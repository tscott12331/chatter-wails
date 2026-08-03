import SearchIcon from "@/components/svg/search-icon";
import { SearchTwitchStreams } from "@wailsjs/chatter-wails/services/native/emoteservice";
import { ApiGetStreamsRes } from "@wailsjs/chatter-wails/internal/api/nativeApi";
import { useContext, useEffect, useRef, useState } from "react";
import { isDefined } from "@/util/assert";
import SearchResultCard from "@/components/search/search-result-card";
import { createTab, TabContext } from "@/contexts/tab-context";

export default function SearchPage() {
    const inputChangeSearchDelay = 500; // ms

    const inputChangeTimeoutId = useRef<number|undefined>(undefined);
    const searchAbortController = useRef<AbortController|null>(null);
    const [searchData, setSearchData] = useState<ApiGetStreamsRes|null>(null);

    const { addTab } = useContext(TabContext);

    const search = (query: string, abortSignal?: AbortSignal) => {
        const promise = SearchTwitchStreams(query)
            .then((d) => {
                setSearchData(d);
            });
        if(abortSignal) {
            promise.cancelOn(abortSignal);
        }
    }

    const cancelSearch = () => {
        if(isDefined(searchAbortController.current)) {
            searchAbortController.current.abort();
        }
        clearTimeout(inputChangeTimeoutId.current);
    }

    const immediateSearch = (query: string) => {
        cancelSearch();
        searchAbortController.current = new AbortController();
        search(query, searchAbortController.current.signal)
    }

    const searchAfterDelay = (query: string) => {
        cancelSearch();

        searchAbortController.current = new AbortController();

        inputChangeTimeoutId.current = setTimeout(() => {
            search(query, searchAbortController.current?.signal);
        }, inputChangeSearchDelay);
    }

    const handleSearchInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        searchAfterDelay(e.target.value);
    }

    const handleSubmit = (formData: FormData) => {
        const query = formData.get('query') as string;
        if(!query) return;

        immediateSearch(query);
    }

    useEffect(() => {
        immediateSearch("");
    }, []);

    return (
        <div className="h-full p-5 flex flex-col items-center gap-5 scroller-y *:w-[min(672px,100%)] relative isolate scrollbar-gutter-both">
            <div className="w-full shrink-0 h-37.5 sticky -top-28 xs:-top-24 sm:-top-25 md:-top-26 flex flex-col justify-evenly items-center gap-5 z-1">
                <h1 className="font-extrabold text-5xl sm:text-6xl md:text-7xl text-center tracking-tight">Browse Channels</h1>
                <form 
                    className="w-full p-1 flex items-center bg-input-bg/70 backdrop-blur-sm outline outline-input-outline rounded-md"
                    action={handleSubmit}
                >
                    <input 
                        className="h-5 grow border-0! bg-transparent!"
                        placeholder="search"
                        name="query"
                        onChange={handleSearchInput}
                    />
                    <button 
                        className="w-6.25 h-6.25 bg-transparent! hover:bg-input-bg-3! border-0! scale-80 hover:scale-100 will-change-[scale] transition-[transform_background] cursor-pointer"
                        type="submit"
                    >
                        <SearchIcon className="fill-text-1" />
                    </button>
                </form>
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
