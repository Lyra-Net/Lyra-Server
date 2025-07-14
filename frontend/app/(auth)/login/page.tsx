'use client';

import { getDeviceId } from '@/utils/device';
import axios from 'axios';
import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

export default function LoginPage() {
  const [showForm, setShowForm] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isPending, setIsPending] = useState(false);

  const router = useRouter();
  useEffect(() => {
    setTimeout(() => setShowForm(true), 100);
  }, []);

  const handleLogin = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (isPending) return;

      setIsPending(true);
      const device_id = getDeviceId();

      try {
        const res = await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/login`, {
          username,
          password,
          device_id,
        });

        const { access_token } = res.data;

        localStorage.setItem('access_token', access_token);
        localStorage.setItem('username', username);
        toast.success('Logged in successfuly!');
        router.push('/dashboard');
      } catch (err: any) {
        toast.error(err.response?.data?.message || 'Login failure');
        console.error('Login failed', err.response?.data || err.message);
      } finally {
        setIsPending(false);
      }
    },
    [username, password, isPending]
  );

  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-transparent">
      <div
        className={`
          w-full max-w-md
          p-8 rounded-xl shadow-xl
          bg-white/50 dark:bg-gray-900/50
          text-gray-800 dark:text-gray-100
          backdrop-blur-md
          transform transition-all duration-700 ease-out
          ${showForm ? 'translate-y-0 opacity-100' : '-translate-y-32 opacity-0'}
        `}
      >
        <h2 className="text-2xl font-bold mb-6 text-center">Login</h2>
        <form onSubmit={handleLogin}>
          <div className="relative z-0 w-full mb-6 group">
            <input
              type="text"
              name="username"
              id="username"
              value={username}
              onChange={e => setUsername(e.target.value)}
              disabled={isPending}
              className="block py-2.5 px-0 w-full text-sm text-gray-900 dark:text-white bg-transparent border-0 border-b-2 border-gray-300 dark:border-gray-600 appearance-none focus:outline-none focus:ring-0 focus:border-blue-600 peer"
              placeholder=" "
              required
            />
            <label
              htmlFor="username"
              className="absolute text-sm text-gray-500 dark:text-gray-400 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-placeholder-shown:scale-100 
              peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6"
            >
              Username
            </label>
          </div>
          <div className="relative z-0 w-full mb-6 group">
            <input
              type="password"
              name="password"
              id="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              disabled={isPending}
              className="block py-2.5 px-0 w-full text-sm text-gray-900 dark:text-white bg-transparent border-0 border-b-2 border-gray-300 dark:border-gray-600 appearance-none focus:outline-none focus:ring-0 focus:border-blue-600 peer"
              placeholder=" "
              required
            />
            <label
              htmlFor="password"
              className="absolute text-sm text-gray-500 dark:text-gray-400 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-placeholder-shown:scale-100 
              peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6"
            >
              Password
            </label>
          </div>
          <button
            type="submit"
            disabled={isPending}
            className={`w-full py-2 rounded-md transition text-white ${
              isPending ? 'bg-blue-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'
            }`}
          >
            {isPending ? 'Signing in...' : 'Sign In'}
          </button>
        </form>
      </div>
    </div>
  );
}
