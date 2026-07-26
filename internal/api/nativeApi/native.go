package nativeApi

import (
	"chatter-wails/internal/api"
	"net/http"
	"net/url"
	"time"
)

var (
	VALIDATE_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "id.twitch.tv",
		Path: "/oauth2/validate",
	}

	SUBSCRIPTIONS_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/eventsub/subscriptions",
	}

	USERS_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/users",
	}

	MESSAGES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/chat/messages",
	}

	BADGES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/chat/badges",
	}

	GLOBAL_BADGES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/chat/badges/global",
	}

	EMOTES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/chat/emotes",
	}

	USER_EMOTES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/chat/emotes/user",
	}

	GLOBAL_EMOTES_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/chat/emotes/global",
	}

	STREAMS_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/streams",
	}

	SHARED_CHAT_SESSION_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "api.twitch.tv",
		Path: "/helix/shared_chat/session",
	}
)

var apiSubscriptionHeaders = api.ApiHeaders;

var apiUsersHeaders = api.ApiHeaders;

var apiMessagesHeaders = api.ApiHeaders;

var apiBadgesHeaders = api.ApiHeaders;

var apiEmotesHeaders = api.ApiHeaders;

var apiStreamsHeaders = api.ApiHeaders;

func apiValidateHeaders(access_token string) *http.Header {
    return &http.Header{
        "Authorization" : {"OAuth " + access_token},
    }
}



func ApiGetNativeValidate(
        access_token string,
    ) (*api.ApiResponse[any], error) {
    return api.ApiGet[any](
		VALIDATE_ENDPOINT,
		apiValidateHeaders(access_token),
		map[string][]string{},
	)
}



type ApiSubscriptionTransport struct{
	Method string			`json:"method"`
	Callback string			`json:"callback"`
	Session_id string		`json:"session_id"`
}

type ApiSubscription struct{
	Id string							`json:"id"`
	Status string						`json:"status"`
	Sub_type string						`json:"type"`
	Version string						`json:"version"`
	Condition map[string]any			`json:"condition"`
	Created_at string					`json:"created_at"`
    Transport ApiSubscriptionTransport	`json:"transport"`
	Cost int							`json:"cost"`
}


type ApiPagination struct{
	Cursor string			`json:"cursor"`
}

type ApiGetSubscriptionsRes struct{
	Data []ApiSubscription		`json:"data"`
	Connected_at string			`json:"connected_at"`
	Disconnected_at string		`json:"disconnected_at"`
	Total int					`json:"total"`
	Total_cost int				`json:"total_cost"`
	Max_total_cost int			`json:"max_total_cost"`
	Pagination ApiPagination 	`json:"pagination"`
}

func ApiGetNativeSubscriptions(
    access_token string,
    params map[string][]string,
	) (*api.ApiResponse[ApiGetSubscriptionsRes], error) {
    return api.ApiGet[ApiGetSubscriptionsRes](
			SUBSCRIPTIONS_ENDPOINT,
			apiSubscriptionHeaders(access_token),
			params,
    )
}



type ApiPostSubscriptionsRes struct{
	Data []struct{
		Id string 				`json:"id"`
		Status string			`json:"status"`
		Sub_type string			`json:"type"`
		Version string			`json:"version"`
		Condition struct{}		`json:"condition"`
		Created_at string		`json:"created_at"`
		Transport struct{
			Method string 		`json:"method"`
			Session_id string 	`json:"session_id"`
			Connected_at string	`json:"connected_at"`
		}						`json:"transport"`
		Cost int				`json:"cost"`

	}							`json:"data"`
	Total int					`json:"total"`
	Total_cost int				`json:"total_cost"`
	Max_total_cost int			`json:"max_total_cost"`
}

type ApiPostSubscriptionsBody struct {
	Sub_type string								`json:"type"`
	Version string								`json:"version"`
	Condition map[string]string					`json:"condition"`
	Transport ApiPostSubscriptionsBodyTransport	`json:"transport"`
}

type ApiPostSubscriptionsBodyTransport struct {
	Method string			`json:"method"`
	Session_id string		`json:"session_id"`
}

func ApiPostNativeSubscriptions(
    access_token string,
    body ApiPostSubscriptionsBody,
    params map[string][]string,
    ) (*api.ApiResponse[ApiPostSubscriptionsRes], error) {
    return api.ApiPost[ApiPostSubscriptionsRes](
        SUBSCRIPTIONS_ENDPOINT,
        apiSubscriptionHeaders(access_token),
        body,
        params,
    )
}

