import { useMemo, useState, useEffect } from 'react'
import { Power, Search, X, Network } from 'lucide-react'
import { useMutation } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { clsx } from 'clsx'

import { Button, Badge, Pagination } from '../../../components/ui'
import { pppApi } from '../../../api/ppp'
import { usePPPActiveWS } from '../../../hooks/usePPPActiveWS'

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
  const [search, setSearch] = useState('')
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { connections, isConnected, lastUpdate } = usePPPActiveWS(routerId)

  const kickMutation = useMutation({
    mutationFn: (id: string) => pppApi.disconnectActive(routerId, id),
    onSuccess: () => toast.success('Connection disconnected'),
    onError: (err: any) => toast.error(err.message || 'Failed to disconnect'),
  })

  const filtered = useMemo(() =>
    connections.filter((c) =>
      c.name?.toLowerCase().includes(search.toLowerCase()) ||
      c.address?.toLowerCase().includes(search.toLowerCase()) ||
      c.callerID?.toLowerCase().includes(search.toLowerCase())
    ), [connections, search]
  )

  const pageCount = Math.ceil(filtered.length / pageSize)
  const paginatedRows = filtered.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [search])

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex items-center gap-3 w-full flex-wrap">
          <div className="relative flex-1 min-w-[160px] sm:max-w-72">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search name, IP, caller..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-9 pr-8 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {search && (
              <button onClick={() => setSearch('')} className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
          <div className="flex items-center gap-2 ml-auto">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              {filtered.length} connection{filtered.length !== 1 ? 's' : ''}
              {lastUpdate && ` — ${lastUpdate.toLocaleTimeString()}`}
            </p>
            <Badge variant={isConnected ? 'success' : 'danger'}>
              <span className={clsx('w-1.5 h-1.5 rounded-full mr-1.5', isConnected ? 'bg-white animate-pulse' : 'bg-white')} />
              {isConnected ? 'Live' : 'Reconnecting...'}
            </Badge>
          </div>
        </div>
      </div>

      {/* Table — scrollable */}
      <div className="flex-1 min-h-0 overflow-x-auto overflow-y-auto">
        <table className="w-full text-sm">
          <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-dark-700 border-b border-gray-200 dark:border-dark-700">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Name</th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Service</th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">IP Address</th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Caller ID</th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Uptime</th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Traffic</th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
            {paginatedRows.map((c) => (
              <tr key={c.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <div className="w-7 h-7 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center shrink-0">
                      <span className="text-xs font-bold text-primary-600 dark:text-primary-400">{c.name?.[0]?.toUpperCase()}</span>
                    </div>
                    <span className="font-semibold text-gray-900 dark:text-white text-sm">{c.name}</span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  {c.service ? <Badge variant="info">{c.service}</Badge> : <span className="text-gray-400 text-xs">—</span>}
                </td>
                <td className="px-4 py-3 font-mono text-xs text-gray-600 dark:text-gray-300">{c.address || '—'}</td>
                <td className="px-4 py-3 font-mono text-xs text-gray-500 dark:text-gray-400">{c.callerID || '—'}</td>
                <td className="px-4 py-3 text-xs text-gray-600 dark:text-gray-300">{c.uptime || '—'}</td>
                <td className="px-4 py-3 text-xs space-y-0.5">
                  <div className="text-success-600 dark:text-success-400">↓ {formatBytes(c.bytesIn || 0)}</div>
                  <div className="text-primary-600 dark:text-primary-400">↑ {formatBytes(c.bytesOut || 0)}</div>
                </td>
                <td className="px-4 py-3">
                  <Button
                    variant="ghost"
                    size="xs"
                    leftIcon={<Power className="w-3.5 h-3.5 text-danger-500" />}
                    onClick={() => { if (confirm(`Disconnect "${c.name}"?`)) kickMutation.mutate(c.id) }}
                  >
                    Kick
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {filtered.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400 flex flex-col items-center gap-3">
            <Network className="w-10 h-10 text-gray-300 dark:text-gray-600" />
            <span>{isConnected ? 'No active PPP connections' : 'Connecting to WebSocket...'}</span>
          </div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filtered.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />
    </div>
  )
}
