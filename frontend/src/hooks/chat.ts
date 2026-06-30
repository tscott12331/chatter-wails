import { TChatroomEmotes } from "@/api/emote";
import { IAppEmote } from "@/api/native-emote";
import { TUser } from "@/App";
import { TChatMessage } from "@/components/chat/chat-message";
import { IViewcountData } from "@/components/chat/viewcount";

import { ConnectToChatroom, EnableSevenTV, SendChatMessage } from '@wailsjs/chatter-wails/appservice';
import { Events } from "@wailsio/runtime";

import { useEffect, useState } from "react";

interface IChatOpenData {
    channel: string;
    accessToken: string;
    open: boolean;
}

interface IChatroomData {
    subId: string|null;
    broadcasterId: string|null;
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
        channelEmotes: {},
    })

    const [emotes, setEmotes] = useState<TChatroomEmotes>(emoteRecord);

    const [viewcountData, setViewcountData] = useState<IViewcountData>({live: false, viewCount: 0});
    

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
            // @ts-ignore
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
            // @ts-ignore
            newEmotes['seventv'] = new Map(Object.entries(e));

            return newEmotes;
        })).catch(e => console.error(e));

        setEmotes(curEmotes => {
            const newEmotes: TChatroomEmotes = {};
            Object.assign(newEmotes, curEmotes);
            newEmotes['channel'] = new Map(Object.entries(chatroomData.channelEmotes));

            return newEmotes;
        });

        setChatMessages([]);

        Events.On(subId, (d) => appendChatMessage(d.data));
        const viewcountEventName = `viewcount:${channel?.toLowerCase()}`;
        if(channel) {
            Events.On(viewcountEventName, (d) => setViewcountData(d.data));

            const chatOpenData: IChatOpenData = {
                channel,
                accessToken: user.access_token,
                open: true
            }
            Events.Emit("chatopen",  chatOpenData);
        }

        
        return () => {
            Events.Off(subId);
            if(channel) {
                Events.Off(viewcountEventName);

                const chatOpenData: IChatOpenData = {
                    channel,
                    accessToken: user.access_token,
                    open: false,
                }
                Events.Emit("chatopen",  chatOpenData);

            }
        }
    }, [chatroomData]);

    return { chatMessages, sendChatMessage, emotes, viewcountData }
}
