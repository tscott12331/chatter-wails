package ffzapi

import (
	"chatter-wails/internal/api"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	FFZ_ROOM_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.frankerfacez.com",

		// /v1/room/id/:twitchId
		Path: "/v1/room/id/",
	}

	FFZ_GLOBAL_EMOTE_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.frankerfacez.com",

		// https://api.frankerfacez.com/v1/set/global/ids
		Path: "/v1/set/global/ids",
	}
)


type FFZRoom struct{
    Obj_id int			`json:"_id"`
    Twitch_id int			`json:"twitch_id"`
    Youtube_id *string			`json:"youtube_id"`
    Id string			`json:"id"`
    Is_group bool			`json:"is_group"`
    Display_name *string			`json:"display_name"`
    Set int			`json:"set"`
    Moderator_badge *string			`json:"moderator_badge"`
    Vip_badge *map[string]string			`json:"vip_badge"`
    Mod_urls *map[string]string			`json:"mod_urls"`
    // "user_badges": {},
    // "user_badge_ids": {},
    Css *string			`json:"css"`
}

type FFZOwner struct{
	Obj_id int			`json:"_id"`
	Name string			`json:"name"`
	Display_name *string			`json:"display_name"`
}

type FFZEmote struct{
	  Id int			`json:"id"`
	  Name string			`json:"name"`
	  Height int			`json:"height"`
	  Width int			`json:"width"`
	  Public bool			`json:"public"`
	  Hidden bool			`json:"hidden"`
	  Modifier bool			`json:"modifier"`
	  Modifier_flags int			`json:"modifier_flags"`
	  Offset *string			`json:"offset"`
	  Margins *string			`json:"margins"`
	  Css *string			`json:"css"`
	  Owner *FFZOwner			`json:"owner"`
	  Artist *FFZOwner			`json:"artist"`
	  Urls map[string]string			`json:"urls"`
	  Status int			`json:"status"`
	  Usage_count int			`json:"usage_count"`
	  Created_at string			`json:"created_at"`
	  Last_updated string			`json:"last_updated"`
}

type FFZEmoteSet struct{
      Id int			`json:"id"`
      Obj_type int			`json:"_type"`
      Icon *string			`json:"icon"`
      Title string			`json:"title"`
      Css *string			`json:"css"`
      Emoticons []FFZEmote			`json:"emoticons"`
}

type GetFFZRoomRes struct{
  Room FFZRoom			`json:"room"`
  Sets map[string]FFZEmoteSet			`json:"sets"`
}


func GetFFZRoom(provider, broadcasterId string) (*GetFFZRoomRes, error) {
	if provider == cache.TWITCH_KEY {
		return GetFFZTwitchRoom(broadcasterId)
	} else {
		return nil, errors.New("WARNING: fetching FFZ channel emotes for provider %s is not implemented")
	}
}

func GetFFZTwitchRoom(broadcasterId string) (*GetFFZRoomRes, error) {
	endpoint := FFZ_ROOM_ENDPOINT.JoinPath(fmt.Sprintf("/%s", broadcasterId))
	res, err := api.ApiGet[GetFFZRoomRes](*endpoint, nil, map[string][]string{})
	if err != nil {
		return nil, err
	}

	if res.Status != 200 {
		return nil, &api.StatusError[GetFFZRoomRes]{ Res: res }
	}

	return &res.Body, nil
}



type GetFFZGlobalEmoteSetRes struct{
	Default_sets []int			`json:"default_sets"`
	Sets map[string]FFZEmoteSet			`json:"sets"`

	// not gonna populate this
	// user_ids map[string][]int
}

func GetFFZGlobalEmoteSet() (*GetFFZGlobalEmoteSetRes, error) {
	res, err := api.ApiGet[GetFFZGlobalEmoteSetRes](FFZ_GLOBAL_EMOTE_ENDPOINT, nil, map[string][]string{})
	if err != nil {
		return nil, err
	}

	if res.Status != 200 {
		return nil, &api.StatusError[GetFFZGlobalEmoteSetRes]{ Res: res }
	}

	return &res.Body, nil
}


// TODO: consider redesigning cache to allow for more sections, 
// FFZ may have multiple sets under "global" or "channel"J
func GetAppEmoteSetFromFFZEmoteSet(emoteSet *FFZEmoteSet, sectionId string) *types.AppEmoteSet {
	emoteMap := types.AppEmoteMap{}

	for _, emote := range emoteSet.Emoticons {
		appEmote := FFZEmoteToAppEmote(&emote, sectionId)
		emoteMap[appEmote.Name] = appEmote
	}

	appSet := &types.AppEmoteSet{
		Id: FFZIdToStringId(emoteSet.Id),
		Provider: cache.FFZ_KEY,
		Section: sectionId,
		Emotes: emoteMap,
	}

	return appSet
}

func FFZEmoteToAppEmote(emote *FFZEmote, sectionId string) *types.AppEmote {
	// TODO: add modifier data to AppEmote
	srcSet := ExtractFFZEmoteSrcSet(emote)
	appEmote := &types.AppEmote{
		Id: FFZIdToStringId(emote.Id),
		Name: emote.Name,
		LightSrcSet: srcSet,
		DarkSrcSet: srcSet,
		Provider: cache.FFZ_KEY,
		Section: sectionId,
		ZeroWidth: emote.Modifier,
	}

	return appEmote
}

func ExtractFFZEmoteSrcSet(emote *FFZEmote) string {
	if len(emote.Urls) == 0 {
		return ""
	}

	srcs := make([]string, len(emote.Urls))
	i := 0
	for scale, u := range emote.Urls {
		srcs[i] = fmt.Sprintf("%s %sx", u, scale)
		i += 1
	}

	return strings.Join(srcs, ", ")
}

type ffzid interface {
	int|string
}
func FFZIdToStringId[T ffzid](id T) string {
	return fmt.Sprintf("%s:%v", cache.FFZ_KEY, id)
}
