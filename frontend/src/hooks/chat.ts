import { TChatroomEmotes } from "@/api/emote";

import { ConnectToChatroom, EnableSevenTV, SendChatMessage } from '@wailsjs/chatter-wails/appservice';
import { Events } from "@wailsio/runtime";

import { useContext, useEffect, useState } from "react";
import { AppUser } from "@wailsjs/chatter-wails/services/auth";
import { ESChatMessage, StreamData } from "@wailsjs/chatter-wails/services/eventsub";
import { AppEmote } from "@wailsjs/chatter-wails/services/emote";
import { assertDefined, isDefined } from "@/util/assert";
import { GlobalContext } from "@/contexts/global-context";

export type TBanTypeInfo = 
    | {
        isPermanent: false,
        duration: number
      }
    | {
        isPermanent: true,
      }

export type TBanInfo = 
    | {
        isBanned: true,
        banTypeInfo: TBanTypeInfo,
      }
    | {
        isBanned: false,
      }

export interface IAppChatMessage extends ESChatMessage {
    banInfo: TBanInfo;
    deleted: boolean;
}

export default function useChat({ channel, user, emoteRecord, maxMessages = 200 }: {
    channel: string|undefined,
    user: AppUser,
    emoteRecord: TChatroomEmotes,
    maxMessages?: number,
}) {

    const [chatMessages, setChatMessages] = useState<IAppChatMessage[]>([]);

    const [emotes, setEmotes] = useState<TChatroomEmotes>(emoteRecord);

    const [streamData, setStreamData] = useState<StreamData>({channel: "", live: false, viewCount: 0, title: "", gameName: "" });

    const { broadcastError } = useContext(GlobalContext);
    

    const appendChatMessage = (message: ESChatMessage) => {
        const appMessage: IAppChatMessage = {
            ...message,
            banInfo: { isBanned: false },
            deleted: false,
        }
        setChatMessages(curMessages => {
            const numExtraMessages = curMessages.length - maxMessages;
            if(numExtraMessages >= 0) {
                return [...curMessages.slice(numExtraMessages + 1), appMessage];
            } else {
                return [...curMessages, appMessage];
            }
        })
    }

    const handleChatMessageEvent = (event: Events.WailsEvent<"common:chat-message">) => {
        if(event.data?.channel && channel === event.data.channel) {
            appendChatMessage(event.data);
        }
    }

    const handleBanEvent = (event: Events.WailsEvent<"common:ban">) => {
        if(event.data.channel !== channel) return;

        const banTypeInfo: TBanTypeInfo = event.data.isPermanent
        ? {
            isPermanent: true,
          }
        : {
            isPermanent: false,
            duration: event.data.duration ?? 0,
          }
        const banInfo: TBanInfo = {
            isBanned: true,
            banTypeInfo,
        }

        setChatMessages(cur => {
            const matchedIndices = cur.flatMap((mes, i) => mes.username.toLowerCase() === event.data.userLogin ? [i] : []);
            const newMessages = [...cur];
            
            for(const i of matchedIndices) {
                newMessages[i] = { ...newMessages[i], banInfo }
            }

            return newMessages;
        })
    }

    const handleClearMsgEvent = (event: Events.WailsEvent<"common:clear-msg">) => {
        if(event.data.channel !== channel) return;

        setChatMessages(cur => {
            const msgToDeleteIndex = cur.findIndex(m => m.id === event.data.messageID);
            if(msgToDeleteIndex === -1) return cur;

            const newMessages = [...cur];
            newMessages[msgToDeleteIndex] = { ...newMessages[msgToDeleteIndex], deleted: true };
            return newMessages;
        })
    }

    const sendChatMessage = async (message: string, replyId: string|null|undefined) => {
        if(!channel) return false;
        const trimmed = message.trim();
        if(trimmed.length === 0) return false;

        try {
            await SendChatMessage(channel, trimmed, replyId ?? null);
            return true;
        } catch(err) {
            broadcastError(err);
            return false;
        }
    }
    
    const addEmoteSet = (set: Record<string, AppEmote|null|undefined>, setName: string) => {
        setEmotes(curEmotes => {
            const newEmotes: TChatroomEmotes = {};
            Object.assign(newEmotes, curEmotes);

            // always make the set a new map to avoid accumulation between channels
            newEmotes[setName] = new Map();

            for(const [name, emote] of Object.entries(set)) {
                if(!isDefined(emote)) continue;
                newEmotes[setName].set(name, emote);
            }

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
        offFns.push(Events.On('common:ban', handleBanEvent));
        offFns.push(Events.On('common:clear-msg', handleClearMsgEvent));

        return () => {
            offFns.forEach(fn => fn());
        }
    }

    useEffect(() => {
        if(!channel) return;
        const abortController = new AbortController();

        setChatMessages([]);

        ConnectToChatroom(channel)
            .then(d => {
                assertDefined(d?.channelEmotes);
                addEmoteSet(d.channelEmotes, 'channel');
            })
            .then(() => emitChatOpenState(channel, user.access_token, true))
            // TODO: make optional?
            .then(() => EnableSevenTV(channel))
            .then(e => {
                assertDefined(e);
                addEmoteSet(e, 'seventv');
            })
            .cancelOn(abortController.signal)
            .catch(broadcastError);

        const listenersOff = listenersOn();
        
        return () => {
            abortController.abort("Channel/user changed");
            emitChatOpenState(channel, user.access_token, false);
            listenersOff();
        }
    }, [channel, user]);

    return { chatMessages, sendChatMessage, emotes, streamData }
}
