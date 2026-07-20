interface ITextTooltip {
    type: "text";
    text: string
}
interface IImageTooltip {
    type: "image";
    imageSrcSet: string;
    imageDesc: string
}

export type TTooltipData = (ITextTooltip | IImageTooltip) & {
    posX: number;
    posY: number;
};

interface ITooltipProps {
    data: TTooltipData;
}

const tooltipMaxImageWidth = 50;
const tooltipMaxTextWidth = 300;

export default function Tooltip({
    data,
}: ITooltipProps) {
    function determinePosition(posX: number, posY: number, type: TTooltipData["type"]) {
        const windowWidth = window.innerWidth;

        const posStyles = {};
        const padding = (type == "text" 
                        ? tooltipMaxTextWidth 
                        : tooltipMaxImageWidth) / 2;

        if(posX + padding > windowWidth) {
            posStyles['right'] = padding;
        } else if(posX < padding) {
            posStyles['left'] = padding;
        } else {
            posStyles['left'] = posX;
        }

        posStyles['top'] = posY

        return posStyles;
    }
    return (
        <div
            className="absolute flex flex-col justify-center items-center z-10000 p-1 bg-bg-2/80 backdrop-blur-xs outline outline-outline-2/80 border border-outline-1 text-sm text-center wrap-break-word -translate-x-1/2 -translate-y-full rounded-[5%] shadow-[0_0_1px_1px] shadow-bg-9/50 italic"
            style={{
                ...determinePosition(data.posX, data.posY, data.type),
            }}
        >
        {data.type === "text"
        ? <p 
            style={{
                maxWidth: `min(100%,${tooltipMaxImageWidth}px)`
            }}
        >{data.text}</p>
        : data.type === "image"
        ? <>
        <img 
            srcSet={data.imageSrcSet}
            style={{
                maxWidth: `${tooltipMaxImageWidth}px`,
            }}
        />
        <div className="w-full border border-bg-9/50 mt-1.5 mb-0.5"></div>
        <p 
            style={{
                maxWidth: `${tooltipMaxImageWidth}px`,
            }}
        >{data.imageDesc}</p>
        </>
        : <></>
        }
        </div>
    )
}
