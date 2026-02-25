import { Sun, Moon } from 'lucide-react'
import { useThemeStore } from '../../stores/themeStore'

export function ThemeToggle() {
  const { isDark, toggleTheme } = useThemeStore()

  return (
    <button
      onClick={toggleTheme}
      className="relative inline-flex h-9 w-16 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 dark:focus:ring-offset-dark-800 bg-gray-200 dark:bg-dark-700"
      aria-label="Toggle theme"
    >
      <span className="sr-only">Toggle theme</span>
      <span
        className={`${
          isDark ? 'translate-x-8' : 'translate-x-1'
        } inline-flex h-7 w-7 transform items-center justify-center rounded-full bg-white shadow-lg transition-transform duration-200`}
      >
        {isDark ? (
          <Moon className="h-4 w-4 text-primary-600" />
        ) : (
          <Sun className="h-4 w-4 text-warning-500" />
        )}
      </span>
    </button>
  )
}
