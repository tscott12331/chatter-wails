package eventsub

import (
	"chatter-wails/internal/api"
	"chatter-wails/services/emote"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
)

type ESSubscriptionCondition map[string]string

type ESBadge struct {
	Set_id string			`json:"set_id"`
	Id string				`json:"id"`
	Info string				`json:"info"`
}


type AppChatMessageFragment struct{
	Fragment_type string							`json:"type"`
	Text string										`json:"text"`
    Cheermote *struct{
		Prefix string								`json:"prefix"`
		Bits int									`json:"bits"`
		Tier int									`json:"tier"`
	}												`json:"cheermote,omitempty"`
	Emote *emote.AppEmote									`json:"emote"`
    Mention *struct{
		User_id string								`json:"user_id"`
		User_name string							`json:"user_name"`
		User_login string							`json:"user_login"`
	}												`json:"mention,omitempty"`
}

type ESChatMessageFragment struct{
	Fragment_type string							`json:"type"`
	Text string										`json:"text"`
    Cheermote *struct{
		Prefix string								`json:"prefix"`
		Bits int									`json:"bits"`
		Tier int									`json:"tier"`
	}												`json:"cheermote,omitempty"`
    Emote *struct{
		Id string									`json:"id"`
		Emote_set_id string							`json:"emote_set_id"`
		Owner_id string								`json:"owner_id"`
		Format []string								`json:"format"`
	}												`json:"emote,omitempty"`
    Mention *struct{
		User_id string								`json:"user_id"`
		User_name string							`json:"user_name"`
		User_login string							`json:"user_login"`
	}												`json:"mention,omitempty"`
}


type ESMessageBadge struct{
	SrcSet string									`json:"srcSet"`
	Info string										`json:"info"`
	Title string									`json:"title"`
}

type ESMessageReply struct{
	Parent_message_body string						`json:"parent_message_body"`
	Parent_message_id string						`json:"parent_message_id"`
	Parent_user_id string							`json:"parent_user_id"`
	Parent_user_login string						`json:"parent_user_login"`
	Parent_user_name string							`json:"parent_user_name"`
	Thread_message_id string						`json:"thread_message_id"`
	Thread_user_id string							`json:"thread_user_id"`
	Thread_user_login string						`json:"thread_user_login"`
	Thread_user_name string							`json:"thread_user_name"`
}

type ESChatMessage struct{
	Id string							`json:"id"`
	Username string						`json:"username"`
	Channel string						`json:"channel"`
	Text string							`json:"text"`
	Fragments []*AppChatMessageFragment	`json:"fragments"`
	Color string						`json:"color"`
	Badges []ESMessageBadge				`json:"badges"`
	Reply *ESMessageReply				`json:"reply,omitempty"`
}



/* EVENTSUB EVENT TYPES */

type ESChatMessageEventMessage struct{
	Text string							`json:"text"`
	Fragments []*ESChatMessageFragment	`json:"fragments"`
}

type ESChatMessageEvent = struct{
	Broadcaster_user_id string				`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Chatter_user_id string					`json:"chatter_user_id"`
	Chatter_user_login string				`json:"chatter_user_login"`
	Chatter_user_name string				`json:"chatter_user_name"`
	Message_id string						`json:"message_id"`
	Message ESChatMessageEventMessage 		`json:"message"`
	Color string							`json:"color"`
	Badges []ESBadge						`json:"badges"`
	Message_type string						`json:"message_type"`
	Cheer *struct {
		Bits int							`json:"bits"`
	}										`json:"cheer,omitempty"`
	Reply *ESMessageReply							`json:"reply,omitempty"`
	Channel_points_custom_reward_id *string			`json:"channel_points_custom_reward_id,omitempty"`
	Source_broadcaster_user_id *string				`json:"source_broadcaster_user_id,omitempty"`
	Source_broadcaster_user_name *string			`json:"source_broadcaster_user_name,omitempty"`
	Source_broadcaster_user_login *string			`json:"source_broadcaster_user_login,omitempty"`
	Source_message_id *string						`json:"source_message_id,omitempty"`
	Source_badges *[]ESBadge						`json:"source_badges,omitempty"`
	Is_source_only *bool							`json:"is_source_only,omitempty"`
}

type ESSharedChatParticipant struct{
	Broadcaster_user_id string			`json:"broadcaster_user_id"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
}

