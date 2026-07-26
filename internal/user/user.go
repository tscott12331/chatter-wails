package user

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/api/nativeApi"
	"errors"
	"log"
)



func GetUserByToken(accessToken string) (*nativeApi.ApiUser, error){
	return getUser(accessToken, nil, nil)
}

func GetUserByLogin(accessToken string, loginName string) (*nativeApi.ApiUser, error) {
	return getUser(accessToken, &loginName, nil)
}

func GetUserById(accessToken string, id string) (*nativeApi.ApiUser, error) {
	return getUser(accessToken, nil, &id)
}

func getUser(accessToken string, loginName *string, id *string) (*nativeApi.ApiUser, error) {
	params := map[string][]string{}
	if loginName != nil {
		params["login"] = []string{*loginName}
	}
	if id != nil {
		params["id"] = []string{*id}
	}

	res, err := nativeApi.GetNativeUsers(accessToken, params)
	if err != nil {
		log.Printf("[getUser]: An error occurred fetching user, aborting\n%+v\n\n", res)
		return nil, err
	}
	if res.Status != 200 {
		log.Printf("[getUser]: Failed to get user, aborting\n%+v\n\n", res)
		return nil, &api.StatusError[nativeApi.ApiGetUsersRes]{ Res: res }
	}
	if len(res.Body.Data) == 0 {
		log.Printf("[getUser]: Request parameters returned no users\n\n")
		return nil, errors.New("No users found")
	}

	return &res.Body.Data[0], nil
}
