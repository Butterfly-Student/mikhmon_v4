import { Activity, Cpu, Server, HardDrive, Radio } from 'lucide-react'

import { Card } from '../../../components/ui/Card'
import type { ResourceStats } from '../../../hooks/useResourceWebSocket'

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

interface ResourceMonitorCardProps {
  resources: ResourceStats | null
  isConnected: boolean
}

export function ResourceMonitorCard({ resources, isConnected }: ResourceMonitorCardProps) {
  return (
    <Card className="h-full">
      <Card.Header>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Activity className="w-5 h-5 text-primary-500" />
            <h3 className="font-semibold text-gray-900 dark:text-white">System Resources</h3>
          </div>
          <div className="flex items-center gap-1.5 text-xs text-gray-400">
            <Radio className={`w-3 h-3 ${isConnected ? 'text-success-500 animate-pulse' : 'text-gray-400'}`} />
            {isConnected ? 'Live' : 'Connecting…'}
          </div>
        </div>
      </Card.Header>
      <Card.Body className="space-y-5">
        {/* CPU */}
        <div>
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
              <Cpu className="w-4 h-4" />
              CPU Load
              {resources?.cpu && (
                <span className="text-xs text-gray-400">({resources.cpu})</span>
              )}
            </div>
            <span className="text-sm font-medium">{resources?.cpuLoad ?? 0}%</span>
          </div>
          <div className="progress-bar">
            <div
              className="progress-bar-fill bg-gradient-to-r from-primary-500 to-primary-600"
              style={{ width: `${resources?.cpuLoad ?? 0}%` }}
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
                width: `${resources?.totalMemory ? ((resources.totalMemory - (resources.freeMemory ?? 0)) / resources.totalMemory * 100) : 0}%`
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
                width: `${resources?.totalHddSpace ? ((resources.totalHddSpace - (resources.freeHddSpace ?? 0)) / resources.totalHddSpace * 100) : 0}%`
              }}
            />
          </div>
        </div>
      </Card.Body>
    </Card>
  )
}
