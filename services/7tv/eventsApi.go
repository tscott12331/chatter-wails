package seventv

import (
	"encoding/json"
	"log"
	"time"
)

type EventsApiMessage struct {
	Op   int             `json:"op"`
	Body json.RawMessage `json:"d"`
}

type EventsApiHelloBody struct {
	HeartbeatInterval uint32 `json:"heartbeat_interval"`
	SessionId         string `json:"session_id"`
	SubscriptionLimit int32  `json:"subscription_limit"`
}

func (stv *SevenTVService) handleEventsApiMessage(m []byte) {
	log.Printf("7TV EVENTS API: %s\n", string(m))

	var message EventsApiMessage
	err := json.Unmarshal(m, &message)
	if err != nil {
		log.Printf("ERROR: %+v", err)
		return
	}

	switch message.Op {
	case EAPI_OP_HELLO:
		stv.handleEventsApiHello(&message)
	case EAPI_OP_RECONNECT:
		stv.handleEventsApiReconnect(&message)
	case EAPI_OP_HEARTBEAT:
		stv.handleEventsApiHeartbeat(&message)
	}
}

func (stv *SevenTVService) handleEventsApiHello(message *EventsApiMessage) {
	var body EventsApiHelloBody
	err := json.Unmarshal(message.Body, &body)
	if err != nil {
		log.Printf("ERROR: %+v", err)
		return
	}

	heartbeatInterval := time.Duration(body.HeartbeatInterval)*time.Millisecond + EAPI_HEARTBEAT_LEV
	deadline := time.Now().Add(heartbeatInterval)

	stv.socket.SetReadDeadline(deadline)
	stv.eventsApiHeartbeat = heartbeatInterval
}

func (stv *SevenTVService) handleEventsApiHeartbeat(message *EventsApiMessage) {
	// ignoring body for now
	stv.socket.SetReadDeadline(time.Now().Add(stv.eventsApiHeartbeat))
}

func (stv *SevenTVService) handleEventsApiReconnect(message *EventsApiMessage) {
	// might have to resubscribe to subscriptions that were lost before eos?
	err := stv.connect()
	if err != nil {
		log.Printf("ERROR: %+v\n", err)
	}
}

func (stv *SevenTVService) listenChannelEvent(subType, platform, broadcasterId string) error {
	return stv.eventsApiSubscribe(subType, map[string]string{
		EAPI_CONDITION_CTX:      EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: platform,
		EAPI_CONDITION_ID:       broadcasterId,
	})
}

func (stv *SevenTVService) eventsApiSubscribe(subType string, condition map[string]string) error {
	data := map[string]any{
		"type":      subType,
		"condition": condition,
	}

	return stv.eventsApiSend(EAPI_OP_SUBSCRIBE, data)
}

func (stv *SevenTVService) eventsApiSend(opCode int, data any) error {
	payload := map[string]any{
		"op": opCode,
		"d":  data,
	}

	enc, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	stv.socket.SendMessage(enc)
	return nil
}
