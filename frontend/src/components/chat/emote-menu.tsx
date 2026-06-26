import { TChatroomEmotes } from "@/api/emote"
import { IAppEmote } from "@/api/native-emote";
import { useEffect, useState } from "react";
import { preload } from "react-dom";
import Tooltip from "../util/tooltip";

interface IEmoteMenuProps {
    emotes: TChatroomEmotes;
    open: boolean;
    handleEmoteSelect: (emote: IAppEmote) => void;
    ref?: React.Ref<HTMLDivElement>;
}

export default function EmoteMenu({
    emotes,
    open,
    handleEmoteSelect,
    ref,
}: IEmoteMenuProps) {
    const tabs = Object.keys(emotes);
    const [tab, setTab] = useState<number>(0);
    
    useEffect(() => {
        for(const emoteMap of Object.values(emotes)) {
            for(const emote of emoteMap.values()) {
                const srcSet = emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet;
                for(const image of srcSet.split(', ')) {
                    preload(image, { as: "image", imageSrcSet: srcSet });
                }
            }
        }
    }, [emotes]);
    return (
        <div 
            className={`${!open && 'invisible'} flex flex-col w-[calc(100%-30px)] h-75 border border-outline-1 rounded-xs m-3.5 absolute left-0 bottom-full bg-bg-2/80 backdrop-blur-xs p-1 z-600`}
            ref={ref}
        >
            <div className="h-10 shrink-0 flex justify-around items-center gap-5 pb-1 ps-1 border-b border-b-outline-1">
                {tabs.map((t, i) => 
                    <button 
                        key={t}
                        className={`${i == tab ? 'bg-bg-3/50!' : 'bg-input-bg/50!'} cursor-pointer text-center grow rounded-xl border border-outline-2 px-2`}
                        onClick={() => setTab(i)}
                    >
                        {t}
                    </button>
                 )}
            </div>
            <div
            className='grow grid grid-cols-[repeat(auto-fill,40px)] auto-rows-min items-start justify-between gap-1 scroller-y'
            >
            {tabs[tab] && [...emotes[tabs[tab]].values()].map(emote =>
                <Tooltip
                 text={emote.name}
                 hoverTime={0}
                 key={emote.id}
                 >
                     <div
                         className='flex justify-center items-center cursor-pointer w-10 h-10 p-0.5 rounded-xs opacity-90 hover:bg-bg-5'
                         onClick={() => handleEmoteSelect(emote)}
                     >
                         <img
                             srcSet={emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet}
                         />
                     </div>
                 </Tooltip>

            )}
            </div>
        </div>
    )
}
