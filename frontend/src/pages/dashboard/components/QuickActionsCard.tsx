import { Activity } from 'lucide-react'

import { Card } from '../../../components/ui/Card'

const actions = [
  { label: 'Add Hotspot User', color: 'primary', href: '/hotspot/users' },
  { label: 'Generate Vouchers', color: 'secondary', href: '/vouchers/generate' },
  { label: 'View Reports', color: 'success', href: '/reports' },
  { label: 'Manage Routers', color: 'warning', href: '/routers' },
]

export function QuickActionsCard() {
  return (
    <Card className="h-full">
      <Card.Header>
        <div className="flex items-center gap-2">
          <Activity className="w-5 h-5 text-success-500" />
          <h3 className="font-semibold text-gray-900 dark:text-white">Quick Actions</h3>
        </div>
      </Card.Header>
      <Card.Body className="space-y-3">
        {actions.map((action) => (
          <button
            key={action.label}
            onClick={() => window.location.href = action.href}
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
  )
}
