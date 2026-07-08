import { TChatroomEmotes } from "@/api/emote";

import { ConnectToChatroom, EnableSevenTV, SendChatMessage } from '@wailsjs/chatter-wails/appservice';
import { Events } from "@wailsio/runtime";

import { useEffect, useState } from "react";
import { AppUser } from "@wailsjs/chatter-wails/services/auth";
import { ESChatMessage, StreamData } from "@wailsjs/chatter-wails/services/eventsub";
import { AppEmote } from "@wailsjs/chatter-wails/services/emote";

export default function useChat({ channel, user, emoteRecord, maxMessages = 200 }: {
    channel: string|undefined,
    user: AppUser,
    emoteRecord: TChatroomEmotes,
    maxMessages?: number,
}) {

    const [chatMessages, setChatMessages] = useState<ESChatMessage[]>([]);

    const [emotes, setEmotes] = useState<TChatroomEmotes>(emoteRecord);

    const [streamData, setStreamData] = useState<StreamData>({channel: "", live: false, viewCount: 0, title: "" });
    

    const appendChatMessage = (message: ESChatMessage) => {
        setChatMessages(curMessages => {
            const numExtraMessages = curMessages.length - maxMessages;
            if(numExtraMessages >= 0) {
                return [...curMessages.slice(numExtraMessages + 1), message];
            } else {
                return [...curMessages, message];
            }
        })
    }

    const handleChatMessageEvent = (event: Events.WailsEvent<"common:chat-message">) => {
        if(channel === event.data.channel) {
            appendChatMessage(event.data);
        }
    }

    const sendChatMessage = async (message: string, replyId?: string) => {
        if(!channel) return false;
        const trimmed = message.trim();
        if(trimmed.length === 0) return false;

        try {
            await SendChatMessage(channel, trimmed, replyId);
            return true;
        } catch(err) {
            console.error(err);
            return false;
        }
    }
    
    const addEmoteSet = (set: Record<string, AppEmote>, setName: string) => {
        setEmotes(curEmotes => {
            const newEmotes: TChatroomEmotes = {};
            Object.assign(newEmotes, curEmotes);
            newEmotes[setName] = new Map(Object.entries(set));

            return newEmotes;
        })
    }

    const emitChatOpenState = (channel: string, accessToken: string, open: boolean) => {
        Events.Emit('common:chat-open', {
            channel,
            accessToken,
            open,
        });
    }

    const listenersOn = () => {
        const offFns: (() => void)[] = [];
        offFns.push(Events.On('common:chat-message', handleChatMessageEvent));
        offFns.push(Events.On('common:stream-data', (e) => setStreamData(e.data)));

        return () => {
            offFns.forEach(fn => fn());
        }
    }

    useEffect(() => {
        if(!channel) return;
        const abortController = new AbortController();

        setChatMessages([]);

        ConnectToChatroom(channel)
            .then(d => addEmoteSet(d.channelEmotes, 'channel'))
            .then(() => emitChatOpenState(channel, user.access_token, true))
            // TODO: make optional?
            .then(() => EnableSevenTV(channel))
            .then(e => addEmoteSet(e, 'seventv'))
            .cancelOn(abortController.signal)
            .catch(err => console.error(err));

        const listenersOff = listenersOn();
        
        return () => {
            abortController.abort("Channel/user changed");
            emitChatOpenState(channel, user.access_token, false);
            listenersOff();
        }
    }, [channel, user]);

    return { chatMessages, sendChatMessage, emotes, streamData }
}
