// Hotspot Page with Tabs
import { useState } from 'react'
import { motion } from 'framer-motion'
import { Users, PieChart, Laptop, Wifi } from 'lucide-react'
import { clsx } from 'clsx'

import { useRouterStore } from '../../stores/routerStore'
import { UsersTab } from './tabs/UsersTab'
import { ProfilesTab } from './tabs/ProfilesTab'
import { ActiveTab } from './tabs/ActiveTab'
import { InactiveTab } from './tabs/InactiveTab'
import { HostsTab } from './tabs/HostsTab'

const tabs = [
  { id: 'users', name: 'Users', icon: Users },
  { id: 'profiles', name: 'User Profile', icon: PieChart },
  { id: 'active', name: 'Active', icon: Wifi },
  { id: 'inactive', name: 'Inactive', icon: Users },
  { id: 'hosts', name: 'Hosts', icon: Laptop },
]

export function HotspotPage() {
  const [activeTab, setActiveTab] = useState('users')
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '1')

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      <div className="rounded-2xl border border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800 shadow-card overflow-hidden flex flex-col h-[calc(100vh-148px)] sm:h-[calc(100vh-152px)] lg:h-[calc(100vh-104px)]">
        {/* Tab Bar — fixed */}
        <div className="flex items-center justify-between px-4 border-b border-gray-100 dark:border-dark-700 bg-gray-50/50 dark:bg-dark-800 flex-shrink-0">
          <div className="flex items-center">
            {tabs.map((tab) => {
              const Icon = tab.icon
              const isActive = activeTab === tab.id
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={clsx(
                    'relative flex items-center gap-1.5 px-3 sm:px-4 py-3.5 text-sm font-semibold transition-colors',
                    isActive
                      ? 'text-primary-600 dark:text-primary-400'
                      : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                  )}
                >
                  <Icon className="w-4 h-4 shrink-0" />
                  <span className="hidden sm:inline">{tab.name}</span>
                  {isActive && (
                    <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-500 rounded-full" />
                  )}
                </button>
              )
            })}
          </div>
          <div className="text-xs text-gray-400 dark:text-gray-500 font-medium hidden sm:block">
            {selectedRouter ? selectedRouter.name : 'No router selected'}
          </div>
        </div>

        {/* Tab Content — fills remaining height */}
        <div className="flex-1 min-h-0 flex flex-col overflow-hidden">
          {activeTab === 'users' && <UsersTab routerId={routerId} />}
          {activeTab === 'profiles' && <ProfilesTab routerId={routerId} />}
          {activeTab === 'active' && <ActiveTab routerId={routerId} />}
          {activeTab === 'inactive' && <InactiveTab routerId={routerId} />}
          {activeTab === 'hosts' && <HostsTab routerId={routerId} />}
        </div>
      </div>
    </motion.div>
  )
}
