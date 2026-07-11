package eventsub

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"

	"chatter-wails/internal/api"
	"chatter-wails/internal/message"
	"chatter-wails/internal/user"
	"chatter-wails/internal/util"

	"chatter-wails/services/auth"
	"chatter-wails/services/badge"
	"chatter-wails/services/emote"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	// maxMessageSize = 512


)

const CHAT_SUB_TYPE = "channel.chat.message"


var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	twitchESURL = url.URL{
		Scheme: "wss",
		Host: "eventsub.wss.twitch.tv",
		Path: "/ws",
	}

	twitchURL = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
	}

	twitchApiURL = url.URL{
		Scheme: twitchURL.Scheme,
		Host: twitchURL.Host,
		Path: "/helix",
	}

	twitchOauthURL = url.URL{
		Scheme: twitchURL.Scheme,
		Host: "id.twitch.tv",
		Path: "/oauth2",
	}
)

var esDarkTheme = "dark"
var esLightTheme = "light"
func esChatMessageFragmentToAppMessageFragment(cmf *ESChatMessageFragment, emoteMap map[string]*emote.AppEmote) []*AppChatMessageFragment {
	var appEmote *emote.AppEmote = nil
	if cmf.Fragment_type == "text" {
		return parseTextFragmentEmotes(cmf, emoteMap)
	}

	if cmf.Emote != nil {
		appEmote = &emote.AppEmote{
			Id: cmf.Emote.Id,
			Name: cmf.Text,
			LightSrcSet: emote.GetEmoteSrcSet(cmf.Emote.Id, cmf.Emote.Format, &esLightTheme, nil),
			DarkSrcSet: emote.GetEmoteSrcSet(cmf.Emote.Id, cmf.Emote.Format, &esDarkTheme, nil),
		}
	}

	return []*AppChatMessageFragment{{
		Fragment_type: cmf.Fragment_type,
		Text: cmf.Text,
		Cheermote: cmf.Cheermote,
		Emote: appEmote,
		Mention: cmf.Mention,
	}}
}

func consumeWhitespace(text string, i int) (int, string) {
	nextIndex := i
	for nextIndex < len(text) && unicode.IsSpace(rune(text[nextIndex])) {
		nextIndex += 1
	}

	return nextIndex, text[i:nextIndex]
}

func parseWord(text string, i int, emoteMap map[string]*emote.AppEmote) (int, string, *emote.AppEmote) {
	nextIndex := i
	for nextIndex < len(text) && !unicode.IsSpace(rune(text[nextIndex])) {
		nextIndex += 1
	}
	
	word := text[i:nextIndex]
	emote := emoteMap[word]

	return nextIndex, word, emote
}

func parseTextFragmentEmotes(cmf *ESChatMessageFragment, emoteMap map[string]*emote.AppEmote) []*AppChatMessageFragment {
	fragments := []*AppChatMessageFragment{}
	
	text := cmf.Text
	i := 0

	var sb strings.Builder

	for i < len(text) {
		nextIndex, spaces := consumeWhitespace(text, i)
		i = nextIndex
		sb.WriteString(spaces)

		nextIndex, word, appEmote := parseWord(text, i, emoteMap)
		i = nextIndex
		if appEmote == nil {
			sb.WriteString(word)
		} else {
			textFragment := &AppChatMessageFragment{
				Fragment_type: cmf.Fragment_type,
				Text: sb.String(),
				Cheermote: cmf.Cheermote,
				Emote: nil,
				Mention: cmf.Mention,
			}

			i, appEmote.EmoteStack = parseZeroWidthEmotes(text, i, emoteMap)

			emoteFragment := &AppChatMessageFragment{
				Fragment_type: "emote",
				Text: word,
				Cheermote: cmf.Cheermote,
				Emote: appEmote,
				Mention: cmf.Mention,
			}

			fragments = append(fragments, textFragment)
			fragments = append(fragments, emoteFragment)

			sb.Reset()
		}
	}

	if sb.Len() > 0 {
		textFragment := &AppChatMessageFragment{
			Fragment_type: cmf.Fragment_type,
			Text: sb.String(),
			Cheermote: cmf.Cheermote,
			Emote: nil,
			Mention: cmf.Mention,
		}
		fragments = append(fragments, textFragment)
	}

	return fragments
}


