package main

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/api/seventv"
	"chatter-wails/internal/message"
	"chatter-wails/internal/user"
	"chatter-wails/internal/util"
	"chatter-wails/services"
	"context"
	"errors"
	"log"
	"sync"
)

// App struct
type App struct {
	ctx context.Context
	esService *services.EventSubService
	emoteService *services.EmoteService
	badgeService *services.BadgeService
	authService *services.AuthService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		esService: services.NewEventSubService(),
		emoteService: services.NewEmoteService(),
		badgeService: services.NewBadgeService(),
		authService: services.NewAuthService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.esService.Ctx = ctx
	a.emoteService.Ctx = ctx
	a.badgeService.Ctx = ctx
	a.authService.Ctx = ctx

	go a.esService.Connect()
}

type NotLoggedInError struct{}
func (nli *NotLoggedInError) Error() string {
	return "User is not yet logged in"
}


const CHAT_SUB_TYPE = "channel.chat.message"


type ChatroomData struct{
	SubId string					`json:"subId"`
	BroadcasterId string			`json:"broadcasterId"`
	BadgeSets []api.ApiBadgeSet		`json:"badgeSets"`
	ChannelEmotes map[string]*services.AppEmote `json:"channelEmotes"`
	SevenTVEmotes map[string]*services.AppEmote `json:"sevenTVEmotes"`
}

func (a *App) ConnectToChatroom(channelName string) (*ChatroomData, error) {
	if a.authService.User == nil {
		return nil, &NotLoggedInError{}
	}

	accessToken := a.authService.User.Access_token
	
	broadcaster, err := user.GetUserByLogin(accessToken, channelName)
	if err != nil {
		log.Printf("[ConnectToChatroom]: An error occurred fetching the broadcaster info, aborting\n\n")
		return nil, err
	}

	var badgeSetsDone chan *[]api.ApiBadgeSet = make(chan *[]api.ApiBadgeSet)
	var badgeSetsErr chan error = make(chan error)
	var channelEmotesDone chan map[string]*services.AppEmote = make(chan map[string]*services.AppEmote)
	var channelEmotesErr chan error = make(chan error)
	var sevenTVUserReqDone chan *seventv.ApiGetSevenTVUserRes = make(chan *seventv.ApiGetSevenTVUserRes)
	var sevenTVErr chan error = make(chan error)
	defer func() {
		close(badgeSetsDone)
		close(badgeSetsErr)
		close(channelEmotesDone)
		close(channelEmotesErr)
	}()

	bsCtx, bsCancel := context.WithCancel(a.ctx)
	defer bsCancel()

	var wg sync.WaitGroup

	// fetch channel badge sets in goroutine
	wg.Add(1)
	go a.goGetChannelBadgeSets(bsCtx, badgeSetsDone, badgeSetsErr, &wg, accessToken, broadcaster.Id)

	// fetch channel emotes in goroutine
	wg.Add(1)
	go a.goGetChannelEmotes(bsCtx, channelEmotesDone, channelEmotesErr, &wg, accessToken, broadcaster.Id)

	wg.Add(1)
	go a.goGetSevenTVEmotes(bsCtx, sevenTVUserReqDone, sevenTVErr, &wg, accessToken, broadcaster.Id)

	condition := map[string]string{
		"broadcaster_user_id": broadcaster.Id,
		"user_id": a.authService.User.Id,
	}
	subId, err := a.esService.CreateSubscription(accessToken, condition, CHAT_SUB_TYPE)
	if err != nil {
		log.Printf("[ConnectToChatroom]: An error occurred creating the chat subscription, aborting\n\n")
		wg.Wait()
		return nil, err
	}

	chatroomData := &ChatroomData{
		SubId: subId,
		BroadcasterId: broadcaster.Id,
	}

	select {
	case channelBadgeSets := <-badgeSetsDone:
		badgeSets := services.CombineChannelGlobalSets(channelBadgeSets, a.badgeService.GlobalBadgeSets)
		chatroomData.BadgeSets = *badgeSets
	case err := <-badgeSetsErr:
		log.Printf("[ConnectToChatroom]: Failed to get channel badge sets\n%+v\n\n", err)
	}

	select {
	case channelEmotes := <-channelEmotesDone:
		chatroomData.ChannelEmotes = channelEmotes
	case err := <-channelEmotesErr:
		log.Printf("[ConnectToChatroom]: Failed to get channel emotes\n%+v\n\n", err)
	}

	select {
	case sevenTVUserRes := <-sevenTVUserReqDone:
		chatroomData.SevenTVEmotes = seventv.GetAppEmotesFromSevenTVUserRes(sevenTVUserRes)
	case err := <-sevenTVErr:
		log.Printf("[ConnectToChatroom]: Failed to get 7tv emotes\n%+v\n\n", err)
	}

	newSub := services.ESSubscription[services.ESChatSubscriptionData]{
		SubType: CHAT_SUB_TYPE,
		Data: services.ESChatSubscriptionData{
			BroadcasterId: broadcaster.Id,
			Channel: channelName,
			ChannelBadgeSets: chatroomData.BadgeSets,
			ChannelEmotes: chatroomData.ChannelEmotes,
			SevenTVEmotes: chatroomData.SevenTVEmotes,
		},
	}
	a.esService.Client.ChatSubscriptions[subId] = newSub
	
	return chatroomData, nil

}

