interface ISearchResultCardProps {
    thumbnailSrc: string;
    title: string;
    channel: string;
    category: string;
    viewcount: number

    onClick?: () => void;
}

export default function SearchResultCard({
    thumbnailSrc,
    title,
    channel,
    category,
    viewcount,

    onClick,
}: ISearchResultCardProps) {
    return (
        <div onClick={onClick} className="bg-bg-09/70 w-full flex items-center gap-3 p-3 drop-shadow-sm drop-shadow-outline-1/50 hover:scale-102 hover:bg-bg-09 transition-[scale] will-change-[scale] cursor-pointer">
            <div className="aspect-video w-40 h-22.5 relative">
                <img src={thumbnailSrc} alt="No thumbnail" />
                <p className="mix-blend-hard-light bg-bg-09/70 text-text-3 p-0.5 text-xs absolute top-1 right-1 drop-shadow-sm drop-shadow-outline-2 backdrop-blur-xs">🔴{viewcount}</p>
            </div>
            <div className="flex flex-col gap-1 overflow-hidden *:text-ellipsis *:whitespace-nowrap *:overflow-hidden *:subpixel-antialiased">
                <h3>{title}</h3>
                <p>{channel}</p>
                <p>{category}</p>
            </div>
        </div>
    )
}
