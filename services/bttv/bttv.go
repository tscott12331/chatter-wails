package bttv

import (
	"chatter-wails/internal/api/bttvApi"
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type BTTVService struct {
	Ctx context.Context

	app *application.App
}

// TODO: app emote?
func (bttv *BTTVService) GetGlobalEmotes() []bttvApi.BTTVEmote {
	panic("not implemented")
}

func (bttv *BTTVService) GetUser(platform, channel string) *bttvApi.BTTVUser {
	panic("not implemented")
}
