package seventv

import (
	"chatter-wails/internal/api/seventvApi"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"context"
	"errors"
	"fmt"
	"sync"

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

func (stv *SevenTVService) RequestSevenTVEmotes(broadcasterId string) error {
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Go(func(){
		errChan <- stv.RequestSevenTVChannelEmotes(broadcasterId)
	})
	wg.Go(func(){
		errChan <- stv.RequestSevenTVGlobalEmotes()
	})

	wg.Wait()
	close(errChan)


	var err error
	for e := range errChan {
		err = errors.Join(err, e)
	}

	return err
}

func (stv *SevenTVService) GetSevenTVChannelEmotes(broadcasterId string) (*types.AppEmoteSet, error) {
	if set, exists := cache.GetEmoteSet(cache.STV_KEY, broadcasterId); exists {
		return set, nil
	}

	userRes, err := seventvApi.GetSevenTVUser("twitch", broadcasterId)
	if err != nil {
		return nil, err
	}

	set := seventvApi.GetAppEmotesFromSevenTVUserRes(userRes)
	cache.SetEmoteSet(cache.STV_KEY, broadcasterId, set)

	return set, nil
}

func (stv *SevenTVService) RequestSevenTVChannelEmotes(broadcasterId string) error {
	set, err := stv.GetSevenTVChannelEmotes(broadcasterId)
	if err != nil {
		return err
	}

	shared.EmitNewSet(stv.app, set, true, broadcasterId)

	// TODO: add emote set udpate listeners

	return nil
}

func (stv *SevenTVService) GetSevenTVGlobalEmotes() (*types.AppEmoteSet, error) {
	if set, exists := cache.GetEmoteSet(cache.STV_KEY, cache.GLOBAL_EMOTE_SECTION); exists {
		return set, nil
	}

	res, err := seventvApi.GetGlobalEmotes()
	if err != nil {
		return nil, err
	}

	set := seventvApi.GetAppEmotesFromSevenTVEmotes(res.Emotes, fmt.Sprintf("%s:%s", cache.STV_KEY, cache.GLOBAL_EMOTE_SECTION), cache.GLOBAL_EMOTE_SECTION)
	cache.SetEmoteSet(cache.STV_KEY, cache.GLOBAL_EMOTE_SECTION, set)

	return set, nil
}

func (stv *SevenTVService) RequestSevenTVGlobalEmotes() error {
	set, err := stv.GetSevenTVGlobalEmotes()
	if err != nil {
		return err
	}

	shared.EmitNewSet(stv.app, set, false, "")
	return nil
}
