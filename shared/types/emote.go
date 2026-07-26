package types

type AppEmote struct{
	Id string					`json:"id"`
	Name string					`json:"name"`
	LightSrcSet string		`json:"lightSrcSet"`
	DarkSrcSet string		`json:"darkSrcSet"`
	Provider string				`json:"provider"`
	Section string				`json:"section"`
	ZeroWidth bool			`json:"zeroWidth"`
	EmoteStack []*AppEmote  `json:"emoteStack"`
}


type AppEmoteMap map[string]*AppEmote

type AppEmoteSet struct{
	Id string
	Provider string
	Section string
	Emotes AppEmoteMap
}
