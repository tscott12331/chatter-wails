package bttv

import (
	"chatter-wails/internal/api/bttvApi"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"context"
	"errors"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type BTTVService struct {
	Ctx context.Context

	app *application.App
}

func NewBTTVService(app *application.App) *BTTVService {
	return &BTTVService{
		app: app,
	}
}

func (bttv *BTTVService) RequestBTTVEmotes(broadcasterId string) error {
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Go(func(){
		errChan <- bttv.RequestBTTVGlobalEmotes()
	})
	wg.Go(func() {
		errChan <- bttv.RequestBTTVChannelEmotes(broadcasterId)
	})

	wg.Wait()
	close(errChan)


	var err error
	for e := range errChan {
		err = errors.Join(err, e)
	}

	return err
}

func (bttv *BTTVService) GetBTTVGlobalEmotes() (*types.AppEmoteSet, error) {
	if set, exists := cache.GetEmoteSet(cache.BTTV_KEY, cache.GLOBAL_EMOTE_SECTION); exists {
		return set, nil
	}
	
	res, err := bttvApi.GetGlobalEmotes()
	if err != nil {
		return nil, err
	}

	set := bttvApi.BTTVEmotesToAppEmoteSet(res, cache.GLOBAL_EMOTE_SECTION)
	cache.SetEmoteSet(cache.BTTV_KEY, cache.GLOBAL_EMOTE_SECTION, set)

	return set, nil
}

func (bttv *BTTVService) RequestBTTVGlobalEmotes() error {
	set, err := bttv.GetBTTVGlobalEmotes()
	if err != nil {
		return err
	}

	shared.EmitNewSet(bttv.app, set, false, "")
	return nil
}


func (bttv *BTTVService) GetBTTVChannelEmotes(broadcasterId string) (*types.AppEmoteSet, error) {
	if set, exists := cache.GetChannelEmoteSet(cache.BTTV_KEY, broadcasterId); exists {
		return set, nil
	}

	res, err := bttvApi.GetUser(bttvApi.BTTV_TWITCH_PROVIDER, broadcasterId)
	if err != nil {
		return nil, err
	}

	set := bttvApi.BTTVEmotesToAppEmoteSet(append(res.ChannelEmotes, res.SharedEmotes...), cache.CHANNEL_EMOTE_SECTION)
	cache.SetChannelEmoteSet(cache.BTTV_KEY, broadcasterId, set)

	return set, nil
}

func (bttv *BTTVService) RequestBTTVChannelEmotes(broadcasterId string) error {
	set, err := bttv.GetBTTVChannelEmotes(broadcasterId)
	if err != nil {
		return err
	}

	// add listeners for emote updates

	shared.EmitNewSet(bttv.app, set, true, broadcasterId)
	return nil
}
