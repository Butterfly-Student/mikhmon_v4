import { useState, useEffect, useRef } from 'react'
import { Menu, Bell, Activity, Clock, ChevronDown, Check, Router } from 'lucide-react'
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
  const { selectedRouter, routers, setSelectedRouter } = useRouterStore()
  const [currentTime, setCurrentTime] = useState(new Date())
  const { toggleSidebar } = useLayoutStore()
  const [routerDropOpen, setRouterDropOpen] = useState(false)
  const dropRef = useRef<HTMLDivElement>(null)

  const { latency } = usePingWebSocket(selectedRouter?.id, {
    address: '8.8.8.8',
    interval: 1,
    count: 0,
    size: 64,
  })

  const getLatencyColor = (ms: number | null) => {
    if (ms === null) return 'text-gray-400'
    if (ms < 50) return 'text-success-500'
    if (ms < 100) return 'text-warning-500'
    return 'text-danger-500'
  }

  const getLatencyBg = (ms: number | null) => {
    if (ms === null) return 'bg-gray-100 dark:bg-dark-700 border-gray-200 dark:border-dark-600'
    if (ms < 50) return 'bg-success-50 dark:bg-success-900/20 border-success-200 dark:border-success-800'
    if (ms < 100) return 'bg-warning-50 dark:bg-warning-900/20 border-warning-200 dark:border-warning-800'
    return 'bg-danger-50 dark:bg-danger-900/20 border-danger-200 dark:border-danger-800'
  }

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(new Date()), 1000)
    return () => clearInterval(timer)
  }, [])

  // Close dropdown when clicking outside
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (dropRef.current && !dropRef.current.contains(e.target as Node)) {
        setRouterDropOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  const initials = user?.username?.[0]?.toUpperCase() || 'A'

  return (
    <header className="sticky top-0 z-40 bg-white/90 dark:bg-dark-800/90 backdrop-blur-md border-b border-gray-200/70 dark:border-dark-700/70 h-14">
      <div className="flex items-center justify-between h-full px-3 sm:px-5 gap-2">
        {/* Left — hamburger + clock */}
        <div className="flex items-center gap-2 min-w-0">
          <button
            onClick={toggleSidebar}
            className="lg:hidden p-2 rounded-xl text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors shrink-0"
          >
            <Menu className="w-5 h-5" />
          </button>

          {/* Clock — always visible, even on mobile */}
          <div className="flex items-center gap-1.5">
            <Clock className="w-3.5 h-3.5 text-gray-400 shrink-0" />
            <div className="flex items-center gap-1 text-xs leading-none">
              <span className="font-bold text-gray-800 dark:text-gray-200 tabular-nums">
                {format(currentTime, 'HH:mm:ss')}
              </span>
              <span className="hidden xs:inline text-gray-300 dark:text-gray-600">•</span>
              <span className="hidden xs:inline text-gray-500 dark:text-gray-400 whitespace-nowrap">
                {format(currentTime, 'EEE dd/MM', { locale: id })}
              </span>
            </div>
          </div>
        </div>

        {/* Right */}
        <div className="flex items-center gap-1.5 shrink-0">
          {/* Ping badge — always visible */}
          {selectedRouter && (
            <div className={clsx(
              'flex items-center gap-1 px-2 py-1 rounded-lg border text-xs font-mono transition-colors',
              getLatencyBg(latency)
            )}>
              <Activity className={clsx('w-3 h-3 shrink-0', getLatencyColor(latency))} />
              <span className={clsx('font-semibold', getLatencyColor(latency))}>
                {latency !== null ? `${Math.round(latency)}ms` : '--'}
              </span>
            </div>
          )}

          {/* Router switcher dropdown */}
          <div ref={dropRef} className="relative">
            <button
              onClick={() => setRouterDropOpen(!routerDropOpen)}
              className={clsx(
                'flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border text-xs font-medium transition-all',
                selectedRouter
                  ? 'bg-primary-50 dark:bg-primary-900/20 border-primary-200 dark:border-primary-800 text-primary-700 dark:text-primary-300'
                  : 'bg-gray-100 dark:bg-dark-700 border-gray-200 dark:border-dark-600 text-gray-500 dark:text-gray-400'
              )}
              title="Switch Router"
            >
              {selectedRouter && (
                <div className="w-1.5 h-1.5 rounded-full bg-success-500 animate-pulse shrink-0" />
              )}
              <Router className="w-3.5 h-3.5 shrink-0" />
              <span className="hidden sm:inline truncate max-w-[90px]">
                {selectedRouter?.name || 'Select Router'}
              </span>
              <ChevronDown className={clsx('w-3 h-3 shrink-0 transition-transform', routerDropOpen && 'rotate-180')} />
            </button>

            {routerDropOpen && (
              <div className="absolute right-0 mt-2 w-52 bg-white dark:bg-dark-800 rounded-xl shadow-lg border border-gray-200 dark:border-dark-700 z-50 py-1 overflow-hidden">
                {routers.length === 0 ? (
                  <div className="px-4 py-3 text-xs text-gray-400 text-center">No routers configured</div>
                ) : (
                  routers.map((router) => (
                    <button
                      key={router.id}
                      onClick={() => { setSelectedRouter(router); setRouterDropOpen(false) }}
                      className="w-full px-3 py-2 flex items-center gap-2.5 text-sm hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors text-left"
                    >
                      <div className={clsx(
                        'w-1.5 h-1.5 rounded-full shrink-0',
                        selectedRouter?.id === router.id ? 'bg-success-500' : 'bg-gray-300 dark:bg-dark-600'
                      )} />
                      <span className="flex-1 truncate text-gray-800 dark:text-gray-200 font-medium">{router.name}</span>
                      {router.host && (
                        <span className="text-xs text-gray-400 font-mono shrink-0">{router.host}</span>
                      )}
                      {selectedRouter?.id === router.id && (
                        <Check className="w-3.5 h-3.5 text-primary-500 shrink-0" />
                      )}
                    </button>
                  ))
                )}
              </div>
            )}
          </div>

          {/* Theme toggle — hidden on very small screens */}
          <div className="hidden sm:block">
            <ThemeToggle />
          </div>

          {/* Notifications */}
          <button className="relative p-2 rounded-xl text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors">
            <Bell className="w-4 h-4" />
            <span className="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-danger-500 rounded-full" />
          </button>

          {/* User avatar */}
          <div className="flex items-center gap-2 pl-2 border-l border-gray-200 dark:border-dark-700">
            <div className="hidden md:block text-right">
              <p className="text-xs font-semibold text-gray-800 dark:text-white leading-tight">{user?.username}</p>
              <p className="text-[10px] text-gray-400 dark:text-gray-500 leading-tight">Admin</p>
            </div>
            <div className="w-8 h-8 rounded-xl bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center text-white text-sm font-bold shadow-sm">
              {initials}
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}
