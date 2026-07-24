package seventv

import (
	"chatter-wails/internal/api/seventv"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type SevenTVService struct {
	Ctx context.Context

	app *application.App
	eventAPISubs map[string]string
}

func NewSevenTVService(app *application.App) *SevenTVService{
	return &SevenTVService{
		app: app,
		eventAPISubs: map[string]string{},
	}
}

func (stv *SevenTVService) EnableSevenTV(broadcasterId string) error {
	if set, exists := cache.GetEmoteSet(cache.STV_KEY, broadcasterId); exists {
		shared.EmitNewSet(stv.app, set, broadcasterId)
		return nil
	}

	userRes, err := seventv.GetSevenTVUser("twitch", broadcasterId)
	if err != nil {
		return err
	}

	set := seventv.GetAppEmotesFromSevenTVUserRes(userRes)
	cache.SetEmoteSet(cache.STV_KEY, broadcasterId, set)
	shared.EmitNewSet(stv.app, set, broadcasterId)

	// TODO: add emote set udpate listeners

	return nil
}
