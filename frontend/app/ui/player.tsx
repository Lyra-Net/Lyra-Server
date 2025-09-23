'use client';

import { useEffect, useRef, useState } from 'react';
import { usePlayerStore } from '@/stores/player';
import {
  CirclePause,
  CirclePlay,
  Repeat,
  Shuffle,
  SkipBack,
  SkipForward,
} from 'lucide-react';

export default function Player() {
  const audioRef = useRef<HTMLAudioElement>(null);

  const {
    currentSong,
    isPlaying,
    progress,
    duration,
    volume,
    repeat,
    shuffle,
    restartSong,
    setProgress,
    setDuration,
    setPlaying,
    setVolume,
    playNext,
    playPrev,
    toggleRepeat,
    toggleShuffle,
    setRestartSong,
    canShuffleRepeat,
  } = usePlayerStore();

  const [hasHydrated, setHasHydrated] = useState(false);

  useEffect(() => {
    setHasHydrated(usePlayerStore.persist.hasHydrated());
  }, []);

  // Restart song
  useEffect(() => {
    if (!restartSong || !audioRef.current) return;
    const audio = audioRef.current;
    audio.currentTime = 0;
    audio.play().catch(() => {});
    setRestartSong(false);
  }, [restartSong, setRestartSong]);

  // Load new song
  useEffect(() => {
    const audio = audioRef.current;
    if (currentSong && audio) {
      audio.src = `${process.env.NEXT_PUBLIC_API_URL}/stream/${currentSong.song_id}.mp3`;
      audio.load();
      audio
        .play()
        .then(() => setPlaying(true))
        .catch((err) => {
          console.warn('Playback failed:', err);
          setPlaying(false);
        });
    }
    if (!currentSong && audio) {
      audio.src = '';
      setPlaying(false);
    }
  }, [currentSong, setPlaying]);

  // Play / Pause
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;
    if (isPlaying) {
      audio.play().catch(() => {});
    } else {
      audio.pause();
    }
  }, [isPlaying]);

  // Progress + Ended
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const updateProgress = () => {
      setProgress(audio.currentTime);
      setDuration(audio.duration || 0);
    };

    const handleEnded = () => {
      playNext();
    };

    audio.addEventListener('timeupdate', updateProgress);
    audio.addEventListener('loadedmetadata', updateProgress);
    audio.addEventListener('ended', handleEnded);

    return () => {
      audio.removeEventListener('timeupdate', updateProgress);
      audio.removeEventListener('loadedmetadata', updateProgress);
      audio.removeEventListener('ended', handleEnded);
    };
  }, [setProgress, setDuration, playNext]);

  // Seek
  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!audioRef.current) return;
    const value = Number(e.target.value);
    audioRef.current.currentTime = value;
    setProgress(value);
  };

  // Volume
  const handleVolume = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Number(e.target.value);
    if (audioRef.current) audioRef.current.volume = value;
    setVolume(value);
  };

  const formatTime = (time: number) => {
    if (isNaN(time)) return '0:00';
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60).toString().padStart(2, '0');
    return `${minutes}:${seconds}`;
  };

  const canUseShuffle = hasHydrated && canShuffleRepeat();
  const canUseRepeat =
    hasHydrated &&
    (repeat === 'one' || canShuffleRepeat());

  return (
    <div className="h-16 px-4 bg-black flex items-center text-sm text-gray-300 gap-x-4">
      <div className="flex items-center gap-x-3 w-1/4 truncate">
        {currentSong && (
          <img
            src={`https://i.ytimg.com/vi/${currentSong.song_id}/default.jpg`}
            className="w-12 h-12 rounded"
          />
        )}
        <span className="truncate">
          {currentSong
            ? `🎵 ${currentSong.title} - ${currentSong.artists
                .map((a) => a.name)
                .join(', ')}`
            : 'No song selected'}
        </span>
      </div>
      <div className="flex-1 flex flex-col items-center justify-center gap-y-1">
        <div className="flex items-center gap-x-4">
          <Shuffle
            color={shuffle ? 'green' : 'white'}
            onClick={canUseShuffle ? toggleShuffle : undefined}
            className={
              canUseShuffle ? 'cursor-pointer' : 'opacity-50 cursor-not-allowed'
            }
          />
          <SkipBack onClick={playPrev} />
          {isPlaying ? (
            <CirclePause
              size={32}
              color="white"
              onClick={() => setPlaying(false)}
            />
          ) : (
            <CirclePlay
              size={32}
              color="white"
              onClick={() => {
                if (currentSong) {
                  setPlaying(true);
                } else {
                  playNext();
                }
              }}
            />
          )}
          <SkipForward onClick={playNext} />
          <div className="relative flex items-center justify-center w-6 h-6">
            <Repeat
              color={repeat ? 'green' : 'white'}
              onClick={canUseRepeat ? toggleRepeat : undefined}
              className={
                canUseRepeat ? 'cursor-pointer' : 'opacity-50 cursor-not-allowed'
              }
            />
            {repeat === 'one' && (
              <span className="absolute top-[-8px] right-[8px] text-[10px] text-green-500 font-bold">
                1
              </span>
            )}
          </div>
        </div>

        <div className="flex items-center gap-x-2 w-full">
          <span className="text-xs">{formatTime(progress)}</span>
          <input
            type="range"
            min={0}
            max={duration || 0}
            value={progress}
            onChange={handleSeek}
            className="w-full"
          />
          <span className="text-xs">{formatTime(duration)}</span>
        </div>
      </div>
      <div className="w-1/4 flex items-center justify-end gap-x-2">
        <span>🔊</span>
        <input
          type="range"
          min={0}
          max={1}
          step={0.01}
          value={volume}
          onChange={handleVolume}
        />
      </div>
      <audio ref={audioRef} autoPlay hidden />
    </div>
  );
}
