import { useState, useMemo, useEffect } from 'react'
import { Search, X } from 'lucide-react'
import { clsx } from 'clsx'

import { Badge, Pagination } from '../../../components/ui'
import { useHotspotInactiveWS } from '../../../hooks/useHotspotInactiveWS'

interface InactiveTabProps {
  routerId: string
}

export function InactiveTab({ routerId }: InactiveTabProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { users: inactiveUsers, isConnected, lastUpdate } = useHotspotInactiveWS(routerId)

  const filteredInactive = useMemo(() =>
    inactiveUsers.filter((u) =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (u.comment || '').toLowerCase().includes(searchQuery.toLowerCase())
    ), [inactiveUsers, searchQuery]
  )

  const pageCount = Math.ceil(filteredInactive.length / pageSize)
  const paginatedRows = filteredInactive.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [searchQuery])

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Total Inactive: <span className="font-medium text-gray-900 dark:text-white">{inactiveUsers.length}</span>
            {lastUpdate && ` — updated ${lastUpdate.toLocaleTimeString()}`}
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant={isConnected ? 'success' : 'danger'}>
              <span className={clsx('w-1.5 h-1.5 rounded-full mr-1.5', isConnected ? 'bg-white animate-pulse' : 'bg-white')} />
              {isConnected ? 'Live' : 'Reconnecting...'}
            </Badge>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Filter..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 pr-8 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
              {searchQuery && (
                <button onClick={() => setSearchQuery('')} className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                  <X className="w-4 h-4" />
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Table — scrollable */}
      <div className="flex-1 min-h-0 overflow-x-auto overflow-y-auto">
        <table className="w-full text-sm">
          <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-dark-700 border-b border-gray-200 dark:border-dark-700">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Server</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">User</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Profile</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">MAC Address</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Last Uptime</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Status</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Comment</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
            {paginatedRows.map((user) => (
              <tr key={user.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{user.server || 'all'}</td>
                <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{user.name}</td>
                <td className="px-4 py-3 text-sm">
                  <Badge variant="primary">{typeof user.profile === 'string' ? user.profile : user.profile?.name || '-'}</Badge>
                </td>
                <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{user.macAddress || '-'}</td>
                <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{user.uptime || '-'}</td>
                <td className="px-4 py-3 text-sm">
                  <Badge variant={user.disabled ? 'danger' : 'secondary'}>
                    {user.disabled ? 'Disabled' : 'Enabled'}
                  </Badge>
                </td>
                <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{user.comment || '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {filteredInactive.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">
            {isConnected ? 'No inactive users' : 'Connecting to WebSocket...'}
          </div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filteredInactive.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />
    </div>
  )
}
