package auth

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/user"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppUser struct {
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
    Access_token string			`json:"access_token"`
}

type AuthService struct {
	app *application.App

	Ctx context.Context
	User *AppUser

	cancelValidation chan struct{}
	validationRunning bool
	validationMutex sync.Mutex
}

func NewAuthService(app *application.App) *AuthService {
	return &AuthService{
		app: app,
		cancelValidation: make(chan struct{}),
	}
}


func apiUserToAppUser(user *api.ApiUser, accessToken string) *AppUser {
	return &AppUser{
		Id: user.Id,
		Login: user.Login,
		Display_name: user.Display_name,
		User_type: user.User_type,
		Broadcaster_type: user.Broadcaster_type,
		Description: user.Description,
		Profile_image_url: user.Profile_image_url,
		Offline_image_url: user.Offline_image_url,
		View_count: user.View_count,
		Email: user.Email,
		Created_at: user.Created_at,
		Access_token: accessToken,
	}
}

func (as *AuthService) validate() error {
	if as.User == nil {
		log.Printf("[validate]: Cannot validate a nil user, aborting\n\n")
		return errors.New("Cannot validate a nil user")
	}

	res, err := api.ApiGetValidate(as.User.Access_token)
	if err != nil {
		log.Printf("[validate]: An error occurred validating token, signing out and aborting\n%+v\n\n", err)
		return err
	}
	if res.Status != 200 {
		log.Printf("[validate]: Failed to validate token, signing out and aborting\n\n")
		return &api.StatusError[any]{ Res: res }
	}

	return nil
}

func (as *AuthService) validateLoop() {
	log.Printf("[validateLoop]: Starting validate loop\n\n")

	err := as.validate()
	if err != nil {
		log.Printf("[validateLoop]: Failed to validate user access token, aborting\n\n")
		goto errorOccurred
	}

	for {
		select {
		case <-time.After(time.Hour * 1):
			log.Printf("[validateLoop]: Validating access token\n\n")
			err := as.validate()
			if err != nil {
				log.Printf("[validateLoop]: Failed to validate user access token, aborting\n\n")
				goto errorOccurred
			}
		case <-as.cancelValidation:
			log.Printf("[validateLoop]: Validation loop has been cancelled, aborting\n\n")
			return
		}
	}

errorOccurred: // if set validation running to false when an error occurs
	as.validationMutex.Lock()
	as.validationRunning = false
	as.validationMutex.Unlock()
}

func (as *AuthService) Login(accessToken string) (*AppUser, error) {
	apiUser, err := user.GetUserByToken(accessToken)
	if err != nil {
		log.Printf("[Login]: An error occurred while logging in, aborting\n%+v\n\n", err)
		return nil, err
	}

	as.validationMutex.Lock() // ensure only one validation loop runs
	if as.validationRunning {
		// a user was previously signed in
		as.cancelValidation <- struct{}{} // cancel validation loop for previous user's access token
	} else {
		as.validationRunning = true
	}

	go as.validateLoop()

	as.validationMutex.Unlock()

	user := apiUserToAppUser(apiUser, accessToken)
	as.User = user
	as.app.Event.Emit("common:user-login", as.User)

	return user, nil
}

func (as *AuthService) Logout() {
	if as.User != nil {
		as.cancelValidation <- struct{}{}

		as.User = nil
	}
}
