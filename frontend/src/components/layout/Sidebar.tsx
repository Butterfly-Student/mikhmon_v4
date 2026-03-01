import { NavLink, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  Wifi,
  Ticket,
  BarChart3,
  Router,
  Settings,
  LogOut,
  WifiIcon,
  X,
  ChevronRight,
  ChevronDown,
} from 'lucide-react'
import { clsx } from 'clsx'
import { useAuthStore } from '../../stores/authStore'
import { useLayoutStore } from '../../stores/layoutStore'
import { useRouterStore } from '../../stores/routerStore'
import { useState } from 'react'

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard, color: 'from-primary-500 to-indigo-600' },
  { name: 'Hotspot', href: '/hotspot', icon: Wifi, color: 'from-cyan-500 to-blue-600' },
  { name: 'Vouchers', href: '/vouchers', icon: Ticket, color: 'from-secondary-500 to-pink-600' },
  { name: 'Reports', href: '/reports', icon: BarChart3, color: 'from-success-500 to-teal-600' },
  { name: 'Routers', href: '/routers', icon: Router, color: 'from-warning-500 to-orange-600' },
  { name: 'Settings', href: '/settings', icon: Settings, color: 'from-gray-500 to-gray-600' },
]

export function Sidebar() {
  const location = useLocation()
  const logout = useAuthStore((state) => state.logout)
  const { sidebarOpen, setSidebarOpen } = useLayoutStore()

  return (
    <>
      {/* Mobile Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <aside
        className={clsx(
          'fixed inset-y-0 left-0 z-50 w-64 flex flex-col',
          'bg-white dark:bg-dark-800',
          'border-r border-gray-100 dark:border-dark-700',
          'shadow-xl lg:shadow-none',
          'transform transition-transform duration-300 ease-in-out lg:translate-x-0',
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        {/* Brand */}
        <div className="flex items-center justify-between px-3 py-3 sm:px-5 sm:py-4 border-b border-gray-100 dark:border-dark-700">
          <div className="flex items-center gap-2 sm:gap-3">
            <div className="w-8 h-8 sm:w-9 sm:h-9 rounded-xl bg-gradient-to-br from-primary-500 via-purple-500 to-secondary-500 flex items-center justify-center shadow-glow">
              <WifiIcon className="w-4 h-4 sm:w-5 sm:h-5 text-white drop-shadow" />
            </div>
            <div>
              <h1 className="text-sm sm:text-base font-extrabold text-gray-900 dark:text-white tracking-tight">
                Mikhmon<span className="text-primary-500">v4</span>
              </h1>
              <p className="text-[9px] sm:text-[10px] text-gray-400 dark:text-gray-500 uppercase tracking-widest">Hotspot Manager</p>
            </div>
          </div>
          <button
            onClick={() => setSidebarOpen(false)}
            className="lg:hidden p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700 hover:text-gray-600 dark:hover:text-gray-200 transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-2 py-3 sm:px-3 sm:py-4 space-y-0.5 overflow-y-auto scrollbar-hide">
          <p className="px-2 sm:px-3 mb-2 text-[9px] sm:text-[10px] font-semibold text-gray-400 dark:text-gray-600 uppercase tracking-widest">
            Main Menu
          </p>
          {navigation.map((item) => {
            const Icon = item.icon
            const isActive = location.pathname.startsWith(item.href)

            return (
              <NavLink
                key={item.name}
                to={item.href}
                onClick={() => setSidebarOpen(false)}
                className={clsx('sidebar-link group', isActive && 'sidebar-link-active')}
              >
                <div className={clsx(
                  'w-7 h-7 rounded-lg flex items-center justify-center shrink-0 transition-all duration-200',
                  isActive
                    ? `bg-gradient-to-br ${item.color} shadow-sm`
                    : 'bg-gray-100 dark:bg-dark-700 group-hover:bg-gray-200 dark:group-hover:bg-dark-600'
                )}>
                  <Icon className={clsx('w-4 h-4', isActive ? 'text-white' : 'text-gray-500 dark:text-gray-400')} />
                </div>
                <span className="flex-1">{item.name}</span>
                {isActive && <ChevronRight className="w-3.5 h-3.5 text-primary-400 opacity-60" />}
              </NavLink>
            )
          })}
        </nav>

        {/* Router Switcher */}
        <RouterSwitcher onSwitch={() => setSidebarOpen(false)} />

        {/* User + Logout */}
        <div className="p-3 border-t border-gray-100 dark:border-dark-700 space-y-1">
          <button
            onClick={logout}
            className="flex w-full items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-xl text-danger-600 dark:text-danger-400 hover:bg-danger-50 dark:hover:bg-danger-900/20 transition-colors"
          >
            <div className="w-7 h-7 rounded-lg bg-danger-100 dark:bg-danger-900/30 flex items-center justify-center">
              <LogOut className="w-4 h-4" />
            </div>
            Logout
          </button>
        </div>
      </aside>
    </>
  )
}

// ─── Router Switcher Sub-component ───────────────────────────────────────────
function RouterSwitcher({ onSwitch }: { onSwitch: () => void }) {
  const { routers, selectedRouter, setSelectedRouter } = useRouterStore()
  const [open, setOpen] = useState(true)

  if (routers.length === 0) return null

  return (
    <div className="mx-3 mb-2">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold text-gray-400 dark:text-gray-600 uppercase tracking-widest hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors"
      >
        <div className="flex items-center gap-1.5">
          <Router className="w-3.5 h-3.5" />
          Routers
        </div>
        <ChevronDown className={clsx('w-3.5 h-3.5 transition-transform duration-200', open && 'rotate-180')} />
      </button>

      {open && (
        <div className="mt-1 space-y-0.5">
          {routers.map((router) => {
            const isActive = selectedRouter?.id === router.id
            return (
              <button
                key={router.id}
                onClick={() => { setSelectedRouter(router); onSwitch() }}
                className={clsx(
                  'w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-sm font-medium transition-all',
                  isActive
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300'
                    : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700'
                )}
              >
                <div className={clsx(
                  'w-2 h-2 rounded-full shrink-0 transition-colors',
                  isActive ? 'bg-success-500' : 'bg-gray-300 dark:bg-dark-600'
                )} />
                <span className="flex-1 truncate text-left">{router.name}</span>
                {router.host && (
                  <span className="text-[10px] font-mono text-gray-400 dark:text-gray-600 truncate max-w-[80px]">
                    {router.host}
                  </span>
                )}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