func parseZeroWidthEmotes(text string, i int, emoteMap map[string]*emote.AppEmote) (int, []*emote.AppEmote){
	emotes := []*emote.AppEmote{}

	indexAfterLast := i
	nextIndex, _ := consumeWhitespace(text, indexAfterLast)

	for {
		ni, _, emote := parseWord(text, nextIndex, emoteMap)
		if emote == nil || !emote.ZeroWidth {
			break
		}

		emotes = append(emotes, emote)
		indexAfterLast = ni
		nextIndex, _ = consumeWhitespace(text, indexAfterLast)

		if nextIndex >= len(text) { break }
	}

	return indexAfterLast, emotes
}


type ESChatSubscriptionSevenTVData struct{
	SevenTVEmotes map[string]*emote.AppEmote
	Enabled bool
}

type SharedChatParticipant struct{
	Name string					`json:"name"`
	ProfileImageURL string		`json:"profileImageURL"`
}

type ESChatSubscriptionData struct{
	ChatOpen bool
	PollCancel context.CancelFunc

	BroadcasterId string
	Channel string
	ChannelBadgeSets *util.SingleWriteMutex[[]api.ApiBadgeSet]
	ChannelEmotes map[string]*emote.AppEmote

	SharedChatSubId string
	SharedChatParticipants map[string]*SharedChatParticipant
	
	SevenTV *ESChatSubscriptionSevenTVData
}
type ESSubscription[T any] struct{
	SubId string
	SubType string
	Data T
}

type ESSubscriptionMap[T any] map[string]*ESSubscription[T]

type ChatSubscriptions struct{
	fromSubId ESSubscriptionMap[*ESChatSubscriptionData]
	fromChannelName ESSubscriptionMap[*ESChatSubscriptionData]
}

func NewChatSubscriptions() *ChatSubscriptions {
	return &ChatSubscriptions{
		fromSubId: make(ESSubscriptionMap[*ESChatSubscriptionData]),
		fromChannelName: make(ESSubscriptionMap[*ESChatSubscriptionData]),
	}
}

func (cs *ChatSubscriptions) SetSubData(data *ESSubscription[*ESChatSubscriptionData]) {
	cs.fromSubId[data.SubId] = data
	cs.fromChannelName[strings.ToLower(data.Data.Channel)] = data
}

func (cs *ChatSubscriptions) GetSubFromChannelName(channelName string) (*ESSubscription[*ESChatSubscriptionData], bool) {
	d, e := cs.fromChannelName[strings.ToLower(channelName)]
	return d, e
}

func (cs *ChatSubscriptions) GetSubFromId(subId string) (*ESSubscription[*ESChatSubscriptionData], bool) {
	d, e := cs.fromSubId[subId]
	return d, e
}

func (cs *ChatSubscriptions) GetChatroomData(channelName string) (*ChatroomData, bool) {
	sub, exists := cs.GetSubFromChannelName(channelName)
	if !exists {
		return nil, false
	}

	data := &ChatroomData{
		SubId: sub.SubId,
		BroadcasterId: sub.Data.BroadcasterId,
		ChannelEmotes: sub.Data.ChannelEmotes,
	}

	return data, true
}


type Client struct {
	app *application.App
	ctx context.Context

	user *auth.AppUser

	socket *util.Socket

	connected bool

	sessionId *string

	ChatSubscriptions util.MutexValue[*ChatSubscriptions]
	
	waitingClient chan struct{}
	sessionIdChan chan *string
	newSessionIdChan chan *string

	ready chan bool
	done chan struct{}
}

func (c *Client) AddChatSubscription(data *ESSubscription[*ESChatSubscriptionData]) {
	c.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		(*cs).SetSubData(data)
	})
}

func (c *Client) ToggleChatSubscriptionFromChannelName(data *ChatOpenData) {
	c.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		sub, exists := (*cs).GetSubFromChannelName(data.Channel)
		if !exists { return }

		if sub.Data.PollCancel != nil {
			sub.Data.PollCancel()
		}

		if data.Open && !sub.Data.ChatOpen {
			pollCxt, pollCancel := context.WithCancel(c.ctx)

			go c.pollChatSubscription(pollCxt, data.AccessToken, data.Channel)

			sub.Data.PollCancel = pollCancel
		}

		sub.Data.ChatOpen = data.Open
	})
}

