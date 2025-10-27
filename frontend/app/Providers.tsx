'use client';

import { PlayerProvider } from "./context/PlayerContext";
import ToasterClient from "@/app/components/ToasterClient";
import Player from "@/app/ui/player";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <PlayerProvider>
      <div className="flex flex-col min-h-screen">
        <main className="flex-1">{children}</main>
        <footer className="fixed z-10 bottom-0 bg-gray-900/80 min-w-screen">
          <Player />
        </footer>
      </div>
      <ToasterClient />
    </PlayerProvider>
  );
}
