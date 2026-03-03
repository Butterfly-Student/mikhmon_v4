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
  RefreshCw,
  AlertTriangle,
  Gauge,
  Network,
} from 'lucide-react'

import { StatCard } from '../components/common/StatCard'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { QueueMonitor, InterfaceMonitor } from '../components/monitor'
import { dashboardApi } from '../api/dashboard'
import { hotspotApi } from '../api/hotspot'
import { useRouterStore } from '../stores/routerStore'
import { toggleApiDebug } from '../api/axios'
import { useResourceWebSocket } from '../hooks/useResourceWebSocket'
import { ResourceMonitorCard } from './dashboard/components/ResourceMonitorCard'
import { SystemInfoCard } from './dashboard/components/SystemInfoCard'
import { QuickActionsCard } from './dashboard/components/QuickActionsCard'

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.1 }
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
  const routerId = selectedRouter ? String(selectedRouter.id) : '1'

  if (import.meta.env.DEV) {
    console.log('[DashboardPage] Rendered - routerId:', routerId, 'selectedRouter:', selectedRouter)
  }

  const { stats: resources, isConnected } = useResourceWebSocket(routerId)
  console.log('[DashboardPage] Resource stats:', resources, 'WebSocket connected:', isConnected)

  const { data: identityInfo } = useQuery({
    queryKey: ['identity', routerId],
    queryFn: () => dashboardApi.getIdentity(routerId),
    refetchInterval: 60000,
    enabled: !!routerId,
    retry: 2,
  })

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

  const activeUsersCount = activeUsers?.length || 0
  const isOnline = isConnected || !!totalUsersCount || !!activeUsers

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(value)
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
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
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

      {/* WS Disconnected Banner */}
      {!isConnected && (
        <motion.div variants={itemVariants}>
          <Card className="bg-warning-50 dark:bg-warning-900/20 border-warning-200 dark:border-warning-800">
            <Card.Body className="flex items-start gap-3">
              <AlertTriangle className="w-5 h-5 text-warning-600 dark:text-warning-400 flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <h3 className="font-semibold text-warning-900 dark:text-warning-100 mb-1">
                  Resource Monitor Disconnected
                </h3>
                <p className="text-sm text-warning-700 dark:text-warning-300">
                  Waiting for realtime data from router. Reconnecting automatically...
                </p>
                <div className="mt-3 flex gap-2">
                  <Button variant="secondary" size="sm" onClick={() => window.location.reload()}>
                    <RefreshCw className="w-4 h-4 mr-1" />
                    Refresh Page
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => toggleApiDebug(true)}>
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
        <StatCard title="Active Users" value={activeUsersCount} subtitle="Currently connected" icon={Wifi} gradient="cyan" />
        <StatCard title="Total Users" value={totalUsersCount || 0} subtitle="Registered users" icon={Users} gradient="indigo" />
        <StatCard title="Monthly Income" value={formatCurrency(0)} subtitle="This month" icon={Wallet} gradient="emerald" />
        <StatCard title="Daily Income" value={formatCurrency(0)} subtitle="Today" icon={TrendingUp} gradient="amber" />
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
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
            <ResourceMonitorCard resources={resources} isConnected={isConnected} />
          </motion.div>
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
            <SystemInfoCard resources={resources} identityName={identityInfo?.name} />
          </motion.div>
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2 }}>
            <QuickActionsCard />
          </motion.div>
        </div>
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
                <p>Queue monitoring allows you to track bandwidth usage for specific queues in real-time.</p>
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
                <p>Interface monitoring allows you to track network traffic on specific interfaces in real-time.</p>
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
