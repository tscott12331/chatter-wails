import { useEffect, useState } from 'react';
import ChatMessage, { TChatMessage } from './chat-message';
import styles from './user-popup.module.css';
import PlusIcon from '../svg/plus-icon';

export interface IPopupUser {
    username: string;
    messages: TChatMessage[];
}

interface IUserPopupProps {
    user?: IPopupUser;
    initPos: { x: number, y: number };
    showChatReplyButton?: boolean;
    getChatterColor: (username: string) => string;
    showUserPopup: (username: string|undefined, mouseX: number, mouseY: number) => void;
    onChatReplyClick?: (message: TChatMessage) => void;
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

    useEffect(() => {
        setPos(initPos);
    }, [initPos]);

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
        <div className={styles.wrapper + ' absolute z-2000'}
            style={{
                left: `${pos.x}px`,
                top: `${pos.y}px`,
            }}
        >
            <section
                className={styles.infoSection + ' flex-column'}
                onMouseDown={handleMouseDown}
            >
                <div className={styles.controls + ' flex-justify-end'}>
                    <PlusIcon
                        onClick={onPopupClose}
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
            <section className={styles.chatSection + ' flex-col scroller-y'}>
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
                    <p className={styles.noMessages}>no recent messages</p>
                }
            </section>
        </div>
    )
}
