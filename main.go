package main

import (
	"chatter-wails/services/eventsub"
	"chatter-wails/shared/types"
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Register events
	// TODO: rename events to be under chatter namespace
	application.RegisterEvent[*eventsub.ESChatMessage]("common:chat-message")
	application.RegisterEvent[eventsub.StreamData]("common:stream-data")
	application.RegisterEvent[eventsub.ChatOpenData]("common:chat-open")
	application.RegisterEvent[*types.AppUser]("common:user-login")
	application.RegisterEvent[eventsub.SharedChatBeginEventData]("common:shared-chat-begin")
	application.RegisterEvent[eventsub.SharedChatUpdateEventData]("common:shared-chat-update")
	application.RegisterEvent[eventsub.SharedChatEndEventData]("common:shared-chat-end")
	application.RegisterEvent[eventsub.BanEventData]("common:ban")
	application.RegisterEvent[eventsub.ClearMsgEventData]("common:clear-msg")
	application.RegisterEvent[types.NewEmoteSetEvent]("chatter:emote:new-set")

	// Create an instance of the app structure
	// Create application with options
	app := application.New(application.Options{
		Name: "Chatter",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	appServiceRaw := NewAppService(app)

	appService := application.NewService(appServiceRaw)
	esService := application.NewService(appServiceRaw.esService)
	emoteService := application.NewService(appServiceRaw.emoteService)
	badgeService := application.NewService(appServiceRaw.badgeService)
	authService := application.NewService(appServiceRaw.authService)
	seventvService := application.NewService(appServiceRaw.seventvService)

	app.RegisterService(appService)
	app.RegisterService(esService)
	app.RegisterService(emoteService)
	app.RegisterService(badgeService)
	app.RegisterService(authService)
	app.RegisterService(seventvService)


	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "chatter-wails",
		Width:  1024,
		Height: 768,
		Frameless: true,
		BackgroundColour: application.NewRGBA(27, 38, 54, 1),
		BackgroundType: application.BackgroundTypeSolid,
	})

	err := app.Run()

	if err != nil {
		println("ERROR:", err.Error())
	}
}
