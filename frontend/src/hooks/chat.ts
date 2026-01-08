import { FailedApiRequest } from "@/api/api-response";
import { esBadgesToMessageBadges, IBadgeSet } from "@/api/badges";
import { IESNotification, TESMessage } from "@/api/eventsub";
import { sendMessage } from "@/api/messages";
import { TUser } from "@/App";
import { TChatMessage } from "@/components/chat/chat-message";

import { ConnectToChatroom } from '@wailsjs/go/main/App'
import { api } from "@wailsjs/go/models";
import { DeleteSubscription } from "@wailsjs/go/services/EventSubService"
import { EventsOn, EventsOff } from "@wailsjs/runtime/runtime";

import { useEffect, useState } from "react";

interface IChatroomData {
    subId: string|null;
    broadcasterId: string|null;
    badgeSets: api.ApiBadgeSet[]
}

export default function useChat({ channel, user, maxMessages = 200 }: {
    channel: string|undefined,
    user: TUser
    maxMessages?: number,
}) {

    const [chatMessages, setChatMessages] = useState<TChatMessage[]>([]);

    const [chatroomData, setChatroomData] = useState<IChatroomData>({
        subId: null,
        broadcasterId: null,
        badgeSets: [],
    })
    

    const appendChatMessage = (message: TChatMessage) => {
        setChatMessages(curMessages => {
            const numExtraMessages = curMessages.length - maxMessages;
            if(numExtraMessages >= 0) {
                return [...curMessages.slice(numExtraMessages + 1), message];
            } else {
                return [...curMessages, message];
            }
        })
    }

    const handleChatMessage = (message: TESMessage) => {
        const data = message as IESNotification;

        const chatMessage: TChatMessage = {
            id: data.payload.event.message_id,
            username: data.payload.event.chatter_user_name,
            text: data.payload.event.message.text,
            fragments: data.payload.event.message.fragments,
            color: data.payload.event.color,
            badges: esBadgesToMessageBadges(data.payload.event.badges, chatroomData.badgeSets),
            reply: data.payload.event.reply,
        }
        appendChatMessage(chatMessage);
    }

    const sendChatMessage = async (message: string, replyId?: string) => {
        const trimmed = message.trim();
        if(trimmed.length === 0 || !chatroomData.broadcasterId) return FailedApiRequest();

        return await sendMessage(user, user.access_token, chatroomData.broadcasterId, trimmed, replyId);
    }
    
    useEffect(() => {
        if(!channel) return;

        ConnectToChatroom(channel)
            .then(d => setChatroomData(d))
            .catch(err => console.error(err))
    }, [channel, user]);

    useEffect(() => {
        const subId = chatroomData.subId;
        if(!subId) return;

        // subscription.addEventListener('message', handleNotificationMessage);
        setChatMessages([]);

        EventsOn(subId, (m) => handleChatMessage(m));
        
        return () => {
            EventsOff(subId);
            DeleteSubscription(user.access_token, subId);
        }
    }, [chatroomData]);

    return { chatMessages, sendChatMessage }
}
