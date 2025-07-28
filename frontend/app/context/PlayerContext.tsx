'use client';

import React, { createContext, useContext, useState } from 'react';

interface PlayerContextType {
  currentSong: { id: string; title: string } | null;
  isPlaying: boolean;
  play: () => void;
  pause: () => void;
  setCurrentSong: (song: { id: string; title: string }) => void;
}

const PlayerContext = createContext<PlayerContextType>({
  currentSong: null,
  isPlaying: false,
  play: () => {},
  pause: () => {},
  setCurrentSong: () => {},
});

export const usePlayer = () => useContext(PlayerContext);

export const PlayerProvider = ({ children }: { children: React.ReactNode }) => {
  const [currentSong, setCurrentSong] = useState<{ id: string; title: string } | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);

  const play = () => setIsPlaying(true);
  const pause = () => setIsPlaying(false);

  return (
    <PlayerContext.Provider value={{ currentSong, setCurrentSong, isPlaying, play, pause }}>
      {children}
    </PlayerContext.Provider>
  );
};
