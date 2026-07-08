import { StreamData } from "@wailsjs/chatter-wails/services/eventsub";

interface IViewcountProps {
    streamData: StreamData;
}

export default function StreamDataIsland({
    streamData,
}: IViewcountProps) {
    return (
        <div className="text-text-2 drop-shadow-sm drop-shadow-outline-2 p-2 rounded-sm bg-bg-2/50 backdrop-blur-xs">
        {streamData.live
        ? 
            <>
            <span className="text-xs align-middle">LIVE 🔴 </span>
            <span className="align-middle">{streamData.viewCount}</span>
            </>
        :
            <span>offline</span>
        }
        </div>
    )
}
