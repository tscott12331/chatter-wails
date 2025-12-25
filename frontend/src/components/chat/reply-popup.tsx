import ChatMessage, { TChatMessage } from './chat-message';
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
        <div className='flex flex-col bg-bg-1 p-1.5 border border-outline-1 gap-1'>
            <div className='flex justify-between items-center h-6 [&_svg]:h-4/5 [&_svg]:rotate-45 [&_svg]:fill-red-500 [&_svg]:hover:fill-red-700'>
                <i>replying to:</i>
                <PlusIcon
                    onClick={onCloseClicked}
                />
            </div>
            <div>
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
