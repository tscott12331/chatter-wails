import { TTooltipData } from "@/components/util/tooltip";
import { createContext, useState } from "react";

interface ITooltipContext {
    currentTooltip?: TTooltipData;
    tooltipOn: (data: TTooltipData, id: string) => void;
    tooltipOff: (id: string) => void;
}

export const TooltipContext = createContext<ITooltipContext>({
    tooltipOn: (_data, _id) => {},
    tooltipOff: (_id) => {},
});

type TUniqueTooltip = TTooltipData & { id: string }

export function TooltipContextProvider({
    children
}: { children: React.ReactNode }) {
    const [currentTooltip, setCurrentTooltip] = useState<TUniqueTooltip|undefined>();

    const tooltipOn = (data: TTooltipData, id: string) => {
        setCurrentTooltip({...data, id});
    }

    const tooltipOff = (id: string) => {
        setCurrentTooltip(cur => {
            if(cur?.id !== id) return cur;

            return undefined;
        })
    }

    return (
        <TooltipContext.Provider value={{
            currentTooltip,
            tooltipOn,
            tooltipOff,
        }}>
            {children}
        </TooltipContext.Provider>
    )
}
