package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"chatter-wails/internal/api"
	"chatter-wails/internal/user"
	"chatter-wails/internal/util"
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

type ESSubscriptionCondition map[string]string

type ESBadge struct {
	Set_id string			`json:"set_id"`
	Id string				`json:"id"`
	Info string				`json:"info"`
}



type AppChatMessageFragment struct{
	Fragment_type string							`json:"type"`
	Text string										`json:"text"`
    Cheermote *struct{
		Prefix string								`json:"prefix"`
		Bits int									`json:"bits"`
		Tier int									`json:"tier"`
	}												`json:"cheermote,omitempty"`
	Emote *AppEmote									`json:"emote"`
    Mention *struct{
		User_id string								`json:"user_id"`
		User_name string							`json:"user_name"`
		User_login string							`json:"user_login"`
	}												`json:"mention,omitempty"`
}

type ESChatMessageFragment struct{
	Fragment_type string							`json:"type"`
	Text string										`json:"text"`
    Cheermote *struct{
		Prefix string								`json:"prefix"`
		Bits int									`json:"bits"`
		Tier int									`json:"tier"`
	}												`json:"cheermote,omitempty"`
    Emote *struct{
		Id string									`json:"id"`
		Emote_set_id string							`json:"emote_set_id"`
		Owner_id string								`json:"owner_id"`
		Format []string								`json:"format"`
	}												`json:"emote,omitempty"`
    Mention *struct{
		User_id string								`json:"user_id"`
		User_name string							`json:"user_name"`
		User_login string							`json:"user_login"`
	}												`json:"mention,omitempty"`
}

