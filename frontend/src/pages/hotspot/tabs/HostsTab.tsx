import { useState, useMemo, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Search, X, RefreshCw, Edit2, Trash2 } from 'lucide-react'
import { clsx } from 'clsx'

import { Pagination } from '../../../components/ui'
import { hotspotApi } from '../../../api/hotspot'

interface HostsTabProps {
  routerId: string
}

export function HostsTab({ routerId }: HostsTabProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [filterType, setFilterType] = useState<'all' | 'authorized' | 'bypassed'>('all')
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { data: hosts = [], isLoading, refetch } = useQuery({
    queryKey: ['hosts', routerId],
    queryFn: () => hotspotApi.getHosts(routerId),
  })

  const filteredHosts = useMemo(() => hosts.filter((host) => {
    const matchesSearch =
      host.macAddress.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (host.address || '').toLowerCase().includes(searchQuery.toLowerCase())
    const matchesFilter =
      filterType === 'all' ||
      (filterType === 'authorized' && host.authorized) ||
      (filterType === 'bypassed' && host.bypassed)
    return matchesSearch && matchesFilter
  }), [hosts, searchQuery, filterType])

  const pageCount = Math.ceil(filteredHosts.length / pageSize)
  const paginatedRows = filteredHosts.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [searchQuery, filterType])

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Total Hosts: <span className="font-medium text-gray-900 dark:text-white">{hosts.length}</span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <button onClick={() => refetch()} className="p-1.5 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700" title="Refresh">
              <RefreshCw className="w-4 h-4" />
            </button>
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
            <div className="flex items-center gap-1">
              {([
                { key: 'all', label: 'All', activeClass: 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400' },
                { key: 'authorized', label: 'Auth', activeClass: 'bg-success-100 dark:bg-success-900/30 text-success-600 dark:text-success-400' },
                { key: 'bypassed', label: 'Pass', activeClass: 'bg-warning-100 dark:bg-warning-900/30 text-warning-600 dark:text-warning-400' },
              ] as const).map(({ key, label, activeClass }) => (
                <button
                  key={key}
                  onClick={() => setFilterType(key)}
                  className={clsx(
                    'px-3 py-1.5 text-xs font-medium rounded-lg transition-colors',
                    filterType === key ? activeClass : 'bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600'
                  )}
                >
                  {label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Table — scrollable */}
      <div className="flex-1 min-h-0 overflow-x-auto overflow-y-auto">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-dark-700 border-b border-gray-200 dark:border-dark-700">
              <tr>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase w-10"></th>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase w-10"></th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">MAC Address</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Address</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">To Address</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Server</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Comment</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
              {paginatedRows.map((host) => (
                <tr key={host.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                  <td className="px-4 py-3 text-center">
                    <button className="text-gray-400 hover:text-gray-600"><Edit2 className="w-4 h-4" /></button>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <button className="text-gray-400 hover:text-gray-600"><Trash2 className="w-4 h-4" /></button>
                  </td>
                  <td className="px-4 py-3 text-sm font-mono text-gray-900 dark:text-white">{host.macAddress}</td>
                  <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{host.address}</td>
                  <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{host.toAddress || '-'}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{host.server || '-'}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{host.comment || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {!isLoading && filteredHosts.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">No hosts found</div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filteredHosts.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />
    </div>
  )
}