const CHAT_SUB_POLL_DURATION = 30 * time.Second
func (c *Client) pollChatSubscription(ctx context.Context, accessToken, channel string) {
	ticker := time.NewTicker(CHAT_SUB_POLL_DURATION)

	// poll initial data
	c.pollStreamData(accessToken, channel)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.pollStreamData(accessToken, channel)
		}
	}
}

type StreamData struct{
	Channel string 		`json:"channel"`
	Live bool			`json:"live"`
	ViewCount int		`json:"viewCount"`
	Title string		`json:"title"`
	GameName string 	`json:"gameName"`
}

func (c *Client) pollStreamData(accessToken, channel string) {
	res, err := api.ApiGetStreams(accessToken, map[string][]string{
		"user_login": {channel},
		"first": {"1"},
	})
	if err != nil {
		log.Printf("ERROR: recieved while polling viewcount %v", err)
		return
	}

	var streamData StreamData
	if len(res.Body.Data) > 0 {
		streamDataRes := res.Body.Data[0]
		streamData = StreamData{
			Channel: channel,
			Live: streamDataRes.Type == "live",
			ViewCount: streamDataRes.Viewer_count,
			Title: streamDataRes.Title,
			GameName: streamDataRes.Game_name,
		}
	}

	c.app.Event.Emit("common:stream-data", streamData)
}

type ChatroomData struct{
	SubId string					`json:"subId"`
	BroadcasterId string			`json:"broadcasterId"`
	ChannelEmotes map[string]*emote.AppEmote `json:"channelEmotes"`
}

type EventSubService struct {
	app *application.App
	Ctx context.Context
	Client Client
}

func NewEventSubService(app *application.App) *EventSubService {
	es := &EventSubService{app: app}
	es.Client = Client{
		app: app,

		ChatSubscriptions: *util.NewMutexValue(NewChatSubscriptions()),

		waitingClient: make(chan struct{}),
		sessionIdChan: make(chan *string),
		newSessionIdChan: make(chan *string),

		ready: make(chan bool),
		done: make(chan struct{}),
	}

	return es
}

type ChatOpenData struct{
	Channel string		`json:"channel"`
	AccessToken string 	`json:"accessToken"`
	Open bool			`json:"open"`
}

func (es *EventSubService) handleChatOpenEvent(event *application.CustomEvent) {
	data, ok := event.Data.(ChatOpenData)
	if !ok {
		log.Printf("ERROR: failed to cast chat open data\nData: %+v\n", data)
		return
	}

	es.Client.ToggleChatSubscriptionFromChannelName(&data)
}

func (es *EventSubService) handleUserLoginEvent(event *application.CustomEvent) {
	data, ok := event.Data.(*auth.AppUser)
	if !ok {
		log.Printf("ERROR: failed to cast user login data\nData: %+v\n", data)
		return
	}

	es.Client.user = data
}

func (es *EventSubService) Connect() {
	ready := es.Client.ready

	es.Client.ctx = es.Ctx

	esCtx, cancel := context.WithCancel(es.Ctx)
	defer cancel()
	
	newSessionIdChan := es.Client.newSessionIdChan

	es.app.Event.On("common:user-login", es.handleUserLoginEvent)

	for {
		select {
		case r, ok := <-ready:
			if !ok { return }

			if r {
				if !es.Client.connected {
					log.Printf("[Connect]: Connecting to eventsub web server\n\n")
					var err error
					es.app.Event.On("common:chat-open", es.handleChatOpenEvent)
					es.Client.socket, err = util.NewSocket(esCtx, twitchESURL.String(), es.Client.handleESMessage)
					if err != nil {
						log.Fatal(err.Error())
					}

					es.Client.connected = true

					log.Printf("[Connect]: Starting read and write goroutines\n\n")
				}
			} else {
				log.Printf("[Connect]: Ready is false, setting connected to false\n\n")
				es.Client.connected = false
				es.app.Event.Off("common:chat-open")
			}

		case id := <-newSessionIdChan:
			log.Printf("[Connect]: New session ID given, send out id to waiting clients\n\n")
			for range es.Client.waitingClient {
				log.Printf("[Connect]: Sent out session ID to waiting client\n\n")
				es.Client.sessionIdChan <- id
			}

		case <-es.Ctx.Done():
			log.Printf("[Connect]: Parent context cancelled, aborting\n\n")
			es.Client.connected = false
			es.app.Event.Off("chatopen")
			return
		case <-esCtx.Done():
			es.Client.ready <- false
		}
	}

}

