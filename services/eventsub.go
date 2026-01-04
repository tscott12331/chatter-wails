package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	
	"chatter-wails/internal/api"
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

type User struct {
	Id string					`json:"id"`
    Login string				`json:"login"`
    Display_name string			`json:"display_name"`
    User_type string			`json:"type"`
    Broadcaster_type string		`json:"broadcaster_type"`
    Description string			`json:"description"`
    Profile_image_url string	`json:"profile_image_url"`
    Offline_image_url string	`json:"offline_image_url"`
    View_count int				`json:"view_count"`
    Created_at time.Time		`json:"created_at"`
    Access_token string			`json:"access_token"`
}

type ESSubscriptionCondition map[string]string

type ESBadge struct {
	Set_id string			`json:"set_id"`
	Id string				`json:"id"`
	Info string				`json:"info"`
}

type ESChatMessageFragment struct {
	Message_type string			`json:"type"`
	Text string					`json:"text"`
	Cheermote struct {
		Prefix string			`json:"prefix"`
		Bits int				`json:"bits"`
		Tier int				`json:"tier"`
	}							`json:"cheermote,omitempty"`
	Emote *struct{
		Id string						`json:"id"`
		Emote_set_id string				`json:"emote_set_id"`
		Owner_id string					`json:"owner_id"`
		Format []string					`json:"format"`
	}									`json:"emote,omitempty"`
	Mention *struct{
		User_id string			`json:"user_id"`
		User_name string		`json:"user_name"`
		User_login string		`json:"user_login"`
	}							`json:"mention,omitempty"`

}

type ESMessageReply struct {
	Parent_message_body string		`json:"parent_message_body"`
	Parent_message_id string		`json:"parent_message_id"`
	Parent_user_id string			`json:"parent_user_id"`
	Parent_user_login string		`json:"parent_user_login"`
	Parent_user_name string			`json:"parent_user_name"`
	Thread_message_id string		`json:"thread_message_id"`
	Thread_user_id string			`json:"thread_user_id"`
	Thread_user_login string		`json:"thread_user_login"`
	Thread_user_name string			`json:"thread_user_name"`
}

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
		Fragments ESChatMessageFragment		`json:"fragments"`
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

type ESMessage map[string]map[string]any

// type ESNotification struct {
// 	Metadata struct {
// 		Message_id string			`json:"message_id"`
// 		Message_type string			`json:"message_type"`
// 		Message_timestamp string	`json:"message_timestamp"`
// 		Subscription_type string	`json:"subscription_type"`
// 		Subscription_version string `json:"subscription_version"`
// 	}								`json:"metadata"`
// 	Payload struct {
// 		Subscription struct {
// 			Id string					`json:"id"`
// 			Status string				`json:"status"`
// 			Sub_type string				`json:"type"`
// 			Version string				`json:"version"`
// 			Cost int					`json:"cost"`
// 			Condition map[string]string `json:"condition"`
// 			Transport struct {
// 				Method string			`json:"method"`
// 				Session_id string		`json:"session_id"`
// 			}							`json:"transport"`
// 			Created_at string			`json:"created_at"`
// 		}								`json:"subscription"`
// 		Event ESEvent					`json:"event"`
// 	}
// }

// type ESWelcomeMessage struct {
// 	Metadata map[string]string			`json:"metadata"`
// 	Payload struct {
// 		Session struct {
// 			Id string							`json:"id"`
// 			Status string						`json:"status"`
// 			Keepalive_timeout_seconds string	`json:"keepalive_timeout_seconds"`
// 			Reconnect_url string 				`json:"reconnect_url"`
// 			Connected_at string 				`json:"connected_at"`
// 		}										`json:"session"`
// 	}											`json:"payload"`
// }

// type ESWelcomeMessage struct {
// 	Metadata struct {
// 		Message_id string			`json:"message_id"`
// 		Message_type string			`json:"message_type"`
// 		Message_timestamp string 	`json:"message_timestamp"`
// 	}								`json:"metadata"`
// 	Payload struct {
// 		Session struct {
// 			Id string							`json:"id"`
// 			Status string						`json:"status"`
// 			Keepalive_timeout_seconds string	`json:"keepalive_timeout_seconds"`
// 			Reconnect_url string 				`json:"reconnect_url"`
// 			Connected_at string 				`json:"connected_at"`
// 		}										`json:"session"`
// 	}											`json:"payload"`
// }


type ESSubscriptionHandler func(message string)

type Client struct {
	conn *websocket.Conn
	connected bool

	sessionId *string
	subscriptions map[string][]string

	waitingClient chan struct{}
	sessionIdChan chan *string
	newSessionIdChan chan *string

	send chan []byte

	ready chan bool
	done chan struct{}
}

type EventSubService struct {
	Ctx context.Context
	Client Client
}

