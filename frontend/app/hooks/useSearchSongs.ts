import { useEffect, useState } from 'react';
import { searchSongs } from '@/app/actions/search';
import { YouTubeSearchResults } from 'youtube-search';

export const useSearchSongs = (query: string) => {
  const [data, setData] = useState<YouTubeSearchResults[]>([]);
  const [isFetching, setIsFetching] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let cancelled = false;

    const fetchSongs = async () => {
      if (!query) return;

      setIsFetching(true);
      setError(null);

      try {
        const results = await searchSongs(query);
        if (!cancelled) {
          setData(results);
        }
      } catch (err: any) {
        if (!cancelled) {
          setError(err);
        }
      } finally {
        if (!cancelled) {
          setIsFetching(false);
        }
      }
    };

    fetchSongs();

    return () => {
      cancelled = true;
    };
  }, [query]);

  return { data, isFetching, error };
};
