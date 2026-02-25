import { NavLink, useLocation } from 'react-router-dom'
import { LayoutDashboard, Wifi, Ticket, BarChart3, Settings } from 'lucide-react'
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
    <nav className="bottom-nav">
      <div className="flex items-center justify-around max-w-lg mx-auto">
        {navigation.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname.startsWith(item.href)
          
          return (
            <NavLink
              key={item.name}
              to={item.href}
              className={clsx(
                'bottom-nav-item',
                isActive && 'bottom-nav-item-active'
              )}
            >
              <Icon className={clsx('w-5 h-5 mb-1', isActive && 'text-current')} />
              <span>{item.name}</span>
            </NavLink>
          )
        })}
      </div>
    </nav>
  )
}
