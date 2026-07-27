package message

import (
	"chatter-wails/internal/api"
	"chatter-wails/internal/api/nativeApi"
	"errors"
	"fmt"
	"log"
)

func SendMessage(userId string, access_token string, brdId string, message string, replyId *string) (*nativeApi.ApiPostMessagesData, error) {
	body := nativeApi.ApiPostMessagesBody{
		Sender_id:      userId,
		Broadcaster_id: brdId,
		Message: message,
		Reply_parent_message_id: replyId,
	}

	res, err := nativeApi.PostMessage(access_token, body, map[string][]string{})
	if err != nil {
		log.Printf("[SendMessage]: An error occurred trying to send a message, aborting\n%+v\n\n", err)
		return nil, err
	}
	if res.Status != 200 {
		log.Printf("[SendMessage]: Failed to send message, aborting\n%+v\n\n", res)
		return nil, &api.StatusError[nativeApi.ApiPostMessagesRes]{ Res: res }
	}
	if len(res.Body.Data) == 0 {
		log.Printf("[SendMessage]: API returned no message data, aborting\n%+v\n\n", res)
		return nil, errors.New("API returned no message data")
	}

	messageData := res.Body.Data[0]

	if messageData.Drop_reason != nil {
		log.Printf("[SendChatMessage]: Message was dropped\n%+v\n\n", messageData.Drop_reason)
		return nil, fmt.Errorf("Message was dropped %+v", messageData.Drop_reason)
	}

	return &messageData, nil
}
