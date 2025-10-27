'use client';

import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import DashboardLayout from '@/app/ui/dashboardLayout';
import { toast } from 'sonner';
import { playlistApi } from '@/lib/playlistApi';
import AddSongToPlaylist from '@/app/ui/addSongToPlaylist';
import { Playlist, Song } from '@/declarations/playlists';
import { usePlayerStore } from '@/stores/player';
import { FaPlayCircle } from 'react-icons/fa';
import { FaCirclePause } from "react-icons/fa6";

import {
  DndContext,
  closestCenter,
  DragEndEvent
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

// draggable song item
function SortableSong({ song, onRemove }: {
  song: Song;
  onRemove: (id: string) => void;
}) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: song.song_id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <li
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className="p-3 rounded flex items-center justify-between shadow-sm hover:bg-[var(--background-highlight)] cursor-grab"
    >
      <div className="flex items-center gap-3">
        <img
          src={`https://i.ytimg.com/vi/${song.song_id}/default.jpg`}
          alt={song.title}
          className="w-16 h-14 rounded"
        />
        <div>
          <p className="font-medium text-gray-900 dark:text-gray-100">{song.title}</p>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            {song.artists?.map(a => a.name).join(", ")}
          </p>
        </div>
      </div>

      <button
        onClick={() => onRemove(song.song_id)}
        className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
      >
        ✕
      </button>
    </li>
  );
}

export default function PlaylistDetailPage() {
  const params = useParams();
  const router = useRouter();
  const playlistId = params.playlistId as string;
  const [playlist, setPlaylist] = useState<Playlist>();
  const [loading, setLoading] = useState(true);
  const [renaming, setRenaming] = useState(false);
  const [newName, setNewName] = useState('');
  const [showAddSong, setShowAddSong] = useState(false);

  const currentSong = usePlayerStore(s => s.currentSong);
  const isPlaying = usePlayerStore(s => s.isPlaying);
  const source = usePlayerStore(s => s.source);
  const addToQueueUnique = usePlayerStore(s => s.addToQueueUnique);
  const setQueue = usePlayerStore(s => s.setQueue);
  const setPlaying = usePlayerStore(s => s.setPlaying);

  const isCurrentPlaylist = currentSong &&
    source?.type === "playlist" && source?.id === playlist?.playlist_id;

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
        is_public: playlist?.is_public ?? false,
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
  const handleMove = async (songId: string, new_position: number) => {
    try {
      await playlistApi.moveSong(playlistId, songId, new_position);
      fetchPlaylist();
    } catch (err) {
      console.error("Reorder failed", err);
      toast.error("Failed to move song");
    }
  };

  // play playlist
  const handlePlayNow = () => {
    if (!playlist?.songs.length) return;
    setQueue(playlist.songs, {
      type: "playlist",
      id: playlist.playlist_id,
      name: playlist.playlist_name,
    });
    toast.success(`Playing playlist: ${playlist.playlist_name}`);
  };

  const handleAddToQueue = () => {
    if (!playlist?.songs.length) return;
    playlist.songs.forEach((song) => addToQueueUnique(song));
    toast.success(`Added ${playlist.songs.length} songs to queue`);
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
        <div className="flex gap-6 items-start bg-gray-500 w-full h-48">
          <div>
            {playlist.songs.length && <img
              src={`https://i.ytimg.com/vi/${playlist.songs[0].song_id}/hqdefault.jpg`}
              className="m-2 w-80 h-[180px] object-cover rounded"
            />}
          </div>
          <div className="flex mt-10 flex-col justify-center gap-2">
            <span className="text-sm text-gray-400">
              {playlist.is_public ? "public" : "private"} playlist
            </span>
            <span className="text-white text-5xl font-bold">
              {playlist.playlist_name}
            </span>
            <span className="text-sm text-gray-400">
              {playlist.songs.length}{" "}
              {playlist.songs.length === 1 ? "song" : "songs"}
            </span>
          </div>
        </div>
      </div>
      <div className="flex items-center gap-4">
        {isCurrentPlaylist ? (
          isPlaying ? (
            <FaCirclePause
              size={48}
              color='green'
              className="text-primary cursor-pointer"
              onClick={() => setPlaying(false)}
            />
          ) : (
            <FaPlayCircle
              size={48}
              color='green'
              className="text-primary cursor-pointer"
              onClick={() => setPlaying(true)}
            />
          )
        ) : (
          <FaPlayCircle
            size={48}
            color='green'
            className="text-primary cursor-pointer"
            onClick={handlePlayNow}
          />
        )}
      </div>

      <ul className="space-y-1 p-2">
        {playlist.songs.length ? (
          <DndContext
            collisionDetection={closestCenter}
            onDragEnd={async (event: DragEndEvent) => {
              const { active, over } = event;
              if (over && active.id !== over.id) {
                const oldIndex = playlist.songs.findIndex(s => s.song_id === active.id);
                const newIndex = playlist.songs.findIndex(s => s.song_id === over.id);

                // re-order local UI
                const newSongs = arrayMove(playlist.songs, oldIndex, newIndex);
                setPlaylist({ ...playlist, songs: newSongs });

                // gọi API update position
                await handleMove(active.id as string, newIndex);
              }
            }}
          >
            <SortableContext
              items={playlist.songs.map(s => s.song_id)}
              strategy={verticalListSortingStrategy}
            >
              {playlist.songs.map((song) => (
                <SortableSong
                  key={song.song_id}
                  song={song}
                  onRemove={handleRemoveSong}
                />
              ))}
            </SortableContext>
          </DndContext>
        ) : (
          <p className="text-gray-500">No songs in this playlist.</p>
        )}
      </ul>

      {showAddSong && (
        <AddSongToPlaylist
          setShowAddSong={setShowAddSong}
          refreshPlaylist={fetchPlaylist}
          playlist={playlist}
        />
      )}
    </DashboardLayout>
  );
}
