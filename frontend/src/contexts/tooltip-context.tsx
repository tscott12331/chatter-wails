import Tooltip, { TPartialTooltipData, TTooltipData } from "@/components/util/tooltip";
import { AppEmote } from "@wailsjs/chatter-wails/shared/types";
import { createContext, useCallback, useMemo, useState } from "react";

interface ITooltipContext {
    currentTooltip?: TTooltipData;
    tooltipOn: (data: TTooltipData, id: string) => void;
    tooltipOff: (id: string) => void;
    tooltipOnEmote: (emote: AppEmote, element: HTMLElement) => void;
    tooltipOffEmote: (emote: AppEmote) => void;
    tooltipOnPartial: (data: TPartialTooltipData, id: string, element: HTMLElement) => void;
}

export const TooltipContext = createContext<ITooltipContext>({
    tooltipOn: (_data, _id) => {},
    tooltipOff: (_id) => {},
    tooltipOnEmote: () => {},
    tooltipOffEmote: () => {},
    tooltipOnPartial: () => {},
});

type TUniqueTooltip = TTooltipData & { id: string }

export function TooltipContextProvider({
    children
}: { children: React.ReactNode }) {
    const [currentTooltip, setCurrentTooltip] = useState<TUniqueTooltip|undefined>();

    const tooltipOn = useCallback((data: TTooltipData, id: string) => {
        setCurrentTooltip({...data, id});
    }, []);

    const tooltipOff = useCallback((id: string) => {
        setCurrentTooltip(cur => {
            if(cur?.id !== id) return cur;

            return undefined;
        })
    }, []);

    const tooltipOffEmote = useCallback((emote: AppEmote) => {
        tooltipOff(emote.id);
    }, [])

    const tooltipOnEmote = useCallback((emote: AppEmote, element: HTMLElement) => {
        const rect = element.getBoundingClientRect();
        const imageSubDesc = `${emote.provider}${
                                emote.section.length > 0
                                ? ` · ${emote.section}`
                                : ''
                                }`
        tooltipOn({
            type: "image",
            imageSrcSet: emote.darkSrcSet.length > 0 ? emote.darkSrcSet : emote.lightSrcSet,
            imageDesc: [emote.name, imageSubDesc],
            posX: rect.x + rect.width/2,
            posY: rect.y,
        }, emote.id);
    }, []);

    const tooltipOnPartial = useCallback((data: TPartialTooltipData, id: string, element: HTMLElement) => {
        const rect = element.getBoundingClientRect();
        tooltipOn({
            ...data,
            posX: rect.x + rect.width/2,
            posY: rect.y,
        }, id);
    }, []);

    

    const contextValue = useMemo<ITooltipContext>(() => ({
        currentTooltip,
        tooltipOn,
        tooltipOff,
        tooltipOnEmote,
        tooltipOffEmote,
        tooltipOnPartial,
    }), [currentTooltip]);

    return (
        <TooltipContext.Provider value={contextValue}>
            {children}
            {currentTooltip && <Tooltip data={currentTooltip} />}
        </TooltipContext.Provider>
    )
}
