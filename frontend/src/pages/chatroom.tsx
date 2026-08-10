import { useNavigate, useParams } from 'react-router-dom';
import ChatMessage from '@components/chat/chat-message';
import { useCallback, useContext, useEffect, useState } from 'react';
import AutoScroller from '@components/util/auto-scroller';
import JumpToRecentPopup from '@components/chat/jump-to-recent-popup';
import UserPopup, { IPopupUser } from '@components/chat/user-popup';
import useChat, { IAppChatMessage } from '@/hooks/chat';
import { GlobalContext } from '@/contexts/global-context';
import ChatControls from '@/components/chat/chat-controls';
import StreamDataIsland from '@/components/chat/stream-data-island';

interface IChatroomProps {
}

export default function Chatroom({
}: IChatroomProps) {
    const { user } = useContext(GlobalContext);
    const navigate = useNavigate()

    if(!user) {
        navigate('/');
        return <></>;
    }
    
    const { channel } = useParams();

    const MAX_MESSAGES = 200;
    const { chatMessages, sendChatMessage, emotes, streamData } = useChat({
                         channel,
                         user,
                         maxMessages: MAX_MESSAGES,
                    });

    const [isReplying, setIsReplying] = useState<boolean>(false);
    const [replyingToMessage, setReplyingToMessage] = useState<IAppChatMessage|undefined>();


    const [shouldShowUserPopup, setShouldShowUserPopup] = useState<boolean>(false);
    const [curPopupUser, setCurPopupUser] = useState<IPopupUser>();
    const [initUserPopupPos, setInitUserPopupPos] = useState<{x:number, y:number}>({x: 0, y: 0});


    useEffect(() => {
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

    }, [chatMessages.length, curPopupUser]);

    const handleChatReplyClick = useCallback((message: IAppChatMessage) => {
        setIsReplying(true);
        setReplyingToMessage(message);
    }, []);

    const handleChatReplyClose = useCallback(() => {
        setIsReplying(false);
        setReplyingToMessage(undefined);
    }, []);

    const handleSendMessage = useCallback(async (message: string) => {
        const res = await sendChatMessage(message, replyingToMessage?.id);
        if(res) {
            handleChatReplyClose();
        }

        return res;
    }, [replyingToMessage]);

    const getChatterColor = useCallback((username: string): string => {
        const chatterMessage = chatMessages.find(m => m.username === username);

        return chatterMessage?.color ?? "var(--color-chatter-text-secondary)";
    }, [chatMessages]);

    const showUserPopup = useCallback(async (username: string|undefined, mouseX: number, mouseY: number) => {
        if(!username) return;

        const recentMessages = chatMessages.filter(m => m.username === username);
        setCurPopupUser({
            username,
            messages: recentMessages,
        })
        setInitUserPopupPos({x: mouseX, y: mouseY});

        setShouldShowUserPopup(true);
    }, [chatMessages]);

    return (
        <div
            className={'flex flex-col w-full h-full grow relative'}
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
            <ChatControls 
                isReplying={isReplying}
                replyingToMessage={replyingToMessage}
                emotes={emotes}
                getChatterColor={getChatterColor}
                onSendMessage={handleSendMessage}
                onShowUserPopup={showUserPopup}
                onReplyClosed={handleChatReplyClose}
            />
            <div className="absolute right-5 top-2.5 z-4999">
                <StreamDataIsland 
                    streamData={streamData}
                />
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
