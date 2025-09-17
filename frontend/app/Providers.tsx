'use client';

import { SessionProvider } from "next-auth/react";
import { PlayerProvider } from "./context/PlayerContext";
import ToasterClient from "@/app/components/ToasterClient";
import Player from "@/app/ui/player";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <PlayerProvider>
        {children}
        <div className="fixed block bottom-0 left-1/2 transform -translate-x-1/2 z-50 w-full max-w-2xl px-4">
          <Player />
        </div>
        <ToasterClient />
      </PlayerProvider>
    </SessionProvider>
  );
}