func (c *Client) handleESNotification(message ESMessage) {
	notification := esMessageToESNotification(&message)
	if notification == nil {
		log.Printf("[handleESNotification]: Converted welcome message is nil, aborting\n\n")
		return
	}

	sub_type := notification.Metadata.Subscription_type

	switch sub_type {
	case "channel.chat.message":
		subId := notification.Payload.Subscription.Id

		sub, ok := c.ChatSubscriptions.Read().GetSubFromId(subId)
		if !ok || !sub.Data.ChatOpen {
			return
		}

		chatMessage := esNotificationToEsChatMessage(notification, sub.Data)
		c.app.Event.Emit("common:chat-message", chatMessage)
	case "channel.shared_chat.begin":
		c.handleSharedChatBegin(notification)
	}
}

func (c *Client) handleSharedChatBegin(notification *ESNotification) {
		var sharedChatBeginEvent ESSharedChatBeginEvent
		json.Unmarshal(*notification.Payload.Event, &sharedChatBeginEvent)
		c.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
			sub, exists := (*cs).GetSubFromChannelName(sharedChatBeginEvent.Host_broadcaster_user_login)
			if !exists { return }
			
			sub.Data.SharedChatSubId = notification.Payload.Subscription.Id
			go c.fetchSharedChatProfileImages(&sharedChatBeginEvent)
		})
}

func (c *Client) fetchSharedChatProfileImages(sharedChatBeginEvent *ESSharedChatBeginEvent) {
	var wg sync.WaitGroup
	for _, participant := range sharedChatBeginEvent.Participants {
		wg.Add(1)
		go c.fetchAndSetSharedParticipantProfileImage(sharedChatBeginEvent.Host_broadcaster_user_login, participant.Broadcaster_user_id, &wg)
	}

	wg.Wait()

	sub, exists := c.ChatSubscriptions.Read().GetSubFromChannelName(sharedChatBeginEvent.Host_broadcaster_user_name)
	if !exists { return }

	data := sub.Data.SharedChatParticipants
	c.app.Event.Emit("common:shared-chat-begin", data)
}

func (c *Client) fetchAndSetSharedParticipantProfileImage(mainBroadcaster string, participant_user_id string, wg *sync.WaitGroup) {
	defer wg.Done()

	if c.user == nil { return }
	user, err := user.GetUserById(c.user.Access_token, participant_user_id)
	if err == nil {
		c.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
			sub, exists := (*cs).GetSubFromChannelName(mainBroadcaster)
			if !exists { return }

			sub.Data.SharedChatParticipants[user.Login] = &SharedChatParticipant{
				Name: user.Display_name,
				ProfileImageURL: user.Profile_image_url,
			}
		})
	}
}

func (c *Client) handleESWelcome(message ESMessage) {
	welcome := esMessageToESWelcome(&message)
	if welcome == nil {
		log.Printf("[handleESWelcome]: Converted welcome message is nil, aborting\n\n")
		return
	}

	id := welcome.Payload.Session.Id
	log.Printf("[handleESWelcome]: Setting session ID to %v\n\n", id)
	c.sessionId = &id
	c.newSessionIdChan <- c.sessionId
}


func (c *Client) handleESMessage(message []byte) {
	var message_obj ESMessage
	err := json.Unmarshal(message, &message_obj)
	if err != nil {
		log.Printf("[handleESMessage]: An error occurred parsing the eventsub message, aborting\n%v\n\n", err)
		return
	}
	// if message_obj["metadata"] == nil || message_obj["payload"] == nil || message_obj["metadata"]["message_type"] == nil {
	// 	log.Printf("[handleESMessage]: Unmarshalled object does not contain metadata, payload, or message_type fields, aborting\n\n")
	// 	return
	// }

	mType := message_obj.Metadata.Message_type

	switch mType {
	case "session_welcome":
		c.handleESWelcome(message_obj)
	case "notification":
		c.handleESNotification(message_obj)
	case "session_keepalive":
		log.Printf("[handleESMessage]: session_keepalive recieved\n\n")
	default:
		log.Printf("[handleESMessage]: not handling this message type\n\n")
	}

}

type CastError struct {}
func (ce *CastError) Error() string {
	return "An error occurred casting data types"
}

type ESNoConnError struct {}
func (ese *ESNoConnError) Error() string {
	return "A connection to the websocket server has not yet been established"
}

