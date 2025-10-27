'use client';

import { useRouter } from "next/navigation";


export default function Header() {

  const router = useRouter();

  return (
    <header className="h-16 px-6 flex items-center justify-between bg-gray-900:0">
      <div className="text-lg font-semibold">search</div>
      {/* <div className="text-sm">{displayName ? `Welcome back, ${displayName}!` : 'Welcome!'}</div> */}
    <div onClick={() => router.push('/profile')} className="cursor-pointer">
      <div className=
      'flex items-center justify-center rounded-full overflow-hidden w-[50px] h-[50px] bg-[url(https://picsum.photos/seed/picsum/128/128)]'>
        <div className='flex items-center justify-center rounded-full overflow-hidden w-[50px] h-[50px] bg-gray-800/40'>
          <div className='rounded-full overflow-hidden w-[33px] h-[33px]'>
            <img src="https://picsum.photos/seed/picsum/64/64" alt="random photo" />
          </div>
        </div>
      </div>
    </div>
    </header>
  );
}
