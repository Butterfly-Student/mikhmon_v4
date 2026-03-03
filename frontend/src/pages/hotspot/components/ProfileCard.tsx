import { Edit2, Trash2, Users, Gauge, Clock, DollarSign, Lock, Server } from 'lucide-react'

import { Card, Button, Badge } from '../../../components/ui'
import type { UserProfile } from '../../../types'

interface ProfileCardProps {
  profile: UserProfile
  onEdit: (profile: UserProfile) => void
  onDelete: (id: string) => void
}

const getExpireModeBadge = (mode?: string): 'default' | 'danger' | 'warning' => {
  const variants: Record<string, 'default' | 'danger' | 'warning'> = {
    '0': 'default',
    rem: 'danger',
    ntf: 'warning',
    remc: 'danger',
    ntfc: 'warning',
  }
  return variants[mode || '0'] || 'default'
}

export function ProfileCard({ profile, onEdit, onDelete }: ProfileCardProps) {
  return (
    <Card hover>
      <Card.Body>
        <div className="flex items-start justify-between mb-4">
          <div>
            <h3 className="font-semibold text-lg text-gray-900 dark:text-white">
              {profile.name}
            </h3>
            <Badge variant={getExpireModeBadge(profile.expireMode)} className="mt-1">
              {profile.expireMode === '0' ? 'No Expiry' : profile.expireMode?.toUpperCase()}
            </Badge>
          </div>
          <div className="flex gap-1">
            <Button variant="ghost" size="sm" onClick={() => onEdit(profile)}>
              <Edit2 className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                if (confirm('Are you sure?')) {
                  onDelete(profile.id)
                }
              }}
            >
              <Trash2 className="w-4 h-4 text-danger-500" />
            </Button>
          </div>
        </div>

        <div className="space-y-3">
          <div className="flex items-center gap-3 text-sm">
            <Users className="w-4 h-4 text-gray-400" />
            <span className="text-gray-600 dark:text-gray-300">
              {profile.sharedUsers} shared users
            </span>
          </div>

          {profile.rateLimit && (
            <div className="flex items-center gap-3 text-sm">
              <Gauge className="w-4 h-4 text-gray-400" />
              <span className="text-gray-600 dark:text-gray-300">{profile.rateLimit}</span>
            </div>
          )}

          {profile.validity && (
            <div className="flex items-center gap-3 text-sm">
              <Clock className="w-4 h-4 text-gray-400" />
              <span className="text-gray-600 dark:text-gray-300">{profile.validity}</span>
            </div>
          )}

          <div className="flex items-center gap-3 text-sm">
            <DollarSign className="w-4 h-4 text-gray-400" />
            <span className="text-gray-600 dark:text-gray-300">
              Rp {profile.price?.toLocaleString('id-ID')}
              {profile.sellingPrice > profile.price && (
                <span className="text-success-600 ml-1">
                  → Rp {profile.sellingPrice?.toLocaleString('id-ID')}
                </span>
              )}
            </span>
          </div>

          {(profile.lockUser === 'Enable' || profile.lockServer === 'Enable') && (
            <div className="flex items-center gap-2 pt-2">
              {profile.lockUser === 'Enable' && (
                <Badge variant="info" size="sm">
                  <Lock className="w-3 h-3 mr-1" />
                  MAC Lock
                </Badge>
              )}
              {profile.lockServer === 'Enable' && (
                <Badge variant="warning" size="sm">
                  <Server className="w-3 h-3 mr-1" />
                  Server Lock
                </Badge>
              )}
            </div>
          )}
        </div>
      </Card.Body>
    </Card>
  )
}
