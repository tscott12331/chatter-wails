package main

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/api/seventv"
	"chatter-wails/internal/message"
	"chatter-wails/services"
	"context"
	"errors"
	"log"
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



func (a *App) ConnectToChatroom(channelName string) (*services.ChatroomData, error) {
	if a.authService.User == nil {
		return nil, &NotLoggedInError{}
	}

	accessToken := a.authService.User.Access_token

	return a.esService.CreateChatSubscription(accessToken, a.authService.User.Id, channelName, a.badgeService.GlobalBadgeSets)

	// if data, exists := a.esService.Client.GetChatroomData(channelName); exists {
	// 	return data, nil
	// }
	//
	// broadcaster, err := user.GetUserByLogin(accessToken, channelName)
	// if err != nil {
	// 	log.Printf("[ConnectToChatroom]: An error occurred fetching the broadcaster info, aborting")
	// 	return nil, err
	// }
	//
	// var channelEmotesDone chan map[string]*services.AppEmote = make(chan map[string]*services.AppEmote)
	// var channelEmotesErr chan error = make(chan error)
	// defer func() {
	// 	close(channelEmotesDone)
	// 	close(channelEmotesErr)
	// }()
	//
	// bsCtx, bsCancel := context.WithCancel(a.ctx)
	// defer bsCancel()
	//
	// var wg sync.WaitGroup
	//
	// // fetch channel emotes in goroutine
	// wg.Add(1)
	// go a.goGetChannelEmotes(bsCtx, channelEmotesDone, channelEmotesErr, &wg, accessToken, broadcaster.Id)
	//
	// condition := map[string]string{
	// 	"broadcaster_user_id": broadcaster.Id,
	// 	"user_id": a.authService.User.Id,
	// }
	// subId, err := a.esService.CreateSubscription(accessToken, condition, CHAT_SUB_TYPE)
	// if err != nil {
	// 	log.Printf("[ConnectToChatroom: ERROR: %+v\n", err)
	// 	wg.Wait()
	// 	return nil, err
	// }
	//
	// chatroomData := &services.ChatroomData{
	// 	SubId: subId,
	// 	BroadcasterId: broadcaster.Id,
	// }
	//
	// select {
	// case channelEmotes := <-channelEmotesDone:
	// 	chatroomData.ChannelEmotes = channelEmotes
	// case err := <-channelEmotesErr:
	// 	log.Printf("[ConnectToChatroom]: Failed to get channel emotes\n%+v\n\n", err)
	// }
	//
	// newSub := &services.ESSubscription[*services.ESChatSubscriptionData]{
	// 	SubId: subId,
	// 	SubType: CHAT_SUB_TYPE,
	// 	Data: &services.ESChatSubscriptionData{
	// 		BroadcasterId: broadcaster.Id,
	// 		Channel: channelName,
	// 		ChannelBadgeSets: &util.SingleWriteMutex[[]api.ApiBadgeSet]{},
	// 		ChannelEmotes: chatroomData.ChannelEmotes,
	// 		SevenTV: &services.ESChatSubscriptionSevenTVData{},
	// 	},
	// }
	// a.esService.Client.AddChatSubscription(newSub)
	//
	// // fetch channel badge sets in goroutine
	// go a.goGetChannelBadgeSets(accessToken, broadcaster.Id, subId)
	//
	// return chatroomData, nil

}

func (a *App) EnableSevenTV(subId string) (map[string]*services.AppEmote, error) {
	sub, subExists := a.esService.Client.ChatSubscriptions.Read().GetSubFromId(subId)

	if !subExists {
		return nil, errors.New("Cannot enable 7tv on nonexistent chat subscription")
	}

	if subExists && sub.Data.SevenTV.Enabled {
		return sub.Data.SevenTV.SevenTVEmotes, nil
	}

	userRes, err := seventv.GetSevenTVUser("twitch", sub.Data.BroadcasterId)
	if err != nil {
		return nil, err
	}

	emotes := seventv.GetAppEmotesFromSevenTVUserRes(userRes)
	sub.Data.SevenTV.SevenTVEmotes = emotes
	sub.Data.SevenTV.Enabled = true

	return emotes, nil
}

func (a *App) goGetChannelBadgeSets(
	accessToken string,
	broadcasterId string,
	subId string,
) {
	sub, exists := a.esService.Client.ChatSubscriptions.Read().GetSubFromId(subId)
	if !exists {
		log.Printf("Subscription %v doesn't exist\n", subId)
		return
	}

	// data already fetched
	if sub.Data.ChannelBadgeSets.IsWritten() {
		log.Printf("Badge sets already written")
		return
	}

	badgeSets, err := services.GetChannelBadgeSets(accessToken, broadcasterId)
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}
	combinedSets := services.CombineChannelGlobalSets(badgeSets, a.badgeService.GlobalBadgeSets)

	if !sub.Data.ChannelBadgeSets.Write(*combinedSets) {
		log.Printf("Tried to write badge sets which were already written")
	}
}

func (a *App) SendChatMessage(chatSubId string, messageContent string, replyId *string) (*api.ApiPostMessagesData, error) {
	user := a.authService.User
	if user == nil {
		log.Printf("[SendChatMessage]: User not logged int, cannot send message, aborting\n\n")
		return nil, &NotLoggedInError{}
	}

	subData, ok := a.esService.Client.ChatSubscriptions.Read().GetSubFromId(chatSubId)
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

func (a *App) DisconnectFromChatroom(channelName string) error {
	return a.esService.DeleteChatSubscriptionFromChannelName(a.authService.User.Access_token, channelName)
}
