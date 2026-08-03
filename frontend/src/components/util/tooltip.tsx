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

const tooltipMaxImageWidth = 100;
const tooltipMaxTextWidth = 500;

export default function Tooltip({
    data,
}: ITooltipProps) {
    function determinePosition(data: TTooltipData) {
        const windowWidth = window.innerWidth;

        const posStyles = {};
        const padding = (data.type == "text" 
                        ? Math.min(tooltipMaxTextWidth, data.text.length*14)
                        : tooltipMaxImageWidth) / 2;

        if(data.posX + padding > windowWidth) {
            posStyles['right'] = 0;
        } else if(data.posX < padding) {
            posStyles['left'] = padding;
        } else {
            posStyles['left'] = data.posX;
        }

        posStyles['top'] = data.posY

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
                maxWidth: `min(100%,${tooltipMaxTextWidth}px)`
            }}
        >{displayText(data.text)}</p>
        : data.type === "image"
        ? <>
        <img 
            srcSet={data.imageSrcSet}
            style={{
                maxWidth: `${tooltipMaxImageWidth}px`,
            }}
        />
        <div className="w-full border border-chatter-border-strong/50 mt-1.5 mb-0.5"></div>
        <p 
            style={{
                maxWidth: `${tooltipMaxImageWidth}px`,
            }}
        >{displayText(data.imageDesc)}</p>
        </>
        : <></>
        }
        </div>
    )
}
