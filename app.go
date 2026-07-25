package main

import (
	"chatter-wails/internal/api/nativeApi"
	seventv "chatter-wails/services/7tv"
	"chatter-wails/services/auth"
	"chatter-wails/services/badge"
	"chatter-wails/services/eventsub"
	"chatter-wails/services/native"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"context"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AppService struct
type AppService struct {
	app *application.App
	ctx context.Context
	esService *eventsub.EventSubService
	emoteService *native.EmoteService
	badgeService *badge.BadgeService
	authService *auth.AuthService
	seventvService *seventv.SevenTVService
}

// NewAppService creates a new App application struct
func NewAppService(app *application.App) *AppService {
	return &AppService{
		app: app,
		esService: eventsub.NewEventSubService(app),
		emoteService: native.NewEmoteService(app),
		badgeService: badge.NewBadgeService(app),
		authService: auth.NewAuthService(app),
		seventvService: seventv.NewSevenTVService(app),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.ctx = ctx
	a.esService.Ctx = ctx
	a.emoteService.Ctx = ctx
	a.badgeService.Ctx = ctx
	a.authService.Ctx = ctx
	a.seventvService.Ctx = ctx

	go a.esService.Connect()

	return nil
}

type NotLoggedInError struct{}
func (nli *NotLoggedInError) Error() string {
	return "User is not yet logged in"
}



func (a *AppService) ConnectToChatroom(channelName string) (*eventsub.ChatroomData, error) {
	user := shared.GetUser()
	if user == nil {
		return nil, &NotLoggedInError{}
	}

	return a.esService.CreateChatSubscription(channelName, a.badgeService.GlobalBadgeSets)
}

func (a *AppService) SendChatMessage(channelName string, messageContent string, replyId *string) (*nativeApi.ApiPostMessagesData, error) {
	user := shared.GetUser()
	if user == nil {
		log.Printf("[SendChatMessage]: User not logged in, cannot send message, aborting\n\n")
		return nil, &NotLoggedInError{}
	}

	return a.esService.SendChatMessageFromChannelName(channelName, messageContent, replyId)
}

func (a *AppService) DisconnectFromChatroom(channelName string) error {
	a.esService.Client.IrcListener.PartChannel(channelName)

	if sub, exists := a.esService.Client.ChatSubscriptions.Read().GetSubFromChannelName(channelName); exists {
		cache.RemoveBroadcasterEmoteSets(sub.Data.BroadcasterId)
	}

	return a.esService.DeleteChatSubscriptionFromChannelName(channelName)
}
