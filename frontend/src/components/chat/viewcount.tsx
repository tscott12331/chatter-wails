export interface IViewcountData {
    live: boolean;
    viewCount: number;
}

interface IViewcountProps {
    viewcountData: IViewcountData;
}

export default function Viewcount({
    viewcountData,
}: IViewcountProps) {
    return (
        <div className="text-text-2 drop-shadow-sm drop-shadow-outline-2 p-2 rounded-sm bg-bg-2/50 backdrop-blur-xs">
        {viewcountData.live
        ? 
            <>
            <span className="text-xs align-middle">LIVE 🔴 </span>
            <span className="align-middle">{viewcountData.viewCount}</span>
            </>
        :
            <span>offline</span>
        }
        </div>
    )
}
