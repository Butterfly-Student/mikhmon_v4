import { useState, useMemo, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Edit2, Trash2, UserPlus, Network, Search, X, ToggleLeft, ToggleRight } from 'lucide-react'
import toast from 'react-hot-toast'

import { Button, Input, Badge, Modal, Select, Pagination } from '../../../components/ui'
import { pppApi } from '../../../api/ppp'
import type { PPPSecret, PPPProfile } from '../../../types'

const secretSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  password: z.string().optional(),
  profile: z.string().optional(),
  service: z.string().optional(),
  localAddress: z.string().optional(),
  remoteAddress: z.string().optional(),
  comment: z.string().optional(),
})
type SecretForm = z.infer<typeof secretSchema>

const SERVICE_OPTIONS = [
  { value: '', label: 'Any' },
  { value: 'ppp', label: 'PPP' },
  { value: 'pppoe', label: 'PPPoE' },
  { value: 'pptp', label: 'PPTP' },
  { value: 'l2tp', label: 'L2TP' },
  { value: 'ovpn', label: 'OpenVPN' },
  { value: 'sstp', label: 'SSTP' },
]

interface SecretsTabProps {
  routerId: string
}

export function SecretsTab({ routerId }: SecretsTabProps) {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editingSecret, setEditingSecret] = useState<PPPSecret | null>(null)
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { data: secrets = [], isLoading } = useQuery({
    queryKey: ['ppp-secrets', routerId],
    queryFn: () => pppApi.getSecrets(routerId),
    enabled: !!routerId,
  })

  const { data: profiles = [] } = useQuery<PPPProfile[]>({
    queryKey: ['ppp-profiles', routerId],
    queryFn: () => pppApi.getProfiles(routerId),
    enabled: !!routerId,
  })

  const profileOptions = useMemo(() => [
    { value: '', label: 'Default' },
    ...profiles.map((p) => ({ value: p.name, label: p.name })),
  ], [profiles])

  const { register, handleSubmit, reset, formState: { errors } } = useForm<SecretForm>({
    resolver: zodResolver(secretSchema),
  })

  const createMutation = useMutation({
    mutationFn: (data: SecretForm) => pppApi.createSecret(routerId, data),
    onSuccess: () => { toast.success('Secret created'); queryClient.invalidateQueries({ queryKey: ['ppp-secrets', routerId] }); closeModal() },
    onError: (err: any) => toast.error(err.message || 'Failed to create'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: SecretForm }) => pppApi.updateSecret(routerId, id, data),
    onSuccess: () => { toast.success('Secret updated'); queryClient.invalidateQueries({ queryKey: ['ppp-secrets', routerId] }); closeModal() },
    onError: (err: any) => toast.error(err.message || 'Failed to update'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => pppApi.deleteSecret(routerId, id),
    onSuccess: () => { toast.success('Secret deleted'); queryClient.invalidateQueries({ queryKey: ['ppp-secrets', routerId] }) },
    onError: (err: any) => toast.error(err.message || 'Failed to delete'),
  })

  const toggleMutation = useMutation({
    mutationFn: ({ id, disabled }: { id: string; disabled: boolean }) =>
      disabled ? pppApi.enableSecret(routerId, id) : pppApi.disableSecret(routerId, id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['ppp-secrets', routerId] }),
    onError: (err: any) => toast.error(err.message || 'Failed to toggle'),
  })

  function openCreate() {
    setEditingSecret(null)
    reset({ name: '', password: '', profile: '', service: '', localAddress: '', remoteAddress: '', comment: '' })
    setModalOpen(true)
  }

  function openEdit(secret: PPPSecret) {
    setEditingSecret(secret)
    reset({ name: secret.name, password: '', profile: secret.profile || '', service: secret.service || '', localAddress: secret.localAddress || '', remoteAddress: secret.remoteAddress || '', comment: secret.comment || '' })
    setModalOpen(true)
  }

  function closeModal() { setModalOpen(false); setEditingSecret(null); reset() }

  function onSubmit(data: SecretForm) {
    if (editingSecret) updateMutation.mutate({ id: editingSecret.id, data })
    else createMutation.mutate(data)
  }

  const filtered = useMemo(() =>
    secrets.filter((s) =>
      s.name.toLowerCase().includes(search.toLowerCase()) ||
      (s.profile || '').toLowerCase().includes(search.toLowerCase()) ||
      (s.comment || '').toLowerCase().includes(search.toLowerCase())
    ), [secrets, search]
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
              placeholder="Search name, profile, comment..."
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
          <Button variant="primary" size="sm" leftIcon={<UserPlus className="w-4 h-4" />} onClick={openCreate}>
            Add Secret
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
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Profile</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Service</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Comment</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Status</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
              {paginatedRows.map((secret) => (
                <tr key={secret.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <div className="w-7 h-7 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center shrink-0">
                        <Network className="w-3.5 h-3.5 text-primary-600 dark:text-primary-400" />
                      </div>
                      <span className="font-semibold text-gray-900 dark:text-white text-sm">{secret.name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    {secret.profile ? <Badge variant="primary">{secret.profile}</Badge> : <span className="text-gray-400 text-xs">default</span>}
                  </td>
                  <td className="px-4 py-3">
                    {secret.service ? <Badge variant="info">{secret.service}</Badge> : <span className="text-gray-400 text-xs">any</span>}
                  </td>
                  <td className="px-4 py-3 text-xs text-gray-500 dark:text-gray-400">{secret.comment || '-'}</td>
                  <td className="px-4 py-3">
                    <button
                      onClick={() => toggleMutation.mutate({ id: secret.id, disabled: secret.disabled })}
                      className="flex items-center gap-1.5 text-xs transition-colors"
                    >
                      {secret.disabled ? (
                        <><ToggleLeft className="w-4 h-4 text-gray-400" /><span className="text-gray-400">Disabled</span></>
                      ) : (
                        <><ToggleRight className="w-4 h-4 text-success-500" /><span className="text-success-600 dark:text-success-400">Enabled</span></>
                      )}
                    </button>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1">
                      <Button variant="ghost" size="xs" leftIcon={<Edit2 className="w-3.5 h-3.5" />} onClick={() => openEdit(secret)}>Edit</Button>
                      <Button variant="ghost" size="xs" leftIcon={<Trash2 className="w-3.5 h-3.5 text-danger-500" />}
                        onClick={() => { if (confirm(`Delete secret "${secret.name}"?`)) deleteMutation.mutate(secret.id) }}>Delete</Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {!isLoading && filtered.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400 flex flex-col items-center gap-3">
            <Network className="w-10 h-10 text-gray-300 dark:text-gray-600" />
            <span>No PPP secrets found</span>
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
        title={editingSecret ? 'Edit Secret' : 'Add Secret'}
        subtitle={editingSecret ? `Editing: ${editingSecret.name}` : 'Create a new PPP secret'}
        footer={
          <div className="flex justify-end gap-2">
            <Button variant="outline" size="sm" onClick={closeModal}>Cancel</Button>
            <Button variant="primary" size="sm" isLoading={isSaving} onClick={handleSubmit(onSubmit)}>
              {editingSecret ? 'Save' : 'Create'}
            </Button>
          </div>
        }
      >
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Name" {...register('name')} error={errors.name?.message} placeholder="username" />
          <Input label="Password" type="password" {...register('password')} placeholder="••••••••" />
          <div className="grid grid-cols-2 gap-4">
            <Select label="Profile" options={profileOptions} {...register('profile')} />
            <Select label="Service" options={SERVICE_OPTIONS} {...register('service')} />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <Input label="Local Address" {...register('localAddress')} placeholder="192.168.1.1" />
            <Input label="Remote Address" {...register('remoteAddress')} placeholder="192.168.1.0/24" />
          </div>
          <Input label="Comment" {...register('comment')} placeholder="Optional comment" />
        </form>
      </Modal>
    </div>
  )
}
