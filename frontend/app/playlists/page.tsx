'use client';

import { useEffect, useState } from 'react';
import DashboardLayout from '../ui/dashboardLayout';
import { Playlists } from '@/declarations/playlists';
import AuthRequiredForm from '../components/AuthRequiredForm';
import api from '@/lib/api';

export default function Playlist() {
  const [myPlaylist, setMyPlaylist] = useState<Playlists | null>(null);
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      setAuthenticated(false);
      setLoading(false);
      return;
    }

    api
      .get<Playlists>('/api/v1/playlists')
      .then(res => {
        setMyPlaylist(res.data);
        setAuthenticated(true);
      })
      .catch(err => {
        console.error(err);
        setAuthenticated(false);
      })
      .finally(() => setLoading(false));
  }, []);

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
        <AuthRequiredForm homePath="/" trendsPath="/trends" />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <button>Create new playlist</button>
      {myPlaylist && (
        <div className="mt-4">
          <pre>{JSON.stringify(myPlaylist, null, 2)}</pre>
        </div>
      )}
    </DashboardLayout>
  );
}
