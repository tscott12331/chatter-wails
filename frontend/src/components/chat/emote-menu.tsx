import { TChatroomEmotes } from "@/api/emote"
import { useContext, useEffect, useState } from "react";
import { preload } from "react-dom";
import { AppEmote } from "@wailsjs/chatter-wails/shared/types";
import { TooltipContext } from "@/contexts/tooltip-context";
import { isDefined } from "@/util/assert";

interface IEmoteMenuProps {
    emotes: TChatroomEmotes;
    open: boolean;
    handleEmoteSelect: (emote: AppEmote) => void;
    ref?: React.Ref<HTMLDivElement>;
}

export default function EmoteMenu({
    emotes,
    open,
    handleEmoteSelect,
    ref,
}: IEmoteMenuProps) {
    const tabs = [...emotes.keys()];

    const [tab, setTab] = useState<number>(0);

    const { tooltipOnEmote, tooltipOffEmote } = useContext(TooltipContext);
    
    useEffect(() => {
        const providerMaps = emotes.values();
        const flatEmotes = Array.from(providerMaps.flatMap(sectionMap =>
            sectionMap.values().flatMap(set =>
                set.Emotes
                ? Object.values(set.Emotes).filter(isDefined)
                : []
            )
        ));

        for(const emote of flatEmotes) {
            const srcSet = emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet;
            for(const image of srcSet.split(', ')) {
                preload(image, { as: "image", imageSrcSet: srcSet });
            }
        }
    }, [emotes]);

    return (
        <div 
            className={`${!open && 'invisible opacity-0'} flex flex-col w-[calc(100%-30px)] h-75 border border-outline-1 rounded-xs m-3.5 absolute left-0 bottom-full bg-bg-2/80 backdrop-blur-xs p-1 z-600 *:select-none transition-all duration-75 ease-in-out transition-discrete`}
            ref={ref}
        >
            <div className="h-10 shrink-0 flex justify-around items-center gap-5 pb-1 ps-1 border-b border-b-outline-1">
                {tabs.map((t, i) => 
                    <button 
                        key={t}
                        className={`${i == tab ? 'bg-bg-3/50!' : 'bg-input-bg/50!'} hover:bg-input-bg-3/50! cursor-pointer text-center grow rounded-xl border border-outline-2 px-2 transition-colors`}
                        onClick={() => setTab(i)}
                    >
                        {t}
                    </button>
                 )}
            </div>
            <div
            className='grow grid grid-cols-[repeat(auto-fill,40px)] auto-rows-min items-start justify-between gap-1 scroller-y'
            >
            {tabs[tab] && emotes.has(tabs[tab]) && [...emotes.get(tabs[tab])!.values()].flatMap(set =>
                set.Emotes
                ? Object.values(set.Emotes).filter(isDefined).map(emote =>
                 <div
                     className='flex justify-center items-center cursor-pointer w-10 h-10 p-0.5 rounded-xs opacity-90 hover:bg-bg-5 transition-all duration-150'
                     onClick={() => handleEmoteSelect(emote)}
                     onMouseEnter={e => tooltipOnEmote(emote, e.currentTarget)}
                     onMouseLeave={() => tooltipOffEmote(emote)}
                     key={emote.id}
                 >
                     <img
                         srcSet={emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet}
                     />
                 </div>
                )
                : []
            )}
            </div>
        </div>
    )
}
