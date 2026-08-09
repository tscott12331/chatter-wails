import { TChatroomEmotes } from "@/api/emote"
import { useContext, useEffect, useRef, useState } from "react";
import { AppEmote, AppEmoteSet } from "@wailsjs/chatter-wails/shared/types";
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
    const [tab, setTab] = useState<number>(0);
    const [providers, setProviders] = useState<string[]>([]);
    const [sets, setSets] = useState<AppEmoteSet[]>([]);

    const sectionsRef = useRef<Record<string, HTMLDivElement>>({});

    const { tooltipOnEmote, tooltipOffEmote, tooltipOnPartial, tooltipOff } = useContext(TooltipContext);

    useEffect(() => {
        setProviders([...emotes.keys()]);

    }, [emotes]);

    useEffect(() => {
        const sectionMap = emotes.get(providers[tab]);
        sectionsRef.current = {};
        if(sectionMap) {
            setSets([...sectionMap.values()]);
        }

    }, [emotes, providers, tab]);

    return (
        <div 
            className={`${open ? 'visible opacity-100 translate-y-0 scale-100' : 'invisible opacity-0 translate-y-1/5 scale-80'} flex flex-col w-[calc(100%-30px)] h-75 border border-chatter-border rounded-xs m-3.5 absolute left-0 bottom-full bg-chatter-surface-elevated/90 backdrop-blur-xs p-1 z-600 *:select-none transition-[opacity_visibility_translate] duration-150 ease-in-out`}
            ref={ref}
        >
            <div className="h-10 shrink-0 flex justify-around items-center gap-5 pb-1 ps-1 border-b border-b-chatter-border">
                {providers.map((t, i) => 
                    <button 
                        key={t}
                        className={`${i == tab ? 'bg-chatter-accent/15! border-chatter-accent!' : 'bg-chatter-surface-inset/50! border-chatter-border!'} hover:bg-chatter-surface/70! cursor-pointer text-center grow rounded-xl border px-2 transition-colors`}
                        onClick={() => setTab(i)}
                    >
                        {t}
                    </button>
                )}
            </div>
            <div className="flex gap-1 overflow-hidden">
                <div className="scroller-y grow">
                    {sets.map(set =>
                        set.Emotes && Object.keys(set.Emotes).length > 0 &&
                            <div
                            key={set.Id}
                            >
                                <h3
                                    className="bg-linear-150 from-chatter-bg/50 to-chatter-surface/30 p-2"
                                    ref={(el) => {
                                        if(!el) return;
                                        sectionsRef.current[set.Section] = el;
                                    }}
                                >{set.Section}</h3>
                                <div
                                    className='grow grid grid-cols-[repeat(auto-fill,40px)] auto-rows-min items-start justify-between gap-1'
                                >
                                    {
                                        Object.values(set.Emotes).filter(isDefined).map(emote =>
                                            <div
                                                className='flex justify-center items-center cursor-pointer w-10 h-10 p-0.5 rounded-xs opacity-90 hover:bg-chatter-surface transition-all duration-150'
                                                onClick={() => handleEmoteSelect(emote)}
                                                onMouseEnter={e => tooltipOnEmote(emote, e.currentTarget)}
                                                onMouseLeave={() => tooltipOffEmote(emote)}
                                                key={emote.id}
                                            >
                                                <img
                                                    className="[content-visibility:auto] [contain-intrinsic-size:40px_40px]"
                                                    srcSet={emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet}
                                                    loading="lazy"
                                                />
                                            </div>
                                        )
                                    }
                                </div>
                            </div>
                    )}
                </div>
                <div className="scroller-y w-11 border-l border-chatter-border flex flex-col gap-2 items-center py-2 ps-1">
                    {sets.map(set => {
                        if(!isDefined(set.Emotes)) return;
                        const first = Object.values(set.Emotes)[0];
                        if(!isDefined(first)) return;

                        return <a
                            key={set.Id}
                            onClick={() => sectionsRef.current[set.Section]?.scrollIntoView()}
                            className='flex justify-center items-center cursor-pointer w-10 h-10 p-0.5 rounded-xs opacity-90 hover:bg-chatter-surface transition-all duration-150'
                            onMouseEnter={e => tooltipOnPartial({
                                type: "text",
                                text: set.Section,
                            }, set.Id, e.currentTarget)}
                            onMouseLeave={() => tooltipOff(set.Id)}
                        >
                            <img
                                srcSet={first.darkSrcSet.length > 0 ? first.darkSrcSet : first.lightSrcSet}
                            />
                        </a>
                    })}
                </div>
            </div>
        </div>
    )
}
