package main

import (
	"context"
	"chatter-wails/services"
	"github.com/joho/godotenv"
)

// App struct
type App struct {
	ctx context.Context
	esService *services.EventSubService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		esService: services.NewEventSubService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	godotenv.Load()

	a.ctx = ctx
	a.esService.Ctx = ctx

	a.esService.Connect()
}