func ApiDeleteNativeSubscriptions(
    access_token string,
    params map[string][]string,
    ) (*api.ApiResponse[any], error) {
    return api.ApiDelete[any](
        SUBSCRIPTIONS_ENDPOINT,
        apiSubscriptionHeaders(access_token),
        params,
    )
}



type ApiUser struct {
	Id string					`json:"id"`
    Login string				`json:"login"`
    Display_name string			`json:"display_name"`
    User_type string			`json:"type"`
    Broadcaster_type string		`json:"broadcaster_type"`
    Description string			`json:"description"`
    Profile_image_url string	`json:"profile_image_url"`
    Offline_image_url string	`json:"offline_image_url"`
    View_count int				`json:"view_count"`
	Email string				`json:"email"`
    Created_at time.Time		`json:"created_at"`
    // Access_token string			`json:"access_token"`
}

type ApiGetUsersRes struct{
	Data []ApiUser				`json:"data"`
}

func GetNativeUsers(
    access_token string,
    params map[string][]string,
    ) (*api.ApiResponse[ApiGetUsersRes], error) {
    return api.ApiGet[ApiGetUsersRes](
        USERS_ENDPOINT,
        apiUsersHeaders(access_token),
        params,
    )
}

func PostNativeUsers(
    access_token string,
    body any,
    params map[string][]string,
    ) (*api.ApiResponse[any], error) {
    return api.ApiPost[any](
        USERS_ENDPOINT,
        apiUsersHeaders(access_token),
        body,
        params,
    )
}


type ApiPostMessagesBody struct{
	Broadcaster_id string				`json:"broadcaster_id"`
	Sender_id string					`json:"sender_id"`
	Message string 						`json:"message"`
	Reply_parent_message_id *string		`json:"reply_parent_message_id"`
}

type ApiMessageDropReason struct{
	Code string				`json:"code"`
	Message string			`json:"message"`
}

type ApiPostMessagesData struct{
	Message_id string						`json:"message_id"`
	Is_sent bool							`json:"is_sent"`
	Drop_reason *ApiMessageDropReason		`json:"drop_reason,omitempty"`
};

type ApiPostMessagesRes struct{
	Data []ApiPostMessagesData
}

func PostMessage(
    access_token string,
    body ApiPostMessagesBody,
    params map[string][]string,
    ) (*api.ApiResponse[ApiPostMessagesRes], error) {
    return api.ApiPost[ApiPostMessagesRes](
        MESSAGES_ENDPOINT,
        apiMessagesHeaders(access_token),
        body,
        params,
    )
}


type ApiGetChannelBadgesRes = ApiGetGlobalBadgesRes

func GetNativeChannelBadges(
    access_token string,
    params map[string][]string,
    ) (*api.ApiResponse[ApiGetChannelBadgesRes], error) {
    return api.ApiGet[ApiGetChannelBadgesRes](
        BADGES_ENDPOINT,
        apiBadgesHeaders(access_token),
        params,
    )
}


type ApiBadgeSetVersions struct{
	Id string				`json:"id"`
	Image_url_1x string		`json:"image_url_1x"`
	Image_url_2x string		`json:"image_url_2x"`
	Image_url_4x string		`json:"image_url_4x"`
	Title string			`json:"title"`
	Description string		`json:"description"`
	Click_action *string	`json:"click_action,omitempty"`
	Click_url *string		`json:"click_url,omitempty"`
}
type ApiBadgeSet struct{
	Set_id string						`json:"set_id"`
    Versions []ApiBadgeSetVersions		`json:"versions"`
}

type ApiGetGlobalBadgesRes struct{
	Data []ApiBadgeSet
}

func GetNativeGlobalBadges(
    access_token string,
    ) (*api.ApiResponse[ApiGetGlobalBadgesRes], error) {
    return api.ApiGet[ApiGetGlobalBadgesRes](
        GLOBAL_BADGES_ENDPOINT,
        apiBadgesHeaders(access_token),
		map[string][]string{},
    )
}

  // "data": [
  //   {
  //     "emote_set_id": "",
  //     "emote_type": "hypetrain",
  //     "format": [
  //       "static"
  //     ],
  //     "id": "304420818",
  //     "name": "HypeLol",
  //     "owner_id": "477339272",
  //     "scale": [
  //       "1.0",
  //       "2.0",
  //       "3.0"
  //     ],
  //     "theme_mode": [
  //       "light",
  //       "dark"
  //     ]
  //   }
  // ],
  // "template": "https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}",
  // "pagination": {
  //   "cursor": "eyJiIjpudWxsLJxhIjoiIn0gf5"
  // }

// important api emote fields (more are commented out below)
type ApiEmote struct{
	Id string				`json:"id"`
	Name string				`json:"name"`
	Format []string			`json:"format"`
	Scale []string			`json:"scale"`
	Theme_mode []string 	`json:"theme_mode"`
}

