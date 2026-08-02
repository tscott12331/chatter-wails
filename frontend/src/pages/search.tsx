import SearchIcon from "@/components/svg/search-icon";
import { SearchTwitchStreams } from "@wailsjs/chatter-wails/services/native/emoteservice";
import { ApiGetStreamsRes } from "@wailsjs/chatter-wails/internal/api/nativeApi";
import { useRef, useState } from "react";
import { isDefined } from "@/util/assert";
import SearchResultCard from "@/components/search/search-result-card";

export default function SearchPage() {
    const inputChangeSearchDelay = 500; // ms

    const inputChangeTimeoutId = useRef<number|undefined>(undefined);
    const searchAbortController = useRef<AbortController|null>(null);
    const [searchData, setSearchData] = useState<ApiGetStreamsRes|null>(null);

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
        <div className="h-full p-5 flex flex-col items-center gap-5 scroller-y *:max-w-[min(672px,100%)]">
            <div className="w-full h-[max(15%,100px)] flex flex-col justify-evenly items-center gap-5">
                <h1 className="text-5xl sm:text-6xl md:text-7xl text-center basis-0">Browse Channels</h1>
                <div className="w-full p-1 flex bg-input-bg outline outline-input-outline">
                    <input 
                        className="h-5 grow border-0! focus:bg-input-bg!"
                        placeholder="search"
                        id="query"
                        onChange={handleSearchInput}
                    />
                    <button className="w-5 h-5 border-0! hover:scale-125 transition-transform cursor-pointer">
                        <SearchIcon className="fill-text-1" />
                    </button>
                </div>
            </div>
            <div className="flex flex-col gap-5 items-center shrink-0 *:max-w-full">
                {searchData?.data?.map(stream => 
                    <SearchResultCard 
                        thumbnailSrc={stream.thumbnail_url}
                        title={stream.title}
                        channel={stream.user_name}
                        category={stream.game_name}
                        viewcount={stream.viewer_count}
                        key={stream.id}
                    />
                )}
            </div>
        </div>
    )
}
