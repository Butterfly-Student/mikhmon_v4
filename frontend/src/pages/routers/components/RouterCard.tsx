import { Router, Globe, Shield, Edit2, Trash2, RefreshCw } from 'lucide-react'

import { Card, Button, Badge } from '../../../components/ui'
import type { Router as RouterType } from '../../../types'

interface RouterCardProps {
  router: RouterType
  isSelected: boolean
  testingId: string | number | null
  onEdit: (router: RouterType) => void
  onDelete: (id: string | number) => void
  onSelect: (router: RouterType) => void
  onTest: (id: string | number) => void
}

export function RouterCard({ router, isSelected, testingId, onEdit, onDelete, onSelect, onTest }: RouterCardProps) {
  return (
    <Card className={isSelected ? 'ring-2 ring-primary-500' : ''}>
      <Card.Body>
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center gap-3">
            <div
              className={`w-12 h-12 rounded-xl flex items-center justify-center ${
                router.isActive
                  ? 'bg-success-100 dark:bg-success-900/30'
                  : 'bg-gray-100 dark:bg-dark-700'
              }`}
            >
              <Router
                className={`w-6 h-6 ${router.isActive ? 'text-success-600' : 'text-gray-400'}`}
              />
            </div>
            <div>
              <h3 className="font-semibold text-lg text-gray-900 dark:text-white">
                {router.name}
              </h3>
              <div className="flex items-center gap-2 mt-1">
                <Badge variant={router.isActive ? 'success' : 'default'}>
                  {router.isActive ? 'Online' : 'Offline'}
                </Badge>
                {router.useSsl && (
                  <Badge variant="primary">
                    <Shield className="w-3 h-3 mr-1" />
                    SSL
                  </Badge>
                )}
              </div>
            </div>
          </div>
          <div className="flex gap-1">
            <Button variant="ghost" size="sm" onClick={() => onEdit(router)}>
              <Edit2 className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                if (confirm('Are you sure?')) {
                  onDelete(router.id)
                }
              }}
            >
              <Trash2 className="w-4 h-4 text-danger-500" />
            </Button>
          </div>
        </div>

        <div className="space-y-2 text-sm">
          <div className="flex items-center gap-2 text-gray-600 dark:text-gray-300">
            <Globe className="w-4 h-4 text-gray-400" />
            <span className="font-mono">
              {router.host}:{router.port}
            </span>
          </div>
          {router.description && (
            <p className="text-gray-500">{router.description}</p>
          )}
          {router.lastConnected && (
            <p className="text-xs text-gray-400">
              Last connected: {new Date(router.lastConnected).toLocaleString()}
            </p>
          )}
        </div>

        <div className="mt-4 flex gap-2">
          <Button
            variant={isSelected ? 'primary' : 'ghost'}
            size="sm"
            className="flex-1"
            onClick={() => onSelect(router)}
          >
            {isSelected ? 'Selected' : 'Select'}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            isLoading={testingId === router.id}
            onClick={() => onTest(router.id)}
          >
            <RefreshCw className="w-4 h-4" />
          </Button>
        </div>
      </Card.Body>
    </Card>
  )
}
