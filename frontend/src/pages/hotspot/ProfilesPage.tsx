import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import {
  Plus,
  Edit2,
  Trash2,
  Users,
  Gauge,
  Clock,
  DollarSign,
  Lock,
  Server,
  AlertTriangle,
  Activity,
  RefreshCw,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { Card, Button, Input, Badge, Modal, Select } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'
import type { UserProfile } from '../../types'
import { Link } from 'react-router-dom'
import { toggleApiDebug } from '../../api/axios'

const profileSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  sharedUsers: z.coerce.number().min(1).default(1),
  rateLimit: z.string().optional(),
  addressPool: z.string().optional(),
  parentQueue: z.string().optional(),
  expireMode: z.string().default('0'),
  validity: z.string().optional(),
  price: z.coerce.number().min(0).default(0),
  sellingPrice: z.coerce.number().min(0).default(0),
  lockUser: z.string().default('Disable'),
  lockServer: z.string().default('Disable'),
})

type ProfileForm = z.infer<typeof profileSchema>

const expireModes = [
  { value: '0', label: 'None' },
  { value: 'rem', label: 'Remove' },
  { value: 'ntf', label: 'Notice' },
  { value: 'remc', label: 'Remove & Record' },
  { value: 'ntfc', label: 'Notice & Record' },
]

const lockOptions = [
  { value: 'Disable', label: 'Disable' },
  { value: 'Enable', label: 'Enable' },
]

