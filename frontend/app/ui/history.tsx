'use client'
import { usePlayerStore } from '@/stores/player'

export default function History() {
  const { history } = usePlayerStore()

  return (
    <div className="p-1">
      <h2 className="text-lg font-semibold mb-4">History</h2>
      {history.length ? (
        <ul className="space-y-2">
          {history
            .slice()
            .reverse()
            .map((song) => (
              <li
                key={song.song_id}
                className="flex justify-between items-center"
              >
                <span>{song.title}</span>
              </li>
            ))}
        </ul>
      ) : (
        <p className="text-gray-400">No songs played yet</p>
      )}
    </div>
  )
}
