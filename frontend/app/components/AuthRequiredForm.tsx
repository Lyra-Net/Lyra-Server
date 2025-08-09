'use client';

import React from 'react';
import { useRouter } from 'next/navigation';

interface AuthRequiredFormProps {
  homePath?: string;
  trendsPath?: string;
}

export default function AuthRequiredForm({
  homePath = '/',
  trendsPath = '/trends',
}: AuthRequiredFormProps) {
  const router = useRouter();

  return (
    <div className="flex items-center justify-center min-h-[60vh]">
      <div className="flex flex-col items-center justify-center p-6 bg-gray-100 dark:bg-gray-800 rounded-lg shadow-lg max-w-sm w-full animate-fade-slide-down">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 text-center">
          Login to use this feature
        </h2>
        <div className="flex flex-col gap-3 w-full">
          <button
            onClick={() => router.push('/login')}
            className="w-full px-4 py-2 rounded bg-blue-500 hover:bg-blue-600 text-white transition"
          >
            Log in
          </button>
          <button
            onClick={() => router.push(homePath)}
            className="w-full px-4 py-2 rounded bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 transition"
          >
            Go to Home
          </button>
          <button
            onClick={() => router.push(trendsPath)}
            className="w-full px-4 py-2 rounded bg-green-500 hover:bg-green-600 text-white transition"
          >
            See Trending Playlists
          </button>
        </div>
      </div>
    </div>
  );
}
