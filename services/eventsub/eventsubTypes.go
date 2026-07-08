package eventsub

import (
	"chatter-wails/internal/api"
	"chatter-wails/services/emote"
	"fmt"
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

type ESEvent struct {
	Broadcaster_user_id string				`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Chatter_user_id string					`json:"chatter_user_id"`
	Chatter_user_login string				`json:"chatter_user_login"`
	Chatter_user_name string				`json:"chatter_user_name"`
	Message_id string						`json:"message_id"`
	Message struct {
		Text string							`json:"text"`
		Fragments []ESChatMessageFragment	`json:"fragments"`
	}										`json:"message"`
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
	Source_braodcaster_user_login *string			`json:"source_broadcaster_user_login,omitempty"`
	Source_message_id *string						`json:"source_message_id,omitempty"`
	Source_badges *[]ESBadge						`json:"source_badges,omitempty"`
	Is_source_only *bool							`json:"is_source_only,omitempty"`

}

type ESChatMessageEventMessage struct{
	Text string							`json:"text"`
	Fragments []*AppChatMessageFragment	`json:"fragments"`
}
type ESChatMessageEvent = struct{
	Broadcaster_user_id string				`json:"broadcaster_user_id"`
	Broadcaster_user_login string			`json:"broadcaster_user_login"`
	Broadcaster_user_name string			`json:"broadcaster_user_name"`
	Chatter_user_id string					`json:"chatter_user_id"`
	Chatter_user_login string				`json:"chatter_user_login"`
	Chatter_user_name string				`json:"chatter_user_name"`
	Message_id string						`json:"message_id"`
	Message ESChatMessageEventMessage		`json:"message"`
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
	Source_braodcaster_user_login *string			`json:"source_broadcaster_user_login,omitempty"`
	Source_message_id *string						`json:"source_message_id,omitempty"`
	Source_badges *[]ESBadge						`json:"source_badges,omitempty"`
	Is_source_only *bool							`json:"is_source_only,omitempty"`
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
	Event *ESEvent									`json:"event,omitempty"`
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
	Event *ESEvent					`json:"event,omitempty"`
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

func esNotificationToEsChatMessage(notification *ESNotification, chatSubscriptionData *ESChatSubscriptionData) *ESChatMessage {
	channelBadges, _ := chatSubscriptionData.ChannelBadgeSets.Read()
	seventvEmotes := chatSubscriptionData.SevenTV.SevenTVEmotes
	var fragments = []*AppChatMessageFragment{}
	for _, fragment := range notification.Payload.Event.Message.Fragments {
		fragments = append(fragments, esChatMessageFragmentToAppMessageFragment(&fragment, seventvEmotes)...)
	}

	return &ESChatMessage{
		Id: notification.Payload.Event.Message_id,
		Username: notification.Payload.Event.Chatter_user_name,
		Channel: notification.Payload.Event.Broadcaster_user_login,
		Text: notification.Payload.Event.Message.Text,
		Fragments: fragments,
		Color: notification.Payload.Event.Color,
		Badges: esBadgesToMessageBadges(notification.Payload.Event.Badges, channelBadges),
		Reply: notification.Payload.Event.Reply,
	}
}

func esNotificationToEsChatMessageNotification(notification *ESNotification) *ESChatMessageNotification {
	var fragments = []*AppChatMessageFragment{}
	for _, fragment := range notification.Payload.Event.Message.Fragments {
		fragments = append(fragments, esChatMessageFragmentToAppMessageFragment(&fragment, map[string]*emote.AppEmote{})...)
	}

	return &ESChatMessageNotification{
		Metadata: notification.Metadata,
		Payload: ESChatMessageNotificationPayload{
			Subscription: notification.Payload.Subscription,
			Event: &ESChatMessageEvent{
				Broadcaster_user_id: notification.Payload.Event.Broadcaster_user_id,
				Broadcaster_user_login: notification.Payload.Event.Broadcaster_user_login,
				Broadcaster_user_name: notification.Payload.Event.Broadcaster_user_name,
				Chatter_user_id: notification.Payload.Event.Chatter_user_id,
				Chatter_user_login: notification.Payload.Event.Chatter_user_login,
				Chatter_user_name: notification.Payload.Event.Chatter_user_name,
				Message_id: notification.Payload.Event.Message_id,
				Message: ESChatMessageEventMessage{
					Text: notification.Payload.Event.Message.Text,
					Fragments: fragments,
				},
				Color: notification.Payload.Event.Color,
				Badges: notification.Payload.Event.Badges,
				Message_type: notification.Payload.Event.Message_type,
				Cheer: notification.Payload.Event.Cheer,
				Reply: notification.Payload.Event.Reply,
				Channel_points_custom_reward_id: notification.Payload.Event.Channel_points_custom_reward_id,
				Source_broadcaster_user_id: notification.Payload.Event.Source_broadcaster_user_id,
				Source_broadcaster_user_name: notification.Payload.Event.Source_broadcaster_user_name,
				Source_braodcaster_user_login: notification.Payload.Event.Source_braodcaster_user_login,
				Source_message_id: notification.Payload.Event.Source_message_id,
				Source_badges: notification.Payload.Event.Source_badges,
				Is_source_only: notification.Payload.Event.Is_source_only,
				
			},
		},
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
