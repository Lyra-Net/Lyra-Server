import { create } from "zustand";


interface AuthState {
  userId: string | null;
  username: string | null;
  email: string | null;
  displayName: string | null;
  accessToken: string | null;
  exp: number | null;
  setAuth: (userId: string, username: string, email: string, displayName: string, accessToken: string, exp: number) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>(( set) => ({
  userId: null,
  username: null,
  email: null,
  displayName: null,
  accessToken: null,
  exp: null,
  setAuth: (userId, username, email, displayName, accessToken, exp) => set({ userId, username, email, displayName, accessToken, exp }),
  clearAuth: () => set({ userId: null, displayName: null, accessToken: null, exp: null }),
}));
