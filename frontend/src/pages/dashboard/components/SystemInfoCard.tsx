import { Server } from 'lucide-react'

import { Card } from '../../../components/ui/Card'
import type { ResourceStats } from '../../../hooks/useResourceWebSocket'

interface SystemInfoCardProps {
  resources: ResourceStats | null
  identityName?: string
}

export function SystemInfoCard({ resources, identityName }: SystemInfoCardProps) {
  const items = [
    { label: 'Uptime',    value: resources?.uptime   || '-' },
    { label: 'Identity',  value: identityName        || '-' },
    { label: 'Board',     value: resources?.boardName || '-' },
    { label: 'Platform',  value: resources?.platform  || '-' },
    { label: 'RouterOS',  value: resources?.version   || '-' },
    { label: 'CPU',       value: resources?.cpu        || '-' },
    { label: 'CPU Freq',  value: resources?.cpuFrequency ? `${resources.cpuFrequency} MHz` : '-' },
  ]

  return (
    <Card className="h-full">
      <Card.Header>
        <div className="flex items-center gap-2">
          <Server className="w-5 h-5 text-secondary-500" />
          <h3 className="font-semibold text-gray-900 dark:text-white">System Info</h3>
        </div>
      </Card.Header>
      <Card.Body>
        <div className="space-y-4">
          {items.map((item) => (
            <div key={item.label} className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-dark-700 last:border-0">
              <span className="text-sm text-gray-500 dark:text-gray-400">{item.label}</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white text-right">{item.value}</span>
            </div>
          ))}
        </div>
      </Card.Body>
    </Card>
  )
}
