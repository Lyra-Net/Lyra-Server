import React from "react";
import { PlayerProvider } from "../context/PlayerContext";
import Player from "../ui/player";

export default function PlayerWapperLayout({ children }: { children:React.ReactNode }) {
  return (
    <PlayerProvider>
      {children}
      <div className="fixed bottom-0 left-0 right-0 z-50">
        <Player />
      </div>
    </PlayerProvider>
  )
}