var esDarkTheme = "dark"
var esLightTheme = "light"
func esChatMessageFragmentToAppMessageFragment(cmf *ESChatMessageFragment, emoteMap map[string]*AppEmote) []*AppChatMessageFragment {
	var appEmote *AppEmote = nil
	if cmf.Fragment_type == "text" {
		return parseTextFragmentEmotes(cmf, emoteMap)
	}

	if cmf.Emote != nil {
		appEmote = &AppEmote{
			Id: cmf.Emote.Id,
			Name: cmf.Text,
			LightSrcSet: GetEmoteSrcSet(cmf.Emote.Id, cmf.Emote.Format, &esLightTheme, nil),
			DarkSrcSet: GetEmoteSrcSet(cmf.Emote.Id, cmf.Emote.Format, &esDarkTheme, nil),
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

func parseWord(text string, i int, emoteMap map[string]*AppEmote) (int, string, *AppEmote) {
	nextIndex := i
	for nextIndex < len(text) && !unicode.IsSpace(rune(text[nextIndex])) {
		nextIndex += 1
	}
	
	// fmt.Printf("nextIndex %d\n", nextIndex)
	word := text[i:nextIndex]
	emote := emoteMap[word]

	return nextIndex, word, emote
}

func parseTextFragmentEmotes(cmf *ESChatMessageFragment, emoteMap map[string]*AppEmote) []*AppChatMessageFragment {
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


type ESMessageBadge struct{
	SrcSet string									`json:"srcSet"`
	Info string										`json:"info"`
	Title string									`json:"title"`
}

type ESMessageReply struct{
	Parent_message_body string						`json:"parent_message_body"`
	Parent_message_id string						`json:"parent_message_id"`
	Parent_user_id string							`json:"parent_user_id"`
	Parent_user_login string						`json:"parent_user_login"`
	Parent_user_name string							`json:"parent_user_name"`
	Thread_message_id string						`json:"thread_message_id"`
	Thread_user_id string							`json:"thread_user_id"`
	Thread_user_login string						`json:"thread_user_login"`
	Thread_user_name string							`json:"thread_user_name"`
}

type ESChatMessage struct{
	Id string							`json:"id"`
	Username string						`json:"username"`
	Text string							`json:"text"`
	Fragments []*AppChatMessageFragment	`json:"fragments"`
	Color string						`json:"color"`
	Badges []ESMessageBadge				`json:"badges"`
	Reply *ESMessageReply				`json:"reply,omitempty"`
}



/* EVENTSUB EVENT TYPES */

type ESEvent struct {
	Broadcaster_user_id string				`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Chatter_user_id string					`json:"chatter_user_id"`
	Chatter_user_login string				`json:"chatter_user_login"`
	Chatter_user_name string				`json:"chatter_user_name"`
	Message_id string						`json:"message_id"`
	Message struct {
		Text string							`json:"text"`
		Fragments []ESChatMessageFragment	`json:"fragments"`
	}										`json:"message"`
	Color string							`json:"color"`
	Badges []ESBadge						`json:"badges"`
	Message_type string						`json:"message_type"`
	Cheer *struct {
		Bits int							`json:"bits"`
	}										`json:"cheer,omitempty"`
	Reply *ESMessageReply							`json:"reply,omitempty"`
	Channel_points_custom_reward_id *string			`json:"channel_points_custom_reward_id,omitempty"`
	Source_broadcaster_user_id *string				`json:"source_broadcaster_user_id,omitempty"`
	Source_broadcaster_user_name *string			`json:"source_broadcaster_user_name,omitempty"`
	Source_braodcaster_user_login *string			`json:"source_broadcaster_user_login,omitempty"`
	Source_message_id *string						`json:"source_message_id,omitempty"`
	Source_badges *[]ESBadge						`json:"source_badges,omitempty"`
	Is_source_only *bool							`json:"is_source_only,omitempty"`

}

type ESChatMessageEventMessage struct{
	Text string							`json:"text"`
	Fragments []*AppChatMessageFragment	`json:"fragments"`
}
type ESChatMessageEvent = struct{
	Broadcaster_user_id string				`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Chatter_user_id string					`json:"chatter_user_id"`
	Chatter_user_login string				`json:"chatter_user_login"`
	Chatter_user_name string				`json:"chatter_user_name"`
	Message_id string						`json:"message_id"`
	Message ESChatMessageEventMessage		`json:"message"`
	Color string							`json:"color"`
	Badges []ESBadge						`json:"badges"`
	Message_type string						`json:"message_type"`
	Cheer *struct {
		Bits int							`json:"bits"`
	}										`json:"cheer,omitempty"`
	Reply *ESMessageReply							`json:"reply,omitempty"`
	Channel_points_custom_reward_id *string			`json:"channel_points_custom_reward_id,omitempty"`
	Source_broadcaster_user_id *string				`json:"source_broadcaster_user_id,omitempty"`
	Source_broadcaster_user_name *string			`json:"source_broadcaster_user_name,omitempty"`
	Source_braodcaster_user_login *string			`json:"source_broadcaster_user_login,omitempty"`
	Source_message_id *string						`json:"source_message_id,omitempty"`
	Source_badges *[]ESBadge						`json:"source_badges,omitempty"`
	Is_source_only *bool							`json:"is_source_only,omitempty"`
}



/* WELCOME MESSAGE TYPES */

type ESWelcomeMetadata struct {
	Message_id string			`json:"message_id"`
	Message_type string			`json:"message_type"`
	Message_timestamp string 	`json:"message_timestamp"`
}

type ESWelcomePayloadSession struct {
	Id string							`json:"id"`
	Status string						`json:"status"`
	Keepalive_timeout_seconds int		`json:"keepalive_timeout_seconds"`
	Reconnect_url string 				`json:"reconnect_url"`
	Connected_at string 				`json:"connected_at"`
}

type ESWelcomePayload struct {
	Session ESWelcomePayloadSession 	`json:"session"`
}

type ESWelcome struct {
	Metadata ESWelcomeMetadata					`json:"metadata"`
	Payload ESWelcomePayload					`json:"payload"`
}



/* NOTIFICATION TYPES */

type ESNotificationMetadata = ESMessageMetadata

type ESNotificationPayloadSubscription = ESMessagePayloadSubscription

type ESNotificationPayload struct {
	Subscription ESNotificationPayloadSubscription 	`json:"subscription"`
	Event *ESEvent									`json:"event,omitempty"`
}

type ESNotification struct {
	Metadata ESNotificationMetadata		`json:"metadata"`
	Payload ESNotificationPayload		`json:"payload"`
}



type ESChatMessageNotificationPayload struct {
	Subscription ESNotificationPayloadSubscription 	`json:"subscription"`
	Event *ESChatMessageEvent						`json:"event,omitempty"`
}

type ESChatMessageNotification struct {
	Metadata ESNotificationMetadata				`json:"metadata"`
	Payload ESChatMessageNotificationPayload	`json:"payload"`
}



/* EVENTSUB MESSAGE TYPES */

type ESMessageMetadata struct {
	Message_id string			`json:"message_id"`
	Message_type string			`json:"message_type"`
	Message_timestamp string	`json:"message_timestamp"`
	Subscription_type string	`json:"subscription_type"`
	Subscription_version string `json:"subscription_version"`
}

type ESMessagePayloadSubscriptionTransport struct {
	Method string			`json:"method"`
	Session_id string		`json:"session_id"`
}

type ESMessagePayloadSubscription struct {
	Id string												`json:"id"`
	Status string											`json:"status"`
	Sub_type string											`json:"type"`
	Version string											`json:"version"`
	Cost int												`json:"cost"`
	Condition map[string]string								`json:"condition"`
	Transport ESMessagePayloadSubscriptionTransport			`json:"transport"`
	Created_at string										`json:"created_at"`
}

type ESMessagePayload struct {
	Subscription ESMessagePayloadSubscription 	`json:"subscription"`
	Event *ESEvent					`json:"event,omitempty"`
	Session ESWelcomePayloadSession				`json:"session"`
}

// type ESMessage map[string]map[string]any
type ESMessage struct {
	Metadata ESMessageMetadata 			`json:"metadata"`
	Payload ESMessagePayload			`json:"payload"`
}



func esMessageToESNotification(message *ESMessage) *ESNotification {
	return &ESNotification{
		Metadata: message.Metadata,
		Payload: ESNotificationPayload{
			Subscription: message.Payload.Subscription,
			Event: message.Payload.Event,
		},
	}
}




func esBadgesToMessageBadges(esBadges []ESBadge, badgeSets []api.ApiBadgeSet) []ESMessageBadge {
    var messageBadges = []ESMessageBadge{}

	for _, badge := range esBadges {
		setIndex := slices.IndexFunc(badgeSets, func(bs api.ApiBadgeSet) bool {
			return bs.Set_id == badge.Set_id
		})
        if(setIndex == -1) {
			continue
		}

		var set = badgeSets[setIndex]


		versionIndex := slices.IndexFunc(set.Versions, func(v api.ApiBadgeSetVersions) bool {
			return v.Id == badge.Id
		})
        if(versionIndex == -1) {
			continue
		}

        var version = set.Versions[versionIndex]

        var urls = []string{version.Image_url_1x, version.Image_url_2x, version.Image_url_4x};

		
        var srcSet strings.Builder
		srcSet.Grow(252)

		for i, url := range urls {
			var scale int
			if i == 0 {
				scale = 1
			} else {
				scale = 2 << i
			}

			fmt.Fprintf(&srcSet, "%s %dx", url, scale)
			if i < 2 {
				fmt.Fprint(&srcSet, ", ")
			}
		}

		messageBadges = append(messageBadges, ESMessageBadge{
			SrcSet: srcSet.String(),
			Info: badge.Info,
			Title: version.Title,
		})
    }

    return messageBadges;
}

func esNotificationToEsChatMessage(notification *ESNotification, chatSubscriptionData *ESChatSubscriptionData) *ESChatMessage {
	channelBadges, _ := chatSubscriptionData.ChannelBadgeSets.Read()
	seventvEmotes := chatSubscriptionData.SevenTV.SevenTVEmotes
	var fragments = []*AppChatMessageFragment{}
	for _, fragment := range notification.Payload.Event.Message.Fragments {
		fragments = append(fragments, esChatMessageFragmentToAppMessageFragment(&fragment, seventvEmotes)...)
	}

	return &ESChatMessage{
		Id: notification.Payload.Event.Message_id,
		Username: notification.Payload.Event.Chatter_user_name,
		Text: notification.Payload.Event.Message.Text,
		Fragments: fragments,
		Color: notification.Payload.Event.Color,
		Badges: esBadgesToMessageBadges(notification.Payload.Event.Badges, channelBadges),
		Reply: notification.Payload.Event.Reply,
	}
}

func esNotificationToEsChatMessageNotification(notification *ESNotification) *ESChatMessageNotification {
	var fragments = []*AppChatMessageFragment{}
	for _, fragment := range notification.Payload.Event.Message.Fragments {
		fragments = append(fragments, esChatMessageFragmentToAppMessageFragment(&fragment, map[string]*AppEmote{})...)
	}

	return &ESChatMessageNotification{
		Metadata: notification.Metadata,
		Payload: ESChatMessageNotificationPayload{
			Subscription: notification.Payload.Subscription,
			Event: &ESChatMessageEvent{
				Broadcaster_user_id: notification.Payload.Event.Broadcaster_user_id,
				Broadcaster_user_login: notification.Payload.Event.Broadcaster_user_login,
				Broadcaster_user_name: notification.Payload.Event.Broadcaster_user_name,
				Chatter_user_id: notification.Payload.Event.Chatter_user_id,
				Chatter_user_login: notification.Payload.Event.Chatter_user_login,
				Chatter_user_name: notification.Payload.Event.Chatter_user_name,
				Message_id: notification.Payload.Event.Message_id,
				Message: ESChatMessageEventMessage{
					Text: notification.Payload.Event.Message.Text,
					Fragments: fragments,
				},
				Color: notification.Payload.Event.Color,
				Badges: notification.Payload.Event.Badges,
				Message_type: notification.Payload.Event.Message_type,
				Cheer: notification.Payload.Event.Cheer,
				Reply: notification.Payload.Event.Reply,
				Channel_points_custom_reward_id: notification.Payload.Event.Channel_points_custom_reward_id,
				Source_broadcaster_user_id: notification.Payload.Event.Source_broadcaster_user_id,
				Source_broadcaster_user_name: notification.Payload.Event.Source_broadcaster_user_name,
				Source_braodcaster_user_login: notification.Payload.Event.Source_braodcaster_user_login,
				Source_message_id: notification.Payload.Event.Source_message_id,
				Source_badges: notification.Payload.Event.Source_badges,
				Is_source_only: notification.Payload.Event.Is_source_only,
				
			},
		},
	}
}

func esMessageToESWelcome(message *ESMessage) *ESWelcome {
	return &ESWelcome{
		Metadata: ESWelcomeMetadata{
			Message_id: message.Metadata.Message_id,
			Message_type: message.Metadata.Message_type,
			Message_timestamp: message.Metadata.Message_timestamp,
		},
		Payload: ESWelcomePayload{
			Session: message.Payload.Session,
		},
	}
}


type ESChatSubscriptionSevenTVData struct{
	SevenTVEmotes map[string]*AppEmote
	Enabled bool
}

type ESChatSubscriptionData struct{
	BroadcasterId string
	Channel string
	ChannelBadgeSets *util.SingleWriteMutex[[]api.ApiBadgeSet]
	ChannelEmotes map[string]*AppEmote
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

	conn *websocket.Conn
	connected bool

	sessionId *string

	ChatSubscriptions util.MutexValue[*ChatSubscriptions]
	
	waitingClient chan struct{}
	sessionIdChan chan *string
	newSessionIdChan chan *string

	send chan []byte

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
	ChannelEmotes map[string]*AppEmote `json:"channelEmotes"`
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

		send: make(chan []byte),

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
					client, _, err := websocket.DefaultDialer.Dial(twitchESURL.String(), nil)
					if err != nil {
						log.Fatal(err.Error())
					}
					es.Client.conn = client
					es.Client.connected = true

					log.Printf("[Connect]: Starting read and write goroutines\n\n")
					go es.Client.readPump(esCtx)
					go es.Client.writePump(esCtx)
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
func (c *Client) readPump(ctx context.Context) {
	log.Printf("[readPump]: Started read goroutine for eventsub\n\n")
	defer func() {
		log.Printf("[readPump]: Setting ready to false\n\n")
		c.ready <- false
	}()

	// c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		select {
		case <-ctx.Done():
			log.Printf("[readPump]: Context canceled, closing\n\n")
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("[readPump]: unexpected close, %v\n\n", err)
				} else {
					log.Printf("[readPump]: expected close, %v\n\n", err)
				}
				return
			}
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

			c.handleESMessage(message)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(ctx context.Context) {
	log.Printf("[writePump]: Started write goroutine for eventsub\n\n")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Printf("[writePump]: Closing connection\n\n")
		ticker.Stop()
		c.conn.Close()
	}()
	
	ready := c.ready
	for {
		select {
		case message, ok := <-c.send:
			log.Printf("[writePump]: Attempting to write message to eventsub\n\n")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Printf("[writePump]: Message channel was closed, aborting\n\n")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[writePump]: An error occurred creating the message writer, aborting\n\n")
				return
			}
			w.Write(message)
			log.Printf("[writePump]: Added message to write: %s\n\n", message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for range n {
				queued_message := <-c.send
				log.Printf("[writePump]: Added queued message to write: %s\n\n", queued_message)
				w.Write(newline)
				w.Write(queued_message)
			}

			if err := w.Close(); err != nil {
				log.Printf("[writePump]: An error occurred sending the message, aborting\n\n")
				return
			}
		case <-ticker.C:
			log.Printf("[writePump]: Attempting to send ping message\n\n")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[writePump]: An error occurred sending ping message, aborting\n\n")
				return
			}
		case <-c.done:
			log.Printf("[writePump]: Connection has been closed, closing\n\n")
			return
		case r := <-ready:
			if !r {
				log.Printf("[writePump]: Ready is false, closing\n\n")
				return
			}
		case <-ctx.Done():
			log.Printf("[writePump]: Context canceled, closing\n\n")
			return
		}
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

func (es *EventSubService) CreateChatSubscription(accessToken, userId, channelName string, globalBadgeSets *[]api.ApiBadgeSet) (*ChatroomData, error) {
	// maybe this could return (sub, exists, fetching)
	var chatroomData *ChatroomData
	var err error
	es.Client.ChatSubscriptions.Update(func(cs **ChatSubscriptions) {
		log.Printf("in update, %+v\n", *cs)
		data, exists := (*cs).GetChatroomData(channelName)
		log.Printf("channel %+v already fetched? %+v\n", channelName, exists)
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
		// TODO: move this to similar pattern as badge sets
		var channelEmotesDone chan map[string]*AppEmote = make(chan map[string]*AppEmote)
		var channelEmotesErr chan error = make(chan error)
		defer func() {
			close(channelEmotesDone)
			close(channelEmotesErr)
		}()

		ceCtx, ceCancel := context.WithCancel(es.Ctx)
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

	badgeSets, err := GetChannelBadgeSets(accessToken, broadcasterId)
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}
	combinedSets := CombineChannelGlobalSets(badgeSets, globalBadgeSets)

	if !sub.Data.ChannelBadgeSets.Write(*combinedSets) {
		log.Printf("Tried to write badge sets which were already written")
	}
}

func (es *EventSubService) goGetChannelEmotes(
	ctx context.Context, 
	channelEmotesDone chan map[string]*AppEmote, 
	channelEmotesErr chan error, 
	wg *sync.WaitGroup,
	accessToken string,
	broadcasterId string,
) {
	defer wg.Done()
	channelEmotes, err := GetChannelEmotes(accessToken, broadcasterId)
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
		channelEmoteMap := util.ArrToMap(*channelEmotes, func(item AppEmote) (string, *AppEmote) {
			return item.Name, &item
		})
		channelEmotesDone <- channelEmoteMap
		return
	}

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
