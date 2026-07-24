import { TooltipContext } from '@/contexts/tooltip-context';
import { IAppChatMessage } from '@/hooks/chat';
import ReplyIcon from '@components/svg/reply-icon';
import { AppEmote } from '@wailsjs/chatter-wails/shared/types';
import { AppChatMessageFragment } from '@wailsjs/chatter-wails/services/eventsub';
import React, { useContext } from 'react';

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
    message: IAppChatMessage;
    onChatReplyClick?: (message: IAppChatMessage) => void;
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
    const { tooltipOn, tooltipOff } = useContext(TooltipContext);
    const fragmentToNode = (fragment: AppChatMessageFragment, index: number): React.ReactNode => {
        switch(fragment.type) {
            case 'emote':
                if(!fragment.emote) return fragment.text;

                return (
                    <div className="align-middle inline-grid place-items-center grid-cols-1 grid-rows-1" key={index}>
                        <img
                            className="row-1 col-1"
                            srcSet={fragment.emote.darkSrcSet.length > 0 ? fragment.emote.darkSrcSet : fragment.emote.lightSrcSet}
                            alt={fragment.text}
                            onMouseEnter={e => {
                                const rect = e.currentTarget.getBoundingClientRect();
                                tooltipOn({
                                    type: "image",
                                    imageSrcSet: fragment.emote ? (fragment.emote.darkSrcSet.length > 0 ? fragment.emote.darkSrcSet : fragment.emote.lightSrcSet) : '',
                                    imageDesc: fragment.text,
                                    posX: rect.x + rect.width/2,
                                    posY: rect.y,
                                }, fragment.text)
                            }}
                            onMouseLeave={() => tooltipOff(fragment.text)}
                        />
                        {
                        fragment.emote.emoteStack?.map(e =>
                            e &&
                            <img
                                className="row-1 col-1"
                                srcSet={e.darkSrcSet.length > 0 ? e.darkSrcSet : e.lightSrcSet}
                                key={e.id}
                                alt={fragment.text}
                            />
                         )
                        }
                    </div>
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
        return <img 
                    className="inline mr-1 max-w-4.5"
                    srcSet={badge.srcSet}
                    onMouseEnter={e => {
                        const rect = e.currentTarget.getBoundingClientRect();
                        tooltipOn({
                            type: "image",
                            imageSrcSet: badge.srcSet,
                            imageDesc: badge.title,
                            posX: rect.x + rect.width/2,
                            posY: rect.y,
                        }, badge.title);
                    }}
                    onMouseLeave={() => tooltipOff(badge.title)}
                    key={index} 
                />
    }

    const badgesToNodes = (badges: IMessageBadge[]): React.ReactNode => {
        return badges.map((badge, i) =>
            badgeToNode(badge, i)
        );
    }

    return (
        <div className={`p-1.5 relative hover:bg-bg-1 hover:[&_.chat-controls]:visible ${message.deleted && 'line-through hover:no-underline'}`}>
            {message.reply &&
            <p 
                className='w-full text-text-3 text-sm ellipsis'
                onMouseEnter={e => {
                    const rect = e.currentTarget.getBoundingClientRect();
                    tooltipOn({
                        type: "text",
                        text: message.reply?.parent_message_body ?? "",
                        posX: rect.x + rect.width/2,
                        posY: rect.y,
                    }, message.id + message.reply?.parent_message_id)
                }}
                onMouseLeave={() => tooltipOff(message.id + message.reply?.parent_message_id)}
            >replying to @{message.reply.parent_user_name}: {message.reply.parent_message_body}
            </p>
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
            {(message.banInfo.isBanned || message.deleted) &&
            <>
            <div className="absolute inset-0 bg-bg-09/30 pointer-events-none">
            </div>
            <p className="text-right text-text-3 pe-2">
                {
                message.banInfo.isBanned
                ? message.banInfo.banTypeInfo.isPermanent
                    ? "permanently banned"
                    : `timed out (${message.banInfo.banTypeInfo.duration}s)`
                : 'deleted'
                }
            </p>
            </>
            }
        </div>
    )
}
