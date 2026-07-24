package shared

import (
	"chatter-wails/internal/util"
	"chatter-wails/shared/types"
)

var user util.RWValue[*types.AppUser] = *util.NewRWValue[*types.AppUser]()

func GetUser() *types.AppUser {
	return user.Get()
}

func SetUser(newUser *types.AppUser) {
	user.Set(newUser)
}

func UpdateUser(update func(**types.AppUser)) {
	user.Update(update)
}
