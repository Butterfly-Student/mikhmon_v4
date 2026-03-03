import { useState, useMemo, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Search, X, RefreshCw, UserPlus, ShieldCheck, Ticket, Edit2, Trash2 } from 'lucide-react'
import toast from 'react-hot-toast'

import { Button, Input, Modal, Pagination } from '../../../components/ui'
import { hotspotApi } from '../../../api/hotspot'
import type { UserProfile } from '../../../types'

const profileSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  sharedUsers: z.string().optional(),
  rateLimit: z.string().optional(),
})
type ProfileForm = z.infer<typeof profileSchema>

interface ProfilesTabProps {
  routerId: string
}

export function ProfilesTab({ routerId }: ProfilesTabProps) {
  const queryClient = useQueryClient()
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProfile, setEditingProfile] = useState<UserProfile | null>(null)
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { data: profiles = [], isLoading, refetch } = useQuery<UserProfile[]>({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
  })

  const createMutation = useMutation({
    mutationFn: (data: ProfileForm) => hotspotApi.createProfile(routerId, {
      ...data,
      sharedUsers: data.sharedUsers ? parseInt(data.sharedUsers) : undefined,
    }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profiles', routerId] })
      toast.success('Profile created successfully')
      setIsModalOpen(false)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProfileForm }) =>
      hotspotApi.updateProfile(routerId, id, {
        ...data,
        sharedUsers: data.sharedUsers ? parseInt(data.sharedUsers) : undefined,
      }),
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

  const setupExpireMonitorMutation = useMutation({
    mutationFn: () => hotspotApi.setupExpireMonitor(routerId),
    onSuccess: (result) => {
      toast.success(result.status === 'existing' ? 'Expire Monitor already active' : 'Expire Monitor activated')
    },
    onError: (error: any) => toast.error(error.message || 'Failed to setup Expire Monitor'),
  })

  const { register, handleSubmit, reset, formState: { errors } } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
  })

  const openModal = (profile?: UserProfile) => {
    if (profile) {
      setEditingProfile(profile)
      reset({ name: profile.name, sharedUsers: profile.sharedUsers?.toString(), rateLimit: profile.rateLimit })
    } else {
      setEditingProfile(null)
      reset({ name: '', sharedUsers: undefined, rateLimit: '' })
    }
    setIsModalOpen(true)
  }

  const onSubmit = (data: ProfileForm) => {
    if (editingProfile) updateMutation.mutate({ id: editingProfile.id, data })
    else createMutation.mutate(data)
  }

  const filteredProfiles = useMemo(() =>
    profiles.filter((p) => p.name.toLowerCase().includes(searchQuery.toLowerCase())),
    [profiles, searchQuery]
  )

  const pageCount = Math.ceil(filteredProfiles.length / pageSize)
  const paginatedRows = filteredProfiles.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [searchQuery])

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex flex-col lg:flex-row lg:items-center gap-3">
          <div className="flex items-center gap-2">
            <span className="px-3 py-1.5 bg-gray-100 dark:bg-dark-700 rounded-lg text-xs font-medium text-gray-700 dark:text-gray-300">
              Total: {profiles.length}
            </span>
            <button onClick={() => refetch()} className="p-1.5 rounded-lg bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600 transition-colors" title="Refresh">
              <RefreshCw className="w-4 h-4" />
            </button>
          </div>

          <div className="relative flex-1 min-w-[160px] sm:max-w-64">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Filter profiles..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-9 pr-8 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {searchQuery && (
              <button onClick={() => setSearchQuery('')} className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                <X className="w-4 h-4" />
              </button>
            )}
          </div>

          <div className="flex items-center gap-2 flex-wrap">
            <Button onClick={() => openModal()} size="sm">
              <UserPlus className="w-4 h-4 mr-1" />Add Profile
            </Button>
            <Button variant="secondary" size="sm" isLoading={setupExpireMonitorMutation.isPending} onClick={() => setupExpireMonitorMutation.mutate()}>
              <ShieldCheck className="w-4 h-4 mr-1" />Expire Monitor
            </Button>
            <Button variant="secondary" size="sm" onClick={() => window.location.href = '/vouchers/generate'}>
              <Ticket className="w-4 h-4 mr-1" />Generate
            </Button>
          </div>
        </div>
      </div>

      {/* Table — scrollable */}
      <div className="flex-1 min-h-0 overflow-x-auto overflow-y-auto">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="sticky top-0 z-10 bg-gray-50 dark:bg-dark-700 border-b border-gray-200 dark:border-dark-700">
              <tr>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase w-20">Actions</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Name</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Shared Users</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Rate Limit</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Expire Mode</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Validity</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Price</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Selling Price</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">User Lock</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Server Lock</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
              {paginatedRows.map((profile) => (
                <tr key={profile.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                  <td className="px-4 py-3 text-center">
                    <div className="flex items-center justify-center gap-1">
                      <button onClick={() => openModal(profile)} className="p-1.5 rounded-lg text-warning-600 hover:bg-warning-50 dark:hover:bg-warning-900/20">
                        <Edit2 className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => { if (confirm('Delete this profile?')) deleteMutation.mutate(profile.id) }}
                        className="p-1.5 rounded-lg text-danger-600 hover:bg-danger-50 dark:hover:bg-danger-900/20"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{profile.name}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{profile.sharedUsers || '1'}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{profile.rateLimit || '-'}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">-</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">-</td>
                  <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">-</td>
                  <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">-</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">-</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">-</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {!isLoading && filteredProfiles.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">No profiles found</div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filteredProfiles.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingProfile ? 'Edit Profile' : 'Add Profile'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Name" {...register('name')} error={errors.name?.message} />
          <Input label="Shared Users" {...register('sharedUsers')} />
          <Input label="Rate Limit" {...register('rateLimit')} placeholder="e.g. 1M/1M" />
          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>Cancel</Button>
            <Button type="submit" isLoading={createMutation.isPending || updateMutation.isPending}>
              {editingProfile ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
