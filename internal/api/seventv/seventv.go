package seventv

import (
	"chatter-wails/internal/api"
	"fmt"
	"net/url"
)

var (
	
	USER_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "7tv.io",
		// "/v3/users/{platform}/{platform_id}",
		Path: "/v3/users",
	}
)

func constructUserEndpointURL(platform string, platform_id string) *url.URL {
	var endpoint url.URL
	endpoint.Scheme = USER_ENDPOINT.Scheme
	endpoint.Host = USER_ENDPOINT.Host
	endpoint.Path = fmt.Sprintf("%s/%s/%s", USER_ENDPOINT.Path, platform, platform_id)

	return &endpoint
}


type SevenTVUser struct{
	
}

func GetSevenTVUser(platform string, platform_id string)*SevenTVUser {
	endpoint := constructUserEndpointURL(platform, platform_id)
	
	res, err := api.ApiGet[any](*endpoint, nil, map[string][]string{})
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("%+v", res)

	return &SevenTVUser{}
}
