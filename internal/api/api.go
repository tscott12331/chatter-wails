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
	"github.com/wailsapp/wails/v3/pkg/application"
)

type StatusError[T any] struct {
	Res *ApiResponse[T]
}
func (se *StatusError[any]) Error() string {
	return fmt.Sprintf("%v, %v", se.Res.Status, se.Res)
}


type ApiResponse[T any] struct{
	Status int
	Body T
}

// type Config struct{
// 	ClientId string			`env:"VITE_CLIENT_ID,required"`
// }

func ApiHeaders(access_token string) *http.Header {
	// defer dialog until application is initialized
	if len(ClientId) == 0 {
		application.Get().Dialog.Error().
			SetTitle("App is unauthorized").
			SetMessage("Please open an issue at https://github.com/tscott12331/gel/issues").
			Show()
	}

	return &http.Header{
		"Authorization": {"Bearer " + access_token},
		"Client-Id": {ClientId},
	}
}

var ClientId string

func init() {
	if len(ClientId) == 0 {
		// fallback for development
		godotenv.Load()
		var found bool
		ClientId, found = os.LookupEnv("VITE_CLIENT_ID")
		if !found {
			log.Printf("ERROR: could not find VITE_CLIENT_ID")
		}
	}
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
