package seventv

import (
	"chatter-wails/internal/api/seventv"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
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
		stv.emitNewSet(set, broadcasterId)
		return nil
	}

	userRes, err := seventv.GetSevenTVUser("twitch", broadcasterId)
	if err != nil {
		return err
	}

	set := seventv.GetAppEmotesFromSevenTVUserRes(userRes)
	cache.SetEmoteSet(cache.STV_KEY, broadcasterId, set)
	stv.emitNewSet(set, broadcasterId)

	// TODO: add emote set udpate listeners

	return nil
}

func (stv *SevenTVService) emitNewSet(set *types.AppEmoteSet, broadcasterId string) {
	stv.app.Event.Emit("chatter:emote:new-set", types.NewEmoteSetEvent{
		BroadcasterId: broadcasterId,
		AppEmoteSet: *set,
	})
}
