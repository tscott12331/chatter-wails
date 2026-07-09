import { StreamData } from "@wailsjs/chatter-wails/services/eventsub";

interface IViewcountProps {
    streamData: StreamData;
}

export default function StreamDataIsland({
    streamData,
}: IViewcountProps) {
    console.log(streamData);
    return (
        <div className="grid grid-cols-[0fr_1fr] hover:auto-rows-min hover:grid-cols-[1fr_1fr] hover:scale-115 max-w-100 origin-top-right transition-[scale] ease-out text-text-2 text-right drop-shadow-sm drop-shadow-outline-2 p-2 rounded-sm bg-bg-2/50 backdrop-blur-xs overflow-hidden group">
        {streamData.live
        ? 
            <>
            <span className="group-hover:inline hidden text-xs text-left self-center">{streamData.gameName}</span>
            <span className="col-start-2 self-center">
                <span className="text-xs align-middle">LIVE 🔴 </span>
                <span className="align-middle">{streamData.viewCount}</span>
            </span>
            <span className="group-hover:col-span-2 group-hover:inline hidden text-left text-sm self-center">{streamData.title}</span>
            </>
        :
            <span className="col-start-2 group-hover:col-span-2">offline</span>
        }
        </div>
    )
}
