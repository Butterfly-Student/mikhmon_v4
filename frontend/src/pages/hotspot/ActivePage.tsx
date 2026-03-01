import { useMemo } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import { Wifi, Power, Globe, Smartphone, RefreshCw } from 'lucide-react'
import toast from 'react-hot-toast'
import type { ColumnDef } from '@tanstack/react-table'

import { Card, Button, Badge, DataTable } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'

type ActiveUser = { id: string; user: string; address: string; macAddress: string; uptime: string; bytesIn: number; bytesOut: number; sessionTimeLeft?: string; loginBy: string }

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024, s = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + s[i]
}

export function ActivePage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '')

  const { data: activeUsers, isLoading } = useQuery({
    queryKey: ['active', routerId],
    queryFn: () => hotspotApi.getActive(routerId),
    refetchInterval: 5000,
    enabled: !!selectedRouter,
  })

  const columns = useMemo<ColumnDef<ActiveUser, any>[]>(() => [
    {
      accessorKey: 'user',
      header: 'User',
      cell: ({ getValue }) => (
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center shrink-0">
            <span className="text-xs font-bold text-primary-600 dark:text-primary-400">
              {(getValue() as string)[0]?.toUpperCase()}
            </span>
          </div>
          <span className="font-semibold text-gray-900 dark:text-white text-sm">{getValue()}</span>
        </div>
      ),
    },
    {
      accessorKey: 'address',
      header: 'IP Address',
      cell: ({ getValue }) => <span className="font-mono text-xs text-gray-600 dark:text-gray-300">{getValue()}</span>,
    },
    {
      accessorKey: 'macAddress',
      header: 'MAC Address',
      cell: ({ getValue }) => <span className="font-mono text-xs text-gray-500 dark:text-gray-400">{getValue()}</span>,
    },
    {
      accessorKey: 'uptime',
      header: 'Uptime',
      cell: ({ getValue }) => <span className="text-xs text-gray-600 dark:text-gray-300">{getValue()}</span>,
    },
    {
      id: 'bytes',
      header: 'Traffic',
      accessorFn: (r) => r.bytesIn + r.bytesOut,
      cell: ({ row }) => (
        <div className="text-xs space-y-0.5">
          <div className="text-success-600 dark:text-success-400">↓ {formatBytes(row.original.bytesIn)}</div>
          <div className="text-primary-600 dark:text-primary-400">↑ {formatBytes(row.original.bytesOut)}</div>
        </div>
      ),
    },
    {
      accessorKey: 'sessionTimeLeft',
      header: 'Time Left',
      cell: ({ getValue }) => {
        const v = getValue()
        return v ? <Badge variant="warning">{v}</Badge> : <span className="text-gray-400">∞</span>
      },
    },
    {
      accessorKey: 'loginBy',
      header: 'Login By',
      cell: ({ getValue }) => {
        const v = (getValue() as string)?.toLowerCase()
        return (
          <div className="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300">
            {v === 'mac' ? <Smartphone className="w-3.5 h-3.5" /> : <Globe className="w-3.5 h-3.5" />}
            <span className="capitalize">{getValue()}</span>
          </div>
        )
      },
    },
    {
      id: 'actions',
      header: '',
      enableSorting: false,
      cell: ({ row }) => (
        <Button
          variant="ghost"
          size="xs"
          leftIcon={<Power className="w-3.5 h-3.5 text-danger-500" />}
          onClick={() => toast.success(`Kick ${row.original.user} (demo)`)}
        >
          Kick
        </Button>
      ),
    },
  ], [])

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <Wifi className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">Silahkan pilih router untuk melihat active sessions.</p>
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
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Active Sessions</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {activeUsers?.length || 0} users currently connected &mdash; auto-refreshes every 5s
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="success">
            <span className="w-1.5 h-1.5 bg-white rounded-full mr-1.5 animate-pulse" />
            Live
          </Badge>
          <Button variant="ghost" size="sm" leftIcon={<RefreshCw className="w-4 h-4" />}
            onClick={() => queryClient.invalidateQueries({ queryKey: ['active', routerId] })}>
            Refresh
          </Button>
        </div>
      </div>

      <Card>
        <Card.Body>
          <DataTable
            data={(activeUsers as ActiveUser[]) || []}
            columns={columns}
            isLoading={isLoading}
            searchPlaceholder="Search by user, IP, MAC..."
            emptyMessage="No active sessions"
            emptyIcon={<Wifi className="w-10 h-10 text-gray-300 dark:text-gray-600" />}
          />
        </Card.Body>
      </Card>
    </motion.div>
  )
}

