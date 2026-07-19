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

export default function Tooltip({
    data,
}: ITooltipProps) {
    return (
        <div
            className="absolute flex flex-col justify-center items-center z-10000 p-1 bg-bg-2/80 backdrop-blur-xs outline outline-outline-2/80 border border-outline-1 text-sm text-center wrap-break-word -translate-x-1/2 -translate-y-full rounded-[5%] shadow-[0_0_1px_1px] shadow-bg-9/50 italic"
            style={{
                left: data.posX,
                top: data.posY,
                
            }}
        >
        {data.type === "text"
        ? <p className="max-w-[min(100%,300px)]">{data.text}</p>
        : data.type === "image"
        ? <>
        <img srcSet={data.imageSrcSet} className="max-w-25" />
        <div className="w-full border border-bg-9/50 mt-1.5 mb-0.5"></div>
        <p className="max-w-full">{data.imageDesc}</p>
        </>
        : <></>
        }
        </div>
    )
}
