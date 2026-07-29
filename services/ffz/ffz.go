package ffz

import (
	ffzapi "chatter-wails/internal/api/ffzApi"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"context"
	"errors"
	"maps"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type FFZService struct {
	Ctx context.Context

	app *application.App
}

func NewFFZService(app *application.App) *FFZService{
	return &FFZService{
		app: app,
	}
}

func (ffz *FFZService) RequestFFZEmotes(broadcasterId string) error {
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Go(func(){
		errChan <- ffz.RequestFFZGlobalEmotes()
	})
	wg.Go(func() {
		errChan <- ffz.RequestFFZChannelEmotes(broadcasterId)
	})

	wg.Wait()
	close(errChan)


	var err error
	for e := range errChan {
		err = errors.Join(err, e)
	}

	return err
}

func (ffz *FFZService) GetFFZGlobalEmotes() (*types.AppEmoteSet, error) {
	if set, exists := cache.GetEmoteSet(cache.FFZ_KEY, cache.GLOBAL_EMOTE_SECTION); exists {
		return set, nil
	}

	res, err := ffzapi.GetFFZGlobalEmoteSet()
	if err != nil {
		return nil, err
	}
	sets := []*types.AppEmoteSet{}
	for _, setId := range res.Default_sets {
		ffzSet := res.Sets[strconv.Itoa(setId)]
		appSet := ffzapi.GetAppEmoteSetFromFFZEmoteSet(&ffzSet, cache.GLOBAL_EMOTE_SECTION)
		sets = append(sets, appSet)
	}

	if len(sets) == 0 {
		return nil, errors.New("No global sets exist")
	}

	baseSet := sets[0]
	for i := 1; i < len(sets); i++ {
		maps.Copy(baseSet.Emotes, sets[i].Emotes)
	}

	cache.SetEmoteSet(cache.FFZ_KEY, cache.GLOBAL_EMOTE_SECTION, baseSet)
	
	return baseSet, nil
}

func (ffz *FFZService) RequestFFZGlobalEmotes() error {
	set, err := ffz.GetFFZGlobalEmotes()
	if err != nil {
		return err
	}

	shared.EmitNewSet(ffz.app, set, false, "")
	return nil
}


func (ffz *FFZService) GetFFZChannelEmotes(broadcasterId string) (*types.AppEmoteSet, error) {
	if set, exists := cache.GetChannelEmoteSet(cache.FFZ_KEY, broadcasterId); exists {
		return set, nil
	}

	res, err := ffzapi.GetFFZRoom(cache.TWITCH_KEY, broadcasterId)
	if err != nil {
		return nil, err
	}

	set, exists := res.Sets[strconv.Itoa(res.Room.Set)]
	if !exists {
		return nil, errors.New("Fetched channel set does not exist")
	}

	appSet := ffzapi.GetAppEmoteSetFromFFZEmoteSet(&set, cache.CHANNEL_EMOTE_SECTION)
	cache.SetChannelEmoteSet(cache.FFZ_KEY, broadcasterId, appSet)

	return appSet, nil
}

func (ffz *FFZService) RequestFFZChannelEmotes(broadcasterId string) error {
	set, err := ffz.GetFFZChannelEmotes(broadcasterId)
	if err != nil {
		return err
	}

	shared.EmitNewSet(ffz.app, set, true, broadcasterId)
	return nil
}
