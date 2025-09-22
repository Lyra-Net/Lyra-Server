import { create } from "zustand";
import { persist } from "zustand/middleware";
import { Song } from "@/declarations/playlists";

type PlayerSource = { type: "playlist"; id: string; name: string } | null;

function shuffleArray<T>(array: T[]): T[] {
  const arr = [...array];
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [arr[i], arr[j]] = [arr[j], arr[i]];
  }
  return arr;
}

interface PlayerState {
  currentSong: Song | null;
  queue: Song[]; 
  history: Song[];
  currentTime: number;
  duration: number;
  progress: number;
  isPlaying: boolean;
  volume: number;
  repeat: boolean;
  shuffle: boolean;
  restartSong: boolean;
  source: PlayerSource;

  // shuffle state
  originalQueue: Song[];
  shuffledQueue: Song[];
  shuffledIndex: number;
  playedSongs: Song[];

  setCurrentSong: (song: Song) => void;
  addToQueue: (song: Song) => void;
  addToQueueUnique: (song: Song) => void;
  setQueue: (songs: Song[], src?: PlayerSource) => void;
  removeFromQueue: (songId: string) => void;
  clearQueue: () => void;
  playNext: () => void;
  playPrev: () => void;

  setCurrentTime: (time: number) => void;
  setDuration: (time: number) => void;
  setProgress: (progress: number) => void;
  setPlaying: (playing: boolean) => void;
  setVolume: (vol: number) => void;
  toggleRepeat: () => void;
  toggleShuffle: () => void;
  setRestartSong: (flag: boolean) => void;
  canShuffleRepeat: () => boolean;
  setSource: (src: PlayerSource) => void;
}

export const usePlayerStore = create<PlayerState>()(
  persist(
    (set, get) => ({
      currentSong: null,
      queue: [],
      history: [],
      currentTime: 0,
      duration: 0,
      progress: 0,
      isPlaying: false,
      volume: 1,
      repeat: false,
      shuffle: false,
      restartSong: false,
      source: null as PlayerSource,

      originalQueue: [],
      shuffledQueue: [],
      shuffledIndex: 0,
      playedSongs: [],

      setCurrentSong: (song) =>
        set((state) => ({
          currentSong: song,
          queue: state.queue.filter((s) => s.song_id !== song.song_id),
          currentTime: 0,
          progress: 0,
          isPlaying: true,
        })),

      addToQueue: (song) => set((state) => ({ queue: [...state.queue, song] })),

      addToQueueUnique: (song) =>
        set((state) => {
          if (state.queue.some((s) => s.song_id === song.song_id)) return {};
          return { queue: [...state.queue, song] };
        }),

      setQueue: (songs: Song[], src?: PlayerSource) =>
        set(() => ({
          queue: songs.slice(1),
          currentSong: songs[0] || null,
          currentTime: 0,
          progress: 0,
          src: src || null,
          isPlaying: songs.length > 0,
          history: [],
          originalQueue: songs,
          shuffledQueue: [],
          shuffledIndex: 0,
          playedSongs: songs[0] ? [songs[0]] : [],
        })),

      removeFromQueue: (songId) =>
        set((state) => ({ queue: state.queue.filter((s) => s.song_id !== songId) })),

      clearQueue: () =>
        set({ queue: [], history: [], originalQueue: [], shuffledQueue: [], shuffledIndex: 0, playedSongs: [] }),

      playNext: () =>
        set((state) => {
          const { shuffle, repeat, currentSong, history, source } = state;

          if (shuffle) {
            const nextIndex = state.shuffledIndex + 1;
            if (nextIndex < state.shuffledQueue.length) {
              const nextSong = state.shuffledQueue[nextIndex];
              return {
                currentSong: nextSong,
                shuffledIndex: nextIndex,
                playedSongs: [...state.playedSongs, nextSong],
                history: currentSong ? [...history, currentSong] : history,
                currentTime: 0,
                progress: 0,
                isPlaying: true,
              };
            }
            // handle if repeat is on
            return { isPlaying: false };
          }

          if (state.queue.length > 0) {
            const [next, ...rest] = state.queue;
            return {
              currentSong: next,
              queue: rest,
              history: currentSong ? [...history, currentSong] : history,
              currentTime: 0,
              progress: 0,
              isPlaying: true,
            };
          }

          if (repeat && source) {
            // if source exists, restart the queue
            return { currentTime: 0, progress: 0, isPlaying: true };
          }

          return { isPlaying: false };
        }),

      playPrev: () =>
        set((state) => {
          const { history, currentSong, queue, shuffle, shuffledIndex, shuffledQueue } = state;

          if (history.length > 0) {
            const prev = history[history.length - 1];
            const newHistory = history.slice(0, -1);

            if (shuffle) {
              const idx = shuffledQueue.findIndex((s) => s.song_id === prev.song_id);
              return {
                currentSong: prev,
                history: newHistory,
                shuffledIndex: idx >= 0 ? idx : shuffledIndex,
                currentTime: 0,
                progress: 0,
                isPlaying: true,
              };
            }

            const newQueue = currentSong
              ? [currentSong, ...queue.filter((s) => s.song_id !== prev.song_id)]
              : queue.filter((s) => s.song_id !== prev.song_id);

            return {
              currentSong: prev,
              history: newHistory,
              queue: newQueue,
              currentTime: 0,
              progress: 0,
              isPlaying: true,
            };
          }

          return { currentTime: 0, progress: 0, restartSong: true, isPlaying: true };
        }),

      setCurrentTime: (time) => set({ currentTime: time }),
      setDuration: (time) => set({ duration: time }),
      setProgress: (progress) => set({ progress }),
      setPlaying: (playing) => set({ isPlaying: playing }),
      setVolume: (vol) => set({ volume: vol }),
      toggleRepeat: () => set((s) => ({ repeat: !s.repeat })),

      toggleShuffle: () =>
        set((state) => {
          if (!state.shuffle) {
            const shuffled = shuffleArray([...state.originalQueue]);
            const currentIdx = shuffled.findIndex((s) => s.song_id === state.currentSong?.song_id);
            let shuffledQueue = shuffled;

            if (currentIdx > 0) {
              const [current] = shuffled.splice(currentIdx, 1);
              shuffledQueue = [current, ...shuffled];
            }

            return {
              shuffle: true,
              shuffledQueue,
              shuffledIndex: 0,
              playedSongs: state.currentSong ? [state.currentSong] : [],
            };
          } else {
            const remaining = state.originalQueue.filter(
              (s) => !state.playedSongs.some((p) => p.song_id === s.song_id)
            );
            return {
              shuffle: false,
              queue: remaining,
              shuffledQueue: [],
              shuffledIndex: 0,
            };
          }
        }),

      setRestartSong: (flag) => set({ restartSong: flag }),
      canShuffleRepeat: () => !!get().source,
      setSource: (src: PlayerSource) => set({ source: src }),
    }),
    {
      name: "player-storage",
      partialize: (state) => ({
        currentSong: state.currentSong,
        queue: state.queue,
        history: state.history,
        currentTime: state.currentTime,
        duration: state.duration,
        progress: state.progress,
        isPlaying: state.isPlaying,
        volume: state.volume,
        repeat: state.repeat,
        shuffle: state.shuffle,
        restartSong: state.restartSong,
        source: state.source,
        originalQueue: state.originalQueue,
        shuffledQueue: state.shuffledQueue,
        shuffledIndex: state.shuffledIndex,
        playedSongs: state.playedSongs,
      }),
    }
  )
);
