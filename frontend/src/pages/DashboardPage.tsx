// Dashboard Page
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import {
  Wifi,
  Users,
  Wallet,
  TrendingUp,
  Activity,
  Server,
  Cpu,
  HardDrive,
  RefreshCw,
  AlertTriangle,
  Gauge,
  Network,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { StatCard } from '../components/common/StatCard'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { QueueMonitor, InterfaceMonitor } from '../components/monitor'
import { dashboardApi } from '../api/dashboard'
import { hotspotApi } from '../api/hotspot'
import { useRouterStore } from '../stores/routerStore'
import { toggleApiDebug } from '../api/axios'
import { useResourceWebSocket } from '../hooks/useResourceWebSocket'

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1
    }
  }
}

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 }
}

type MonitorTab = 'overview' | 'queue' | 'interface'

export function DashboardPage() {
  const [activeTab, setActiveTab] = useState<MonitorTab>('overview')
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  // Ensure routerId is always a string and fallback to '1' only if we don't have a selected router
  const routerId = selectedRouter ? String(selectedRouter.id) : '1'

  if (import.meta.env.DEV) {
    console.log('[DashboardPage] Rendered - routerId:', routerId, 'selectedRouter:', selectedRouter)
  }

  // Use WebSocket for realtime system resources instead of polling
  const { stats: resources } = useResourceWebSocket(routerId)

  const { data: systemInfo, error: systemInfoError, isLoading: systemInfoLoading } = useQuery({
    queryKey: ['systemInfo', routerId],
    queryFn: () => dashboardApi.getSystemInfo(routerId),
    refetchInterval: 30000,
    enabled: !!routerId,
    retry: 2,
  })

  // Individual Queries
  const { data: totalUsersCount } = useQuery({
    queryKey: ['usersCount', routerId],
    queryFn: () => hotspotApi.getUsersCount(routerId),
    refetchInterval: 30000,
    enabled: !!routerId,
    retry: 2,
  })

  const { data: activeUsers } = useQuery({
    queryKey: ['activeUsers', routerId],
    queryFn: () => hotspotApi.getActive(routerId),
    refetchInterval: 10000,
    enabled: !!routerId,
    retry: 2,
  })

  const { data: routerBoardInfo } = useQuery({
    queryKey: ['routerBoard', routerId],
    queryFn: () => dashboardApi.getRouterBoard(routerId),
    refetchInterval: 60000,
    enabled: !!routerId,
    retry: 2,
  })

  const { data: identityInfo } = useQuery({
    queryKey: ['identity', routerId],
    queryFn: () => dashboardApi.getIdentity(routerId),
    refetchInterval: 60000,
    enabled: !!routerId,
    retry: 2,
  })

  const activeUsersCount = activeUsers?.length || 0
  const isOnline = !!systemInfo || !!totalUsersCount || !!activeUsers

  // Show error toast if there's a connection error
  if (systemInfoError && !systemInfoLoading) {
    const errorMsg = systemInfoError instanceof Error ? systemInfoError.message : 'Failed to load router info'
    if (errorMsg.includes('connection') || errorMsg.includes('Network Error') || errorMsg.includes('timeout')) {
      toast.error('Cannot connect to router. Please check your network settings.', { id: 'dashboard-error', duration: 5000 })
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(value)
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-6"
    >
      {/* Header */}
      <motion.div variants={itemVariants} className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Router: <span className="font-medium text-primary-600 dark:text-primary-400">{selectedRouter?.name || 'Loading...'}</span>
            {selectedRouter?.id && <span className="text-xs text-gray-400 ml-2">(ID: {selectedRouter.id})</span>}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <span className={`w-2 h-2 rounded-full ${isOnline ? 'bg-success-500 animate-pulse' : 'bg-gray-400'}`} />
            {isOnline ? 'Online' : 'Connecting...'}
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => toggleApiDebug()}
            title="Toggle API Debug Logging"
          >
            <Activity className="w-4 h-4" />
          </Button>
        </div>
      </motion.div>

      {/* Error Banner */}
      {systemInfoError && !systemInfoLoading && (
        <motion.div variants={itemVariants}>
          <Card className="bg-danger-50 dark:bg-danger-900/20 border-danger-200 dark:border-danger-800">
            <Card.Body className="flex items-start gap-3">
              <AlertTriangle className="w-5 h-5 text-danger-600 dark:text-danger-400 flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <h3 className="font-semibold text-danger-900 dark:text-danger-100 mb-1">
                  Connection Error
                </h3>
                <p className="text-sm text-danger-700 dark:text-danger-300">
                  {systemInfoError.message ||
                    'Failed to connect to MikroTik router. Please verify the router is online and credentials are correct.'}
                </p>
                <div className="mt-3 flex gap-2">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => window.location.reload()}
                  >
                    <RefreshCw className="w-4 h-4 mr-1" />
                    Refresh Page
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => toggleApiDebug(true)}
                  >
                    <Activity className="w-4 h-4 mr-1" />
                    Enable Debug
                  </Button>
                </div>
              </div>
            </Card.Body>
          </Card>
        </motion.div>
      )}

      {/* Stats Grid */}
      <motion.div variants={itemVariants} className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Active Users"
          value={activeUsersCount}
          subtitle="Currently connected"
          icon={Wifi}
          gradient="cyan"
        />
        <StatCard
          title="Total Users"
          value={totalUsersCount || 0}
          subtitle="Registered users"
          icon={Users}
          gradient="indigo"
        />
        <StatCard
          title="Monthly Income"
          value={formatCurrency(0 /* Fetched from specific endpoint if available */)}
          subtitle="This month"
          icon={Wallet}
          gradient="emerald"
        />
        <StatCard
          title="Daily Income"
          value={formatCurrency(0 /* Fetched from specific endpoint if available */)}
          subtitle="Today"
          icon={TrendingUp}
          gradient="amber"
        />
      </motion.div>

      {/* Monitor Tabs */}
      <motion.div variants={itemVariants} className="border-b border-gray-200 dark:border-dark-700">
        <div className="flex gap-1">
          {[
            { id: 'overview', label: 'Overview', icon: Activity },
            { id: 'queue', label: 'Queue Monitor', icon: Gauge },
            { id: 'interface', label: 'Interface Monitor', icon: Network },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as MonitorTab)}
              className={`flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 transition-colors
                ${activeTab === tab.id
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
                }`}
            >
              <tab.icon className="w-4 h-4" />
              {tab.label}
            </button>
          ))}
        </div>
      </motion.div>

      {/* Tab Content */}
      {activeTab === 'overview' && (
        <>
          {/* Main Content Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* System Resources */}
            <motion.div variants={itemVariants}>
              <Card className="h-full">
                <Card.Header>
                  <div className="flex items-center gap-2">
                    <Activity className="w-5 h-5 text-primary-500" />
                    <h3 className="font-semibold text-gray-900 dark:text-white">System Resources</h3>
                  </div>
                </Card.Header>
                <Card.Body className="space-y-5">
                  {/* CPU */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
                        <Cpu className="w-4 h-4" />
                        CPU Load
                      </div>
                      <span className="text-sm font-medium">{resources?.cpuUsed ?? 0}%</span>
                    </div>
                    <div className="progress-bar">
                      <div
                        className="progress-bar-fill bg-gradient-to-r from-primary-500 to-primary-600"
                        style={{ width: `${resources?.cpuUsed ?? 0}%` }}
                      />
                    </div>
                  </div>

                  {/* Memory */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
                        <Server className="w-4 h-4" />
                        Memory
                      </div>
                      <span className="text-sm font-medium">
                        {formatBytes((resources?.totalMemory ?? 0) - (resources?.freeMemory ?? 0))} / {formatBytes(resources?.totalMemory ?? 0)}
                      </span>
                    </div>
                    <div className="progress-bar">
                      <div
                        className="progress-bar-fill bg-gradient-to-r from-secondary-500 to-secondary-600"
                        style={{
                          width: `${resources?.totalMemory ? ((resources?.totalMemory - (resources?.freeMemory ?? 0)) / resources?.totalMemory * 100) : 0}%`
                        }}
                      />
                    </div>
                  </div>

                  {/* HDD */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
                        <HardDrive className="w-4 h-4" />
                        Storage
                      </div>
                      <span className="text-sm font-medium">
                        {formatBytes((resources?.totalHddSpace ?? 0) - (resources?.freeHddSpace ?? 0))} / {formatBytes(resources?.totalHddSpace ?? 0)}
                      </span>
                    </div>
                    <div className="progress-bar">
                      <div
                        className="progress-bar-fill bg-gradient-to-r from-warning-500 to-warning-600"
                        style={{
                          width: `${resources?.totalHddSpace ? ((resources?.totalHddSpace - (resources?.freeHddSpace ?? 0)) / resources?.totalHddSpace * 100) : 0}%`
                        }}
                      />
                    </div>
                  </div>
                </Card.Body>
              </Card>
            </motion.div>

            {/* System Info */}
            <motion.div variants={itemVariants}>
              <Card className="h-full">
                <Card.Header>
                  <div className="flex items-center gap-2">
                    <Server className="w-5 h-5 text-secondary-500" />
                    <h3 className="font-semibold text-gray-900 dark:text-white">System Info</h3>
                  </div>
                </Card.Header>
                <Card.Body>
                  <div className="space-y-4">
                    {[
                      { label: 'Uptime', value: systemInfo?.uptime || resources?.uptime || '-' },
                      { label: 'Board Name', value: systemInfo?.boardName || routerBoardInfo?.boardName || '-' },
                      { label: 'Model', value: systemInfo?.model || routerBoardInfo?.model || '-' },
                      { label: 'RouterOS', value: systemInfo?.version || routerBoardInfo?.version || '-' },
                      { label: 'Identity', value: identityInfo?.name || '-' },
                    ].map((item) => (
                      <div key={item.label} className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-dark-700 last:border-0">
                        <span className="text-sm text-gray-500 dark:text-gray-400">{item.label}</span>
                        <span className="text-sm font-medium text-gray-900 dark:text-white text-right">{item.value}</span>
                      </div>
                    ))}
                  </div>
                </Card.Body>
              </Card>
            </motion.div>

            {/* Quick Actions */}
            <motion.div variants={itemVariants}>
              <Card className="h-full">
                <Card.Header>
                  <div className="flex items-center gap-2">
                    <Activity className="w-5 h-5 text-success-500" />
                    <h3 className="font-semibold text-gray-900 dark:text-white">Quick Actions</h3>
                  </div>
                </Card.Header>
                <Card.Body className="space-y-3">
                  {[
                    { label: 'Add Hotspot User', color: 'primary', href: '/hotspot/users' },
                    { label: 'Generate Vouchers', color: 'secondary', href: '/vouchers/generate' },
                    { label: 'View Reports', color: 'success', href: '/reports' },
                    { label: 'Manage Routers', color: 'warning', href: '/routers' },
                  ].map((action) => (
                    <button
                      key={action.label}
                      className={`w-full py-3 px-4 rounded-xl text-sm font-medium transition-all duration-200
                    ${action.color === 'primary' && 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 hover:bg-primary-100 dark:hover:bg-primary-900/30'}
                    ${action.color === 'secondary' && 'bg-secondary-50 dark:bg-secondary-900/20 text-secondary-600 dark:text-secondary-400 hover:bg-secondary-100 dark:hover:bg-secondary-900/30'}
                    ${action.color === 'success' && 'bg-success-50 dark:bg-success-900/20 text-success-600 dark:text-success-400 hover:bg-success-100 dark:hover:bg-success-900/30'}
                    ${action.color === 'warning' && 'bg-warning-50 dark:bg-warning-900/20 text-warning-600 dark:text-warning-400 hover:bg-warning-100 dark:hover:bg-warning-900/30'}
                  `}
                    >
                      {action.label}
                    </button>
                  ))}
                </Card.Body>
              </Card>
            </motion.div>
          </div>
        </>
      )}

      {activeTab === 'queue' && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="grid grid-cols-1 lg:grid-cols-2 gap-6"
        >
          <QueueMonitor />
          <Card>
            <Card.Header>
              <div className="flex items-center gap-2">
                <Gauge className="w-5 h-5 text-primary-500" />
                <h3 className="font-semibold text-gray-900 dark:text-white">Queue Information</h3>
              </div>
            </Card.Header>
            <Card.Body className="space-y-4">
              <div className="text-sm text-gray-600 dark:text-gray-300 space-y-2">
                <p>
                  Queue monitoring allows you to track bandwidth usage for specific queues in real-time.
                </p>
                <ul className="list-disc list-inside space-y-1 text-gray-500">
                  <li>Select a queue from the dropdown to start monitoring</li>
                  <li>RX (Download) shows incoming traffic rate</li>
                  <li>TX (Upload) shows outgoing traffic rate</li>
                  <li>Rates are updated every second</li>
                </ul>
              </div>
            </Card.Body>
          </Card>
        </motion.div>
      )}

      {activeTab === 'interface' && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="grid grid-cols-1 lg:grid-cols-2 gap-6"
        >
          <InterfaceMonitor />
          <Card>
            <Card.Header>
              <div className="flex items-center gap-2">
                <Network className="w-5 h-5 text-secondary-500" />
                <h3 className="font-semibold text-gray-900 dark:text-white">Interface Information</h3>
              </div>
            </Card.Header>
            <Card.Body className="space-y-4">
              <div className="text-sm text-gray-600 dark:text-gray-300 space-y-2">
                <p>
                  Interface monitoring allows you to track network traffic on specific interfaces in real-time.
                </p>
                <ul className="list-disc list-inside space-y-1 text-gray-500">
                  <li>Select an interface (ether1, wlan1, etc.) to monitor</li>
                  <li>RX (Download) shows incoming traffic</li>
                  <li>TX (Upload) shows outgoing traffic</li>
                  <li>Packets per second (pps) shows packet rate</li>
                </ul>
              </div>
            </Card.Body>
          </Card>
        </motion.div>
      )}
    </motion.div>
  )
}
