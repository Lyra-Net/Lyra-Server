'use client';

import { useEffect, useRef } from 'react';
import { usePlayer } from '@/app/context/PlayerContext';

export default function Player() {
  const { currentSong, isPlaying, play, pause } = usePlayer();
  const audioRef = useRef<HTMLAudioElement>(null);

  useEffect(() => {
    if (currentSong && audioRef.current) {
      const audio = audioRef.current;
      audio.src = `${process.env.NEXT_PUBLIC_API_URL}/stream/${currentSong.id}.mp3`;
      audio
        .play()
        .then(() => play())
        .catch(err => {
          console.warn('Playback failed:', err);
        });
    }
  }, [currentSong]);

  useEffect(() => {
    if (!audioRef.current) return;

    if (isPlaying) {
      audioRef.current.play().catch(() => {});
    } else {
      audioRef.current.pause();
    }
  }, [isPlaying]);

  return (
    <div className="h-16 bg-gray-900/10 px-6 flex items-center justify-between text-sm text-gray-300">
      <div>Now playing: 🎵 {currentSong?.title || 'No song selected'}</div>
      <div className="space-x-4">
        <button>⏮</button>
        <button onClick={() => play()}>▶️</button>
        <button onClick={() => pause()}>⏸</button>
      </div>

      <audio ref={audioRef} autoPlay controls hidden={!currentSong?.id} />
    </div>
  );
}
