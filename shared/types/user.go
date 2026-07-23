package types

import "time"

type AppUser struct {
	Id                string   `json:"id"`
	Login             string   `json:"login"`
	Display_name      string   `json:"display_name"`
	User_type         string   `json:"type"`
	Broadcaster_type  string   `json:"broadcaster_type"`
	Description       string   `json:"description"`
	Profile_image_url string   `json:"profile_image_url"`
	Offline_image_url string   `json:"offline_image_url"`
	View_count        int      `json:"view_count"`
	Email             string   `json:"email"`
	Created_at        time.Time `json:"created_at"`
	Access_token      string   `json:"access_token"`
}
