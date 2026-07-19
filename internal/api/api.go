package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type StatusError[T any] struct {
	Res *ApiResponse[T]
}
func (se *StatusError[any]) Error() string {
	return fmt.Sprintf("%v, %v", se.Res.Status, se.Res)
}


type ApiSessionReqTimeout struct{}
func (esrt *ApiSessionReqTimeout) Error() string {
	return "The request for a session ID timed out"
}


type ApiResponse[T any] struct{
	Status int
	Body T
}



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

var CLIENT_ID string


func apiHeaders(access_token string) *http.Header {
    return &http.Header{
        "Authorization": {"Bearer " + access_token},
        "Client-Id": {CLIENT_ID},
    }
}

var apiSubscriptionHeaders = apiHeaders;

var apiUsersHeaders = apiHeaders;

var apiMessagesHeaders = apiHeaders;

var apiBadgesHeaders = apiHeaders;

var apiEmotesHeaders = apiHeaders;

var apiStreamsHeaders = apiHeaders;

var apiSharedChatHeaders = apiHeaders;

var apiPollHeaders = apiHeaders;

func apiValidateHeaders(access_token string) *http.Header {
    return &http.Header{
        "Authorization" : {"OAuth " + access_token},
    }
}

func init() {
	godotenv.Load()

	CLIENT_ID, _ = os.LookupEnv("VITE_CLIENT_ID")
}


func ApiFetch[T any](
	method string,
    endpoint url.URL,
    headers *http.Header,
	body any,
    params map[string][]string,
	) (*ApiResponse[T], error) {
		var req_body *bytes.Buffer = nil

		var hasBody = body != nil

		if hasBody {
			req_body_json, err := json.Marshal(body)
			if err != nil {
				log.Printf("[ApiFetch]: An error occurred marshaling the request body, aborting\n\n")
				return nil, err
			}

			req_body = bytes.NewBuffer(req_body_json)
		}

		encParams := url.Values{}
		encParams = params
		
		endpoint.RawQuery = encParams.Encode()
		

		var req *http.Request
		var err error

		if req_body != nil {
			req, err = http.NewRequest(method, endpoint.String(), req_body)
		} else {
			req, err = http.NewRequest(method, endpoint.String(), nil)
		}

		if err != nil {
			log.Printf("[ApiFetch]: An error occurred creating the request, aborting\n\n")
			return nil, err
		}

		if headers != nil {
			req.Header = *headers
		}
		if hasBody {
			req.Header.Set("Content-Type", "application/json")
		}

		httpClient := &http.Client{Timeout: time.Second * 10}

		log.Printf("[ApiFetch]: %s %s\n\n", method, endpoint.String())
		res, err := httpClient.Do(req)
		if err != nil {
			log.Printf("[ApiFetch]: An error occurred making the post request\n\n")
			return nil, err
		}
		defer res.Body.Close()

		res_body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("[ApiFetch]: An error occurred reading the response body, aborting\n\n")
			return nil, err
		}

		var res_body_obj T
		if len(res_body) > 0 {
			err = json.Unmarshal(res_body, &res_body_obj)
			if err != nil {
				log.Printf("[ApiFetch]: An error occurred parsing the response body, aborting\n\n")
				return nil, err
			}
		}

		return &ApiResponse[T]{
			Status: res.StatusCode,
			Body: res_body_obj,
		}, nil
}


func ApiDelete[T any](
    endpoint url.URL,
    headers *http.Header,
    params map[string][]string,
    ) (*ApiResponse[T], error) {
		return ApiFetch[T]("DELETE", endpoint, headers, nil, params)
}

func ApiGet[T any](
    endpoint url.URL,
    headers *http.Header,
    params map[string][]string,
    ) (*ApiResponse[T], error) {
		return ApiFetch[T]("GET", endpoint, headers, nil, params)
}

func ApiPost[T any](
    endpoint url.URL,
    headers *http.Header,
	body any,
    params map[string][]string,
    ) (*ApiResponse[T], error) {
		return ApiFetch[T]("POST", endpoint, headers, body, params)
}



