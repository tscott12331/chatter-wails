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

type StatusError struct {
	Res *http.Response
}

func (se *StatusError) Error() string {
	return fmt.Sprintf("%s, %v", se.Res.Status, se.Res)
}

type APISessionReqTimeout struct{}
func (esrt *APISessionReqTimeout) Error() string {
	return "The request for a session ID timed out"
}


type APIResponse struct{
	Status int
	Body any
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


func ApiFetch(
	method string,
    endpoint url.URL,
    headers *http.Header,
	body any,
    params map[string][]string,
	) (*APIResponse, error) {
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
		

		req, err := http.NewRequest(method, endpoint.String(), req_body)
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

		var res_body_obj any
		err = json.Unmarshal(res_body, &res_body_obj)
		if err != nil {
			log.Printf("[ApiFetch]: An error occurred parsing the response body, aborting\n\n")
			return nil, err
		}

		
		return &APIResponse{
			Status: res.StatusCode,
			Body: res_body_obj,
		}, nil
}


func ApiDelete(
    endpoint url.URL,
    headers *http.Header,
    params map[string][]string,
    ) (*APIResponse, error) {
		return ApiFetch("DELETE", endpoint, headers, nil, params)
}

func ApiGet(
    endpoint url.URL,
    headers *http.Header,
    params map[string][]string,
    ) (*APIResponse, error) {
		return ApiFetch("GET", endpoint, headers, nil, params)
}

func ApiPost(
    endpoint url.URL,
    headers *http.Header,
	body any,
    params map[string][]string,
    ) (*APIResponse, error) {
		return ApiFetch("POST", endpoint, headers, body, params)
}



func ApiGetValidate(
        access_token string,
    ) (*APIResponse, error) {
    return ApiGet(
		VALIDATE_ENDPOINT,
		apiValidateHeaders(access_token),
		map[string][]string{},
	)
}


func ApiGetSubscriptions(
    access_token string,
    params map[string][]string,
	) (*APIResponse, error) {
    return ApiGet(
			SUBSCRIPTIONS_ENDPOINT,
			apiSubscriptionHeaders(access_token),
			params,
    )
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
    ) (*APIResponse, error) {
    return ApiPost(
        SUBSCRIPTIONS_ENDPOINT,
        apiSubscriptionHeaders(access_token),
        body,
        params,
    )
}

func ApiDeleteSubscriptions(
    access_token string,
    params map[string][]string,
    ) (*APIResponse, error) {
    return ApiDelete(
        SUBSCRIPTIONS_ENDPOINT,
        apiSubscriptionHeaders(access_token),
        params,
    )
}


func ApiGetUsers(
    access_token string,
    params map[string][]string,
    ) (*APIResponse, error) {
    return ApiGet(
        USERS_ENDPOINT,
        apiUsersHeaders(access_token),
        params,
    )
}

func ApiPostUsers(
    access_token string,
    body any,
    params map[string][]string,
    ) (*APIResponse, error) {
    return ApiPost(
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
    ) (*APIResponse, error) {
    return ApiPost(
        MESSAGES_ENDPOINT,
        apiMessagesHeaders(access_token),
        body,
        params,
    )
}


func ApiGetChannelBadges(
    access_token string,
    params map[string][]string,
    ) (*APIResponse, error) {
    return ApiGet(
        BADGES_ENDPOINT,
        apiBadgesHeaders(access_token),
        params,
    )
}

func ApiGetGlobalBadges(
    access_token string,
    ) (*APIResponse, error) {
    return ApiGet(
        GLOBAL_BADGES_ENDPOINT,
        apiBadgesHeaders(access_token),
		map[string][]string{},
    )
}


func ApiGetUserEmotes(
    access_token string,
    params map[string][]string,
    ) (*APIResponse, error) {
    return ApiGet(
		USER_EMOTES_ENDPOINT,
		apiEmotesHeaders(access_token),
		params,
	)
}


func ApiGetGlobalEmotes(
    access_token string,
    params map[string][]string,
    ) (*APIResponse, error) {
    return ApiGet(
		GLOBAL_EMOTES_ENDPOINT,
		apiEmotesHeaders(access_token),
		params,
	)
}
