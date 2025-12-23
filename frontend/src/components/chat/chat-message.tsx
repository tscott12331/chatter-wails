import styles from './chat-message.module.css';

import { getEmoteSrcSet } from '@api/native-emote';
import ReplyIcon from '@components/svg/reply-icon';
import Tooltip from '@components/util/tooltip';
import React from 'react';

interface IChatMessageFragment {
    type: 'text'|'cheermote'|'emote'|'mention';
    text: string;
    cheermote: {
        prefix: string;
        bits: number;
        tier: number;
    }|null;
    emote: {
        id: string;
        emote_set_id: string;
        owner_id: string;
        format: ('static'|'animated')[];
    }|null;
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

export type TChatMessage = {
    id: string;
    username: string;
    text: string;
    fragments: IChatMessageFragment[];
    color: string;
    badges: IMessageBadge[];
    reply: IMessageReply|null;
}

interface ChatMessageProps {
    message: TChatMessage;
    onChatReplyClick?: (message: TChatMessage) => void;
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
    const fragmentToNode = (fragment: IChatMessageFragment, index: number): React.ReactNode => {
        switch(fragment.type) {
            case 'emote':
                if(!fragment.emote) return fragment.text;

                return (
                    <Tooltip
                        text={fragment.text}
                        hoverTime={0}
                        key={index}
                        >
                            <img
                                srcSet={getEmoteSrcSet(
                                        fragment.emote.id,
                                        fragment.emote.format,
                                        )}
                                alt={fragment.text}
                            />
                    </Tooltip>
                );
            case 'mention':
                if(!fragment.mention) return fragment.text;
                return (
                    <span
                        className={styles.mention}
                        style={{
                            color: getChatterColor(fragment.mention.user_name),
                        }}
                        onClick={(e: React.MouseEvent<HTMLSpanElement>) =>
                                showUserPopup(fragment.mention?.user_name, e.pageX, e.pageY)}
                        key={index}
                    >{fragment.text}</span>
                );
            case 'text':
            default:
                return fragment.text;
        }
    }

    const defragmentMessage = (fragments: IChatMessageFragment[]): React.ReactNode => {
        return fragments.map((fragment, i) => fragmentToNode(fragment, i));
    }

    const badgeToNode = (badge: IMessageBadge, index: number): React.ReactNode => {
        return <Tooltip text={badge.title} hoverTime={0} key={index}>
                    <img srcSet={badge.srcSet} />
               </Tooltip>
    }

    const badgesToNodes = (badges: IMessageBadge[]): React.ReactNode => {
        return badges.map((badge, i) =>
            badgeToNode(badge, i)
        );
    }

    return (
        <div className={styles.wrapper + ' flex-column'}>
            {message.reply &&
            <Tooltip
                hoverTime={0}
                text={message.reply.parent_message_body}
            >
                <p className={styles.reply + ' ellipsis'}>replying to @{message.reply.parent_user_name}: {message.reply.parent_message_body}</p>
            </Tooltip>
            }
            <p>
                {badgesToNodes(message.badges)}
                <span
                    className={styles.username}
                    style={{
                        color: message.color
                    }}
                    onClick={(e: React.MouseEvent<HTMLSpanElement>) => showUserPopup(message.username, e.pageX, e.pageY)}
                >
                    {message.username}
                </span>
                <span className={styles.colon}>: </span>
                {defragmentMessage(message.fragments)}
            </p>
            <div
                className={styles.chatControls + ' hidden flex flex-justify-end flex-align-center'}
            >
                {showChatReplyButton &&
                <div className={styles.replyWrapper}
                    onClick={() => onChatReplyClick(message)}
                >
                    <ReplyIcon />
                </div>
                }
            </div>
        </div>
    )
}
