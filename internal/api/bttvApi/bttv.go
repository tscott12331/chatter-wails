package bttvApi

import (
	"chatter-wails/internal/api"
	"fmt"
	"net/url"
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
	UserId    *string          `json:"userId"`
	User      *BTTVPartialUser `json:"user"`
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

