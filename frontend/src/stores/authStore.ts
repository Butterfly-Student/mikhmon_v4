import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { UserInfo } from '../types'

interface AuthState {
  user: UserInfo | null
  token: string | null
  isAuthenticated: boolean

  // Actions
  setAuth: (token: string, user: UserInfo) => void
  logout: () => void
  setUser: (user: UserInfo) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      setAuth: (token, user) => {
        if (import.meta.env.DEV) {
          console.log('[AuthStore] setAuth called')
          console.log('[AuthStore] Token:', token ? `***${token.slice(-4)}` : 'none')
          console.log('[AuthStore] User:', user?.username)
        }
        set({
          token,
          user,
          isAuthenticated: true,
        })
      },

      logout: () => {
        if (import.meta.env.DEV) {
          console.log('[AuthStore] logout called')
        }
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        })
      },

      setUser: (user) => set({ user }),
    }),
    {
      name: 'mikhmon-auth',
      onRehydrateStorage: () => (state) => {
        if (import.meta.env.DEV && state) {
          console.log('[AuthStore] Storage rehydrated')
          console.log('[AuthStore] Is authenticated:', state.isAuthenticated)
          console.log('[AuthStore] Token present:', !!state.token)
          console.log('[AuthStore] User:', state.user?.username)
        }
      },
    }
  )
)
