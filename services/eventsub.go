package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512


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

type ESNotification struct {
	Metadata struct {
		Message_id string			`json:"message_id"`
		Message_type string			`json:"message_type"`
		Message_timestamp string	`json:"message_timestamp"`
		Subscription_type string	`json:"subscription_type"`
		Subscription_version string `json:"subscription_version"`
	}								`json:"metadata"`
	Payload struct {
		Subscription struct {
			Id string					`json:"id"`
			Status string				`json:"status"`
			Sub_type string				`json:"type"`
			Version string				`json:"version"`
			Cost int					`json:"cost"`
			Condition map[string]string `json:"condition"`
			Transport struct {
				Method string			`json:"method"`
				Session_id string		`json:"session_id"`
			}							`json:"transport"`
			Created_at string			`json:"created_at"`
		}								`json:"subscription"`
		Event ESEvent					`json:"event"`
	}
}

type ESSubscriptionHandler func(message string)

type ESSubscription struct {
	SubId string				`json:"subId"`
	SubType string				`json:"subType"`
}

type Client struct {
	conn *websocket.Conn

	sessionId *string
	subscriptions map[string][]string

	send chan []byte

	done chan struct{}
}

type EventSubService struct {
	Ctx context.Context
	Client Client
}

func NewEventSubService() *EventSubService {
	es := &EventSubService{}
	es.Client.initClient(twitchESURL)
	return es
}

func (es *EventSubService) Connect() {
	log.Printf("Connecting to eventsub web server\n")
	go es.Client.readPump()
	go es.Client.writePump()
}

func (c *Client) initClient(u url.URL) {
	log.Printf("Initializing eventsub client\n")
	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.conn = client
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	log.Printf("Started read goroutine for eventsub\n")
	defer close(c.done)

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Printf("Recieved message: %s\n", message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	log.Printf("Started write goroutine for eventsub\n")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			log.Printf("Attempting to write message to eventsub\n")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Printf("Message channel was closed, aborting\n")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("An error occurred creating the message writer, aborting\n")
				return
			}
			w.Write(message)
			log.Printf("Added message to write: %s\n", message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				queued_message := <-c.send
				log.Printf("Added queued message to write: %s\n", queued_message)
				w.Write(newline)
				w.Write(queued_message)
			}

			if err := w.Close(); err != nil {
				log.Printf("An error occurred sending the message, aborting\n")
				return
			}
		case <-ticker.C:
			log.Printf("Attempting to send ping message\n")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("An error occurred sending ping message, aborting\n")
				return
			}
		case <-c.done:
			log.Printf("Connection has been closed, closing\n")
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

type StatusError struct {
	Code string
}

func (se *StatusError) Error() string {
	return fmt.Sprintf("BAD REQUEST STATUS: %s", se.Code)
}

type ESNoConnError struct {}
func (ese *ESNoConnError) Error() string {
	return "A connection to the websocket server has not yet been established"
}

func (es *EventSubService) CreateSubscription(user User, condition ESSubscriptionCondition, subType string) (*ESSubscription, error) {
	if es.Client.sessionId == nil {
		log.Printf("EventSub Service not yet connected to websocket server, aborting\n\n")
		return nil, &ESNoConnError{}
	}

	// log.Printf("condition: %v\n\n", condition)
	sub_url := url.URL{
		Scheme: twitchApiURL.Scheme,
		Host: twitchApiURL.Host,
		Path: twitchApiURL.Path + "/eventsub/subscriptions",
	}

	log.Printf("POST %s\n\n", sub_url.String())
	
	req_body := struct {
		Sub_type string				`json:"type"`
		Version string				`json:"version"`
		Condition map[string]string	`json:"condition"`
		Transport struct {
			Method string			`json:"method"`
			Session_id string		`json:"session_id"`
		}							`json:"transport"`
	}{
		Sub_type: subType,
		Version: "1",
		Condition: condition,
		Transport: struct{Method string `json:"method"`; Session_id string `json:"session_id"`}{
			Method: "websocket",
			Session_id: *es.Client.sessionId,
		},
	}

	req_body_json, err := json.Marshal(req_body)
	if err != nil {
		log.Printf("An error occurred marshaling the request body, aborting\n")
		return nil, err
	}

	req, err := http.NewRequest("POST", sub_url.String(), bytes.NewBuffer(req_body_json))
	if err != nil {
		log.Printf("An error occurred creating the request, aborting\n")
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + user.Access_token)
	clientId := os.Getenv("VITE_CLIENT_ID")
	req.Header.Set("Client-Id", clientId)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: time.Second * 10}
	
	log.Printf("Req: %v\n\n", req)
	log.Printf("Req header: %v\n\n", req.Header)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to make subscription\n\n")
		return nil, err
	}
	defer res.Body.Close()

	log.Printf("res: %v\n\n", res)

	if res.StatusCode != 202 {
		log.Printf("Failed to make subscription")
		return nil, &StatusError{Code: res.Status}
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("An error occurred reading the response body, aborting\n\n")
		return nil, err
	}

	var res_body_obj CreateSubscriptionRes
	err = json.Unmarshal(res_body, &res_body_obj)
	if err != nil {
		log.Printf("An error occurred parsing the response body, aborting\n\n")
		return nil, err
	}
	log.Printf("Read response body: %v\n\n", res_body_obj)

	if len(res_body_obj.Data) == 0 {
		log.Panic("Twitch API and chatter data types are out of sync")
	}

	subId := res_body_obj.Data[0].Id

	subList, ok := es.Client.subscriptions[subType]
	if ok {
		// list for this subscription type already exists, append
		subList = append(subList, subId)
	} else {
		subList = []string{subId}
	}

	return &ESSubscription{
		SubId: subId,
		SubType: subType,
	}, nil
	
}