export function ProfilesPage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '')

  if (import.meta.env.DEV) {
    console.log('[ProfilesPage] Rendered - routerId:', routerId, 'selectedRouter:', selectedRouter)
  }

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProfile, setEditingProfile] = useState<UserProfile | null>(null)

  const { data: profiles, isLoading, error: profilesError, refetch } = useQuery({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
    enabled: !!selectedRouter,
    retry: 2,
  })

  // Show error toast if there's a connection error
  if (profilesError && !isLoading) {
    const errorMsg = profilesError instanceof Error ? profilesError.message : 'Failed to load profiles'
    if (errorMsg.includes('connection') || errorMsg.includes('Network Error') || errorMsg.includes('timeout')) {
      toast.error('Cannot connect to router. Please check your network settings.', { id: 'profiles-error', duration: 5000 })
    }
  }

  const createMutation = useMutation({
    mutationFn: (data: ProfileForm) => hotspotApi.createProfile(routerId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profiles', routerId] })
      toast.success('Profile created successfully')
      setIsModalOpen(false)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProfileForm }) =>
      hotspotApi.updateProfile(routerId, id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profiles', routerId] })
      toast.success('Profile updated successfully')
      setIsModalOpen(false)
      setEditingProfile(null)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => hotspotApi.deleteProfile(routerId, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profiles', routerId] })
      toast.success('Profile deleted successfully')
    },
    onError: (error: any) => toast.error(error.message),
  })

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      sharedUsers: 1,
      expireMode: '0',
      lockUser: 'Disable',
      lockServer: 'Disable',
    },
  })

  const expireMode = watch('expireMode')

  const onSubmit = (data: ProfileForm) => {
    if (editingProfile) {
      updateMutation.mutate({ id: editingProfile.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const openModal = (profile?: UserProfile) => {
    if (profile) {
      setEditingProfile(profile)
      reset({
        name: profile.name,
        sharedUsers: profile.sharedUsers,
        rateLimit: profile.rateLimit,
        addressPool: profile.addressPool,
        parentQueue: profile.parentQueue,
        expireMode: profile.expireMode || '0',
        validity: profile.validity,
        price: profile.price,
        sellingPrice: profile.sellingPrice,
        lockUser: profile.lockUser || 'Disable',
        lockServer: profile.lockServer || 'Disable',
      })
    } else {
      setEditingProfile(null)
      reset({
        sharedUsers: 1,
        expireMode: '0',
        lockUser: 'Disable',
        lockServer: 'Disable',
      })
    }
    setIsModalOpen(true)
  }

  const getExpireModeBadge = (mode?: string) => {
    const variants: Record<string, any> = {
      '0': 'default',
      rem: 'danger',
      ntf: 'warning',
      remc: 'danger',
      ntfc: 'warning',
    }
    return variants[mode || '0'] || 'default'
  }

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <Users className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">
          Silahkan pilih router terlebih dahulu untuk melihat dan mengelola user profiles.
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">User Profiles</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Manage hotspot user profiles and pricing
          </p>
          {selectedRouter && (
            <p className="text-xs text-gray-400 mt-1">
              Router: {selectedRouter.name} (ID: {selectedRouter.id})
            </p>
          )}
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => toggleApiDebug()}
            title="Toggle API Debug Logging"
          >
            <Activity className="w-4 h-4" />
          </Button>
          <Button onClick={() => openModal()} leftIcon={<Plus className="w-4 h-4" />}>
            Add Profile
          </Button>
        </div>
      </div>

      {/* Error Banner */}
      {profilesError && !isLoading && (
        <Card className="bg-danger-50 dark:bg-danger-900/20 border-danger-200 dark:border-danger-800">
          <Card.Body className="flex items-start gap-3">
            <AlertTriangle className="w-5 h-5 text-danger-600 dark:text-danger-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <h3 className="font-semibold text-danger-900 dark:text-danger-100 mb-1">
                Failed to Load Profiles
              </h3>
              <p className="text-sm text-danger-700 dark:text-danger-300">
                {profilesError instanceof Error ? profilesError.message : 'An unknown error occurred'}
              </p>
              <div className="mt-3 flex gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => refetch()}
                >
                  <RefreshCw className="w-4 h-4 mr-1" />
                  Retry
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
      )}

      {/* Profiles Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <Card key={i} className="h-64 animate-pulse" />
          ))}
        </div>
      ) : profiles?.length === 0 ? (
        <Card>
          <Card.Body className="text-center py-12">
            <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mx-auto mb-4">
              <Users className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-500">No profiles found</p>
          </Card.Body>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {profiles?.map((profile) => (
            <Card key={profile.id} hover>
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
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => openModal(profile)}
                    >
                      <Edit2 className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        if (confirm('Are you sure?')) {
                          deleteMutation.mutate(profile.id)
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
          ))}
        </div>
      )}

      {/* Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={editingProfile ? 'Edit Profile' : 'Add Profile'}
        size="lg"
        footer={
          <div className="flex justify-end gap-3">
            <Button variant="ghost" onClick={() => setIsModalOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmit(onSubmit)}
              isLoading={createMutation.isPending || updateMutation.isPending}
            >
              {editingProfile ? 'Update' : 'Create'}
            </Button>
          </div>
        }
      >
        <form className="space-y-6">
          {/* General Settings */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">General</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input label="Name" {...register('name')} error={errors.name?.message} />
              <Input
                label="Shared Users"
                type="number"
                {...register('sharedUsers', { valueAsNumber: true })}
              />
              <Input label="Rate Limit" {...register('rateLimit')} placeholder="e.g. 1M/1M" />
              <Input label="Address Pool" {...register('addressPool')} />
              <Input label="Parent Queue" {...register('parentQueue')} />
            </div>
          </div>

          {/* Expiration Settings */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Expiration</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Select
                label="Expire Mode"
                options={expireModes}
                {...register('expireMode')}
              />
              {expireMode !== '0' && (
                <Input
                  label="Validity"
                  {...register('validity')}
                  placeholder="e.g. 30d, 12h"
                />
              )}
            </div>
          </div>

          {/* Pricing */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Pricing</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Price"
                type="number"
                {...register('price', { valueAsNumber: true })}
              />
              <Input
                label="Selling Price"
                type="number"
                {...register('sellingPrice', { valueAsNumber: true })}
              />
            </div>
          </div>

          {/* Lock Settings */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Lock Settings</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Select label="Lock User (MAC)" options={lockOptions} {...register('lockUser')} />
              <Select label="Lock Server" options={lockOptions} {...register('lockServer')} />
            </div>
          </div>
        </form>
      </Modal>
    </motion.div>
  )
}
