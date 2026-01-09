package main

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/user"
	"chatter-wails/services"
	"context"
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


const CHAT_SUB_TYPE = "channel.chat.message"


type ChatroomData struct{
	SubId string					`json:"subId"`
	BroadcasterId string			`json:"broadcasterId"`
	BadgeSets []api.ApiBadgeSet		`json:"badgeSets"`
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
	defer func() {
		close(badgeSetsDone)
		close(badgeSetsErr)
	}()

	bsCtx, bsCancel := context.WithCancel(a.ctx)
	defer bsCancel()

	go func(ctx context.Context) {
		badgeSets, err := a.badgeService.GetChannelBadgeSets(accessToken, broadcaster.Id)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				badgeSetsErr <- err
			}
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			badgeSetsDone<- badgeSets
		}

	}(bsCtx)

	condition := map[string]string{
		"broadcaster_user_id": broadcaster.Id,
		"user_id": a.authService.User.Id,
	}
	subId, err := a.esService.CreateSubscription(accessToken, condition, CHAT_SUB_TYPE)
	if err != nil {
		log.Printf("[ConnectToChatroom]: An error occurred creating the chat subscription, aborting\n\n")
		bsCancel()
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

	newSub := services.ESSubscription[services.ESChatSubscriptionData]{
		SubType: CHAT_SUB_TYPE,
		Data: services.ESChatSubscriptionData{
			ChannelBadgeSets: chatroomData.BadgeSets,
		},
	}
	a.esService.Client.ChatSubscriptions[subId] = newSub
	

	return chatroomData, nil

}
