import { create } from "zustand";

interface TwoFAState {
  sessionId: string | null;
  setSessionId: (sessionId: string) => void;
  clearSessionId: () => void;
}

export const use2FASessionStore = create<TwoFAState>(( set) => ({
  sessionId: null,
  setSessionId: (sessionId) => set({ sessionId }),
  clearSessionId: () => set({ sessionId: null }),
}));