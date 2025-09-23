'use client'

import { useState } from 'react'
import UpNext from '@/app/ui/upNext'
import History from '@/app/ui/history'

export default function RightPanel() {
  const [tab, setTab] = useState<'upnext' | 'history'>('upnext')

  return (
    <div className="">
      <div className="flex space-x-4 mb-4">
        <button
          onClick={() => setTab('upnext')}
          className={`px-2 py-1 rounded ${tab === 'upnext' ? 'bg-gray-700' : 'bg-gray-800'}`}
        >
          Up Next
        </button>
        <button
          onClick={() => setTab('history')}
          className={`px-2 py-1 rounded ${tab === 'history' ? 'bg-gray-700' : 'bg-gray-800'}`}
        >
          History
        </button>
      </div>

      {tab === 'upnext' ? <UpNext /> : <History />}
    </div>
  )
}
