package main

import (
	"chatter-wails/internal/api"
	seventv "chatter-wails/services/7tv"
	"chatter-wails/services/auth"
	"chatter-wails/services/badge"
	"chatter-wails/services/emote"
	"chatter-wails/services/eventsub"
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
	emoteService *emote.EmoteService
	badgeService *badge.BadgeService
	authService *auth.AuthService
	seventvService *seventv.SevenTVService
}

// NewAppService creates a new App application struct
func NewAppService(app *application.App) *AppService {
	return &AppService{
		app: app,
		esService: eventsub.NewEventSubService(app),
		emoteService: emote.NewEmoteService(app),
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
	if a.authService.User == nil {
		return nil, &NotLoggedInError{}
	}

	accessToken := a.authService.User.Access_token

	return a.esService.CreateChatSubscription(accessToken, a.authService.User.Id, channelName, a.badgeService.GlobalBadgeSets)
}

func (a *AppService) goGetChannelBadgeSets(
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

	badgeSets, err := badge.GetChannelBadgeSets(accessToken, broadcasterId)
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}
	combinedSets := badge.CombineChannelGlobalSets(badgeSets, a.badgeService.GlobalBadgeSets)

	if !sub.Data.ChannelBadgeSets.Write(*combinedSets) {
		log.Printf("Tried to write badge sets which were already written")
	}
}

func (a *AppService) SendChatMessage(channelName string, messageContent string, replyId *string) (*api.ApiPostMessagesData, error) {
	user := a.authService.User
	if user == nil {
		log.Printf("[SendChatMessage]: User not logged int, cannot send message, aborting\n\n")
		return nil, &NotLoggedInError{}
	}

	return a.esService.SendChatMessageFromChannelName(user.Access_token, user.Id, channelName, messageContent, replyId)
}

func (a *AppService) DisconnectFromChatroom(channelName string) error {
	a.esService.Client.IrcListener.PartChannel(channelName)

	if sub, exists := a.esService.Client.ChatSubscriptions.Read().GetSubFromChannelName(channelName); exists {
		cache.RemoveBroadcasterEmoteSets(sub.Data.BroadcasterId)
	}

	return a.esService.DeleteChatSubscriptionFromChannelName(a.authService.User.Access_token, channelName)
}
