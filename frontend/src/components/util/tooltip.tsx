interface ITextTooltip {
    type: "text";
    text: string | string[];
}
interface IImageTooltip {
    type: "image";
    imageSrcSet: string;
    imageDesc: string | string[];
}

export type TPartialTooltipData = ITextTooltip|IImageTooltip

export type TTooltipData = TPartialTooltipData & {
    posX: number;
    posY: number;
};

interface ITooltipProps {
    data: TTooltipData;
}

const TOOLTIP_MAX_IMG_W = 100;
const TOOLTIP_MAX_TXT_W = 500;
const TOOLTIP_UPWARD_OFF = 5;

export default function Tooltip({
    data,
}: ITooltipProps) {
    function determinePosition(data: TTooltipData) {
        const windowWidth = window.innerWidth;

        const posStyles = {};
        const padding = (data.type == "text" 
                        ? Math.min(TOOLTIP_MAX_TXT_W, data.text.length*14)
                        : TOOLTIP_MAX_IMG_W) / 2;

        if(data.posX + padding > windowWidth) {
            posStyles['right'] = 0;
        } else if(data.posX < padding) {
            posStyles['left'] = padding;
        } else {
            posStyles['left'] = data.posX;
        }

        posStyles['top'] = data.posY - TOOLTIP_UPWARD_OFF;

        return posStyles;
    }

    const displayText = (text: string|string[]) => {
        return Array.isArray(text)
            ? text.map((t, i) => <span key={t+i} className={`block ${i !== 0 ? 'text-xs text-chatter-text-secondary' : 'mb-1'}`}>{t}</span>)
            : <span>{text}</span>
    }

    return (
        <div
            className="absolute flex flex-col justify-center items-center z-10000 p-1 bg-chatter-surface-elevated/90 backdrop-blur-xs outline outline-chatter-border-strong/80 border border-chatter-border text-sm text-center text-chatter-text-primary wrap-break-word -translate-x-1/2 -translate-y-full rounded-sm shadow-[0_0_1px_1px] shadow-chatter-border-strong/50 italic"
            style={{
                ...determinePosition(data),
            }}
        >
        {data.type === "text"
        ? <p 
        className="text-left"
            style={{
                maxWidth: `min(100%,${TOOLTIP_MAX_TXT_W}px)`
            }}
        >{displayText(data.text)}</p>
        : data.type === "image"
        ? <>
        <img 
            className="h-10"
            srcSet={data.imageSrcSet}
            style={{
                maxWidth: `${TOOLTIP_MAX_IMG_W}px`,
            }}
        />
        <div className="w-full border border-chatter-border-strong/50 mt-1.5 mb-0.5"></div>
        <p 
            style={{
                maxWidth: `${TOOLTIP_MAX_IMG_W}px`,
            }}
        >{displayText(data.imageDesc)}</p>
        </>
        : <></>
        }
        </div>
    )
}