type ESSessionReqTimeout struct{}
func (esrt *ESSessionReqTimeout) Error() string {
	return "The request for a session ID timed out"
}

func (es *EventSubService) GetSubscriptions(accessToken string) ([]api.ApiSubscription, error) {
	res, err := api.ApiGetSubscriptions(accessToken, map[string][]string{})
	if err != nil {
		log.Printf("[GetSubscriptions]: An error occurred trying to get subscriptions, aborting\n\n")
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetSubscriptions]: Failed to get subscriptions, aborting\n\n")
		return nil, &api.StatusError[api.ApiGetSubscriptionsRes]{Res: res}
	}

	return res.Body.Data, nil
}

func (es *EventSubService) CreateSubscription(accessToken string, condition ESSubscriptionCondition, subType string) (string, error) {
	if !es.Client.connected {
		log.Printf("[CreateSubscription]: Client not yet connected, signal ready to connect\n\n")
		es.Client.ready <- true // signal ready to connect
	}

	var sessionId string
	if es.Client.sessionId == nil {
		sessionId = ""
	} else {
		sessionId = *es.Client.sessionId
	}

	if sessionId == "" {
		es.Client.waitingClient <- struct{}{} // tell main thread that we are waiting
		log.Printf("[CreateSubscription]: EventSub Service not yet connected to websocket server, waiting\n\n")

		var idPtr *string
		select {
		case idPtr = <-es.Client.sessionIdChan:
			if idPtr == nil {
				return "", &ESNoConnError{}
			}

			sessionId = *idPtr
		case <-time.After(5 * time.Second):
			log.Printf("[CreateSubscription]: Request for sessionId timed out, aborting\n\n")
			return "", &ESSessionReqTimeout{}
		}

		log.Printf("[CreateSubscription]: EventSub Service has connected\n\n")
	}


	req_body := api.ApiPostSubscriptionsBody{
		Sub_type: subType,
		Version: "1",
		Condition: condition,
		Transport: api.ApiPostSubscriptionsBodyTransport{
			Method: "websocket",
			Session_id: sessionId,
		},
	}

	res, err := api.ApiPostSubscriptions(accessToken, req_body, map[string][]string{})
	if err != nil {
		log.Printf("[CreateSubscription]: An error occurred while making request, aborting\n\n")
		return "", err
	}
	
	if res.Status != 202 {
		log.Printf("[CreateSubscription]: Failed to make subscription\n\n")
		return "", &api.StatusError[api.ApiPostSubscriptionsRes]{Res: res}
	}

	if len(res.Body.Data) == 0 {
		log.Panic("[CreateSubscription]: Twitch API and chatter data types are out of sync\n\n")
	}

	subId := res.Body.Data[0].Id
	log.Printf("[CreateSubscription]: Recieved subscription ID %v\n\n", subId)

	return subId, nil
	
}


func (es *EventSubService) DeleteSubscription(accessToken string, subId string) (error) {
	res, err := api.ApiDeleteSubscriptions(accessToken, map[string][]string{
		"id": {subId},
	})

	if err != nil {
		log.Printf("[DeleteSubscription]: An error occurred trying to delete a subscription, aborting\n\n")
		return err
	}

	if res.Status != 204 {
		log.Printf("[DeleteSubscription]: Failed to delete subscription\n\n")
		return &api.StatusError[any]{Res: res}
	}

	return nil
}

func (es *EventSubService) DeleteAllSubscriptions(accessToken string) error {
	res, err := es.GetSubscriptions(accessToken)
	if err != nil {
		log.Printf("[DeleteAllSubscriptions]: An error occurred getting the current subscriptions, aborting\n\n")
		return err
	}

	for _, v := range res {
		err := es.DeleteSubscription(accessToken, v.Id)
		if err != nil {
			log.Printf("[DeleteAllSubscriptions]: An error occured deleting subscription %v", v.Id)
		}
	}

	return err
}

func (es *EventSubService) DeleteChatSubscriptionFromSubId(accessToken, subId string) error {
	sub, exists := es.Client.ChatSubscriptions.Read().GetSubFromId(subId)
	if !exists {
		return nil
	}

	channelName := sub.Data.Channel
	return es.deleteChatSubscription(accessToken, subId, channelName)
}