type ApiUserEmote ApiEmote

// type ApiUserEmote struct{
// 	Emote_set_id string 	`json:"emote_set_id"`
// 	Emote_type string		`json:"emote_type"`
// 	Format []string			`json:"format"`
// 	Id string				`json:"id"`
// 	Name string				`json:"name"`
// 	Owner_id string			`json:"owner_id"`
// 	Scale []string			`json:"scale"`
// 	Theme_mode []string 	`json:"theme_mode"`
// }

type ApiGetUserEmotesRes struct{
	Data []ApiUserEmote				`json:"data"`
	Template string					`json:"template"`
	Pagination ApiPagination		`json:"pagination"`
}

func GetNativeUserEmotes(
    access_token string,
    params map[string][]string,
    ) (*api.ApiResponse[ApiGetUserEmotesRes], error) {
    return api.ApiGet[ApiGetUserEmotesRes](
		USER_EMOTES_ENDPOINT,
		apiEmotesHeaders(access_token),
		params,
	)
}




type ApiGlobalEmoteImages struct{
	Url_1x string		`json:"url_1x"`
	Url_2x string		`json:"url_2x"`
	Url_4x string		`json:"url_4x"`
}

type ApiGlobalEmote ApiEmote
// type ApiGlobalEmote struct{
// 	Id string					`json:"id"`
// 	Name string					`json:"name"`
//     Images ApiGlobalEmoteImages	`json:"images"`
// 	Format []string				`json:"format"`
// 	Scale []string				`json:"scale"`
// 	Theme_mode []string			`json:"theme_mode"`
// }

type ApiGetGlobalEmotesRes struct{
	Data []ApiGlobalEmote 	`json:"data"`
	Template string			`json:"template"`
}

func GetNativeGlobalEmotes(
    access_token string,
    params map[string][]string,
    ) (*api.ApiResponse[ApiGetGlobalEmotesRes], error) {
    return api.ApiGet[ApiGetGlobalEmotesRes](
		GLOBAL_EMOTES_ENDPOINT,
		apiEmotesHeaders(access_token),
		params,
	)
}

type ApiChannelEmote ApiEmote

type ApiGetChannelEmotesRes struct{
	Data []ApiChannelEmote 	`json:"data"`
	Template string			`json:"template"`
}

func GetNativeChannelEmotes(
    access_token string,
    params map[string][]string,
    ) (*api.ApiResponse[ApiGetChannelEmotesRes], error) {
    return api.ApiGet[ApiGetChannelEmotesRes](
		EMOTES_ENDPOINT,
		apiEmotesHeaders(access_token),
		params,
	)
}

type ApiStream struct{
	Id string		`json:"id"`
	User_id string		`json:"user_id"`
	User_login string		`json:"user_login"`
	User_name string		`json:"user_name"`
	Game_id string		`json:"game_id"`
	Game_name string		`json:"game_name"`
	Type string		`json:"type"`
	Title string		`json:"title"`
	Tags []string
	Viewer_count int		`json:"viewer_count"`
	Started_at string		`json:"started_at"`
	Language string		`json:"language"`
	Thumbnail_url string		`json:"thumbnail_url"`
}

type ApiGetStreamsRes struct{
	Data []ApiStream				`json:"data"`
	Pagination ApiPagination	`json:"pagination"`
}

func GetStreams(
	access_token string,
	params map[string][]string,
	) (*api.ApiResponse[ApiGetStreamsRes], error) {
	return api.ApiGet[ApiGetStreamsRes](
		STREAMS_ENDPOINT,
		apiStreamsHeaders(access_token),
		params,
	)
}

type ApiSharedChatSessionParticipant struct{
	Broadcaster_id string		`json:"broadcaster_id"`
}

type ApiSharedChatSession struct{
	Session_id string		`json:"session_id"`
	Host_broadcaster_id string		`json:"host_broadcaster_id"`
	Participants []ApiSharedChatSessionParticipant		`json:"participants"`
	Created_at string		`json:"created_at"`
	Updated_at string		`json:"updated_at"`
}

type ApiGetSharedChatSessionRes struct{
	Data []ApiSharedChatSession		`json:"data"`
}

func GetSharedChatSession(
	access_token string,
	params map[string][]string,
	) (*api.ApiResponse[ApiGetSharedChatSessionRes], error) {
	return api.ApiGet[ApiGetSharedChatSessionRes](
		SHARED_CHAT_SESSION_ENDPOINT,
		apiStreamsHeaders(access_token),
		params,
	)
}
