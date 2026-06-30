package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
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

	app.RegisterService(appService)
	app.RegisterService(esService)
	app.RegisterService(emoteService)
	app.RegisterService(badgeService)
	app.RegisterService(authService)


	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "chatter-wails",
		Width:  1024,
		Height: 768,
		BackgroundColour: application.NewRGBA(27, 38, 54, 1),
	})

	err := app.Run()

	if err != nil {
		println("ERROR:", err.Error())
	}
}
