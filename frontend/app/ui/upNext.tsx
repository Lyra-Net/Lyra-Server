'use client'
import { usePlayerStore } from '@/stores/player';

export default function UpNext() {
  const { queue, removeFromQueue } = usePlayerStore();

  return (
    <div className='p-1'>
      <h2 className="text-lg font-semibold mb-4">Up Next</h2>
      {queue.length ? (
        <>
          <ul className="space-y-2">
            {queue.map((song) => (
              <li key={song.song_id} className="flex justify-between items-center">
                <span>{song.title}</span>
                <button
                  className="text-sm text-red-400"
                  onClick={() => removeFromQueue(song.song_id)}
                >
                  Remove
                </button>
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
