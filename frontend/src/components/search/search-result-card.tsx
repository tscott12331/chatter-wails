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
        <div onClick={onClick} className="bg-chatter-surface/70 w-[102%] flex items-center gap-3 p-3 drop-shadow-sm drop-shadow-chatter-border/50 scale-98 hover:scale-100 hover:bg-chatter-surface-elevated transition-[scale] will-change-[scale] cursor-pointer">
            <div className="aspect-video w-40 h-22.5 relative">
                <img src={thumbnailSrc} alt="No thumbnail" />
                <p className="bg-chatter-bg/90 text-chatter-text-tertiary p-0.5 text-xs absolute top-1 right-1 backdrop-blur-sm rounded-sm">🔴{viewcount}</p>
            </div>
            <div className="flex flex-col gap-1 overflow-hidden *:text-ellipsis *:whitespace-nowrap *:overflow-hidden *:antialiased">
                <h3>{title}</h3>
                <p>{channel}</p>
                <p>{category}</p>
            </div>
        </div>
    )
}