func (a *App) goGetChannelEmotes(
	ctx context.Context, 
	channelEmotesDone chan map[string]*services.AppEmote, 
	channelEmotesErr chan error, 
	wg *sync.WaitGroup,
	accessToken string,
	broadcasterId string,
) {
	defer wg.Done()
	channelEmotes, err := a.emoteService.GetChannelEmotes(accessToken, broadcasterId)
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
		channelEmoteMap := util.ArrToMap(*channelEmotes, func(item services.AppEmote) (string, *services.AppEmote) {
			return item.Name, &item
		})
		channelEmotesDone <- channelEmoteMap
		return
	}

}

func (a *App) goGetChannelBadgeSets(
	ctx context.Context, 
	badgeSetsDone chan *[]api.ApiBadgeSet, 
	badgeSetsErr chan error, 
	wg *sync.WaitGroup,
	accessToken string,
	broadcasterId string,
) {
	defer wg.Done()
	badgeSets, err := a.badgeService.GetChannelBadgeSets(accessToken, broadcasterId)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		default:
			badgeSetsErr <- err
			return
		}
	}
	select {
	case <-ctx.Done():
		log.Printf("[ConnectToChatroom]: badgeset context closed")
		return
	default:
		badgeSetsDone <- badgeSets
		return
	}

}

func (a *App) goGetSevenTVEmotes(
	ctx context.Context, 
	sevenTVUserReqDone chan *seventv.ApiGetSevenTVUserRes, 
	sevenTVErr chan error, 
	wg *sync.WaitGroup,
	accessToken string,
	broadcasterId string,
) {
	defer wg.Done()
	userRes, err := seventv.GetSevenTVUser("twitch", broadcasterId)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		default:
			sevenTVErr <- err
			return
		}
	}
	select {
	case <-ctx.Done():
		log.Printf("[ConnectToChatroom]: seventv emote context closed\n")
		return
	default:
		sevenTVUserReqDone <- userRes
		return
	}
}

func (a *App) SendChatMessage(chatSubId string, messageContent string, replyId *string) (*api.ApiPostMessagesData, error) {
	user := a.authService.User
	if user == nil {
		log.Printf("[SendChatMessage]: User not logged int, cannot send message, aborting\n\n")
		return nil, &NotLoggedInError{}
	}

	subData, ok := a.esService.Client.ChatSubscriptions[chatSubId]
	if !ok {
		log.Printf("[SendChatMessage]: Failed to find chat subscription data, aborting\n\n")
		return nil, errors.New("Failed to find chat subscription data")
	}

	res, err := message.SendMessage(user.Id, user.Access_token, subData.Data.BroadcasterId, messageContent, replyId)
	if err != nil {
		log.Printf("[SendChatMessage]: An error occurred sending the chat message, aborting\n%+v\n\n", err)
		return nil, err
	}

	return res, nil
}
