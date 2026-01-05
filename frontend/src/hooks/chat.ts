import { FailedApiRequest } from "@/api/api-response";
import { combineChannelGlobalSets, esBadgesToMessageBadges, getChannelBadges, IBadgeSet } from "@/api/badges";
import { ESCon, ESSubscription, IESNotification, TESMessage } from "@/api/eventsub";
import { sendMessage } from "@/api/messages";
import { getUser } from "@/api/user-info";
import { TUser } from "@/App";
import { TChatMessage } from "@/components/chat/chat-message";
import { IPopupUser } from "@/components/chat/user-popup";

import { CreateSubscription, DeleteSubscription } from "@wailsjs/go/services/EventSubService"

import { useEffect, useState } from "react";

export default function useChat({ channel, user, globalBadgeSets, maxMessages = 200 }: {
    channel: string|undefined,
    user: TUser
    globalBadgeSets: IBadgeSet[],
    maxMessages?: number,
}) {
    const [broadcasterId, setBroadcasterId] = useState<string|null>();
    const [subscription, setSubscription] = useState<string|null>();

    const [chatMessages, setChatMessages] = useState<TChatMessage[]>([]);

    const [channelBadgeSets, setChannelBadgeSets] = useState<IBadgeSet[]>([]);

    const badgeSets: IBadgeSet[] = combineChannelGlobalSets(channelBadgeSets, globalBadgeSets);
    

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

    const handleNotificationMessage = (message: TESMessage) => {
        const data = message as IESNotification;
        switch(data.metadata.subscription_type) {
            case 'channel.chat.message':
                const message: TChatMessage = {
                id: data.payload.event.message_id,
                username: data.payload.event.chatter_user_name,
                text: data.payload.event.message.text,
                fragments: data.payload.event.message.fragments,
                color: data.payload.event.color,
                badges: esBadgesToMessageBadges(data.payload.event.badges, badgeSets),
                reply: data.payload.event.reply,
            }

            appendChatMessage(message);
            break;
        }
    }

    const getBroadcasterId = async (channelName: string, access_token: string) => {
        const res = await getUser(access_token, channelName);
        if(res.success) {
            return res.data.user.id;
        } else {
            return null;
        }
    }

    const sendChatMessage = async (message: string, replyId?: string) => {
        const trimmed = message.trim();
        if(trimmed.length === 0 || !broadcasterId) return FailedApiRequest();

        return await sendMessage(user, user.access_token, broadcasterId, trimmed, replyId);
    }
    
    useEffect(() => {
        if(!channel) return;

        getBroadcasterId(channel, user.access_token).then(id => setBroadcasterId(id));
    }, [channel, user]);

    useEffect(() => {
        if(!broadcasterId) return;

        // maybe change escon to connect to chatroom based on channel name
        // ESCon.connectToChatroom(user, broadcasterId).then(s => setSubscription(s.success ? s.data.subscription : null));
        CreateSubscription(user.access_token, {
                'broadcaster_user_id': broadcasterId,
                'user_id': user.id,
            },
            "channel.chat.message")
            .then(s => setSubscription(s))
            .catch(e => {
                console.error(e);
                setSubscription(null);
            });

        getChannelBadges(user.access_token, broadcasterId).then(b => b.success && setChannelBadgeSets(b.data.sets));
    }, [broadcasterId])

    useEffect(() => {
        if(!subscription) return;

        // subscription.addEventListener('message', handleNotificationMessage);
        setChatMessages([]);

        return () => {
            DeleteSubscription(user.access_token, subscription);
            // subscription.removeEventListener('message', handleNotificationMessage);
            // ESCon.deleteSubscription(subscription, user.access_token);
        }
    }, [subscription]);

    return { chatMessages, sendChatMessage }
}
