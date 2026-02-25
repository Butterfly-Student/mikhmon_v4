// Active Page
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import {
  Wifi,
  Power,
  Globe,
  Smartphone,
  RefreshCw,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { Card, Button, Badge } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'

export function ActivePage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '')

  const { data: activeUsers, isLoading } = useQuery({
    queryKey: ['active', routerId],
    queryFn: () => hotspotApi.getActive(routerId),
    refetchInterval: 5000, // Refresh every 5 seconds
    enabled: !!selectedRouter,
  })

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const getLoginByIcon = (loginBy: string) => {
    switch (loginBy?.toLowerCase()) {
      case 'mac':
        return <Smartphone className="w-4 h-4" />
      default:
        return <Globe className="w-4 h-4" />
    }
  }

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <Wifi className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">
          Silahkan pilih router terlebih dahulu untuk melihat active sessions.
        </p>
        <Link
          to="/routers"
          className="px-5 py-2.5 rounded-xl bg-primary-500 text-white font-medium hover:bg-primary-600 transition-colors"
        >
          Manage Routers
        </Link>
      </div>
    )
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Active Sessions</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Currently connected hotspot users
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="success" className="animate-pulse">
            <span className="w-2 h-2 bg-white rounded-full mr-2" />
            Live
          </Badge>
          <Button
            variant="ghost"
            leftIcon={<RefreshCw className="w-4 h-4" />}
            onClick={() => queryClient.invalidateQueries({ queryKey: ['active', routerId] })}
          >
            Refresh
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <Card.Body className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-cyan-100 dark:bg-cyan-900/30 flex items-center justify-center">
              <Wifi className="w-6 h-6 text-cyan-600" />
            </div>
            <div>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {activeUsers?.length || 0}
              </p>
              <p className="text-sm text-gray-500">Active Users</p>
            </div>
          </Card.Body>
        </Card>
      </div>

      {/* Table */}
      <Card>
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                <th>User</th>
                <th>IP Address</th>
                <th>MAC Address</th>
                <th>Uptime</th>
                <th>Bytes In/Out</th>
                <th>Time Left</th>
                <th>Login By</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr>
                  <td colSpan={8} className="text-center py-8">
                    <RefreshCw className="w-6 h-6 animate-spin mx-auto text-gray-400" />
                  </td>
                </tr>
              ) : activeUsers?.length === 0 ? (
                <tr>
                  <td colSpan={8} className="text-center py-12">
                    <div className="flex flex-col items-center">
                      <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mb-4">
                        <Wifi className="w-8 h-8 text-gray-400" />
                      </div>
                      <p className="text-gray-500">No active sessions</p>
                    </div>
                  </td>
                </tr>
              ) : (
                activeUsers?.map((user) => (
                  <tr key={user.id}>
                    <td>
                      <div className="flex items-center gap-2">
                        <div className="w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center">
                          <span className="text-sm font-medium text-primary-600">
                            {user.user[0]?.toUpperCase()}
                          </span>
                        </div>
                        <span className="font-medium">{user.user}</span>
                      </div>
                    </td>
                    <td className="font-mono text-sm">{user.address}</td>
                    <td className="font-mono text-sm">{user.macAddress}</td>
                    <td>{user.uptime}</td>
                    <td>
                      <div className="text-sm space-y-1">
                        <div className="text-success-600">↓ {formatBytes(user.bytesIn)}</div>
                        <div className="text-primary-600">↑ {formatBytes(user.bytesOut)}</div>
                      </div>
                    </td>
                    <td>
                      {user.sessionTimeLeft ? (
                        <Badge variant="warning">{user.sessionTimeLeft}</Badge>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </td>
                    <td>
                      <div className="flex items-center gap-2">
                        {getLoginByIcon(user.loginBy)}
                        <span className="text-sm capitalize">{user.loginBy}</span>
                      </div>
                    </td>
                    <td>
                      <Button
                        variant="ghost"
                        size="sm"
                        leftIcon={<Power className="w-4 h-4 text-danger-500" />}
                        onClick={() => {
                          toast.success(`Kick user ${user.user} (demo)`)
                        }}
                      >
                        Kick
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </motion.div>
  )
}
