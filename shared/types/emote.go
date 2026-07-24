package types

type AppEmote struct{
	Id string					`json:"id"`
	Name string					`json:"name"`
	LightSrcSet string		`json:"lightSrcSet"`
	DarkSrcSet string		`json:"darkSrcSet"`
	// 'global' | 'user' | 'channel' | 'seventv'
	Type string				`json:"type"`
	ZeroWidth bool			`json:"zeroWidth"`
	EmoteStack []*AppEmote  `json:"emoteStack"`
}


type AppEmoteMap map[string]*AppEmote

type AppEmoteSet struct{
	Id string
	Provider string
	Emotes AppEmoteMap
}
