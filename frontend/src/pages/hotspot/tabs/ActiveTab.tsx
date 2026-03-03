import { useState, useMemo, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Search, X, Power } from 'lucide-react'
import toast from 'react-hot-toast'
import { clsx } from 'clsx'

import { Button, Badge, Pagination } from '../../../components/ui'
import { hotspotApi } from '../../../api/hotspot'
import { useHotspotActiveWS } from '../../../hooks/useHotspotActiveWS'

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024, s = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + s[i]
}

interface ActiveTabProps {
  routerId: string
}

export function ActiveTab({ routerId }: ActiveTabProps) {
  const queryClient = useQueryClient()
  const [searchQuery, setSearchQuery] = useState('')
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { users: activeUsers, isConnected, lastUpdate } = useHotspotActiveWS(routerId)

  const kickMutation = useMutation({
    mutationFn: (id: string) => hotspotApi.deleteActiveSession(routerId, id),
    onSuccess: () => {
      toast.success('User disconnected')
      queryClient.invalidateQueries({ queryKey: ['active', routerId] })
    },
    onError: (err: any) => toast.error(err.message || 'Failed to disconnect'),
  })

  const filteredActive = useMemo(() =>
    activeUsers.filter((a) =>
      a.user.toLowerCase().includes(searchQuery.toLowerCase()) ||
      a.address.toLowerCase().includes(searchQuery.toLowerCase()) ||
      a.macAddress.toLowerCase().includes(searchQuery.toLowerCase())
    ), [activeUsers, searchQuery]
  )

  const pageCount = Math.ceil(filteredActive.length / pageSize)
  const paginatedRows = filteredActive.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [searchQuery])

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Total Active: <span className="font-medium text-gray-900 dark:text-white">{activeUsers.length}</span>
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
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Address</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">MAC Address</th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Uptime</th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Traffic</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Time Left</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Login By</th>
              <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
            {paginatedRows.map((active) => (
              <tr key={active.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{active.server}</td>
                <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{active.user}</td>
                <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{active.address}</td>
                <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{active.macAddress}</td>
                <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{active.uptime}</td>
                <td className="px-4 py-3 text-xs text-right space-y-0.5">
                  <div className="text-success-600 dark:text-success-400">↓ {formatBytes(active.bytesIn || 0)}</div>
                  <div className="text-primary-600 dark:text-primary-400">↑ {formatBytes(active.bytesOut || 0)}</div>
                </td>
                <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{active.sessionTimeLeft || '-'}</td>
                <td className="px-4 py-3 text-sm">
                  <Badge variant="success">{active.loginBy}</Badge>
                </td>
                <td className="px-4 py-3 text-center">
                  <Button
                    variant="ghost"
                    size="xs"
                    onClick={() => { if (confirm(`Disconnect ${active.user}?`)) kickMutation.mutate(active.id) }}
                  >
                    <Power className="w-3.5 h-3.5 text-danger-500" />
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {filteredActive.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">
            {isConnected ? 'No active users' : 'Connecting to WebSocket...'}
          </div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filteredActive.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />
    </div>
  )
}
