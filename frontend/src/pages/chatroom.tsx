import { ESCon, ESSubscription, IESNotification, TESMessage } from '@/api/eventsub';
import { useParams } from 'react-router-dom';
import ChatMessage, { TChatMessage } from '@components/chat/chat-message';
import { useEffect, useMemo, useRef, useState } from 'react';
import { TUser } from '@/App';
import { DebugLogger } from '@util/debug';
// import { connectToChatroom, deleteSubscription, ESCon } from '@api/eventsub';
import { getUser } from '@api/user-info';
import AutoScroller from '@components/util/auto-scroller';
import JumpToRecentPopup from '@components/chat/jump-to-recent-popup';
import { sendMessage } from '@api/messages';
import ReplyPopup from '@components/chat/reply-popup';
import { combineChannelGlobalSets, esBadgesToMessageBadges, getChannelBadges, IBadgeSet } from '@api/badges';
import EmoteIcon from '@components/svg/emote-icon';
import { getEmoteSrcSet, IGlobalEmote } from '@api/native-emote';
import Tooltip from '@components/util/tooltip';
import { preload } from 'react-dom';
import { moveCursorToEnd } from '@util/rte';
import UserPopup, { IPopupUser } from '@components/chat/user-popup';
import useChat from '@/hooks/chat';

interface IChatroomProps {
    user: TUser|undefined;
    globalBadgeSets: IBadgeSet[],
    globalEmotes: IGlobalEmote[],
}

const dbLog = new DebugLogger();

