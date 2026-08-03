import SearchIcon from "@/components/svg/search-icon";
import { SearchTwitchStreams } from "@wailsjs/chatter-wails/services/native/emoteservice";
import { ApiGetStreamsRes } from "@wailsjs/chatter-wails/internal/api/nativeApi";
import { useContext, useRef, useState } from "react";
import { isDefined } from "@/util/assert";
import SearchResultCard from "@/components/search/search-result-card";
import { createTab, TabContext } from "@/contexts/tab-context";

export default function SearchPage() {
    const inputChangeSearchDelay = 500; // ms

    const inputChangeTimeoutId = useRef<number|undefined>(undefined);
    const searchAbortController = useRef<AbortController|null>(null);
    const [searchData, setSearchData] = useState<ApiGetStreamsRes|null>(null);

    const { addTab } = useContext(TabContext);

    const handleSearchInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        if(isDefined(searchAbortController.current)) {
            searchAbortController.current.abort();
        }
        clearTimeout(inputChangeTimeoutId.current);

        searchAbortController.current = new AbortController();

        inputChangeTimeoutId.current = setTimeout(() => {
            SearchTwitchStreams(e.target.value)
                .then((d) => {
                    console.log(d);
                    setSearchData(d);
                })
                .cancelOn(searchAbortController.current!.signal);
        }, inputChangeSearchDelay)
    }

    return (
        <div className="h-full p-5 flex flex-col items-center gap-5 scroller-y *:w-[min(672px,100%)] relative isolate">
            <div className="w-full h-[max(15%,100px)] flex flex-col justify-evenly items-center gap-5 sticky top-[min(-7.5%,-50px)] z-1">
                <h1 className="font-extrabold text-5xl sm:text-6xl md:text-7xl text-center tracking-tight">Browse Channels</h1>
                <div className="w-full p-1 flex items-center bg-input-bg/70 backdrop-blur-sm outline outline-input-outline rounded-md">
                    <input 
                        className="h-5 grow border-0! bg-transparent!"
                        placeholder="search"
                        id="query"
                        onChange={handleSearchInput}
                    />
                    <button className="w-6.25 h-6.25 bg-transparent! hover:bg-input-bg-3! border-0! scale-80 hover:scale-100 will-change-[scale] transition-[transform_background] cursor-pointer">
                        <SearchIcon className="fill-text-1" />
                    </button>
                </div>
            </div>
            <div className="flex flex-col gap-5 items-center shrink-0 -z-1">
                {searchData?.data?.map(stream => 
                    <SearchResultCard 
                        thumbnailSrc={stream.thumbnail_url}
                        title={stream.title}
                        channel={stream.user_name}
                        category={stream.game_name}
                        viewcount={stream.viewer_count}
                        key={stream.id}

                        onClick={() => addTab(createTab(stream.user_name))}
                    />
                )}
            </div>
        </div>
    )
}
