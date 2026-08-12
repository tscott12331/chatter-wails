package types

type NewEmoteSetEvent struct{
	BroadcasterId string
	ChannelSpecific bool
	AppEmoteSet
}

type ChatOpenData struct{
	Channel string		`json:"channel"`
	AccessToken string 	`json:"accessToken"`
	Open bool			`json:"open"`
}
