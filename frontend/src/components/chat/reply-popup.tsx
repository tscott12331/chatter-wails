import ChatMessage, { TChatMessage } from './chat-message';
import styles from './reply-popup.module.css';
import PlusIcon from '../svg/plus-icon';

interface IReplyPopupProps {
    message: TChatMessage;
    onCloseClicked?: () => void;
    getChatterColor: (username: string) => string;
    showUserPopup: (username: string|undefined, mouseX: number, mouseY: number) => void;
}

export default function ReplyPopup({
    message,
    onCloseClicked=()=>{},
    getChatterColor,
    showUserPopup,
}: IReplyPopupProps) {
    return (
        <div className={styles.wrapper + ' flex-column'}>
            <div className={styles.topBar + ' flex-justify-space-btw flex-align-center'}>
                <i>replying to:</i>
                <PlusIcon
                    onClick={onCloseClicked}
                />
            </div>
            <div className={styles.messageWrapper}>
                <ChatMessage
                    message={message}
                    showChatReplyButton={false}
                    getChatterColor={getChatterColor}
                    showUserPopup={showUserPopup}
                />
            </div>
        </div>
    )
}
