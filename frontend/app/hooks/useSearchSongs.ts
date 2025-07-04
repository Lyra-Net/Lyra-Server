import { useQuery, type UseQueryOptions } from "@tanstack/react-query";
import { searchSongs } from "@/app/actions/search";
import { YouTubeSearchResults } from "youtube-search";

export const useSearchSongs = (query: string) => {
    return useQuery<
        YouTubeSearchResults[],
        Error,
        YouTubeSearchResults[],
        [string, string]
    >({
        queryKey: ["search", query],
        queryFn: () => searchSongs(query),
        enabled: !!query,
        staleTime: Infinity,
        cacheTime: 1000 * 60 * 60 * 24,
        refetchOnWindowFocus: false,
        refetchOnMount: false,
        refetchOnReconnect: false,
    } as UseQueryOptions<YouTubeSearchResults[], Error, YouTubeSearchResults[], [string, string]>);
};
