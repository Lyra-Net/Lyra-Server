'use client';

import { useEffect, useRef } from 'react';
import { usePlayer } from '@/app/context/PlayerContext';

export default function Player() {
  const { currentSong } = usePlayer();
  const audioRef = useRef<HTMLAudioElement>(null);

  useEffect(() => {
    if (currentSong && audioRef.current) {
      audioRef.current.src = `/api/play/${currentSong}`;
      audioRef.current.play().catch(() => {});
    }
  }, [currentSong]);

  return (
    <div className="h-16 bg-gray-900/10 px-6 flex items-center justify-between text-sm text-gray-300">
      <div>Now playing: 🎵 {currentSong || 'No song'}</div>
      <div className="space-x-4">
        <button>⏮</button>
        <button onClick={() => audioRef.current?.play()}>▶️</button>
        <button onClick={() => audioRef.current?.pause()}>⏸</button>
      </div>
      <audio controls autoPlay src={`/api/play/${currentSong}`} />
    </div>
  );
}