type ESSharedChatBeginEvent struct{
	Session_id string			`json:"session_id"`
	Broadcaster_user_id string			`json:"broadcaster_user_id"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Host_broadcaster_user_id string			`json:"host_broadcaster_user_id"`
	Host_broadcaster_user_name string			`json:"host_broadcaster_user_name"`
	Host_broadcaster_user_login string			`json:"host_broadcaster_user_login"`
	Participants []ESSharedChatParticipant			`json:"participants"`
}

type ESSharedChatUpdateEvent ESSharedChatBeginEvent

type ESSharedChatEndEvent struct{
	Session_id string			`json:"session_id"`
	Broadcaster_user_id string			`json:"broadcaster_user_id"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Host_broadcaster_user_id string			`json:"host_broadcaster_user_id"`
	Host_broadcaster_user_name string			`json:"host_broadcaster_user_name"`
	Host_broadcaster_user_login string			`json:"host_broadcaster_user_login"`
}



type ESBanEvent struct{
	User_id string			`json:"user_id"`
	User_login string			`json:"user_login"`
	User_name string			`json:"user_name"`
	Broadcaster_user_id string			`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Moderator_user_id string			`json:"moderator_user_id"`
	Moderator_user_login string			`json:"moderator_user_login"`
	Moderator_user_name string			`json:"moderator_user_name"`
	Reason string			`json:"reason"`
	Banned_at string			`json:"banned_at"`
	Ends_at string			`json:"ends_at"`
	Is_permanent bool			`json:"is_permanent"`
}


type ESSpecialVoteData struct{
	Is_enabled bool			`json:"is_enabled"`
	Amount_per_vote int			`json:"amount_per_vote"`
}
type ESPollBeginEvent struct{
	Id string			`json:"id"`
	Broadcaster_user_id string			`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Title string			`json:"title"`
	Choices []api.ApiPollChoice			`json:"choices"`
	Bits_voting ESSpecialVoteData			`json:"bits_voting"`
	Channel_points_voting ESSpecialVoteData			`json:"channel_points_voting"`
	Started_at string			`json:"started_at"`
	Ends_at string			`json:"ends_at"`
}

/* WELCOME MESSAGE TYPES */

type ESWelcomeMetadata struct {
	Message_id string			`json:"message_id"`
	Message_type string			`json:"message_type"`
	Message_timestamp string 	`json:"message_timestamp"`
}

type ESWelcomePayloadSession struct {
	Id string							`json:"id"`
	Status string						`json:"status"`
	Keepalive_timeout_seconds int		`json:"keepalive_timeout_seconds"`
	Reconnect_url string 				`json:"reconnect_url"`
	Connected_at string 				`json:"connected_at"`
}

type ESWelcomePayload struct {
	Session ESWelcomePayloadSession 	`json:"session"`
}

type ESWelcome struct {
	Metadata ESWelcomeMetadata					`json:"metadata"`
	Payload ESWelcomePayload					`json:"payload"`
}



/* NOTIFICATION TYPES */

type ESNotificationMetadata = ESMessageMetadata

type ESNotificationPayloadSubscription = ESMessagePayloadSubscription

type ESNotificationPayload struct {
	Subscription ESNotificationPayloadSubscription 	`json:"subscription"`
	Event *json.RawMessage									`json:"event,omitempty"`
}

type ESNotification struct {
	Metadata ESNotificationMetadata		`json:"metadata"`
	Payload ESNotificationPayload		`json:"payload"`
}



type ESChatMessageNotificationPayload struct {
	Subscription ESNotificationPayloadSubscription 	`json:"subscription"`
	Event *ESChatMessageEvent						`json:"event,omitempty"`
}

type ESChatMessageNotification struct {
	Metadata ESNotificationMetadata				`json:"metadata"`
	Payload ESChatMessageNotificationPayload	`json:"payload"`
}



/* EVENTSUB MESSAGE TYPES */

type ESMessageMetadata struct {
	Message_id string			`json:"message_id"`
	Message_type string			`json:"message_type"`
	Message_timestamp string	`json:"message_timestamp"`
	Subscription_type string	`json:"subscription_type"`
	Subscription_version string `json:"subscription_version"`
}

type ESMessagePayloadSubscriptionTransport struct {
	Method string			`json:"method"`
	Session_id string		`json:"session_id"`
}

type ESMessagePayloadSubscription struct {
	Id string												`json:"id"`
	Status string											`json:"status"`
	Sub_type string											`json:"type"`
	Version string											`json:"version"`
	Cost int												`json:"cost"`
	Condition map[string]string								`json:"condition"`
	Transport ESMessagePayloadSubscriptionTransport			`json:"transport"`
	Created_at string										`json:"created_at"`
}

