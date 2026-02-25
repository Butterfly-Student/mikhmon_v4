import { useState, useEffect } from 'react'
import { Menu, Bell, Activity } from 'lucide-react'
import { format } from 'date-fns'
import { id } from 'date-fns/locale'
import { clsx } from 'clsx'

import { ThemeToggle } from '../common/ThemeToggle'
import { useAuthStore } from '../../stores/authStore'
import { useLayoutStore } from '../../stores/layoutStore'
import { useRouterStore } from '../../stores/routerStore'
import { usePingWebSocket } from '../../hooks/usePingWebSocket'

export function TopBar() {
  const user = useAuthStore((state) => state.user)
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const [currentTime, setCurrentTime] = useState(new Date())
  const { toggleSidebar } = useLayoutStore()

  // Real-time ping ke 8.8.8.8 — interval=1s, count=0 (infinite), size=64
  const { latency } = usePingWebSocket(selectedRouter?.id, {
    address: '8.8.8.8',
    interval: 1,
    count: 0,
    size: 64,
  })

  // Determine color based on latency
  const getLatencyColor = (ms: number | null) => {
    if (ms === null) return 'text-gray-400'
    if (ms < 50) return 'text-success-500'
    if (ms < 100) return 'text-warning-500'
    return 'text-danger-500'
  }

  const getLatencyBg = (ms: number | null) => {
    if (ms === null) return 'bg-gray-100 dark:bg-dark-700'
    if (ms < 50) return 'bg-success-50 dark:bg-success-900/20'
    if (ms < 100) return 'bg-warning-50 dark:bg-warning-900/20'
    return 'bg-danger-50 dark:bg-danger-900/20'
  }

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date())
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  return (
    <header className="sticky top-0 z-40 bg-white/80 dark:bg-dark-800/80 backdrop-blur-md border-b border-gray-200 dark:border-dark-700">
      <div className="flex items-center justify-between px-4 lg:px-6 py-3">
        {/* Left */}
        <div className="flex items-center gap-4">
          <button
            onClick={toggleSidebar}
            className="lg:hidden p-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700"
          >
            <Menu className="w-5 h-5" />
          </button>

          {/* Time Display */}
          <div className="hidden sm:flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
            <span className="font-medium">{format(currentTime, 'HH:mm:ss')}</span>
            <span className="text-gray-400">|</span>
            <span>{format(currentTime, 'EEEE, dd MMMM yyyy', { locale: id })}</span>
          </div>
        </div>

        {/* Right */}
        <div className="flex items-center gap-3">
          {/* Ping Display - 8.8.8.8 */}
          {selectedRouter && (
            <div className={clsx(
              "flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors",
              getLatencyBg(latency)
            )}>
              <Activity className={clsx("w-4 h-4", getLatencyColor(latency))} />
              <span className="text-xs text-gray-500 dark:text-gray-400 hidden sm:inline">Ping:</span>
              <span className={clsx("text-sm font-mono font-medium", getLatencyColor(latency))}>
                {latency !== null ? `${Math.round(latency)}ms` : '--'}
              </span>
              <span className="text-xs text-gray-400">8.8.8.8</span>
            </div>
          )}

          {/* Theme Toggle */}
          <div className="hidden sm:block">
            <ThemeToggle />
          </div>

          {/* Notifications */}
          <button className="relative p-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700">
            <Bell className="w-5 h-5" />
            <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-danger-500 rounded-full" />
          </button>

          {/* User */}
          <div className="flex items-center gap-3 pl-3 border-l border-gray-200 dark:border-dark-700">
            <div className="hidden md:block text-right">
              <p className="text-sm font-medium text-gray-900 dark:text-white">{user?.username}</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">Administrator</p>
            </div>
            <div className="w-9 h-9 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center text-white font-medium">
              {user?.username?.[0]?.toUpperCase() || 'A'}
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}
