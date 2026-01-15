import { IAppEmote } from "@/api/native-emote";
import { getCursorPos, moveCursorTo, moveCursorToEnd } from "@/util/rte";
import { useEffect, useMemo, useRef, useState } from "react";
import { preload } from "react-dom";
import EmoteIcon from "../svg/emote-icon";
import Tooltip from "../util/tooltip";
import { TChatMessage } from "./chat-message";
import ReplyPopup from "./reply-popup";

interface IChatControlsProps {
    isReplying: boolean;
    replyingToMessage: TChatMessage|undefined;
    emoteList: IAppEmote[];
    getChatterColor: (username: string) => string;
    onSendMessage: (message: string) => Promise<boolean>;
    onShowUserPopup: (username: string|undefined, mouseX: number, mouseY: number) => void;
    onReplyClosed: () => void;
}

export default function ChatControls({
    isReplying,
    replyingToMessage,
    emoteList,
    getChatterColor,
    onSendMessage,
    onShowUserPopup,
    onReplyClosed,
}: IChatControlsProps) {
    const [shouldShowEmotePopup, setShouldShowEmotePopup] = useState<boolean>(false);

    const messageInputRef = useRef<HTMLDivElement>(null);

    const handleEmoteButtonClick = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation();
        setShouldShowEmotePopup(cur => !cur);
    }

    const handleMessageInputKeydown = (e: React.KeyboardEvent<HTMLDivElement>) => {
        if(e.key === "Enter") {
            e.preventDefault();
            handleSendMessage();
        }
    }

    const handleSendMessage = async () => {
        if(!messageInputRef.current) return;

        const message = inputNodesToText(messageInputRef.current.childNodes);

        const res = await onSendMessage(message);

        if(res) messageInputRef.current.innerHTML = '';
    }

    const inputNodesToText = (nodes: NodeListOf<ChildNode>) => {
        if(!nodes) return "";

        let text = "";
        for(const node of nodes) {
            switch(node.nodeType) {
            case Node.TEXT_NODE:
                const prevSibling = node.previousSibling
                const textContent = node.nodeValue;
                if(prevSibling && prevSibling.nodeType === Node.ELEMENT_NODE
                  && textContent?.charAt(0) !== ' ') text += ' ';
                text += textContent;
                break;
            case Node.ELEMENT_NODE:
                const prevChar = text.at(-1);
                if(prevChar && prevChar !== ' ') {
                    text += ' ';
                }

                text += (node as HTMLImageElement).alt;
                break;
            }
        }

        return text;
    }

    async function processTextNode(node: ChildNode) {
        const text = node.textContent;
        if(!text) return;

        let curTextNodeVal = "";
        let replaceNodes: (Node|string)[] = [];
        const potentialEmotes = text.split(' ');
        for(const potentialEmote of potentialEmotes) {
            if(potentialEmote.length === 0) {
                curTextNodeVal += " ";
                continue;
            }

            const matchedEmote = emoteList.find((e) => e.name === potentialEmote)
            if(!matchedEmote) {
                curTextNodeVal += potentialEmote;
                curTextNodeVal += " ";
                continue;
            }

            if(curTextNodeVal.length > 0 && curTextNodeVal.at(-1) === " ") {
                replaceNodes.push(curTextNodeVal.slice(0, -1));
                curTextNodeVal = "";
            }

            const imgNode = new Image()
            imgNode.srcset = matchedEmote.darkSrcSet;
            imgNode.classList.add("inline");
            imgNode.alt = matchedEmote.name;
            replaceNodes.push(imgNode);
        }

        if(curTextNodeVal.length > 0) {
            if(curTextNodeVal.at(-1) === " ") {
                replaceNodes.push(curTextNodeVal.slice(0, -1));
            } else {
                replaceNodes.push(curTextNodeVal);
            }
        }

        node.replaceWith(...replaceNodes);
    }

    async function processMessageInput(messageNodes: NodeListOf<ChildNode>) {
        const promiseList: Promise<void>[] = []
        for(const node of messageNodes) {
            switch(node.nodeType) {
            case Node.TEXT_NODE:
                promiseList.push(processTextNode(node));
            }
        }

        await Promise.all(promiseList);
        console.log('pmi done');
    }

    async function handleMessageInput(e: React.InputEvent<HTMLDivElement>) {
        const curorPos = getCursorPos(e.currentTarget);
        console.log(curorPos);
        processMessageInput(e.currentTarget.childNodes);
        if(curorPos !== -1) moveCursorTo(e.currentTarget, curorPos + 1);
    }

    function handleEmoteSelect(emote: IAppEmote) {
        if(!messageInputRef.current) return;

        const srcSet = emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet;

        const concat = `<img class="inline" srcset="${srcSet}" alt="${emote.name}"/> `;

        const lastChar = messageInputRef.current.innerHTML.at(-1);
        console.log(messageInputRef.current.innerHTML);
        console.log('last char', lastChar);
        if(lastChar && lastChar === ' ') {
            console.log('last char is space, no need to add space')
            messageInputRef.current.innerHTML += `${concat}`;
        } else {
            console.log('last char is not space, adding space')
            messageInputRef.current.innerHTML += ` ${concat}`;
        }

        messageInputRef.current.focus();
        moveCursorToEnd(messageInputRef.current);
    }

    const filteredEmoteList = useMemo<React.ReactNode>(() => {
        const idHash: Record<string, string> = {};

        const filteredArr: IAppEmote[] = [];

        for(let i = 0; i < emoteList.length; i++) {
            const emote = emoteList[i];
            if(!(emote.id in idHash)) {
                idHash[emote.id] = emote.id;
                filteredArr.push(emote);
            }
        }

        return filteredArr.map(emote => {
                const srcSet = emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet;
                for(const image in srcSet.split(', ')) {
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
    }, [emoteList]);



    useEffect(() => {
        function handleOutsideClick(e: MouseEvent) {
            if(e.target instanceof Node) {
                if(messageInputRef.current && !messageInputRef.current.contains(e.target)) {
                    setShouldShowEmotePopup(false);
                }
            }
        }

        document.addEventListener('click', handleOutsideClick);

        return () => {
            document.removeEventListener('click', handleOutsideClick);
        }
    }, [])


    return (
    <div className='flex flex-col p-1 basis-[max-content] bg-bg-9 relative'>
        {shouldShowEmotePopup &&
        <div
        className='w-[calc(100%-30px)] h-75 border border-outline-1 rounded-xs m-3.5 absolute left-0 bottom-full bg-bg-2/80 backdrop-blur-xs g-1 p-1 z-600 flex flex-wrap justify-between items-center scroller-y [&>span>div]:w-10 [&>span>div]:h-10 [&>span>div]:p-0.5 [&>span>div]:rounded-xs [&>span>div]:opacity-90 [&>span>div]:hover:bg-bg-5'
        onClick={(e: React.MouseEvent<HTMLDivElement>) => e.stopPropagation()}
        >
            {filteredEmoteList}
        </div>
        }
        {isReplying && replyingToMessage &&
        <div>
            <ReplyPopup
                onCloseClicked={onReplyClosed}
                message={replyingToMessage}
                getChatterColor={getChatterColor}
                showUserPopup={onShowUserPopup}
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
    )
}
