'use client';

import { useEffect, useState } from 'react';
import DashboardLayout from '../ui/dashboardLayout';
import { Playlist } from '@/declarations/playlists';
import api from '@/lib/api';
import AuthRequiredForm from "../components/AuthRequiredForm";
import Link from "next/link";
import { useAuthStore } from "@/stores/useAuth";

export default function PlaylistPage() {
  const [myPlaylist, setMyPlaylist] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(true);
  const { userId, accessToken } = useAuthStore();
  console.log({userId, accessToken})
  const [newName, setNewName] = useState("");
  const [isPublic, setIsPublic] = useState(false);
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    if (!accessToken) {
      setAuthenticated(false);
      setLoading(false);
      return;
    }
    api
      .post('/playlist/list')
      .then(res => {
        setMyPlaylist(res.data?.playlists || []);
        setAuthenticated(true);
      })
      .catch((err) => {
        console.error("error fetching playlists", err);
        setAuthenticated(false);
      })
      .finally(() => setLoading(false));
  }, [accessToken]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!accessToken) return;

    setCreating(true);
    try {
      await api.post('/playlist/create', {
        owner_id: userId,
        playlist_name: newName,
        is_public: isPublic,
      });

      const updated = await api.post('/playlist/list', { owner_id: userId });
      setMyPlaylist(updated.data?.playlists || []);
      setNewName("");
      setIsPublic(true);
    } catch (err) {
      console.error("Failed to create playlist", err);
    } finally {
      setCreating(false);
    }
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex justify-center items-center min-h-[60vh]">
          <p className="text-gray-500">Loading playlists...</p>
        </div>
      </DashboardLayout>
    );
  }

  if (!authenticated) {
    return (
      <DashboardLayout>
        <div className="flex justify-center items-center min-h-[60vh]">
          <AuthRequiredForm homePath="/" trendsPath="trends"/>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="mx-20 p-6">
      {/* Create playlist form */}
      <form onSubmit={handleCreate} className="shadow rounded-lg p-4 mb-6">
        <h2 className="text-lg font-semibold mb-3 text-gray-800 dark:text-gray-100">Create Playlist</h2>
        <div className="flex flex-col gap-3">
          <input
            type="text"
            placeholder="Playlist name"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            required
            className="px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600"
          />
          
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={isPublic}
              onChange={(e) => setIsPublic(e.target.checked)}
            />
            <span className="text-gray-700 dark:text-gray-200">Public</span>
          </label>
          <button
            type="submit"
            disabled={creating}
            className="px-4 py-2 rounded bg-blue-500 hover:bg-blue-600 text-white transition disabled:opacity-50"
          >
            {creating ? "Creating..." : "Create Playlist"}
          </button>
        </div>
      </form>

      {/* playlists */}
      {myPlaylist && myPlaylist.length > 0 ? (
        <div className="grid sm:grid-cols-2 md:grid-cols-3 gap-4">
          {myPlaylist.map((pl: Playlist) => {
            return (
            <Link
              key={pl.playlist_id}
              href={`/playlists/${pl.playlist_id}`}
              className="block p-4 rounded-lg shadow bg-white/70 dark:bg-gray-800/70 hover:bg-blue-100/50 dark:hover:bg-gray-700 transition"
            >
              <h3 className="text-lg font-semibold text-gray-800 dark:text-gray-100">
                {pl.playlist_name}
              </h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {pl.is_public ? "Public" : "Private"}
              </p>
              <p className="text-xs text-gray-400 mt-1">
                {pl.songs.length || 0} {pl.songs.length === 1 ? "song" : "songs"}
              </p>
            </Link>
            );} )}
        </div>
      )
       : (
        <div className="m-4">
          <p className="text-gray-500">No playlists found.</p>
        </div>
       )}
      </div>
    </DashboardLayout>
  );
}
