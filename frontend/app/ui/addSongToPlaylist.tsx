'use client';

import { useEffect, useState } from 'react';
import api from '@/lib/api';
import { Playlist } from '@/declarations/playlists';
import { toast } from 'sonner';

type AddSongToPlaylistProps = {
  playlist: Playlist;
  setShowAddSong: (show: boolean) => void;
  refreshPlaylist: () => void;
};

export default function AddSongToPlaylist({ playlist, setShowAddSong, refreshPlaylist }: AddSongToPlaylistProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [typingTimer, setTypingTimer] = useState<NodeJS.Timeout | null>(null);

  // debounce search
  useEffect(() => {
    if (!searchQuery.trim()) {
      setSearchResults([]);
      return;
    }
    if (typingTimer) clearTimeout(typingTimer);

    const timer = setTimeout(async () => {
      try {
        const res = await api.get(`/search?type=song&query=${encodeURIComponent(searchQuery)}`);
        setSearchResults(res.data?.songs || []);
      } catch (err) {
        console.error("Search failed", err);
        toast.error("Search failed");
      }
    }, 200); // debounce 200ms

    setTypingTimer(timer);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  // add song
  const handleAddSong = async (songId: string) => {
    try {
      await api.post('/playlist/add-song', { playlist_id: playlist.playlist_id, song_id: songId });
      toast.success("Song added");
      refreshPlaylist?.();
    } catch (err) {
      console.error("Add song failed", err);
      toast.error("Add song failed");
    }
  };

  // remove song
  const handleRemoveSong = async (songId: string) => {
    try {
      await api.post('/playlist/remove-song', { playlist_id: playlist.playlist_id, song_id: songId });
      toast.success("Song removed");
      refreshPlaylist?.();
    } catch (err) {
      console.error("Remove song failed", err);
      toast.error("Remove song failed");
    }
  };

  const songIdsInPlaylist = new Set(playlist.songs.length ? playlist.songs.map((s: any) => s.id): []);

  return (
    <div className="absolute inset-0 bg-black/20 flex justify-center items-start z-20">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-2xl mt-10 p-4 relative">
        {/* Header */}
        <div className="flex justify-between items-center mb-3">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-100">
            Add or Remove Songs
          </h2>
          <button
            onClick={() => setShowAddSong(false)}
            className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            ✕
          </button>
        </div>

        {/* Search box */}
        <input
          type="text"
          placeholder="Search for a song..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-3 py-2 border rounded mb-4 dark:bg-gray-700 dark:border-gray-600"
        />

        {/* Search results */}
        <ul className="space-y-2 max-h-96 overflow-y-auto">
          {searchResults.length ? searchResults.map((song) => (
            <li
              key={song.song_id}
              className="flex justify-between items-center p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <span>{song.title}</span>
              {songIdsInPlaylist.has(song.song_id) ? (
                "Added"
              ) : (
                <button
                  onClick={() => handleAddSong(song.song_id)}
                  className="px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 text-sm"
                >
                  Add
                </button>
              )}
            </li>
          )) : (
            <p className="text-gray-500">No songs found.</p>
          )}
        </ul>
      </div>
    </div>
  );
}