func NewEventSubService() *EventSubService {
	es := &EventSubService{}
	es.Client = Client{
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
					log.Printf("[Connect]: Initializing eventsub client\n\n")
					client, _, err := websocket.DefaultDialer.Dial(twitchESURL.String(), nil)
					if err != nil {
						log.Fatal(err.Error())
					}
					es.Client.conn = client
					es.Client.connected = true

					log.Printf("[Connect]: Set client connection\n\n")
					log.Printf("[Connect]: Cancelled read and write goroutines\n\n")

					log.Printf("[Connect]: Starting read and write goroutines\n\n")
					go es.Client.readPump(esCtx)
					go es.Client.writePump(esCtx)
				}
			} else {
				log.Printf("[Connect]: Ready is false, cancelling read and write goroutines\n\n")
				cancel()
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
	log.Printf("[handleESNotification]: %v\n\n", message)
}

func (c *Client) handleESWelcome(message ESMessage) {
	if s, ok := message["payload"]["session"].(map[string]any); ok {
		id := s["id"]
		if id, ok := id.(string); ok {
			log.Printf("[handleESWelcome]: Setting session ID to %v\n\n", id)
			c.sessionId = &id
			c.newSessionIdChan <- c.sessionId
		} else {
			log.Printf("[handleESWelcome]: session_welcome message has incorrect type for id, aborting\n\n")
		}
	} else {
		log.Printf("[handleESWelcome]: session_welcome message does not contain session map, aborting\n\n")
	}
}


func (c *Client) handleESMessage(message []byte) {
	var message_obj ESMessage
	err := json.Unmarshal(message, &message_obj)
	if err != nil {
		log.Printf("[handleESMessage]: An error occurred parsing the eventsub message, aborting\n\n")
		return
	}
	if message_obj["metadata"] == nil || message_obj["payload"] == nil || message_obj["metadata"]["message_type"] == nil {
		log.Printf("[handleESMessage]: Unmarshalled object does not contain metadata, payload, or message_type fields, aborting\n\n")
		return
	}

	mType := message_obj["metadata"]["message_type"]

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
		log.Printf("[readPump]: Notifying done channel\n\n")
		c.done <- struct{}{}
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
			log.Printf("[readPump]: RECIEVED MESSAGE\n\n")

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
			for i := 0; i < n; i++ {
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
		case <-ctx.Done():
			log.Printf("[writePump]: Context canceled, closing\n\n")
			return
		}
	}
}

type CreateSubscriptionRes struct{
	Data []struct{
		Id string 				`json:"id"`
		Status string			`json:"status"`
		Sub_type string			`json:"type"`
		Version string			`json:"version"`
		Condition struct{}		`json:"condition"`
		Created_at string		`json:"created_at"`
		Transport struct{
			Method string 		`json:"method"`
			Session_id string 	`json:"session_id"`
			Connected_at string	`json:"connected_at"`
		}						`json:"transport"`
		Cost int				`json:"cost"`

	}							`json:"data"`
	Total int					`json:"total"`
	Total_cost int				`json:"total_cost"`
	Max_total_cost int			`json:"max_total_cost"`
}

type CastError struct {}
func (ce *CastError) Error() string {
	return "An error occurred casting data types"
}

type StatusError struct {
	Res *api.APIResponse
}
func (se *StatusError) Error() string {
	return fmt.Sprintf("%v, %v", se.Res.Status, se.Res)
}

type ESNoConnError struct {}
func (ese *ESNoConnError) Error() string {
	return "A connection to the websocket server has not yet been established"
}

type ESSessionReqTimeout struct{}
func (esrt *ESSessionReqTimeout) Error() string {
	return "The request for a session ID timed out"
}

func (es *EventSubService) CreateSubscription(user User, condition ESSubscriptionCondition, subType string) (*string, error) {
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
				return nil, &ESNoConnError{}
			}

			sessionId = *idPtr
		case <-time.After(5 * time.Second):
			log.Printf("[CreateSubscription]: Request for sessionId timed out, aborting\n\n")
			return nil, &ESSessionReqTimeout{}
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

	res, err := api.ApiPostSubscriptions(user.Access_token, req_body, map[string][]string{})
	if err != nil {
		log.Printf("[CreateSubscription]: An error occurred while making request, aborting\n\n")
		return nil, err
	}
	
	if res.Status != 202 {
		log.Printf("[CreateSubscription]: Failed to make subscription\n\n")
		return nil, &StatusError{Res: res}
	}

	var res_body_obj, ok = res.Body.(CreateSubscriptionRes)
	if !ok {
		log.Printf("[CreateSubscription]: An error occurred casting the response body, aborting\n\n")
		return nil, &CastError{}
	}

	if len(res_body_obj.Data) == 0 {
		log.Panic("[CreateSubscription]: Twitch API and chatter data types are out of sync\n\n")
	}

	subId := res_body_obj.Data[0].Id
	log.Printf("[CreateSubscription]: Recieved subscription ID %v\n\n", subId)

	subList, ok := es.Client.subscriptions[subType]
	if ok {
		// list for this subscription type already exists, append
		subList = append(subList, subId)
	} else {
		subList = []string{subId}
	}

	return &subId, nil
	
}
