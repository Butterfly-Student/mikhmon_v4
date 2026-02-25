import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import {
  Laptop,
  CheckCircle,
  XCircle,
  Shield,
  RefreshCw,
} from 'lucide-react'

import { Card, Button, Badge } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'

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

  const filteredHosts = hosts?.filter((host) => {
    if (filter === 'authorized') return host.authorized
    if (filter === 'bypassed') return host.bypassed
    return true
  })

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <Laptop className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">
          Silahkan pilih router terlebih dahulu untuk melihat data hotspot hosts.
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Hotspot Hosts</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Connected devices on the network
          </p>
        </div>
        <Button
          variant="ghost"
          leftIcon={<RefreshCw className="w-4 h-4" />}
          onClick={() => queryClient.invalidateQueries({ queryKey: ['hosts', routerId] })}
        >
          Refresh
        </Button>
      </div>

      {/* Filter Tabs */}
      <div className="flex gap-2">
        {(['all', 'authorized', 'bypassed'] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${filter === f
                ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400'
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700'
              }`}
          >
            {f.charAt(0).toUpperCase() + f.slice(1)}
          </button>
        ))}
      </div>

      {/* Hosts Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <Card key={i} className="h-48 animate-pulse" />
          ))}
        </div>
      ) : filteredHosts?.length === 0 ? (
        <Card>
          <Card.Body className="text-center py-12">
            <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mx-auto mb-4">
              <Laptop className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-500">No hosts found</p>
          </Card.Body>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredHosts?.map((host) => (
            <Card key={host.id} hover>
              <Card.Body>
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div
                      className={`w-12 h-12 rounded-xl flex items-center justify-center ${host.authorized
                          ? 'bg-success-100 dark:bg-success-900/30'
                          : host.bypassed
                            ? 'bg-warning-100 dark:bg-warning-900/30'
                            : 'bg-gray-100 dark:bg-dark-700'
                        }`}
                    >
                      <Laptop
                        className={`w-6 h-6 ${host.authorized
                            ? 'text-success-600'
                            : host.bypassed
                              ? 'text-warning-600'
                              : 'text-gray-400'
                          }`}
                      />
                    </div>
                    <div>
                      <p className="font-mono text-sm text-gray-900 dark:text-white">
                        {host.macAddress}
                      </p>
                      {host.address && (
                        <p className="text-sm text-gray-500">{host.address}</p>
                      )}
                    </div>
                  </div>
                </div>

                <div className="mt-4 flex flex-wrap gap-2">
                  {host.authorized && (
                    <Badge variant="success">
                      <CheckCircle className="w-3 h-3 mr-1" />
                      Authorized
                    </Badge>
                  )}
                  {host.bypassed && (
                    <Badge variant="warning">
                      <Shield className="w-3 h-3 mr-1" />
                      Bypassed
                    </Badge>
                  )}
                  {host.blocked && (
                    <Badge variant="danger">
                      <XCircle className="w-3 h-3 mr-1" />
                      Blocked
                    </Badge>
                  )}
                </div>

                {host.server && (
                  <p className="mt-3 text-xs text-gray-500">
                    Server: {host.server}
                  </p>
                )}
              </Card.Body>
            </Card>
          ))}
        </div>
      )}
    </motion.div>
  )
}
