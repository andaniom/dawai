import { create } from "zustand";
import { User } from "./types";

interface AuthStore {
  user: User | null;
  accessToken: string | null;
  schoolId: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, token: string, schoolId: string) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  accessToken: null,
  schoolId: null,
  isAuthenticated: false,
  setAuth: (user: User, token: string, schoolId: string) =>
    set({
      user,
      accessToken: token,
      schoolId,
      isAuthenticated: true,
    }),
  clearAuth: () =>
    set({
      user: null,
      accessToken: null,
      schoolId: null,
      isAuthenticated: false,
    }),
}));