func ApiGetValidate(
        access_token string,
    ) (*ApiResponse[any], error) {
    return ApiGet[any](
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

func ApiGetSubscriptions(
    access_token string,
    params map[string][]string,
	) (*ApiResponse[ApiGetSubscriptionsRes], error) {
    return ApiGet[ApiGetSubscriptionsRes](
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

func ApiPostSubscriptions(
    access_token string,
    body ApiPostSubscriptionsBody,
    params map[string][]string,
    ) (*ApiResponse[ApiPostSubscriptionsRes], error) {
    return ApiPost[ApiPostSubscriptionsRes](
        SUBSCRIPTIONS_ENDPOINT,
        apiSubscriptionHeaders(access_token),
        body,
        params,
    )
}

func ApiDeleteSubscriptions(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[any], error) {
    return ApiDelete[any](
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

func ApiGetUsers(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[ApiGetUsersRes], error) {
    return ApiGet[ApiGetUsersRes](
        USERS_ENDPOINT,
        apiUsersHeaders(access_token),
        params,
    )
}

func ApiPostUsers(
    access_token string,
    body any,
    params map[string][]string,
    ) (*ApiResponse[any], error) {
    return ApiPost[any](
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

func ApiPostMessages(
    access_token string,
    body ApiPostMessagesBody,
    params map[string][]string,
    ) (*ApiResponse[ApiPostMessagesRes], error) {
    return ApiPost[ApiPostMessagesRes](
        MESSAGES_ENDPOINT,
        apiMessagesHeaders(access_token),
        body,
        params,
    )
}


type ApiGetChannelBadgesRes = ApiGetGlobalBadgesRes

func ApiGetChannelBadges(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[ApiGetChannelBadgesRes], error) {
    return ApiGet[ApiGetChannelBadgesRes](
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

func ApiGetGlobalBadges(
    access_token string,
    ) (*ApiResponse[ApiGetGlobalBadgesRes], error) {
    return ApiGet[ApiGetGlobalBadgesRes](
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

func ApiGetUserEmotes(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[ApiGetUserEmotesRes], error) {
    return ApiGet[ApiGetUserEmotesRes](
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

func ApiGetGlobalEmotes(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[ApiGetGlobalEmotesRes], error) {
    return ApiGet[ApiGetGlobalEmotesRes](
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

func ApiGetChannelEmotes(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[ApiGetChannelEmotesRes], error) {
    return ApiGet[ApiGetChannelEmotesRes](
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

func ApiGetStreams(
	access_token string,
	params map[string][]string,
	) (*ApiResponse[ApiGetStreamsRes], error) {
	return ApiGet[ApiGetStreamsRes](
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

func ApiGetSharedChatSession(
	access_token string,
	params map[string][]string,
	) (*ApiResponse[ApiGetSharedChatSessionRes], error) {
	return ApiGet[ApiGetSharedChatSessionRes](
		SHARED_CHAT_SESSION_ENDPOINT,
		apiSharedChatHeaders(access_token),
		params,
	)
}

type ApiPollChoice struct {
	Id string			`json:"id"`
	Title string			`json:"title"`
	Votes int			`json:"votes"`
	ChannelPointsVotes int			`json:"channel_points_votes"`
	BitsVotes int			`json:"bits_votes"`
}

type ApiPoll struct{
	Id string			`json:"id"`
	Broadcaster_id string			`json:"broadcaster_id"`
	Broadcaster_name string			`json:"broadcaster_name"`
	Title string			`json:"title"`
	Choices []ApiPollChoice			`json:"choices"`
	Bits_voting_enabled bool			`json:"bits_voting_enabled"`
	Bits_per_vote int			`json:"bits_per_vote"`
	Channel_points_voting_enabled bool			`json:"channel_points_voting_enabled"`
	Channel_points_per_vote int			`json:"channel_points_per_vote"`
	// ACTIVE | COMPLETED | TERMINATED | ARCHIVED | MODERATED | INVALID
	Status string			`json:"status"`
	// seconds			`json:"second"`
	Duration int			`json:"duration"`
	Started_at string			`json:"started_at"`
	Ended_at string			`json:"ended_at"`
}

type ApiGetPollRes struct{
	Data []ApiPoll 		`json:"data"`
	Pagination ApiPagination
}

// broadcaster_id(req), id, first, after
func ApiGetPoll(
	access_token string,
	params map[string][]string,
	) (*ApiResponse[ApiGetPollRes], error) {
	return ApiGet[ApiGetPollRes](
		SHARED_CHAT_SESSION_ENDPOINT,
		apiPollHeaders(access_token),
		params,
	)
}
