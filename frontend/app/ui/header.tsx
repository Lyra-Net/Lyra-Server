'use client';
import { useEffect, useState } from 'react';

export default function Header() {
  const [username, setUsername] = useState<string | null>(null);

  useEffect(() => {
    const storedName = localStorage.getItem('username');
    if (storedName) {
      setUsername(storedName);
    }
  }, []);
  return (
    <header className="h-16 px-6 flex items-center justify-between bg-gray-900:0">
      <div className="text-lg font-semibold">Dashboard</div>
      <div className="text-sm">{username ? `Welcome back, ${username}!` : 'Welcome!'}</div>
    </header>
  );
}
