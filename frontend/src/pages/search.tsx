import SearchIcon from "@/components/svg/search-icon";

export default function SearchPage() {
    return (
        <div className="flex flex-col h-full scroller-y">
            <div className="flex flex-col items-center gap-5 py-[15%] grow-2">
                <h1 className="text-5xl text-center basis-0">Browse Channels</h1>
                <div className="w-4/5 p-1 flex bg-input-bg outline outline-input-outline">
                    <input className="h-5 grow border-0! focus:bg-input-bg!" placeholder="search" id="query" />
                    <button className="w-5 h-5 border-0! hover:scale-125 transition-transform cursor-pointer">
                        <SearchIcon className="fill-text-1" />
                    </button>
                </div>
            </div>
            <div className="flex flex-col gap-1 grow-7 shrink-0">
            </div>
        </div>
    )
}
