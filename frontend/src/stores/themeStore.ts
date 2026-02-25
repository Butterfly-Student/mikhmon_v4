import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type ThemeMode = 'light' | 'dark' | 'system'
type ColorScheme = 'indigo' | 'pink' | 'emerald' | 'amber' | 'cyan' | 'purple'

interface ThemeState {
  mode: ThemeMode
  colorScheme: ColorScheme
  isDark: boolean
  
  // Actions
  setMode: (mode: ThemeMode) => void
  setColorScheme: (scheme: ColorScheme) => void
  toggleTheme: () => void
}

const getInitialDarkMode = (mode: ThemeMode): boolean => {
  if (mode === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return mode === 'dark'
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      mode: 'system',
      colorScheme: 'indigo',
      isDark: getInitialDarkMode('system'),

      setMode: (mode) => {
        const isDark = getInitialDarkMode(mode)
        set({ mode, isDark })
        
        // Apply dark class to document
        if (isDark) {
          document.documentElement.classList.add('dark')
        } else {
          document.documentElement.classList.remove('dark')
        }
      },

      setColorScheme: (colorScheme) => set({ colorScheme }),

      toggleTheme: () => {
        const newMode = get().isDark ? 'light' : 'dark'
        get().setMode(newMode)
      },
    }),
    {
      name: 'mikhmon-theme',
      onRehydrateStorage: () => (state) => {
        // Apply theme on rehydrate
        if (state) {
          const isDark = getInitialDarkMode(state.mode)
          if (isDark) {
            document.documentElement.classList.add('dark')
          } else {
            document.documentElement.classList.remove('dark')
          }
        }
      },
    }
  )
)
