import { create } from "zustand";
import { persist } from "zustand/middleware";
import { Song } from "@/declarations/playlists";

interface PlayerState {
  currentSong: Song | null;
  queue: Song[];
  currentTime: number;
  duration: number;
  isPlaying: boolean;
  volume: number;
  repeat: boolean;
  shuffle: boolean;

  setCurrentSong: (song: Song) => void;
  addToQueue: (song: Song) => void;
  removeFromQueue: (songId: string) => void;
  clearQueue: () => void;
  playNext: () => void;
  playPrev: () => void;

  setCurrentTime: (time: number) => void;
  setDuration: (time: number) => void;
  setPlaying: (playing: boolean) => void;
  setVolume: (vol: number) => void;
  toggleRepeat: () => void;
  toggleShuffle: () => void;
}

export const usePlayerStore = create<PlayerState>()(
  persist(
    (set, get) => ({
      currentSong: null,
      queue: [],
      currentTime: 0,
      duration: 0,
      isPlaying: false,
      volume: 1,
      repeat: false,
      shuffle: false,

      setCurrentSong: (song) =>
        set({ currentSong: song, currentTime: 0, isPlaying: true }),
      addToQueue: (song) => set({ queue: [...get().queue, song] }),
      removeFromQueue: (songId) =>
        set({ queue: get().queue.filter((s) => s.song_id !== songId) }),
      clearQueue: () => set({ queue: [] }),
      playNext: () => {
        const { queue, shuffle, repeat } = get();
        if (shuffle && queue.length > 1) {
          const idx = Math.floor(Math.random() * queue.length);
          const next = queue[idx];
          set({
            currentSong: next,
            queue: queue.filter((_, i) => i !== idx),
            currentTime: 0,
          });
          return;
        }
        if (queue.length > 0) {
          const [next, ...rest] = queue;
          set({ currentSong: next, queue: rest, currentTime: 0 });
        } else if (repeat && get().currentSong) {
          set({ currentTime: 0 });
        } else {
          set({ isPlaying: false });
        }
      },
      playPrev: () => {
        set({ currentTime: 0 });
      },

      setCurrentTime: (time) => set({ currentTime: time }),
      setDuration: (time) => set({ duration: time }),
      setPlaying: (playing) => set({ isPlaying: playing }),
      setVolume: (vol) => set({ volume: vol }),
      toggleRepeat: () => set({ repeat: !get().repeat }),
      toggleShuffle: () => set({ shuffle: !get().shuffle }),
    }),
    {
      name: "player-storage",
      partialize: (state) => ({
        currentSong: state.currentSong,
        queue: state.queue,
        currentTime: state.currentTime,
        isPlaying: state.isPlaying,
        volume: state.volume,
        repeat: state.repeat,
        shuffle: state.shuffle,
      }),
    }
  )
);