func (es *EventSubService) DeleteChatSubscriptionFromChannelName(accessToken, channelName string) error {
	sub, exists := es.Client.ChatSubscriptions.Read().GetSubFromChannelName(channelName)
	if !exists {
		return nil
	}

	subId := sub.SubId
	return es.deleteChatSubscription(accessToken, subId, channelName)
}

func (es *EventSubService) deleteChatSubscription(accessToken, subId, channelName string) error {
	es.Client.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		sub := (*cs).fromSubId[subId]

		// delete auxiliary subscriptions
		go es.DeleteSubscription(accessToken, sub.Data.SharedChatSubId)

		pollCancel := sub.Data.PollCancel
		if pollCancel != nil { pollCancel() }

		delete((*cs).fromSubId, subId)
		delete((*cs).fromChannelName, channelName)
	})

	return es.DeleteSubscription(accessToken, subId)
}

func (es *EventSubService) SendChatMessageFromChannelName(accessToken, userId, channelName string, messageContent string, replyId *string) (*api.ApiPostMessagesData, error) {
	subData, ok := es.Client.ChatSubscriptions.Read().GetSubFromChannelName(channelName)
	if !ok {
		log.Printf("[SendChatMessage]: Failed to find chat subscription data, aborting\n\n")
		return nil, errors.New("Failed to find chat subscription data")
	}

	res, err := message.SendMessage(userId, accessToken, subData.Data.BroadcasterId, messageContent, replyId)
	if err != nil {
		log.Printf("[SendChatMessage]: An error occurred sending the chat message, aborting\n%+v\n\n", err)
		return nil, err
	}

	return res, nil
}


func (es *EventSubService) createSharedChatSubscription(accessToken, broadcasterId, subId string) {
	sharedSubId, err := es.CreateSubscription(accessToken, map[string]string{
		"broadcaster_user_id": broadcasterId,
	}, "channel.shared_chat.begin")
	if err != nil {
		log.Printf("ERROR: Failed to subscribe to shared chat begin event")
		return
	}

	es.Client.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		sub, exists := (*cs).GetSubFromId(subId)
		if !exists {
			log.Printf("ERROR: Cannot create shared chat subscription on non existant chat subscription")
			return
		}

		sub.Data.SharedChatSubId = sharedSubId
	})
}

func (es *EventSubService) createAuxiliaryChatSubscriptions(accessToken, broadcasterId, subId string) {
	go es.createSharedChatSubscription(accessToken, broadcasterId, subId)
}

func (es *EventSubService) CreateChatSubscription(accessToken, userId, channelName string, globalBadgeSets *[]api.ApiBadgeSet) (*ChatroomData, error) {
	var chatroomData *ChatroomData
	var err error
	es.Client.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		data, exists := (*cs).GetChatroomData(channelName)
		if exists {
			chatroomData = data
			return
		}

		broadcaster, broadcasterErr := user.GetUserByLogin(accessToken, channelName)
		if broadcasterErr != nil {
			log.Printf("[CreateChatSubscription]: An error occurred fetching the broadcaster info, aborting")
			err = broadcasterErr
			return
		}

		// fetch channel emotes
		var channelEmotesDone chan map[string]*emote.AppEmote = make(chan map[string]*emote.AppEmote)
		var channelEmotesErr chan error = make(chan error)
		defer func() {
			close(channelEmotesDone)
			close(channelEmotesErr)
		}()

		ceCtx, ceCancel := context.WithTimeout(es.Ctx, 15 * time.Second)
		defer ceCancel()

		var wg sync.WaitGroup

		// fetch channel emotes in goroutine
		wg.Add(1)
		go es.goGetChannelEmotes(ceCtx, channelEmotesDone, channelEmotesErr, &wg, accessToken, broadcaster.Id)

		// fetch
		sub, fetchErr := es.fetchAndInitChatSubscription(accessToken, userId, channelName, broadcaster.Id, globalBadgeSets)
		err = fetchErr
		if err != nil {
			wg.Wait()
			return
		}

		chatroomData = &ChatroomData{
			SubId: sub.SubId,
			BroadcasterId: sub.Data.BroadcasterId,
		}
		select {
		case channelEmotes := <-channelEmotesDone:
			chatroomData.ChannelEmotes = channelEmotes
			sub.Data.ChannelEmotes = channelEmotes
		case err := <-channelEmotesErr:
			log.Printf("[CreateChatSubscription]: Failed to get channel emotes\n%+v\n\n", err)
		}

		(*cs).SetSubData(sub)
	})

	return chatroomData, err
}

