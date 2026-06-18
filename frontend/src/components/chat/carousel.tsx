import { IAppEmote } from "@/api/native-emote";
import Tooltip from "../util/tooltip";

export interface IEmoteCarouselData { 
    emotes: IAppEmote[]; 
    index: number 
}

interface IEmoteCarouselProps {
    data: IEmoteCarouselData
}

export default function EmoteCarousel({
    data
}: IEmoteCarouselProps) {
    let trio: IAppEmote[] = [];
    let highlightIndex = 0;
    if(data.emotes.length < 3) {
        trio.push(data.emotes[data.index]);
    } else {
        trio.push(data.emotes.at(data.index-1)!, data.emotes.at(data.index)!, data.emotes.at((data.index+1) % data.emotes.length)!);
        highlightIndex = 1;
    }

    return (
        <div className="max-w-50 h-10 flex gap-1 justify-around p-2 bg-bg-2/50 backdrop-blur-xs">
            {trio.map((e,i) => 
                <Tooltip
                    text={e.name}
                    hoverTime={0}
                    key={i}
                    >
                        <div className={`w-10 h-10 flex justify-center items-center ${i === highlightIndex && "bg-bg-5"}`}>
                            <img
                                className={`inline p-1`}
                                srcSet={e.darkSrcSet}
                                alt={e.name}
                            />
                        </div>
                </Tooltip>
            )}
        </div>
    )
}
