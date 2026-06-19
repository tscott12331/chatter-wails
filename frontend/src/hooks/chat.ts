import { TChatroomEmotes } from "@/api/emote";
import { IAppEmote } from "@/api/native-emote";
import { TUser } from "@/App";
import { TChatMessage } from "@/components/chat/chat-message";

import { ConnectToChatroom, EnableSevenTV, SendChatMessage } from '@wailsjs/go/main/App'
import { api } from "@wailsjs/go/models";
import { DeleteSubscription } from "@wailsjs/go/services/EventSubService"
import { EventsOn, EventsOff } from "@wailsjs/runtime/runtime";

import { useEffect, useState } from "react";

interface IChatroomData {
    subId: string|null;
    broadcasterId: string|null;
    badgeSets: api.ApiBadgeSet[];
    channelEmotes: Record<string, IAppEmote>;
}

export default function useChat({ channel, user, emoteRecord, maxMessages = 200 }: {
    channel: string|undefined,
    user: TUser,
    emoteRecord: TChatroomEmotes,
    maxMessages?: number,
}) {

    const [chatMessages, setChatMessages] = useState<TChatMessage[]>([]);

    const [chatroomData, setChatroomData] = useState<IChatroomData>({
        subId: null,
        broadcasterId: null,
        badgeSets: [],
        channelEmotes: {},
    })

    const [emotes, setEmotes] = useState<TChatroomEmotes>(emoteRecord);
    

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

    const sendChatMessage = async (message: string, replyId?: string) => {
        if(!chatroomData.subId) return false;
        const trimmed = message.trim();
        if(trimmed.length === 0 || !chatroomData.broadcasterId) return false;

        try {
            await SendChatMessage(chatroomData.subId, trimmed, replyId);
            return true;
        } catch(err) {
            console.error(err);
            return false;
        }
    }
    
    useEffect(() => {
        if(!channel) return;

        ConnectToChatroom(channel)
            .then(d => setChatroomData(d))
            .catch(err => console.error(err));
    }, [channel, user]);

    useEffect(() => {
        const subId = chatroomData.subId;
        if(!subId) return;

        // TODO: make optional?
        EnableSevenTV(subId).then(e => setEmotes(curEmotes => {
            const newEmotes: TChatroomEmotes = {};
            Object.assign(newEmotes, curEmotes);
            newEmotes['seventv'] = new Map(Object.entries(e));

            return newEmotes;
        })).catch(e => console.error(e));

        setEmotes(curEmotes => {
            const newEmotes: TChatroomEmotes = {};
            Object.assign(newEmotes, curEmotes);
            newEmotes['channel'] = new Map(Object.entries(chatroomData.channelEmotes));
            // newEmotes['seventv'] = new Map(Object.entries(chatroomData.sevenTVEmotes));

            return newEmotes;
        });

        setChatMessages([]);

        EventsOn(subId, appendChatMessage);
        
        return () => {
            EventsOff(subId);
            DeleteSubscription(user.access_token, subId);
        }
    }, [chatroomData]);

    return { chatMessages, sendChatMessage, emotes }
}
