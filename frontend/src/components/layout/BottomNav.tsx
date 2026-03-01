import { NavLink, useLocation } from 'react-router-dom'
import {
  LayoutDashboard, Wifi, Ticket, BarChart3, Settings,
} from 'lucide-react'
import { clsx } from 'clsx'

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Hotspot', href: '/hotspot', icon: Wifi },
  { name: 'Vouchers', href: '/vouchers', icon: Ticket },
  { name: 'Reports', href: '/reports', icon: BarChart3 },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export function BottomNav() {
  const location = useLocation()

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 lg:hidden bg-white/90 dark:bg-dark-800/90 backdrop-blur-md border-t border-gray-200 dark:border-dark-700">
      <div className="flex items-center justify-around h-14 sm:h-16 max-w-lg mx-auto px-1">
        {navigation.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname.startsWith(item.href)

          return (
            <NavLink
              key={item.name}
              to={item.href}
              className="bottom-nav-item flex-1"
            >
              <div className={clsx(
                'w-8 h-8 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center transition-all duration-200',
                isActive
                  ? 'bg-primary-500 shadow-glow scale-110'
                  : 'bg-gray-100 dark:bg-dark-700'
              )}>
                <Icon className={clsx('w-4 h-4 sm:w-5 sm:h-5', isActive ? 'text-white' : 'text-gray-500 dark:text-gray-400')} />
              </div>
              <span className={clsx(
                'text-[10px] sm:text-xs font-semibold mt-0.5',
                isActive ? 'text-primary-600 dark:text-primary-400' : 'text-gray-500 dark:text-gray-500'
              )}>
                {item.name}
              </span>
            </NavLink>
          )
        })}
      </div>
    </nav>
  )
}
