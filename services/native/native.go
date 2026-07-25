package native

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/api/nativeApi"
	"chatter-wails/internal/user"
	"chatter-wails/shared"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"context"
	"errors"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type EmoteService struct {
	app *application.App
	Ctx context.Context
	GlobalEmotes *[]types.AppEmote
	UserEmotes *[]types.AppEmote
}

func NewEmoteService(app *application.App) *EmoteService {
	return &EmoteService{app: app}
}

func (es *EmoteService) GetUserEmotes() (*[]types.AppEmote, error) {
	appUser := shared.GetUser()
	if appUser == nil {
		return nil, errors.New("Cannot get user emotes without being logged in")
	}
	user, err := user.GetUserByToken(appUser.Access_token)
	if err != nil {
		return nil, err
	}

	params := map[string][]string{
		"user_id": {user.Id},
	}

	res, err := nativeApi.ApiGetUserEmotes(appUser.Access_token, params)
	if err != nil {
		log.Printf("[GetUserEmotes]: An error occurred fetching user emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetUserEmotes]: Failed to get user emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[nativeApi.ApiGetUserEmotesRes]{ Res: res }
	}

	emotes := []types.AppEmote{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		appEmote := *nativeApi.GetAppEmoteFromApiEmote(nativeApi.ApiEmote(emote), tmpl)
		appEmote.Type = "user"
		emotes = append(emotes, appEmote)
	}

	es.UserEmotes = &emotes

	return &emotes, nil

}

func (es *EmoteService) GetGlobalEmotes() (*[]types.AppEmote, error) {
	appUser := shared.GetUser()
	if appUser == nil {
		return nil, errors.New("Cannot get global emotes without being logged in")
	}

	res, err := nativeApi.ApiGetGlobalEmotes(appUser.Access_token, map[string][]string{})
	if err != nil {
		log.Printf("[GetGlobalEmotes]: An error occurred fetching global emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetGlobalEmotes]: Failed to get global emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[nativeApi.ApiGetGlobalEmotesRes]{ Res: res }
	}


	emotes := []types.AppEmote{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		appEmote := *nativeApi.GetAppEmoteFromApiEmote(nativeApi.ApiEmote(emote), tmpl)
		appEmote.Type = "global"
		emotes = append(emotes, appEmote)
	}

	es.GlobalEmotes = &emotes

	return &emotes, nil
}

func (es *EmoteService) GetChannelEmotes(broadcaster_id string) (*types.AppEmoteSet, error) {
	appUser := shared.GetUser()
	if appUser == nil {
		return nil, errors.New("Cannot get channel emotes without logging in")
	}

	if set, exists := cache.GetEmoteSet(cache.NATIVE_KEY, broadcaster_id); exists {
		shared.EmitNewSet(es.app, set, broadcaster_id)
		return set, nil
	}

	params := map[string][]string{
		"broadcaster_id": {broadcaster_id},
	}

	res, err := nativeApi.ApiGetChannelEmotes(appUser.Access_token, params)
	if err != nil {
		log.Printf("[GetGlobalEmotes]: An error occurred fetching global emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetGlobalEmotes]: Failed to get global emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[nativeApi.ApiGetChannelEmotesRes]{ Res: res }
	}


	emotes := map[string]*types.AppEmote{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		appEmote := *nativeApi.GetAppEmoteFromApiEmote(nativeApi.ApiEmote(emote), tmpl)
		appEmote.Type = "channel"
		emotes[appEmote.Name] = &appEmote
	}

	set := &types.AppEmoteSet{
		Provider: "channel",
		Emotes: emotes,
		Id: broadcaster_id,
	}

	cache.SetEmoteSet(cache.NATIVE_KEY, broadcaster_id, set)

	shared.EmitNewSet(es.app, set, broadcaster_id)

	return set, nil
}
