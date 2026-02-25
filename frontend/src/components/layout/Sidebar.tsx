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
  X
} from 'lucide-react'
import { clsx } from 'clsx'
import { useAuthStore } from '../../stores/authStore'
import { useLayoutStore } from '../../stores/layoutStore'

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Hotspot', href: '/hotspot', icon: Wifi },
  { name: 'Vouchers', href: '/vouchers', icon: Ticket },
  { name: 'Reports', href: '/reports', icon: BarChart3 },
  { name: 'Routers', href: '/routers', icon: Router },
  { name: 'Settings', href: '/settings', icon: Settings },
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
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <aside
        className={clsx(
          'fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-dark-800 border-r border-gray-200/80 dark:border-dark-700 transform transition-transform duration-300 lg:translate-x-0',
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        {/* Logo */}
        <div className="flex items-center justify-between px-6 py-5 border-b border-gray-200/80 dark:border-dark-700">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
              <WifiIcon className="w-5 h-5 text-white" />
            </div>
            <div>
              <h1 className="text-lg font-bold text-gray-900 dark:text-white">Mikhmon v4</h1>
              <p className="text-xs text-gray-500 dark:text-gray-400">Hotspot Manager</p>
            </div>
          </div>
          <button
            onClick={() => setSidebarOpen(false)}
            className="lg:hidden p-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-3 py-6 space-y-1 overflow-y-auto">
          {navigation.map((item) => {
            const Icon = item.icon
            const isActive = location.pathname.startsWith(item.href)
            
            return (
              <NavLink
                key={item.name}
                to={item.href}
                onClick={() => setSidebarOpen(false)}
                className={clsx(
                  'sidebar-link',
                  isActive && 'sidebar-link-active'
                )}
              >
                <Icon className={clsx('w-5 h-5 mr-3', isActive ? 'text-current' : 'text-gray-400')} />
                {item.name}
              </NavLink>
            )
          })}
        </nav>

        {/* Logout */}
        <div className="p-3 border-t border-gray-200/80 dark:border-dark-700">
          <button
            onClick={logout}
            className="flex w-full items-center px-4 py-3 text-sm font-medium text-danger-600 dark:text-danger-400 rounded-lg hover:bg-danger-50 dark:hover:bg-danger-900/20 transition-colors"
          >
            <LogOut className="w-5 h-5 mr-3" />
            Logout
          </button>
        </div>
      </aside>
    </>
  )
}
