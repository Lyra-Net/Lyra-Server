'use client'
import { usePlayerStore } from '@/stores/player';

export default function UpNext() {
  const { queue, shuffledQueue, shuffledIndex, shuffle } = usePlayerStore();

  const upNext = shuffle
    ? shuffledQueue.slice(shuffledIndex + 1)
    : queue;

  return (
    <div className='p-1'>
      <h2 className="text-lg font-semibold mb-4">Up Next</h2>
      {upNext.length ? (
        <>
          <ul className="space-y-2">
            {upNext.map((song) => (
              <li key={song.song_id} className="flex justify-between items-center">
                <span>{song.title}</span>
              </li>
            ))}
          </ul>
      </>
      ) : (
        <p className="text-gray-400">No songs in queue</p>
      )}

      
    </div>
  );
}
