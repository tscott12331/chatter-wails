import { useEffect, useState } from 'react';
import ChatMessage from './chat-message';
import PlusIcon from '../svg/plus-icon';
import { IAppChatMessage } from '@/hooks/chat';

export interface IPopupUser {
    username: string;
    messages: IAppChatMessage[];
}

interface IUserPopupProps {
    user?: IPopupUser;
    initPos: { x: number, y: number };
    showChatReplyButton?: boolean;
    getChatterColor: (username: string) => string;
    showUserPopup: (username: string|undefined, mouseX: number, mouseY: number) => void;
    onChatReplyClick?: (message: IAppChatMessage) => void;
    onPopupClose?: () => any;
}

export default function UserPopup({
    user,
    initPos,
    showChatReplyButton,
    getChatterColor,
    showUserPopup,
    onChatReplyClick,
    onPopupClose,
}: IUserPopupProps) {
    const [pos, setPos] = useState<{x: number, y: number}>(initPos);
    const [isDragging, setIsDragging] = useState<boolean>(false);

    const handleMouseMove = (e: MouseEvent) => {
        setPos(curPos => ({
            x: curPos.x + e.movementX,
            y: curPos.y + e.movementY,
        }))
    }

    const handleMouseUp = (e: MouseEvent) => {
        setIsDragging(false);
    }

    const handleMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
        setIsDragging(true);
    }

    useEffect(() => {
        if(isDragging) {
            document.addEventListener('mousemove', handleMouseMove);

            return () => document.removeEventListener('mousemove', handleMouseMove);
        }

        return;
    }, [isDragging]);

    useEffect(() => {
        document.addEventListener('mouseup', handleMouseUp);

        return () => document.removeEventListener('mouseup', handleMouseUp);
    }, []);

    return (
        <div className='absolute w-75 h-100 bg-chatter-surface-elevated/90 backdrop-blur-xs drop-shadow-[0_0_3px] drop-shadow-chatter-border outline outline-chatter-border-strong/30 rounded-sm z-5000'
            style={{
                left: `${pos.x}px`,
                top: `${pos.y}px`,
            }}
        >
            <section
                className={'flex flex-col h-[15%] border-b border-b-chatter-border-strong/30 cursor-move p-1'}
                onMouseDown={handleMouseDown}
            >
                <div className={'flex justify-end h-1/4 [&_svg]:h-full [&_svg]:rotate-45 [&_svg]:fill-red-500 [&_svg]:cursor-pointer [&_svg]:hover:fill-red-700'}>
                    <PlusIcon
                        onClick={onPopupClose}
                        className="size-full"
                    />
                </div>
                <h3
                style={{
                    color: `${user?.messages.at(0)?.color ?? 'initial'}`,
                }}
                >
                    {user?.username}
                </h3>
            </section>
            <section className={'flex flex-col h-[85%] scroller-y'}>
                {
                    user && user.messages.length > 0
                    ?
                    user?.messages.map(m =>
                    <ChatMessage
                        message={m}
                        showChatReplyButton={showChatReplyButton}
                        getChatterColor={getChatterColor}
                        showUserPopup={showUserPopup}
                        onChatReplyClick={onChatReplyClick}
                        key={m.id}
                    />
                    )
                    :
                    <p className='p-1 text-center'>no recent messages</p>
                }
            </section>
        </div>
    )
}
