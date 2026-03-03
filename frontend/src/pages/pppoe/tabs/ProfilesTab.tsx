import { useState, useMemo, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Edit2, Trash2, Plus, PieChart, Search, X } from 'lucide-react'
import toast from 'react-hot-toast'

import { Button, Input, Badge, Modal, Pagination } from '../../../components/ui'
import { pppApi } from '../../../api/ppp'
import type { PPPProfile } from '../../../types'

const profileSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  localAddress: z.string().optional(),
  remoteAddress: z.string().optional(),
  dnsServer: z.string().optional(),
  rateLimit: z.string().optional(),
  sessionTimeout: z.string().optional(),
  idleTimeout: z.string().optional(),
  comment: z.string().optional(),
})
type ProfileForm = z.infer<typeof profileSchema>

interface ProfilesTabProps {
  routerId: string
}

export function ProfilesTab({ routerId }: ProfilesTabProps) {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editingProfile, setEditingProfile] = useState<PPPProfile | null>(null)
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { data: profiles = [], isLoading } = useQuery({
    queryKey: ['ppp-profiles', routerId],
    queryFn: () => pppApi.getProfiles(routerId),
    enabled: !!routerId,
  })

  const { register, handleSubmit, reset, formState: { errors } } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
  })

  const createMutation = useMutation({
    mutationFn: (data: ProfileForm) => pppApi.createProfile(routerId, data),
    onSuccess: () => { toast.success('Profile created'); queryClient.invalidateQueries({ queryKey: ['ppp-profiles', routerId] }); closeModal() },
    onError: (err: any) => toast.error(err.message || 'Failed to create'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProfileForm }) => pppApi.updateProfile(routerId, id, data),
    onSuccess: () => { toast.success('Profile updated'); queryClient.invalidateQueries({ queryKey: ['ppp-profiles', routerId] }); closeModal() },
    onError: (err: any) => toast.error(err.message || 'Failed to update'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => pppApi.deleteProfile(routerId, id),
    onSuccess: () => { toast.success('Profile deleted'); queryClient.invalidateQueries({ queryKey: ['ppp-profiles', routerId] }) },
    onError: (err: any) => toast.error(err.message || 'Failed to delete'),
  })

  function openCreate() {
    setEditingProfile(null)
    reset({ name: '', localAddress: '', remoteAddress: '', dnsServer: '', rateLimit: '', sessionTimeout: '', idleTimeout: '', comment: '' })
    setModalOpen(true)
  }

  function openEdit(profile: PPPProfile) {
    setEditingProfile(profile)
    reset({ name: profile.name, localAddress: profile.localAddress || '', remoteAddress: profile.remoteAddress || '', dnsServer: profile.dnsServer || '', rateLimit: profile.rateLimit || '', sessionTimeout: profile.sessionTimeout || '', idleTimeout: profile.idleTimeout || '', comment: profile.comment || '' })
    setModalOpen(true)
  }

  function closeModal() { setModalOpen(false); setEditingProfile(null); reset() }

  function onSubmit(data: ProfileForm) {
    if (editingProfile) updateMutation.mutate({ id: editingProfile.id, data })
    else createMutation.mutate(data)
  }

  const filtered = useMemo(() =>
    profiles.filter((p) =>
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      (p.comment || '').toLowerCase().includes(search.toLowerCase())
    ), [profiles, search]
  )

  const pageCount = Math.ceil(filtered.length / pageSize)
  const paginatedRows = filtered.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [search])

  const isSaving = createMutation.isPending || updateMutation.isPending

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex items-center gap-3 flex-wrap">
          <div className="relative flex-1 min-w-[160px] sm:max-w-72">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search profiles..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-9 pr-8 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {search && (
              <button onClick={() => setSearch('')} className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
          <Button variant="primary" size="sm" leftIcon={<Plus className="w-4 h-4" />} onClick={openCreate}>
            Add Profile
          </Button>
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
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Name</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Rate Limit</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Local Address</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Remote Address</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Session Timeout</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Comment</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
              {paginatedRows.map((profile) => (
                <tr key={profile.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <div className="w-7 h-7 rounded-lg bg-secondary-100 dark:bg-secondary-900/30 flex items-center justify-center shrink-0">
                        <PieChart className="w-3.5 h-3.5 text-secondary-600 dark:text-secondary-400" />
                      </div>
                      <span className="font-semibold text-gray-900 dark:text-white text-sm">{profile.name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    {profile.rateLimit ? <Badge variant="warning">{profile.rateLimit}</Badge> : <span className="text-gray-400 text-xs">—</span>}
                  </td>
                  <td className="px-4 py-3 font-mono text-xs text-gray-600 dark:text-gray-300">{profile.localAddress || '—'}</td>
                  <td className="px-4 py-3 font-mono text-xs text-gray-600 dark:text-gray-300">{profile.remoteAddress || '—'}</td>
                  <td className="px-4 py-3 text-xs text-gray-600 dark:text-gray-300">{profile.sessionTimeout || '—'}</td>
                  <td className="px-4 py-3 text-xs text-gray-500 dark:text-gray-400">{profile.comment || '-'}</td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1">
                      <Button variant="ghost" size="xs" leftIcon={<Edit2 className="w-3.5 h-3.5" />} onClick={() => openEdit(profile)}>Edit</Button>
                      <Button variant="ghost" size="xs" leftIcon={<Trash2 className="w-3.5 h-3.5 text-danger-500" />}
                        onClick={() => { if (confirm(`Delete profile "${profile.name}"?`)) deleteMutation.mutate(profile.id) }}>Delete</Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {!isLoading && filtered.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400 flex flex-col items-center gap-3">
            <PieChart className="w-10 h-10 text-gray-300 dark:text-gray-600" />
            <span>No PPP profiles found</span>
          </div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filtered.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />

      <Modal
        isOpen={modalOpen}
        onClose={closeModal}
        title={editingProfile ? 'Edit Profile' : 'Add Profile'}
        subtitle={editingProfile ? `Editing: ${editingProfile.name}` : 'Create a new PPP profile'}
        footer={
          <div className="flex justify-end gap-2">
            <Button variant="outline" size="sm" onClick={closeModal}>Cancel</Button>
            <Button variant="primary" size="sm" isLoading={isSaving} onClick={handleSubmit(onSubmit)}>
              {editingProfile ? 'Save' : 'Create'}
            </Button>
          </div>
        }
      >
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Profile Name" {...register('name')} error={errors.name?.message} placeholder="default" />
          <Input label="Rate Limit" {...register('rateLimit')} placeholder="2M/2M" helperText="Upload/Download e.g. 2M/10M" />
          <div className="grid grid-cols-2 gap-4">
            <Input label="Local Address" {...register('localAddress')} placeholder="192.168.1.1" />
            <Input label="Remote Address" {...register('remoteAddress')} placeholder="192.168.1.0/24" />
          </div>
          <Input label="DNS Server" {...register('dnsServer')} placeholder="8.8.8.8,8.8.4.4" />
          <div className="grid grid-cols-2 gap-4">
            <Input label="Session Timeout" {...register('sessionTimeout')} placeholder="1d 00:00:00" />
            <Input label="Idle Timeout" {...register('idleTimeout')} placeholder="00:30:00" />
          </div>
          <Input label="Comment" {...register('comment')} placeholder="Optional comment" />
        </form>
      </Modal>
    </div>
  )
}
