'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Music, Home, List, Share2, Settings } from 'lucide-react';

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-gray-900/0 p-4">
      <h2 className="text-xl font-bold mb-6">🎵 BypassBeats</h2>
      <nav className="space-y-3">
        <NavItem icon={<Home size={18} />} label="Home" href="/" pathname={pathname} />
        <NavItem icon={<Music size={18} />} label="Search" href="/search" pathname={pathname} />
        <NavItem
          icon={<List size={18} />}
          label="My Playlists"
          href="/playlists"
          pathname={pathname}
        />
        <NavItem
          icon={<Share2 size={18} />}
          label="Trend Playlists"
          href="/trends"
          pathname={pathname}
        />
        <NavItem
          icon={<Settings size={18} />}
          label="Settings"
          href="/settings"
          pathname={pathname}
        />
      </nav>
    </aside>
  );
}

function NavItem({
  icon,
  label,
  href,
  pathname,
}: {
  icon: React.ReactNode;
  label: string;
  href: string;
  pathname: string;
}) {
  const isActive = pathname === href;

  return (
    <Link href={href}>
      <div
        className={`flex items-center space-x-3 px-3 py-2 rounded-md cursor-pointer transition
          ${
            isActive
              ? 'bg-gray-800:20 text-white font-semibold'
              : 'text-gray-300 hover:text-white hover:bg-gray-800:20'
          }`}
      >
        {icon}
        <span>{label}</span>
      </div>
    </Link>
  );
}
