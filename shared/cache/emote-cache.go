package cache

import (
	"chatter-wails/internal/util"
	"chatter-wails/shared/types"
	"maps"
)

// providers

const STV_KEY = "7tv"
const FFZ_KEY = "ffz"
const BTTV_KEY = "bttv"
const TWITCH_KEY = "twitch"

var EMOTE_PROVIDERS = []string{ STV_KEY, FFZ_KEY, BTTV_KEY, TWITCH_KEY }

// sections
const GLOBAL_EMOTE_SECTION = "global"
const USER_EMOTE_SECTION = "user"

const CHANNEL_EMOTE_SECTION = "channel"

var NON_CHANNEL_SPECIFIC_EMOTE_SECTIONS = []string{ GLOBAL_EMOTE_SECTION, USER_EMOTE_SECTION }

// provider -> section_id -> set
var emoteCache map[string]*util.RWMap[string, *types.AppEmoteSet] = map[string]*util.RWMap[string, *types.AppEmoteSet]{
	STV_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
	FFZ_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
	BTTV_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
	TWITCH_KEY: util.NewRWMap[string, *types.AppEmoteSet](),
}

func GetEmoteSet(provider, sectionId string) (*types.AppEmoteSet, bool) {
	if _, exists := emoteCache[provider]; !exists {
		panic("Trying to access emote cache of unknown provider")
	}

	mp := emoteCache[provider]

	val, exists := mp.Get(sectionId)
	return val, exists
}

// TODO: consider prepending "channel" to broadcasterId for clarity on section implication
func GetChannelEmoteSet(provider, broadcasterId string) (*types.AppEmoteSet, bool) {
	return GetEmoteSet(provider, broadcasterId)
}

func SetEmoteSet(provider, sectionId string, set *types.AppEmoteSet) {
	if _, exists := emoteCache[provider]; !exists {
		panic("Trying to set emote cache of unknown provider")
	}

	emoteCache[provider].Set(sectionId, set)
}

func SetChannelEmoteSet(provider, broadcasterId string, set *types.AppEmoteSet) {
	SetEmoteSet(provider, broadcasterId, set)
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

	for _, provider := range EMOTE_PROVIDERS {
		mergeChannelSetIntoMap(provider, broadcasterId, emoteMap)
		for _, sectionId := range NON_CHANNEL_SPECIFIC_EMOTE_SECTIONS {
			mergeSetIntoMap(provider, sectionId, emoteMap)
		}
	}


	return emoteMap
}

func mergeSetIntoMap(provider, sectionId string, emoteMap types.AppEmoteMap) {
	if set, exists := GetEmoteSet(provider, sectionId); exists {
		maps.Copy(emoteMap, set.Emotes)
	}
}

func mergeChannelSetIntoMap(provider, broadcasterId string, emoteMap types.AppEmoteMap) {
	mergeSetIntoMap(provider, broadcasterId, emoteMap)
}
