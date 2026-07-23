package cache

import (
	"chatter-wails/internal/util"
	"chatter-wails/shared/types"
	"maps"
)

const STV_KEY = "7tv"
const FFZ_KEY = "ffz"
const BTTV_KEY = "bttv"
const NATIVE_KEY = "native"

var emoteCache map[string]*util.RWMap[string, *types.AppEmoteSet] = map[string]*util.RWMap[string, *types.AppEmoteSet]{
	STV_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
	FFZ_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
	BTTV_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
	NATIVE_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
}

func GetEmoteSet(provider, broadcasterId string) (*types.AppEmoteSet, bool) {
	if _, exists := emoteCache[provider]; !exists {
		panic("Trying to access emote cache of unknown provider")
	}

	mp := emoteCache[provider]

	val, exists := mp.Get(broadcasterId)
	return val, exists
}

func SetEmoteSet(provider, broadcasterId string, set *types.AppEmoteSet) {
	if _, exists := emoteCache[provider]; !exists {
		panic("Trying to set emote cache of unknown provider")
	}

	emoteCache[provider].Set(broadcasterId, set)
}


func RemoveBroadcasterEmoteSets(broadcasterId string) {
	for provider := range emoteCache {
		RemoveEmoteSet(provider, broadcasterId)
	}
}

func RemoveEmoteSet(provider, broadcasterId string) {
	if _, exists := emoteCache[provider]; !exists {
		panic("Trying to delete emote cache of unknown provider")
	}

	emoteCache[provider].Delete(broadcasterId)
}


func GetBroadcasterEmoteMap(broadcasterId string) types.AppEmoteMap {
	emoteMap := types.AppEmoteMap{}

	// unsure if native should even be here
	// precedence: 7tv > bttv > ffz
	mergeSetIntoMap(STV_KEY, broadcasterId, emoteMap)
	mergeSetIntoMap(BTTV_KEY, broadcasterId, emoteMap)
	mergeSetIntoMap(FFZ_KEY, broadcasterId, emoteMap)

	return emoteMap
}

func mergeSetIntoMap(provider, broadcasterId string, emoteMap types.AppEmoteMap) {
	if set, exists := GetEmoteSet(provider, broadcasterId); exists {
		maps.Copy(emoteMap, set.Emotes)
	}
}
