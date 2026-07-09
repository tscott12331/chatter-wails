import { StreamData } from "@wailsjs/chatter-wails/services/eventsub";

interface IViewcountProps {
    streamData: StreamData;
}

export default function StreamDataIsland({
    streamData,
}: IViewcountProps) {
    return (
        <div className="hover:scale-115 origin-top-right transition-all ease-out text-text-2 text-right drop-shadow-sm drop-shadow-outline-2 p-2 rounded-sm bg-bg-2/50 backdrop-blur-xs overflow-hidden group">
        {streamData.live
        ? 
            <>
            <span className="text-xs align-middle">LIVE 🔴 </span>
            <span className="align-middle">{streamData.viewCount}</span>
            <span className="group-hover:block hidden text-left">{streamData.title}</span>
            </>
        :
            <span>offline</span>
        }
        </div>
    )
}
