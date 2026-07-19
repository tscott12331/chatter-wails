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
            className="absolute flex flex-col justify-center items-center z-10000 p-0.5 bg-bg-2/80 backdrop-blur-xs outline outline-outline-2/80 border border-outline-1 text-nowrap text-sm -translate-x-1/2 -translate-y-full"
            style={{
                left: data.posX,
                top: data.posY,
                
            }}
        >
        {data.type === "text"
        ? data.text
        : data.type === "image"
        ? <>
        <img srcSet={data.imageSrcSet} className="max-w-20"/>
        <p>{data.imageDesc}</p>
        </>
        : <></>
        }
        </div>
    )
}
