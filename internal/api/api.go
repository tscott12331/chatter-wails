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


func ApiGetUsers(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[any], error) {
    return ApiGet[any](
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


type PostMessagesBody struct{
	Broadcaster_id string				`json:"broadcaster_id"`
	Sender_id string					`json:"sender_id"`
	Message string 						`json:"message"`
	Reply_parent_message_id *string		`json:"reply_parent_message_id"`
}

func ApiPostMessages(
    access_token string,
    body PostMessagesBody,
    params map[string][]string,
    ) (*ApiResponse[any], error) {
    return ApiPost[any](
        MESSAGES_ENDPOINT,
        apiMessagesHeaders(access_token),
        body,
        params,
    )
}


func ApiGetChannelBadges(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[any], error) {
    return ApiGet[any](
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


func ApiGetUserEmotes(
    access_token string,
    params map[string][]string,
    ) (*ApiResponse[any], error) {
    return ApiGet[any](
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

type ApiGlobalEmote struct{
	Id string					`json:"id"`
	Name string					`json:"name"`
    Images ApiGlobalEmoteImages	`json:"images"`
	Format []string				`json:"format"`
	Scale []string				`json:"scale"`
	Theme_mode []string			`json:"theme_mode"`
}

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