export default function Chatroom({
    user,
    globalBadgeSets,
    globalEmotes,
}: IChatroomProps) {
    if(!user) return <></>;
    
    const { channel } = useParams();

    const MAX_MESSAGES = 200;
    const { chatMessages, sendChatMessage } = useChat({channel, user, globalBadgeSets, maxMessages: MAX_MESSAGES});

    const [isReplying, setIsReplying] = useState<boolean>(false);
    const [replyingToMessage, setReplyingToMessage] = useState<TChatMessage|undefined>();


    const [shouldShowEmotePopup, setShouldShowEmotePopup] = useState<boolean>(false);

    const [shouldShowUserPopup, setShouldShowUserPopup] = useState<boolean>(false);
    const [curPopupUser, setCurPopupUser] = useState<IPopupUser>();
    const [initUserPopupPos, setInitUserPopupPos] = useState<{x:number, y:number}>({x: 0, y: 0});

    const messageInputRef = useRef<HTMLDivElement>(null);


    const lastMessage = chatMessages.at(-1);
    const lastPopupMessage = curPopupUser?.messages.at(-1);
    if(curPopupUser && lastMessage
       && curPopupUser.username === lastMessage.username
      && (lastMessage.id !== lastPopupMessage?.id)) {
        setCurPopupUser(cur => {
            if(!cur) return cur;
            const numExtraMessages = curPopupUser.messages.length - MAX_MESSAGES;
            return {
                username: cur.username,
                messages: numExtraMessages >= 0
                    ? [...cur.messages.slice(numExtraMessages + 1), lastMessage]
                    : [...cur.messages, lastMessage],
            }
        })
    }
    const inputNodesToText = (nodes: NodeListOf<ChildNode>|undefined) => {
        if(!nodes) return "";

        let text = "";
        for(const node of nodes) {
            switch(node.nodeType) {
                case Node.TEXT_NODE:
                    text += node.nodeValue;
                    break;
                case Node.ELEMENT_NODE:
                    text += (node as HTMLImageElement).alt;
                    break;
            }
        }

        return text;
    }

    const handleSendMessage = async () => {
        if(!messageInputRef.current) return;

        const message = inputNodesToText(messageInputRef.current.childNodes);
        const res = await sendChatMessage(message, replyingToMessage?.id);
        if(res.success) {
            messageInputRef.current.innerHTML = '';
            handleChatReplyClose();
        } else {
            console.log(res.error);
        }
    }

    const handleMessageInputKeydown = (e: React.KeyboardEvent<HTMLDivElement>) => {
        if(e.key === "Enter") {
            e.preventDefault();
            handleSendMessage();
        }
    }

    const handleMessageInput = (e: React.InputEvent<HTMLDivElement>) => {
        const childNodes = e.currentTarget.childNodes;
        const prevNode = childNodes.item(childNodes.length - 2);
        const curNode = childNodes.item(childNodes.length - 1);
        if(prevNode && prevNode.nodeType === Node.ELEMENT_NODE
          && !(curNode.textContent?.charCodeAt(0) === 32
              || curNode.textContent?.charCodeAt(0) === 160)) {
              curNode.textContent = ' ' + curNode.textContent;
              moveCursorToEnd(e.currentTarget);
        }

        const curTextContent = curNode?.textContent;
        if(!curTextContent) return;

        const potentialEmote = curTextContent.split(' ').at(-1);
        if(potentialEmote) {
            for(let i = 0; i < globalEmotes.length; i++) {
                const emote = globalEmotes[i];
                if(potentialEmote === emote.name) {
                    curNode.textContent = curTextContent.slice(0, curTextContent.length - potentialEmote.length);
                    handleEmoteSelect(emote);
                }
            }
        }
    }

    const handleChatReplyClick = (message: TChatMessage) => {
        setIsReplying(true);
        setReplyingToMessage(message);
    }

    const handleChatReplyClose = () => {
        setIsReplying(false);
        setReplyingToMessage(undefined);
    }

    const handleEmoteButtonClick = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation();
        setShouldShowEmotePopup(cur => !cur);
    }

    const handleEmoteSelect = (emote: IGlobalEmote) => {
        if(!messageInputRef.current) return;
        const srcSet = getEmoteSrcSet(emote.id, emote.format) ?? "";
        if(!srcSet) return;

        const concat = `<img class="inline" srcset="${srcSet}" alt="${emote.name}"/>`;

        if(messageInputRef.current.innerHTML.at(-1) === ' ') {
            messageInputRef.current.innerHTML += `${concat} `;
        } else {
            messageInputRef.current.innerHTML += ` ${concat} `;
        }

        messageInputRef.current.focus();
        moveCursorToEnd(messageInputRef.current);
    }

    const getChatterColor = (username: string): string => {
        const chatterMessage = chatMessages.find(m => m.username === username);

        return chatterMessage?.color ?? "var(--c-text-2)";
    }

    const showUserPopup = async (username: string|undefined, mouseX: number, mouseY: number) => {
        if(!username) return;

        const recentMessages = chatMessages.filter(m => m.username === username);
        setCurPopupUser({
            username,
            messages: recentMessages,
        })
        setInitUserPopupPos({x: mouseX, y: mouseY});

        setShouldShowUserPopup(true);
    }

    const filteredGlobalEmotes = useMemo<React.ReactNode>(() => {
        const idHash: Record<string, string> = {};

        const filteredArr: IGlobalEmote[] = [];

        for(let i = 0; i < globalEmotes.length; i++) {
            const emote = globalEmotes[i];
            if(!(emote.id in idHash)) {
                idHash[emote.id] = emote.id;
                filteredArr.push(emote);
            }
        }

        return filteredArr.map(emote => {
                const srcSet = getEmoteSrcSet(emote.id, emote.format);
                for(const image in emote.images) {
                    preload(image, { as: "image", imageSrcSet: srcSet });
                }

                return <Tooltip
                text={emote.name}
                hoverTime={0}
                key={emote.id}
                >
                    <div
                        className='flex-center'
                        onClick={() => handleEmoteSelect(emote)}
                    >
                        <img
                            srcSet={srcSet}
                        />
                    </div>
                </Tooltip>
        }
            );
    }, [globalEmotes]);

    return (
        <div
            className={'flex flex-col w-full h-full grow relative'}
            onClick={() => setShouldShowEmotePopup(false)}
        >
            <AutoScroller
                style={{
                    paddingBottom: "2px",
                }}
                jumpToRecentPopup={JumpToRecentPopup}
            >
                {chatMessages.map(message =>
                    <ChatMessage
                    message={message}
                    onChatReplyClick={handleChatReplyClick}
                    getChatterColor={getChatterColor}
                    showUserPopup={showUserPopup}
                    key={message.id}
                    />
                 )}
            </AutoScroller>
            <div className='flex flex-col p-1 basis-[max-content] bg-bg-9 relative'>
                {shouldShowEmotePopup &&
                <div
                className='w-[calc(100%-30px)] h-75 border border-outline-1 rounded-xs m-3.5 absolute left-0 bottom-full bg-bg-2/80 backdrop-blur-xs g-1 p-1 z-600 flex flex-wrap justify-between items-center scroller-y [&>span>div]:w-10 [&>span>div]:h-10 [&>span>div]:p-0.5 [&>span>div]:rounded-xs [&>span>div]:opacity-90 [&>span>div]:hover:bg-bg-5'
                onClick={(e: React.MouseEvent<HTMLDivElement>) => e.stopPropagation()}
                >
                    {filteredGlobalEmotes}
                </div>
                }
                {isReplying && replyingToMessage &&
                <div>
                    <ReplyPopup
                        onCloseClicked={handleChatReplyClose}
                        message={replyingToMessage}
                        getChatterColor={getChatterColor}
                        showUserPopup={showUserPopup}
                    />
                </div>
                }
                <div className={'flex items-center justify-between flex-wrap g-0.5 h-37.5'}>
                    <div
                    className='w-full h-[calc(100%-46px)]'
                    contentEditable="true"
                    onKeyDown={handleMessageInputKeydown}
                    onInput={handleMessageInput}
                    ref={messageInputRef}
                    onClick={(e: React.MouseEvent<HTMLDivElement>) => e.stopPropagation()}
                    ></div>
                    <div className='flex items-center justify-end gap-1.5 w-full'>
                        <button
                            className='relative flex justify-center items-center w-12.5 p-2.5 bg-none hover:drop-shadow-xs hover:drop-shadow-text-2'
                            onClick={handleEmoteButtonClick}
                        >
                                <EmoteIcon />
                        </button>
                        <button
                            className="h-12.5 p-2.5"
                            onClick={handleSendMessage}
                        >send</button>
                    </div>
                </div>
            </div>
            {shouldShowUserPopup &&
                <UserPopup
                    user={curPopupUser}
                    initPos={initUserPopupPos}
                    onChatReplyClick={handleChatReplyClick}
                    getChatterColor={getChatterColor}
                    showUserPopup={showUserPopup}
                    onPopupClose={() => setShouldShowUserPopup(false)}
                />
            }
        </div>
    )
}
