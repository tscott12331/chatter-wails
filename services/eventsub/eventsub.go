package eventsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"chatter-wails/internal/api"
	"chatter-wails/internal/user"
	"chatter-wails/internal/util"

	"chatter-wails/services/emote"
	"chatter-wails/services/badge"
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

		nextIndex, word, emote := parseWord(text, i, emoteMap)
		i = nextIndex

		if emote == nil {
			sb.WriteString(word)
		} else {
			textFragment := &AppChatMessageFragment{
				Fragment_type: cmf.Fragment_type,
				Text: sb.String(),
				Cheermote: cmf.Cheermote,
				Emote: nil,
				Mention: cmf.Mention,
			}

			emoteFragment := &AppChatMessageFragment{
				Fragment_type: "emote",
				Text: word,
				Cheermote: cmf.Cheermote,
				Emote: emote,
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





type ESChatSubscriptionSevenTVData struct{
	SevenTVEmotes map[string]*emote.AppEmote
	Enabled bool
}

type ESChatSubscriptionData struct{
	BroadcasterId string
	Channel string
	ChannelBadgeSets *util.SingleWriteMutex[[]api.ApiBadgeSet]
	ChannelEmotes map[string]*emote.AppEmote
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
	ctx context.Context

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

type ChatroomData struct{
	SubId string					`json:"subId"`
	BroadcasterId string			`json:"broadcasterId"`
	ChannelEmotes map[string]*emote.AppEmote `json:"channelEmotes"`
}

type EventSubService struct {
	Ctx context.Context
	Client Client
}

func NewEventSubService() *EventSubService {
	es := &EventSubService{}
	es.Client = Client{
		ChatSubscriptions: *util.NewMutexValue(NewChatSubscriptions()),

		waitingClient: make(chan struct{}),
		sessionIdChan: make(chan *string),
		newSessionIdChan: make(chan *string),

		ready: make(chan bool),
		done: make(chan struct{}),
	}

	return es
}

func (es *EventSubService) Connect() {
	ready := es.Client.ready

	es.Client.ctx = es.Ctx

	esCtx, cancel := context.WithCancel(es.Ctx)
	defer cancel()
	
	newSessionIdChan := es.Client.newSessionIdChan

	for {
		select {
		case r, ok := <-ready:
			if !ok { return }

			if r {
				if !es.Client.connected {
					log.Printf("[Connect]: Connecting to eventsub web server\n\n")
					var err error
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
		if !ok {
			return
		}

		chatMessage := esNotificationToEsChatMessage(notification, sub.Data)
		runtime.EventsEmit(c.ctx, notification.Payload.Subscription.Id, chatMessage)
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

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
// func (c *Client) readPump(ctx context.Context) {
// 	log.Printf("[readPump]: Started read goroutine for eventsub\n\n")
// 	defer func() {
// 		log.Printf("[readPump]: Setting ready to false\n\n")
// 		c.ready <- false
// 	}()
//
// 	// c.conn.SetReadLimit(maxMessageSize)
// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Printf("[readPump]: Context canceled, closing\n\n")
// 			return
// 		default:
// 			_, message, err := c.conn.ReadMessage()
// 			if err != nil {
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 					log.Printf("[readPump]: unexpected close, %v\n\n", err)
// 				} else {
// 					log.Printf("[readPump]: expected close, %v\n\n", err)
// 				}
// 				return
// 			}
// 			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
//
// 			c.handleESMessage(message)
// 		}
// 	}
// }

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// func (c *Client) writePump(ctx context.Context) {
// 	log.Printf("[writePump]: Started write goroutine for eventsub\n\n")
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		log.Printf("[writePump]: Closing connection\n\n")
// 		ticker.Stop()
// 		c.conn.Close()
// 	}()
//
// 	ready := c.ready
// 	for {
// 		select {
// 		case message, ok := <-c.send:
// 			log.Printf("[writePump]: Attempting to write message to eventsub\n\n")
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The hub closed the channel.
// 				log.Printf("[writePump]: Message channel was closed, aborting\n\n")
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
//
// 			w, err := c.conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				log.Printf("[writePump]: An error occurred creating the message writer, aborting\n\n")
// 				return
// 			}
// 			w.Write(message)
// 			log.Printf("[writePump]: Added message to write: %s\n\n", message)
//
// 			// Add queued chat messages to the current websocket message.
// 			n := len(c.send)
// 			for range n {
// 				queued_message := <-c.send
// 				log.Printf("[writePump]: Added queued message to write: %s\n\n", queued_message)
// 				w.Write(newline)
// 				w.Write(queued_message)
// 			}
//
// 			if err := w.Close(); err != nil {
// 				log.Printf("[writePump]: An error occurred sending the message, aborting\n\n")
// 				return
// 			}
// 		case <-ticker.C:
// 			log.Printf("[writePump]: Attempting to send ping message\n\n")
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				log.Printf("[writePump]: An error occurred sending ping message, aborting\n\n")
// 				return
// 			}
// 		case <-c.done:
// 			log.Printf("[writePump]: Connection has been closed, closing\n\n")
// 			return
// 		case r := <-ready:
// 			if !r {
// 				log.Printf("[writePump]: Ready is false, closing\n\n")
// 				return
// 			}
// 		case <-ctx.Done():
// 			log.Printf("[writePump]: Context canceled, closing\n\n")
// 			return
// 		}
// 	}
// }


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

	fmt.Printf("Deleted subscription %s\nres: %+v\n", subId, res)

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
		log.Printf("subscription for %s doesn't exist", channelName)
		return nil
	}

	log.Printf("deleting subscription %+v", sub)

	subId := sub.SubId
	return es.deleteChatSubscription(accessToken, subId, channelName)
}

func (es *EventSubService) deleteChatSubscription(accessToken, subId, channelName string) error {
	es.Client.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		delete((*cs).fromSubId, subId)
		delete((*cs).fromChannelName, channelName)
	})

	return es.DeleteSubscription(accessToken, subId)
}



func (es *EventSubService) CreateChatSubscription(accessToken, userId, channelName string, globalBadgeSets *[]api.ApiBadgeSet) (*ChatroomData, error) {
	// maybe this could return (sub, exists, fetching)
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
			log.Printf("[ConnectToChatroom]: An error occurred fetching the broadcaster info, aborting")
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
			log.Printf("[ConnectToChatroom]: Failed to get channel emotes\n%+v\n\n", err)
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
		},
	}

	// fetch channel badge sets in goroutine
	go es.goGetChannelBadgeSets(accessToken, broadcasterId, subId, globalBadgeSets)

	return newSub, nil
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
