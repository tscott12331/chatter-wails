package shared

import (
	"chatter-wails/shared/types"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func EmitNewSet(app *application.App, set *types.AppEmoteSet, broadcasterId string) {
	app.Event.Emit("chatter:emote:new-set", types.NewEmoteSetEvent{
		BroadcasterId: broadcasterId,
		AppEmoteSet:   *set,
	})
}
