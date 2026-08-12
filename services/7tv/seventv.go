package seventv

import (
	"chatter-wails/internal/api/seventvApi"
	"chatter-wails/internal/util"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const STV_WSS = "wss://events.7tv.io/v3"

const EAPI_HEARTBEAT_LEV = 5 * time.Second

const EAPI_OP_HELLO = 1
const EAPI_OP_HEARTBEAT = 2
const EAPI_OP_RECONNECT = 4
const EAPI_OP_SUBSCRIBE = 35

const EAPI_CONDITION_OBJECT_ID = "object_id"
const EAPI_CONDITION_HOST_ID = "host_id"
const EAPI_CONDITION_CONNECTION_ID = "connection_id"
const EAPI_CONDITION_CTX = "ctx"
const EAPI_CONDITION_PLATFORM = "platform"
const EAPI_CONDITION_ID = "id"

const EAPI_CTX_CHANNEL = "channel"
const EAPI_PLATFORM_TWITCH = "TWITCH"

const EAPI_SUB_ENTITLEMENT_CREATE = "entitlement.create"
const EAPI_SUB_ENTITLEMENT_UPDATE = "entitlement.update"
const EAPI_SUB_ENTITLEMENT_DELETE = "entitlement.delete"
const EAPI_SUB_ENTITLEMENT_ALL = "entitlement.*"

const EAPI_SUB_COSMETIC_CREATE = "cosmetic.create"
const EAPI_SUB_COSMETIC_UPDATE = "cosmetic.update"
const EAPI_SUB_COSMETIC_DELETE = "cosmetic.delete"
const EAPI_SUB_COSMETIC_ALL = "cosmetic.*"

type SevenTVService struct {
	Ctx context.Context

	app *application.App

	socket *util.Socket
	eventsApiHeartbeat time.Duration

	eventAPISubs map[string]string
}

func NewSevenTVService(app *application.App) (*SevenTVService, error){
	service := &SevenTVService{
		app: app,
		eventAPISubs: map[string]string{},
	}

	err := service.connect()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (stv *SevenTVService) connect() error {
	socket, err := util.NewSocket(stv.app.Context(), STV_WSS, stv.handleEventsApiMessage, false)
	if err != nil {
		return err
	}

	stv.socket = socket

	return nil
}

type EventsApiMessage struct{
	Op int			`json:"op"`
	Body json.RawMessage		`json:"d"`
}

type EventsApiHelloBody struct{
	HeartbeatInterval uint32			`json:"heartbeat_interval"`
	SessionId string			`json:"session_id"`
	SubscriptionLimit int32			`json:"subscription_limit"`
}

func (stv *SevenTVService) handleEventsApiMessage(m []byte) {
	log.Printf("7TV EVENTS API: %s\n", string(m))

	var message EventsApiMessage
	err := json.Unmarshal(m, &message)
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
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
		fmt.Printf("ERROR: %+v", err)
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

func (stv *SevenTVService) listenAllEntitlement(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_ENTITLEMENT_ALL, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}

func (stv *SevenTVService) listenCreateEntitlement(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_ENTITLEMENT_CREATE, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}

func (stv *SevenTVService) listenUpdateEntitlement(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_ENTITLEMENT_UPDATE, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}

func (stv *SevenTVService) listenDeleteEntitlement(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_ENTITLEMENT_DELETE, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}


func (stv *SevenTVService) listenAllCosmetic(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_COSMETIC_ALL, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}

func (stv *SevenTVService) listenCreateCosmetic(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_COSMETIC_CREATE, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}

func (stv *SevenTVService) listenUpdateCosmetic(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_COSMETIC_UPDATE, map[string]string{
		EAPI_CONDITION_CONNECTION_ID: broadcasterId,
	})
}

func (stv *SevenTVService) listenDeleteCosmetic(broadcasterId string) error {
	return stv.eventsApiSubscribe(EAPI_SUB_COSMETIC_DELETE, map[string]string{
		EAPI_CONDITION_CTX: EAPI_CTX_CHANNEL,
		EAPI_CONDITION_PLATFORM: EAPI_PLATFORM_TWITCH,
		EAPI_CONDITION_ID: broadcasterId,
	})
}



func (stv *SevenTVService) eventsApiSubscribe(subType string, condition map[string]string) error {
	data := map[string]any{
		"type": subType,
		"condition": condition,
	}

	return stv.eventsApiSend(EAPI_OP_SUBSCRIBE, data)
}

func (stv *SevenTVService) eventsApiSend(opCode int, data any) error {
	payload := map[string]any {
		"op": opCode,
		"d": data,
	}

	enc, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	stv.socket.SendMessage(enc)
	return nil
}

func (stv *SevenTVService) RequestSevenTVEmotes(broadcasterId string) error {
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Go(func(){
		errChan <- stv.RequestSevenTVChannelEmotes(broadcasterId)
	})
	wg.Go(func(){
		errChan <- stv.RequestSevenTVGlobalEmotes()
	})
	wg.Go(func() {
		errChan <- stv.listenAllEntitlement(broadcasterId)
	})
	wg.Go(func() {
		errChan <- stv.listenAllCosmetic(broadcasterId)
	})

	wg.Wait()
	close(errChan)


	var err error
	for e := range errChan {
		err = errors.Join(err, e)
	}

	// dispatch subscriptions to main loop?

	return err
}

func (stv *SevenTVService) GetSevenTVChannelEmotes(broadcasterId string) (*types.AppEmoteSet, error) {
	if set, exists := cache.GetEmoteSet(cache.STV_KEY, broadcasterId); exists {
		return set, nil
	}

	userRes, err := seventvApi.GetSevenTVUser("twitch", broadcasterId)
	if err != nil {
		return nil, err
	}

	set := seventvApi.GetAppEmotesFromSevenTVUserRes(userRes)
	cache.SetEmoteSet(cache.STV_KEY, broadcasterId, set)


	return set, nil
}

func (stv *SevenTVService) RequestSevenTVChannelEmotes(broadcasterId string) error {
	set, err := stv.GetSevenTVChannelEmotes(broadcasterId)
	if err != nil {
		return err
	}

	shared.EmitNewSet(stv.app, set, true, broadcasterId)

	// TODO: add emote set udpate listeners

	return nil
}

func (stv *SevenTVService) GetSevenTVGlobalEmotes() (*types.AppEmoteSet, error) {
	if set, exists := cache.GetEmoteSet(cache.STV_KEY, cache.GLOBAL_EMOTE_SECTION); exists {
		return set, nil
	}

	res, err := seventvApi.GetGlobalEmotes()
	if err != nil {
		return nil, err
	}

	set := seventvApi.GetAppEmotesFromSevenTVEmotes(res.Emotes, fmt.Sprintf("%s:%s", cache.STV_KEY, cache.GLOBAL_EMOTE_SECTION), cache.GLOBAL_EMOTE_SECTION)
	cache.SetEmoteSet(cache.STV_KEY, cache.GLOBAL_EMOTE_SECTION, set)

	return set, nil
}

func (stv *SevenTVService) RequestSevenTVGlobalEmotes() error {
	set, err := stv.GetSevenTVGlobalEmotes()
	if err != nil {
		return err
	}

	shared.EmitNewSet(stv.app, set, false, "")
	return nil
}
