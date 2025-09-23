import Sidebar from './sidebar';
import Header from './header';
import UpNext from './upNext';
import RightPanel from '../components/rightPanel';

const PLAYER_HEIGHT = 64; // h-16 ~ 64px

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex flex-col h-screen bg-gray-950:80 text-white">
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <div className="flex flex-col flex-1">
          <Header />
          <div className="flex flex-1 overflow-hidden">
            <main className="flex-1 overflow-y-auto"
             style={{ maxHeight: `calc(100vh - 64px - ${PLAYER_HEIGHT}px)` }}>{children}</main>
            <div className="w-80 bg-gray-900:80 p-4 overflow-y-auto">
              <RightPanel />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
