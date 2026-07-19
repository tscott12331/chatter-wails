import { AppEmote } from "@wailsjs/chatter-wails/services/emote";
import { TooltipContext } from "@/contexts/tooltip-context";
import { useContext } from "react";

export interface IEmoteCarouselData { 
    emotes: AppEmote[]; 
    index: number 
}

interface IEmoteCarouselProps {
    data: IEmoteCarouselData
}

export default function EmoteCarousel({
    data
}: IEmoteCarouselProps) {
    let trio: AppEmote[] = [];
    let highlightIndex = 0;
    if(data.emotes.length < 3) {
        trio.push(data.emotes[data.index]);
    } else {
        trio.push(data.emotes.at(data.index-1)!, data.emotes.at(data.index)!, data.emotes.at((data.index+1) % data.emotes.length)!);
        highlightIndex = 1;
    }

    const { tooltipOn, tooltipOff } = useContext(TooltipContext);

    return (
        <div className="relative max-w-50 flex gap-1 justify-around p-2 bg-bg-2/50 outline outline-outline-2/50 shadow-md shadow-outline-2/50 backdrop-blur-xs">
            <span className="text-xl text-shadow-sm text-shadow-outline-1 text-outline-2 content-center absolute h-full top-0 -left-2">{"<"}</span>
            {trio.map((e,i) => 
                <div 
                    className={`w-10 h-10 flex justify-center items-center ${i === highlightIndex && "bg-bg-5"}`}
                    onMouseEnter={me => {
                        const rect= me.currentTarget.getBoundingClientRect();
                        tooltipOn({
                            type: "image",
                            imageSrcSet: e.darkSrcSet,
                            imageDesc: e.name,
                            posX: rect.x + rect.width/2,
                            posY: rect.y,
                        }, e.id)
                    }}
                    onMouseLeave={() => tooltipOff(e.id)}
                    key={i}
                >
                    <img
                        className={`inline p-1`}
                        srcSet={e.darkSrcSet}
                        alt={e.name}
                    />
                </div>
            )}
            <span className="text-xl text-shadow-sm text-shadow-outline-1 text-outline-2 content-center absolute h-full top-0 -right-2">{">"}</span>
        </div>
    )
}
