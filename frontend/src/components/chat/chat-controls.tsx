import { TChatroomEmotes } from "@/api/emote";
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
    emotes: TChatroomEmotes;
    getChatterColor: (username: string) => string;
    onSendMessage: (message: string) => Promise<boolean>;
    onShowUserPopup: (username: string|undefined, mouseX: number, mouseY: number) => void;
    onReplyClosed: () => void;
}

export default function ChatControls({
    isReplying,
    replyingToMessage,
    emotes,
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

    const handleMessageInputKeydown = async (e: React.KeyboardEvent<HTMLDivElement>) => {
        switch(e.key) {
            case "Enter":
                e.preventDefault();
                handleSendMessage();
                break;
            case "Tab":
                e.preventDefault();
                const curorPos = getCursorPos(e.currentTarget);
                const target = e.currentTarget;
                const cursorOffset = await processMessageInput(e.currentTarget.childNodes);
                if(curorPos !== -1) moveCursorTo(target, curorPos + cursorOffset);
                break;
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

    // TODO: change emoteMap structure to trie to improve lookup on completion
    function getEmoteMatches(potentialEmote: string): IAppEmote[] {
        const matches: IAppEmote[] = [];
        for(const emoteMap of Object.values(emotes)) {
            for(const [emoteName, emote] of emoteMap) {
                if(emoteName.startsWith(potentialEmote)) {
                    matches.push(emote);
                }
            }
        }

        return matches;
    }

    function getWordFromCursorPos(text: string, offset: number) {
        let wordStart = offset - 1;
        let wordEnd = wordStart;
        if(wordStart < 0 || wordStart >= text.length) {
            return {
                word: " ",
                wordStart,
                wordEnd,
            };
        }

        while(wordStart >= 0 && /\S/.test(text[wordStart])) {
            wordStart -= 1;
        }
        wordStart += 1

        while(wordEnd < text.length && /\S/.test(text[wordEnd])) {
            wordEnd += 1;
        }

        return {
            word: text.slice(wordStart, wordEnd),
            wordStart,
            wordEnd,
        }
    }

    interface CompletionOpts {
        node: ChildNode;
        text: string;
        matches: IAppEmote[];
        cursorOffset: number;
        wordStart: number;
        wordEnd: number;
    }
    function completeNodeWord(opts: CompletionOpts) {
        // TODO: add carousel for completion options instead of just completing the first option
        const completion = opts.matches[0];
        const newText = `${opts.text.slice(0, opts.wordStart)}${completion.name}${opts.text.slice(opts.wordEnd)}`;

        opts.node.replaceWith(newText);
        
        return (opts.wordStart + completion.name.length) - opts.cursorOffset;
    }

    async function processTextNode(node: ChildNode, cursorOffset: number) {
        const text = node.textContent;
        if(!text) return 0;

        const { word: potentialEmote, wordStart, wordEnd }  = getWordFromCursorPos(text, cursorOffset);
        if(potentialEmote.length === 0) return 0;

        const matches = getEmoteMatches(potentialEmote);
        if(matches.length === 0) return 0;

        return completeNodeWord({ node, text, matches, cursorOffset, wordStart, wordEnd });
    }

    async function processMessageInput(messageNodes: NodeListOf<ChildNode>) {
        const selection = window.getSelection();
        if(!selection) return 0;
        const range = selection.getRangeAt(0);
        const offset = range.startOffset;

        for(const node of messageNodes) {
            if(node.contains(selection.anchorNode)) {
                const co = await processTextNode(node, offset);
                return co;
            }
        }

        return 0;
    }

    function handleEmoteSelect(emote: IAppEmote) {
        if(!messageInputRef.current) return;

        const childNodes = messageInputRef.current.childNodes;
        const lastNode = childNodes.item(childNodes.length - 1);

        // remove breaks (causes strange behavior)
        if(lastNode && lastNode.nodeName === "BR") {
            lastNode.remove();
        }

        const lastChar = messageInputRef.current.innerHTML.at(-1);
        if(lastChar && lastChar === ' ') {
            messageInputRef.current.innerHTML += `${emote.name} `;
        } else {
            messageInputRef.current.innerHTML += ` ${emote.name} `;
        }

        messageInputRef.current.focus();
        moveCursorToEnd(messageInputRef.current);
    }

    const filteredEmoteList = useMemo<React.ReactNode>(() => {
        const idHash: Record<string, string> = {};

        const filteredArr: IAppEmote[] = [];

        for(const emoteList of Object.values(emotes)) {
            for(const emote of emoteList.values()) {
                if(!(emote.id in idHash)) {
                    idHash[emote.id] = emote.id;
                    filteredArr.push(emote);
                }
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
    }, [emotes]);



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
            className='w-full h-[calc(100%-46px)] overflow-auto'
            contentEditable="true"
            onKeyDown={handleMessageInputKeydown}
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
