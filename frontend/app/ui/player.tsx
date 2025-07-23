'use client';

import { useEffect, useRef } from 'react';
import { usePlayer } from '@/app/context/PlayerContext';

export default function Player() {
  const { currentSong } = usePlayer();
  const audioRef = useRef<HTMLAudioElement>(null);
  useEffect(() => {
    if (currentSong && audioRef.current) {
      audioRef.current.src = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/play/${currentSong}`;
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
      <audio
        ref={audioRef}
        controls
        autoPlay
        src={currentSong ?? `${process.env.NEXT_PUBLIC_API_URL}/api/v1/play/${currentSong}`}
      />
    </div>
  );
}
