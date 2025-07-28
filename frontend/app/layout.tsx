import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import ToasterClient from '@/app/components/ToasterClient';
import { PlayerProvider } from './context/PlayerContext';
import Player from './ui/player';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: 'BypassBeats',
  description: 'A free background music listening platform',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <PlayerProvider>
          {children}
          <div className="fixed block bottom-0 left-1/2 transform -translate-x-1/2 z-50 w-full max-w-2xl px-4">
            <Player />
          </div>
          <ToasterClient />
        </PlayerProvider>
      </body>
    </html>
  );
}
