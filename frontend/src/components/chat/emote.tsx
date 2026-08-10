import { AppEmote } from "@wailsjs/chatter-wails/shared/types";
import React from "react";
import { TPartialTooltipData } from "../util/tooltip";

interface IEmoteMemoProps {
    srcSet: string;
    name: string;
}

interface IEmoteProps {
    emote: AppEmote;

    tooltipOnPartial?: (data: TPartialTooltipData, id: string, element: HTMLElement) => void;
    tooltipOff?: (id: string) => void;
}

const EmoteMemo = React.memo(({
        srcSet,
        name,
    }: IEmoteMemoProps) => {

        return (
            <img
                className="row-1 col-1 max-h-9"
                loading="lazy"
                decoding="async"
                srcSet={srcSet}
                alt={name}
            />
        )
    });

export default function Emote({
    emote,
    tooltipOnPartial,
    tooltipOff,
}: IEmoteProps) {
    const imageSubDesc = `${emote.provider}${
                            emote.section.length > 0
                            ? ` · ${emote.section}`
                            : ''
                            }`
    // TODO: improve stacked emote tooltip
    return (
        <div 
            className="align-middle inline-grid place-items-center grid-cols-1 grid-rows-1"
            onMouseEnter={(e) => {
                tooltipOnPartial?.({
                    type: "image",
                    imageSrcSet: emote.darkSrcSet,
                    imageDesc: [emote.name, imageSubDesc],
                }, emote.id, e.currentTarget);
            }}
            onMouseLeave={() => tooltipOff?.(emote.id)}
        >
            <EmoteMemo
                srcSet={emote.darkSrcSet}
                name={emote.name}
            />
            {
                emote.emoteStack?.map(e =>
                    e &&
                    <EmoteMemo
                        srcSet={e.darkSrcSet}
                        name={e.name}
                    />
                )
            }
        </div>
    )
}
