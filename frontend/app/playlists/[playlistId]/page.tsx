'use client';

import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import DashboardLayout from '@/app/ui/dashboardLayout';
import { toast } from 'sonner';
import { playlistApi } from '@/lib/playlistApi';
import AddSongToPlaylist from '@/app/ui/addSongToPlaylist';

export default function PlaylistDetailPage() {
  const params = useParams();
  const router = useRouter();
  const playlistId = params.playlistId as string;
  const [playlist, setPlaylist] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [renaming, setRenaming] = useState(false);
  const [newName, setNewName] = useState('');
  const [showAddSong, setShowAddSong] = useState(false);

  // fetch playlist
  const fetchPlaylist = async () => {
    try {
      const res = await playlistApi.get(playlistId);
      setPlaylist(res.data);
      setNewName(res.data.playlist_name);
    } catch (err) {
      console.error("Failed to load playlist", err);
      toast.error("Failed to load playlist");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlaylist();
  }, [playlistId]);

  // rename playlist
  const handleRename = async () => {
    if (!newName.trim()) {
      toast.error("Name cannot be empty");
      return;
    }
    try {
      await playlistApi.update({
        playlist_id: playlistId,
        playlist_name: newName,
        is_public: playlist.is_public ?? false,
      });
      toast.success("Playlist renamed");
      setRenaming(false);
      fetchPlaylist();
    } catch (err) {
      console.error("Rename failed", err);
      toast.error("Rename failed");
    }
  };

  // delete playlist
  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this playlist?")) return;
    try {
      await playlistApi.delete(playlistId);
      toast.success("Playlist deleted");
      router.push('/playlists');
    } catch (err) {
      console.error("Delete failed", err);
      toast.error("Delete failed");
    }
  };

  // remove song
  const handleRemoveSong = async (songId: string) => {
    try {
      await playlistApi.removeSong(playlistId, songId);
      toast.success("Song removed");
      fetchPlaylist();
    } catch (err) {
      console.error("Remove song failed", err);
      toast.error("Remove song failed");
    }
  };

  // reorder song
  const handleMove = async (songId: string, direction: 'up' | 'down') => {
    try {
      await playlistApi.moveSong(playlistId, songId, direction);
      fetchPlaylist();
    } catch (err) {
      console.error("Reorder failed", err);
      toast.error("Failed to move song");
    }
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex justify-center items-center min-h-[60vh]">
          <p className="text-gray-500">Loading playlist...</p>
        </div>
      </DashboardLayout>
    );
  }

  if (!playlist) {
    return (
      <DashboardLayout>
        <p className="text-gray-500">Playlist not found.</p>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="flex items-center justify-between mb-6">
        {renaming ? (
          <div className="flex items-center gap-2">
            <input
              type="text"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleRename();
                if (e.key === "Escape") setRenaming(false);
              }}
              className="px-2 py-1 border rounded dark:bg-gray-700 dark:border-gray-600"
              autoFocus
            />
            <button
              onClick={handleRename}
              className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Save
            </button>
            <button
              onClick={() => setRenaming(false)}
              className="px-3 py-1 bg-gray-500 text-white rounded hover:bg-gray-600"
            >
              Cancel
            </button>
          </div>
        ) : (
          <h1 className="text-xl font-bold text-gray-800 dark:text-gray-100">
            {playlist.playlist_name}
          </h1>
        )}

        <div className="flex gap-2">
          {!renaming && (
            <button
              onClick={() => setRenaming(true)}
              className="px-3 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600"
            >
              Rename
            </button>
          )}
          <button
            onClick={handleDelete}
            className="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700"
          >
            Delete
          </button>
        </div>
      </div>

      <div className="mb-4">
        <button
          onClick={() => {
            setShowAddSong(true);
          }}
          className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          + Add Song
        </button>
      </div>

      <ul className="space-y-2">
        {playlist.songs.length ? playlist.songs.map((song: any, idx: number) => (
          <li
            key={song.id}
            className="p-3 bg-white dark:bg-gray-700 rounded flex justify-between items-center"
          >
            <span>{idx + 1}. {song.title}</span>
            <div className="flex gap-2">
              <button
                onClick={() => handleMove(song.id, 'up')}
                disabled={idx === 0}
                className="px-2 py-1 bg-gray-300 dark:bg-gray-600 rounded disabled:opacity-50"
              >
                ↑
              </button>
              <button
                onClick={() => handleMove(song.id, 'down')}
                disabled={idx === playlist.songs.length - 1}
                className="px-2 py-1 bg-gray-300 dark:bg-gray-600 rounded disabled:opacity-50"
              >
                ↓
              </button>
              <button
                onClick={() => handleRemoveSong(song.id)}
                className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
              >
                Remove
              </button>
            </div>
          </li>
        )) : (
          <p className="text-gray-500">No songs in this playlist.</p>
        )}
      </ul>
      {showAddSong && (
        <AddSongToPlaylist
        setShowAddSong={setShowAddSong}
          playlist={playlist} />)}
    </DashboardLayout>
  );
}
