package native

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/api/nativeApi"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"context"
	"errors"
	"log"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type EmoteService struct {
	app *application.App
	Ctx context.Context
}

func NewEmoteService(app *application.App) *EmoteService {
	return &EmoteService{app: app}
}

func (es *EmoteService) RequestTwitchEmotes(broadcasterId string) error {
	// TODO: refactor this pattern somewhere
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Go(func(){
		errChan <- es.RequestTwitchChannelEmotes(broadcasterId)
	})
	wg.Go(func(){
		errChan <- es.RequestTwitchUserEmotes()
	})
	wg.Go(func(){
		errChan <- es.RequestTwitchGlobalEmotes()
	})

	wg.Wait()
	close(errChan)


	var err error
	for e := range errChan {
		err = errors.Join(err, e)
	}

	return err
}

func (es *EmoteService) GetTwitchUserEmotes() (*types.AppEmoteSet, error) {
	appUser := shared.GetUser()
	if appUser == nil {
		return nil, errors.New("Cannot get user emotes without being logged in")
	}

	if set, exists := cache.GetEmoteSet(cache.TWITCH_KEY, cache.USER_EMOTE_SECTION); exists {
		return set, nil
	}

	params := map[string][]string{
		"user_id": {appUser.Id},
	}

	res, err := nativeApi.GetNativeUserEmotes(appUser.Access_token, params)
	if err != nil {
		log.Printf("[GetUserEmotes]: An error occurred fetching user emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	// TODO: move status checks to api method
	if res.Status != 200 {
		log.Printf("[GetUserEmotes]: Failed to get user emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[nativeApi.ApiGetUserEmotesRes]{ Res: res }
	}

	emotes := types.AppEmoteMap{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		appEmote := *nativeApi.GetAppEmoteFromApiEmote(nativeApi.ApiEmote(emote), tmpl)
		appEmote.Section = cache.USER_EMOTE_SECTION
		emotes[appEmote.Name] = &appEmote
	}

	set := &types.AppEmoteSet{
		Provider: cache.TWITCH_KEY,
		Emotes: emotes,
		Id: appUser.Id,
		Section: cache.USER_EMOTE_SECTION,
	}

	cache.SetEmoteSet(cache.TWITCH_KEY, cache.USER_EMOTE_SECTION, set)

	return set, nil

}

func (es *EmoteService) RequestTwitchUserEmotes() error {
	set, err := es.GetTwitchUserEmotes()
	if err != nil {
		return err
	}

	shared.EmitNewSet(es.app, set, false, "")
	return nil
}

func (es *EmoteService) GetTwitchGlobalEmotes() (*types.AppEmoteSet, error) {
	appUser := shared.GetUser()
	if appUser == nil {
		return nil, errors.New("Cannot get global emotes without being logged in")
	}

	if set, exists := cache.GetEmoteSet(cache.TWITCH_KEY, cache.GLOBAL_EMOTE_SECTION); exists {
		return set, nil
	}

	res, err := nativeApi.GetNativeGlobalEmotes(appUser.Access_token, map[string][]string{})
	if err != nil {
		log.Printf("[GetGlobalEmotes]: An error occurred fetching global emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	// TODO: move status checks to api call
	if res.Status != 200 {
		log.Printf("[GetGlobalEmotes]: Failed to get global emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[nativeApi.ApiGetGlobalEmotesRes]{ Res: res }
	}


	emotes := types.AppEmoteMap{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		appEmote := *nativeApi.GetAppEmoteFromApiEmote(nativeApi.ApiEmote(emote), tmpl)
		appEmote.Section = cache.GLOBAL_EMOTE_SECTION
		emotes[appEmote.Name] = &appEmote
	}

	set := &types.AppEmoteSet{
		Section: cache.GLOBAL_EMOTE_SECTION,
		Provider: cache.TWITCH_KEY,
		Emotes: emotes,
		Id: cache.GLOBAL_EMOTE_SECTION,
	}

	cache.SetEmoteSet(cache.TWITCH_KEY, cache.GLOBAL_EMOTE_SECTION, set)

	return set, nil
}

func (es *EmoteService) RequestTwitchGlobalEmotes() error {
	set, err := es.GetTwitchGlobalEmotes()
	if err != nil {
		return err
	}

	shared.EmitNewSet(es.app, set, false, "")
	return nil
}

func (es *EmoteService) GetTwitchChannelEmotes(broadcaster_id string) (*types.AppEmoteSet, error) {
	appUser := shared.GetUser()
	if appUser == nil {
		return nil, errors.New("Cannot get channel emotes without logging in")
	}

	if set, exists := cache.GetChannelEmoteSet(cache.TWITCH_KEY, broadcaster_id); exists {
		shared.EmitNewSet(es.app, set, true, broadcaster_id)
		return set, nil
	}

	params := map[string][]string{
		"broadcaster_id": {broadcaster_id},
	}

	res, err := nativeApi.GetNativeChannelEmotes(appUser.Access_token, params)
	if err != nil {
		log.Printf("[GetGlobalEmotes]: An error occurred fetching global emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetGlobalEmotes]: Failed to get global emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[nativeApi.ApiGetChannelEmotesRes]{ Res: res }
	}


	emotes := types.AppEmoteMap{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		appEmote := *nativeApi.GetAppEmoteFromApiEmote(nativeApi.ApiEmote(emote), tmpl)
		appEmote.Section = cache.CHANNEL_EMOTE_SECTION
		emotes[appEmote.Name] = &appEmote
	}

	set := &types.AppEmoteSet{
		Provider: cache.TWITCH_KEY,
		Emotes: emotes,
		Section: cache.CHANNEL_EMOTE_SECTION,
		Id: broadcaster_id,
	}

	cache.SetChannelEmoteSet(cache.TWITCH_KEY, broadcaster_id, set)

	shared.EmitNewSet(es.app, set, true, broadcaster_id)

	return set, nil
}

func (es *EmoteService) RequestTwitchChannelEmotes(broadcasterId string) error {
	set, err := es.GetTwitchChannelEmotes(broadcasterId)
	if err != nil {
		return err
	}

	shared.EmitNewSet(es.app, set, true, broadcasterId)
	return nil
}
