package bttvApi

import (
	"chatter-wails/internal/api"
	"chatter-wails/shared/cache"
	"chatter-wails/shared/types"
	"fmt"
	"net/url"
	"strings"
)

var (
	BTTV_GLOBAL_EMOTES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.betterttv.net",
		Path: "/3/cached/emotes/global",
	}

	BTTV_USERS_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.betterttv.net",
		Path: "/3/cached/users",
	}

	BTTV_EMOTE_CDN = url.URL{
		Scheme: "https",
		Host: "cdn.betterttv.net",
		Path: "/emote",
	}

	BTTV_TWITCH_PROVIDER = "twitch"
)

type BTTVPartialUser struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	ProviderId  string `json:"providerId"`
}

type BTTVEmote struct {
	Id        string           `json:"id"`
	Code      string           `json:"code"`
	ImageType string           `json:"imageType"`
	Animated  bool             `json:"animated"`
	UserId    *string          `json:"userId,omitempty"`
	User      *BTTVPartialUser `json:"user,omitempty"`
}

type BTTVUser struct {
	Id            string      `json:"id"`
	Bots          []string    `json:"bots"`
	Avatar        string      `json:"avatar"`
	ChannelEmotes []BTTVEmote `json:"channelEmotes"`
	SharedEmotes  []BTTVEmote `json:"sharedEmotes"`
}



func constructBttvUserEndpointURL(provider string, providerId string) *url.URL {
	var endpoint url.URL
	endpoint.Scheme = BTTV_USERS_ENDPOINT.Scheme
	endpoint.Host = BTTV_USERS_ENDPOINT.Host
	endpoint.Path = fmt.Sprintf("%s/%s/%s", BTTV_USERS_ENDPOINT.Path, provider, providerId)

	return &endpoint
}

type ApiGetBTTVUserRes BTTVUser

func GetUser(provider string, providerId string) (*ApiGetBTTVUserRes, error) {
	res, err := api.ApiGet[ApiGetBTTVUserRes](*constructBttvUserEndpointURL(provider, providerId), nil, map[string][]string{})
	if err != nil {
		return nil, err
	}

	return &res.Body, nil
}

type ApiGetBTTVGlobalEmotesRes []BTTVEmote

func GetGlobalEmotes() ([]BTTVEmote, error) {
	res, err := api.ApiGet[ApiGetBTTVGlobalEmotesRes](BTTV_GLOBAL_EMOTES_ENDPOINT, nil, map[string][]string{})
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func BTTVEmotesToAppEmoteSet(emotes []BTTVEmote, sectionId string) *types.AppEmoteSet {
	set := &types.AppEmoteSet{
		Id: fmt.Sprintf("%s:%s", cache.BTTV_KEY, sectionId),
		Provider: cache.BTTV_KEY,
		Section: sectionId,
		Emotes: types.AppEmoteMap{},
	}
	
	for _, emote := range emotes {
		appEmote := BTTVEmoteToAppEmote(&emote, sectionId)
		set.Emotes[appEmote.Name] = appEmote
	}

	return set
}

// sectionId is needed since bttv global emotes still have userId associated with them
func BTTVEmoteToAppEmote(emote *BTTVEmote, sectionId string) *types.AppEmote {
	srcSet := ExtractBTTVSrcSet(emote)

	return &types.AppEmote{
		Id: emote.Id,
		Name: emote.Code,
		LightSrcSet: srcSet,
		DarkSrcSet: srcSet,
		Provider: cache.BTTV_KEY,
		Section: sectionId,
	}
}

func ExtractBTTVSrcSet(emote *BTTVEmote) string {
	// {cdn_base}/{id}/{scale}x.{imagetype}

	srcs := []string{}
	for i := 1; i <= 3; i++ {
		src := fmt.Sprintf("%s/%s/%dx.%s", BTTV_EMOTE_CDN.String(), emote.Id, i, emote.ImageType)
		srcs = append(srcs, src)
	}

	return strings.Join(srcs, ", ")
}