type ESMessagePayload struct {
	Subscription ESMessagePayloadSubscription 	`json:"subscription"`
	Event *json.RawMessage				`json:"event,omitempty"`
	Session ESWelcomePayloadSession				`json:"session"`
}

// type ESMessage map[string]map[string]any
type ESMessage struct {
	Metadata ESMessageMetadata 			`json:"metadata"`
	Payload ESMessagePayload			`json:"payload"`
}


func esMessageToESNotification(message *ESMessage) *ESNotification {
	return &ESNotification{
		Metadata: message.Metadata,
		Payload: ESNotificationPayload{
			Subscription: message.Payload.Subscription,
			Event: message.Payload.Event,
		},
	}
}




func esBadgesToMessageBadges(esBadges []ESBadge, badgeSets []api.ApiBadgeSet) []ESMessageBadge {
    var messageBadges = []ESMessageBadge{}

	for _, badge := range esBadges {
		setIndex := slices.IndexFunc(badgeSets, func(bs api.ApiBadgeSet) bool {
			return bs.Set_id == badge.Set_id
		})
        if(setIndex == -1) {
			continue
		}

		var set = badgeSets[setIndex]


		versionIndex := slices.IndexFunc(set.Versions, func(v api.ApiBadgeSetVersions) bool {
			return v.Id == badge.Id
		})
        if(versionIndex == -1) {
			continue
		}

        var version = set.Versions[versionIndex]

        var urls = []string{version.Image_url_1x, version.Image_url_2x, version.Image_url_4x};

		
        var srcSet strings.Builder
		srcSet.Grow(252)

		for i, url := range urls {
			var scale int
			if i == 0 {
				scale = 1
			} else {
				scale = 2 << i
			}

			fmt.Fprintf(&srcSet, "%s %dx", url, scale)
			if i < 2 {
				fmt.Fprint(&srcSet, ", ")
			}
		}

		messageBadges = append(messageBadges, ESMessageBadge{
			SrcSet: srcSet.String(),
			Info: badge.Info,
			Title: version.Title,
		})
    }

    return messageBadges;
}

func extractSharedChatBadge(chatMessageEvent *ESChatMessageEvent, chatSubscriptionData *ESChatSubscriptionData) []ESMessageBadge {
	res := []ESMessageBadge{}

	broadcaster_login := chatMessageEvent.Broadcaster_user_login

	if chatMessageEvent.Source_broadcaster_user_login != nil {
		broadcaster_login = *chatMessageEvent.Source_broadcaster_user_login
	}

	participant, exists := chatSubscriptionData.SharedChatParticipants[broadcaster_login]
	if !exists {
		return res
	}

	res = append(res, ESMessageBadge{
		SrcSet: fmt.Sprintf("%s 1x", participant.ProfileImageURL),
		Info: "Shared chat",
		Title: participant.Name,
	})

	return res
}

func esNotificationToEsChatMessage(notification *ESNotification, chatSubscriptionData *ESChatSubscriptionData) *ESChatMessage {
	channelBadges, _ := chatSubscriptionData.ChannelBadgeSets.Read()
	seventvEmotes := chatSubscriptionData.SevenTV.SevenTVEmotes

	var chatMessageEvent ESChatMessageEvent
	if notification.Payload.Event == nil { 
		log.Printf("WARNING: recieved empty chat message event payload")
		return nil
	}

	json.Unmarshal(*notification.Payload.Event, &chatMessageEvent)

	var fragments = []*AppChatMessageFragment{}
	for _, fragment := range chatMessageEvent.Message.Fragments {
		fragments = append(fragments, esChatMessageFragmentToAppMessageFragment(fragment, seventvEmotes)...)
	}

	sharedChatBadge := extractSharedChatBadge(&chatMessageEvent, chatSubscriptionData)

	return &ESChatMessage{
		Id: chatMessageEvent.Message_id,
		Username: chatMessageEvent.Chatter_user_name,
		Channel: chatMessageEvent.Broadcaster_user_login,
		Text: chatMessageEvent.Message.Text,
		Fragments: fragments,
		Color: chatMessageEvent.Color,
		Badges: append(sharedChatBadge, esBadgesToMessageBadges(chatMessageEvent.Badges, channelBadges)...),
		Reply: chatMessageEvent.Reply,
	}
}

func esMessageToESWelcome(message *ESMessage) *ESWelcome {
	return &ESWelcome{
		Metadata: ESWelcomeMetadata{
			Message_id: message.Metadata.Message_id,
			Message_type: message.Metadata.Message_type,
			Message_timestamp: message.Metadata.Message_timestamp,
		},
		Payload: ESWelcomePayload{
			Session: message.Payload.Session,
		},
	}
}
