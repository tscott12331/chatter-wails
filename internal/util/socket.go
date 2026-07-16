package util

import (
	"bytes"
	"context"
	"log"
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
	// maxMessageSize = 512


)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Socket struct {
	ctx context.Context
	conn *websocket.Conn
	send chan []byte
}

type SocketMessageHandler func(messageBytes []byte)

func NewSocket(ctx context.Context, connStr string, messageHandler SocketMessageHandler) (*Socket, error) {
	
	client, _, err := websocket.DefaultDialer.Dial(connStr, nil)
	if err != nil {
		return nil, err
	}
	// c.conn.SetReadLimit(maxMessageSize)
	client.SetReadDeadline(time.Now().Add(pongWait))
	client.SetPongHandler(func(string) error { 
		client.SetReadDeadline(time.Now().Add(pongWait)); return nil
	})

	socket := &Socket{
		conn: client,
		ctx: ctx,
		send: make(chan []byte),
	}

	go socket.readPump(messageHandler)
	go socket.writePump()

	return socket, nil
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (s *Socket) readPump(handler SocketMessageHandler) {
	log.Printf("[readPump]: Started read goroutine for socket\n\n")
	defer func() {
		log.Printf("[readPump]: Setting ready to false\n\n")
	}()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[readPump]: Context canceled, closing\n\n")
			return
		default:
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("[readPump]: unexpected close, %v\n\n", err)
				} else {
					log.Printf("[readPump]: expected close, %v\n\n", err)
				}
				return
			}
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

			handler(message)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (s *Socket) writePump() {
	log.Printf("[writePump]: Started write goroutine for socket\n\n")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Printf("[writePump]: Closing connection\n\n")
		ticker.Stop()
		s.conn.Close()
	}()

	for {
		select {
		case message, ok := <-s.send:
			log.Printf("[writePump]: Attempting to write message to socket\n\n")
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Printf("[writePump]: Message channel was closed, aborting\n\n")
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			defer func() {
				if err := w.Close(); err != nil {
					log.Printf("[writePump]: An error occurred sending the message, aborting\n\n")
				}
			}()

			if err != nil {
				log.Printf("[writePump]: An error occurred creating the message writer, aborting\n\n")
				return
			}

			w.Write(message)
			log.Printf("[writePump]: wrote %s\n\n", message)

			// Add queued chat messages to the current websocket message.
			// n := len(s.send)
			// for range n {
			// 	queued_message := <-s.send
			// 	log.Printf("[writePump]: Added queued message to write: %s\n\n", queued_message)
			// 	w.Write(newline)
			// 	w.Write(queued_message)
			// }
		case <-ticker.C:
			log.Printf("[writePump]: Attempting to send ping message\n\n")
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[writePump]: An error occurred sending ping message, aborting\n\n")
				return
			}
		case <-s.ctx.Done():
			log.Printf("[writePump]: Context canceled, closing\n\n")
			return
		}
	}
}

func (s *Socket) SendMessage(message []byte) {
	log.Printf("[SendMessage]: Sending %s", string(message))
	s.send <- message
}
