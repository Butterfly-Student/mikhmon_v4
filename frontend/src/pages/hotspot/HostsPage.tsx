import { useState, useMemo } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import { Laptop, CheckCircle, XCircle, Shield, RefreshCw } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'
import { clsx } from 'clsx'

import { Card, Button, Badge, DataTable } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'

type Host = { id: string; macAddress: string; address?: string; authorized: boolean; bypassed: boolean; blocked: boolean; server?: string }

export function HostsPage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '')
  const [filter, setFilter] = useState<'all' | 'authorized' | 'bypassed'>('all')

  const { data: hosts, isLoading } = useQuery({
    queryKey: ['hosts', routerId],
    queryFn: () => hotspotApi.getHosts(routerId),
    refetchInterval: 10000,
    enabled: !!selectedRouter,
  })

  const filteredHosts = useMemo(() => {
    if (!hosts) return []
    if (filter === 'authorized') return hosts.filter((h: Host) => h.authorized)
    if (filter === 'bypassed') return hosts.filter((h: Host) => h.bypassed)
    return hosts
  }, [hosts, filter])

  const columns = useMemo<ColumnDef<Host, any>[]>(() => [
    {
      accessorKey: 'macAddress',
      header: 'Device (MAC)',
      cell: ({ row }) => {
        const h = row.original
        return (
          <div className="flex items-center gap-3">
            <div className={clsx(
              'w-8 h-8 rounded-lg flex items-center justify-center shrink-0',
              h.authorized ? 'bg-success-100 dark:bg-success-900/30' :
                h.bypassed ? 'bg-warning-100 dark:bg-warning-900/30' :
                  'bg-gray-100 dark:bg-dark-700'
            )}>
              <Laptop className={clsx(
                'w-4 h-4',
                h.authorized ? 'text-success-600' :
                  h.bypassed ? 'text-warning-600' :
                    'text-gray-400'
              )} />
            </div>
            <span className="font-mono text-xs text-gray-900 dark:text-white">{h.macAddress}</span>
          </div>
        )
      },
    },
    {
      accessorKey: 'address',
      header: 'IP Address',
      cell: ({ getValue }) => (
        <span className="font-mono text-xs text-gray-600 dark:text-gray-300">{getValue() || '-'}</span>
      ),
    },
    {
      id: 'status',
      header: 'Status',
      accessorFn: (r) => r.authorized ? 'authorized' : r.bypassed ? 'bypassed' : r.blocked ? 'blocked' : 'unknown',
      cell: ({ row }) => {
        const h = row.original
        return (
          <div className="flex flex-wrap gap-1">
            {h.authorized && <Badge variant="success"><CheckCircle className="w-3 h-3 mr-1" />Authorized</Badge>}
            {h.bypassed && <Badge variant="warning"><Shield className="w-3 h-3 mr-1" />Bypassed</Badge>}
            {h.blocked && <Badge variant="danger"><XCircle className="w-3 h-3 mr-1" />Blocked</Badge>}
            {!h.authorized && !h.bypassed && !h.blocked && <Badge>Unknown</Badge>}
          </div>
        )
      },
    },
    {
      accessorKey: 'server',
      header: 'Server',
      cell: ({ getValue }) => (
        <span className="text-xs text-gray-500 dark:text-gray-400">{getValue() || '-'}</span>
      ),
    },
  ], [])

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <Laptop className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">Silahkan pilih router untuk melihat data hosts.</p>
        <Link to="/routers" className="px-5 py-2.5 rounded-xl bg-primary-500 text-white font-medium hover:bg-primary-600 transition-colors">
          Manage Routers
        </Link>
      </div>
    )
  }

  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="space-y-4">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Hotspot Hosts</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {filteredHosts.length} devices on network
          </p>
        </div>
        <Button variant="ghost" size="sm" leftIcon={<RefreshCw className="w-4 h-4" />}
          onClick={() => queryClient.invalidateQueries({ queryKey: ['hosts', routerId] })}>
          Refresh
        </Button>
      </div>

      {/* Filter chips */}
      <div className="flex gap-2">
        {(['all', 'authorized', 'bypassed'] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={clsx(
              'px-3.5 py-1.5 rounded-xl text-sm font-medium transition-colors',
              filter === f
                ? 'bg-primary-500 text-white shadow-sm'
                : 'bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600'
            )}
          >
            {f.charAt(0).toUpperCase() + f.slice(1)}
          </button>
        ))}
      </div>

      <Card>
        <Card.Body>
          <DataTable
            data={filteredHosts as Host[]}
            columns={columns}
            isLoading={isLoading}
            searchPlaceholder="Search by MAC, IP..."
            emptyMessage="No hosts found"
            emptyIcon={<Laptop className="w-10 h-10 text-gray-300 dark:text-gray-600" />}
          />
        </Card.Body>
      </Card>
    </motion.div>
  )
}

