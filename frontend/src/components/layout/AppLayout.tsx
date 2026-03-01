import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { BottomNav } from './BottomNav'
import { TopBar } from './TopBar'

export function AppLayout() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900 flex">
      {/* Sidebar — fixed, visible on desktop via lg:translate-x-0 */}
      <Sidebar />

      {/* Spacer for desktop so content doesn't sit behind sidebar */}
      <div className="hidden lg:block w-64 shrink-0" />

      {/* Main area */}
      <div className="flex-1 flex flex-col min-h-screen min-w-0">
        <TopBar />
        {/* Mobile: p-3 pb-20 | sm: p-4 | lg: p-6 */}
        <main className="flex-1 p-3 pb-20 sm:p-4 sm:pb-20 lg:p-6 lg:pb-6 overflow-x-hidden">
          <Outlet />
        </main>
      </div>

      {/* Bottom navigation — mobile only */}
      <BottomNav />
    </div>
  )
}
