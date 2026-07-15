import ReplyIcon from '@components/svg/reply-icon';
import Tooltip from '@components/util/tooltip';
import { AppEmote } from '@wailsjs/chatter-wails/services/emote';
import { AppChatMessageFragment, ESChatMessage } from '@wailsjs/chatter-wails/services/eventsub';
import React from 'react';

export interface IChatMessageFragment {
    type: 'text'|'cheermote'|'emote'|'mention';
    text: string;
    cheermote: {
        prefix: string;
        bits: number;
        tier: number;
    }|null;
    emote: AppEmote|null;
    mention: {
        user_id: string;
        user_name: string;
        user_login: string;
    }|null;
}

export interface IMessageBadge {
    srcSet: string;
    info: string;
    title: string;
}

export interface IMessageReply {
    parent_message_body: string;
    parent_message_id: string;
    parent_user_id: string;
    parent_user_login: string;
    parent_user_name: string;
    thread_message_id: string;
    thread_user_id: string;
    thread_user_login: string;
    thread_user_name: string;
}

interface ChatMessageProps {
    message: ESChatMessage;
    onChatReplyClick?: (message: ESChatMessage) => void;
    showChatReplyButton?: boolean;
    getChatterColor: (username: string) => string;
    showUserPopup: (username: string|undefined, mouseX: number, mouseY: number) => void;
}

export default function ChatMessage({
    message,
    onChatReplyClick=() => {},
    showChatReplyButton=true,
    getChatterColor,
    showUserPopup,
}: ChatMessageProps) {
    const fragmentToNode = (fragment: AppChatMessageFragment, index: number): React.ReactNode => {
        switch(fragment.type) {
            case 'emote':
                if(!fragment.emote) return fragment.text;

                return (
                    <Tooltip
                        text={fragment.text}
                        hoverTime={0}
                        key={index}
                        >
                        <div className="align-middle inline-grid place-items-center grid-cols-1 grid-rows-1">
                            <img
                                className="row-1 col-1"
                                srcSet={fragment.emote.darkSrcSet.length > 0 ? fragment.emote.darkSrcSet : fragment.emote.lightSrcSet}
                                alt={fragment.text}
                            />
                            {
                            fragment.emote.emoteStack?.map(e =>
                                e &&
                                <img
                                    className="row-1 col-1"
                                    srcSet={e.darkSrcSet.length > 0 ? e.darkSrcSet : e.lightSrcSet}
                                    alt={fragment.text}
                                />
                             )
                            }
                        </div>
                    </Tooltip>
                );
            case 'mention':
                if(!fragment.mention) return fragment.text;
                return (
                    <span
                        style={{
                            color: getChatterColor(fragment.mention.user_name),
                        }}
                        onClick={(e: React.MouseEvent<HTMLSpanElement>) =>
                                showUserPopup(fragment.mention?.user_name, e.pageX, e.pageY)}
                        className="cursor-pointer"
                        key={index}
                    >{fragment.text}</span>
                );
            case 'text':
            default:
                return fragment.text;
        }
    }

    const defragmentMessage = (fragments: (AppChatMessageFragment|null)[]): React.ReactNode => {
        return fragments.map((fragment, i) => fragment && fragmentToNode(fragment, i));
    }

    const badgeToNode = (badge: IMessageBadge, index: number): React.ReactNode => {
        return <Tooltip text={badge.title} hoverTime={0} key={index}>
                    <img className="inline mr-1 max-w-4.5" srcSet={badge.srcSet} />
               </Tooltip>
    }

    const badgesToNodes = (badges: IMessageBadge[]): React.ReactNode => {
        return badges.map((badge, i) =>
            badgeToNode(badge, i)
        );
    }

    return (
        <div className="p-1.5 relative hover:bg-bg-1 hover:[&_.chat-controls]:visible">
            {message.reply &&
            <Tooltip
                hoverTime={0}
                text={message.reply.parent_message_body}
            >
                <p className='w-full text-text-3 text-sm ellipsis'>replying to @{message.reply.parent_user_name}: {message.reply.parent_message_body}</p>
            </Tooltip>
            }
            <div>
                {message.badges && badgesToNodes(message.badges)}
                <span
                    className="contrast-300 font-medium my-0 mr-0.5 cursor-pointer"
                    style={{
                        color: message.color
                    }}
                    onClick={(e: React.MouseEvent<HTMLSpanElement>) => showUserPopup(message.username, e.pageX, e.pageY)}
                >
                    {message.username}
                </span>
                <span>: </span>
                {message.fragments && defragmentMessage(message.fragments)}
            </div>
            <div
                className='chat-controls absolute -top-1 right-5 invisible flex justify-end items-center'
            >
                {showChatReplyButton &&
                <div className="w-6 h-6 p-1 rounded-xs bg-bg-3/60 backdrop-blur-xs contrast-100 hover:outline hover:outline-outline-2 hover:contrast-200 hover:brightness-120 [&_svg]:fill-text-1 [&_svg]:contrast-200 cursor-pointer"
                    onClick={() => onChatReplyClick(message)}
                >
                    <ReplyIcon />
                </div>
                }
            </div>
        </div>
    )
}
