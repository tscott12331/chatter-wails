package main

import (
	"chatter-wails/services"
	"context"
)

// App struct
type App struct {
	ctx context.Context
	esService *services.EventSubService
	emoteService *services.EmoteService
	badgeService *services.BadgeService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		esService: services.NewEventSubService(),
		emoteService: services.NewEmoteService(),
		badgeService: services.NewBadgeService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.esService.Ctx = ctx
	a.emoteService.Ctx = ctx
	a.badgeService.Ctx = ctx

	go a.esService.Connect()
}
