package types

type NewEmoteSetEvent struct{
	BroadcasterId string
	ChannelSpecific bool
	AppEmoteSet
}