func (es *EventSubService) fetchAndInitChatSubscription(accessToken, userId, channelName, broadcasterId string, globalBadgeSets *[]api.ApiBadgeSet) (*ESSubscription[*ESChatSubscriptionData], error){
	condition := map[string]string{
		"broadcaster_user_id": broadcasterId,
		"user_id": userId,
	}
	subId, err := es.CreateSubscription(accessToken, condition, CHAT_SUB_TYPE)
	if err != nil {
		log.Printf("[ConnectToChatroom]: ERROR: %+v\n", err)
		return nil, err
	}

	newSub := &ESSubscription[*ESChatSubscriptionData]{
		SubId: subId,
		SubType: CHAT_SUB_TYPE,
		Data: &ESChatSubscriptionData{
			BroadcasterId: broadcasterId,
			Channel: channelName,
			ChannelBadgeSets: &util.SingleWriteMutex[[]api.ApiBadgeSet]{},
			// ChannelEmotes: chatroomData.ChannelEmotes,
			SevenTV: &ESChatSubscriptionSevenTVData{},
			SharedChatParticipants: map[string]*SharedChatParticipant{},
		},
	}

	// fetch channel badge sets in goroutine
	go es.goGetChannelBadgeSets(accessToken, broadcasterId, subId, globalBadgeSets)
	go es.goGetSharedChatSession(accessToken, broadcasterId, channelName)
	es.createAuxiliaryChatSubscriptions(accessToken, broadcasterId, subId)

	return newSub, nil
}

func (es *EventSubService) goGetSharedChatSession(
	accessToken,
	broadcasterId,
	channelName string,
) {
	res, err := api.ApiGetSharedChatSession(accessToken, map[string][]string{
		"broadcaster_id": { broadcasterId },
	})

	if err != nil {
		log.Printf("ERROR: failed to fetch initial shared chat session data: %+v", err)
		return
	}

	// no shared chat
	if len(res.Body.Data) == 0 { return }

	data := res.Body.Data[0]

	var wg sync.WaitGroup
	for _, apiParticipant := range data.Participants {
		wg.Add(1)
		go es.Client.fetchAndSetSharedParticipantProfileImage(channelName, apiParticipant.Broadcaster_id, &wg)
	}
	
	wg.Wait()

	sub, exists := es.Client.ChatSubscriptions.Read().GetSubFromChannelName(channelName)
	if !exists { return }

	sharedChatEventData := sub.Data.SharedChatParticipants
	es.app.Event.Emit("common:shared-chat-begin", sharedChatEventData)

}

func (es *EventSubService) goGetChannelBadgeSets(
	accessToken string,
	broadcasterId string,
	subId string,
	globalBadgeSets *[]api.ApiBadgeSet,
) {
	sub, exists := es.Client.ChatSubscriptions.Read().GetSubFromId(subId)
	if !exists {
		log.Printf("Subscription %v doesn't exist\n", subId)
		return
	}

	// data already fetched
	if sub.Data.ChannelBadgeSets.IsWritten() {
		log.Printf("Badge sets already written")
		return
	}

	badgeSets, err := badge.GetChannelBadgeSets(accessToken, broadcasterId)
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}
	combinedSets := badge.CombineChannelGlobalSets(badgeSets, globalBadgeSets)

	if !sub.Data.ChannelBadgeSets.Write(*combinedSets) {
		log.Printf("Tried to write badge sets which were already written")
	}
}

func (es *EventSubService) goGetChannelEmotes(
	ctx context.Context, 
	channelEmotesDone chan map[string]*emote.AppEmote, 
	channelEmotesErr chan error, 
	wg *sync.WaitGroup,
	accessToken string,
	broadcasterId string,
) {
	defer wg.Done()
	channelEmotes, err := emote.GetChannelEmotes(accessToken, broadcasterId)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		default:
			channelEmotesErr <- err
			return
		}
	}
	select {
	case <-ctx.Done():
		log.Printf("[ConnectToChatroom]: channel emote context closed")
		return
	default:
		channelEmoteMap := util.ArrToMap(*channelEmotes, func(item emote.AppEmote) (string, *emote.AppEmote) {
			return item.Name, &item
		})
		channelEmotesDone <- channelEmoteMap
		return
	}

}
