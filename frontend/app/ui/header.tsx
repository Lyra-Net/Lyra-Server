'use client';
import { useSession } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function Header() {
  const [username, setUsername] = useState<string | null>(null);
  const {data : session} = useSession();
  const router = useRouter();
  useEffect(() => {
    if (!session) {
      router.push('/login');
    }
    setUsername(session?.user?.name || null);
  }, []);
  return (
    <header className="h-16 px-6 flex items-center justify-between bg-gray-900:0">
      <div className="text-lg font-semibold">Dashboard</div>
      <div className="text-sm">{username ? `Welcome back, ${username}!` : 'Welcome!'}</div>
    </header>
  );
}
