import { useState, useMemo, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Search, Edit2, Trash2, RefreshCw, Eye, EyeOff, X, Ticket, Filter, UserPlus,
} from 'lucide-react'
import toast from 'react-hot-toast'
import { clsx } from 'clsx'

import { Button, Input, Badge, Modal, Select, Pagination } from '../../../components/ui'
import { hotspotApi } from '../../../api/hotspot'
import type { HotspotUser, UserProfile } from '../../../types'

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024, s = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + s[i]
}

const userSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  password: z.string().optional(),
  profile: z.string().min(1, 'Profile is required'),
  server: z.string().optional(),
  macAddress: z.string().optional(),
  comment: z.string().optional(),
})
type UserForm = z.infer<typeof userSchema>

interface UsersTabProps {
  routerId: string
}

export function UsersTab({ routerId }: UsersTabProps) {
  const queryClient = useQueryClient()
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedProfile, setSelectedProfile] = useState('')
  const [selectedComment, setSelectedComment] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<HotspotUser | null>(null)
  const [showPassword, setShowPassword] = useState<Record<string, boolean>>({})
  const [showPrintMenu, setShowPrintMenu] = useState(false)
  const [pageIndex, setPageIndex] = useState(0)
  const [pageSize, setPageSize] = useState(25)

  const { data: users = [], isLoading, refetch } = useQuery({
    queryKey: ['users', routerId],
    queryFn: () => hotspotApi.getUsers(routerId),
  })

  const { data: profiles = [] } = useQuery<UserProfile[]>({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
  })

  const createMutation = useMutation({
    mutationFn: (data: UserForm) => hotspotApi.createUser(routerId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users', routerId] })
      toast.success('User created successfully')
      setIsModalOpen(false)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UserForm }) =>
      hotspotApi.updateUser(routerId, id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users', routerId] })
      toast.success('User updated successfully')
      setIsModalOpen(false)
      setEditingUser(null)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => hotspotApi.deleteUser(routerId, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users', routerId] })
      toast.success('User deleted successfully')
    },
    onError: (error: any) => toast.error(error.message),
  })

  const { register, handleSubmit, reset, formState: { errors } } = useForm<UserForm>({
    resolver: zodResolver(userSchema),
  })

  const onSubmit = (data: UserForm) => {
    if (editingUser) updateMutation.mutate({ id: editingUser.id, data })
    else createMutation.mutate(data)
  }

  const openModal = (user?: HotspotUser) => {
    if (user) {
      setEditingUser(user)
      reset({
        name: user.name,
        profile: typeof user.profile === 'string' ? user.profile : user.profile?.name || '',
        password: '',
        macAddress: user.macAddress,
        comment: user.comment,
      })
    } else {
      setEditingUser(null)
      reset({ name: '', password: '', profile: '', server: '', macAddress: '', comment: '' })
    }
    setIsModalOpen(true)
  }

  const uniqueComments = useMemo(() => {
    const s = new Set<string>()
    users.forEach((u) => { if (u.comment) s.add(u.comment) })
    return Array.from(s).sort()
  }, [users])

  const filteredUsers = useMemo(() => users.filter((user) => {
    const profileStr = typeof user.profile === 'string' ? user.profile : user.profile?.name || ''
    const q = searchQuery.toLowerCase()
    return (
      (!searchQuery || user.name.toLowerCase().includes(q) || profileStr.toLowerCase().includes(q) || (user.comment || '').toLowerCase().includes(q)) &&
      (!selectedProfile || profileStr === selectedProfile) &&
      (!selectedComment || user.comment === selectedComment)
    )
  }), [users, searchQuery, selectedProfile, selectedComment])

  const pageCount = Math.ceil(filteredUsers.length / pageSize)
  const paginatedUsers = filteredUsers.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize)

  useEffect(() => { setPageIndex(0) }, [searchQuery, selectedProfile, selectedComment])

  const hasFilters = searchQuery || selectedProfile || selectedComment

  const handlePrint = (size: 'small' | 'default') => {
    localStorage.setItem('printUsers', JSON.stringify(filteredUsers))
    localStorage.setItem('printSize', size)
    window.open('/vouchers/print', '_blank')
  }

  return (
    <div className="flex flex-col flex-1 min-h-0">
      {/* Toolbar */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-100 dark:border-dark-700 bg-white dark:bg-dark-800">
        <div className="flex flex-col xl:flex-row gap-3">
          <div className="flex items-center gap-2">
            <span className="px-3 py-1.5 bg-gray-100 dark:bg-dark-700 rounded-lg text-xs font-medium text-gray-700 dark:text-gray-300">
              Total: {users.length}
            </span>
            <button
              onClick={() => refetch()}
              className="p-1.5 rounded-lg bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600 transition-colors"
              title="Refresh"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          </div>

          <div className="flex-1 flex flex-col sm:flex-row gap-2">
            <div className="relative flex-1 min-w-[160px]">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Filter users..."
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
            <Select
              value={selectedProfile}
              onChange={(e) => setSelectedProfile(e.target.value)}
              options={[{ value: '', label: 'All Profiles' }, ...profiles.map((p) => ({ value: p.name, label: p.name }))]}
              className="w-full sm:w-36"
            />
            <Select
              value={selectedComment}
              onChange={(e) => setSelectedComment(e.target.value)}
              options={[{ value: '', label: 'All Comments' }, ...uniqueComments.map((c) => ({ value: c, label: c }))]}
              className="w-full sm:w-36"
            />
            <button
              onClick={() => { setSearchQuery(''); setSelectedProfile(''); setSelectedComment('') }}
              disabled={!hasFilters}
              className={clsx(
                'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors flex items-center gap-1',
                hasFilters
                  ? 'bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600'
                  : 'bg-gray-50 dark:bg-dark-800 text-gray-400 cursor-not-allowed'
              )}
            >
              <Filter className="w-4 h-4" />
            </button>
          </div>

          <div className="flex items-center gap-2">
            <Button onClick={() => openModal()} size="sm">
              <UserPlus className="w-4 h-4 mr-1" />Add
            </Button>
            <Button variant="secondary" size="sm" onClick={() => window.location.href = '/vouchers/generate'}>
              <Ticket className="w-4 h-4 mr-1" />Generate
            </Button>
            <div className="relative">
              <button
                onClick={() => setShowPrintMenu(!showPrintMenu)}
                disabled={filteredUsers.length === 0}
                className={clsx(
                  'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors',
                  filteredUsers.length > 0
                    ? 'bg-success-50 dark:bg-success-900/20 text-success-600 dark:text-success-400 hover:bg-success-100'
                    : 'bg-gray-100 dark:bg-dark-700 text-gray-400 cursor-not-allowed'
                )}
              >
                Print
              </button>
              {showPrintMenu && filteredUsers.length > 0 && (
                <>
                  <div className="fixed inset-0 z-10" onClick={() => setShowPrintMenu(false)} />
                  <div className="absolute right-0 mt-1 w-40 bg-white dark:bg-dark-800 rounded-lg shadow-lg border border-gray-200 dark:border-dark-700 z-20 py-1">
                    <button onClick={() => { handlePrint('small'); setShowPrintMenu(false) }} className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700">Print Small</button>
                    <button onClick={() => { handlePrint('default'); setShowPrintMenu(false) }} className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700">Print Default</button>
                  </div>
                </>
              )}
            </div>
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
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Server</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Name</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Password</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Profile</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">MAC Address</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Uptime</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Bytes In</th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Bytes Out</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Comment</th>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
              {paginatedUsers.map((user) => (
                <tr key={user.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                  <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{user.server || 'all'}</td>
                  <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{user.name}</td>
                  <td className="px-4 py-3 text-sm">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-gray-600 dark:text-gray-400">
                        {showPassword[user.id] ? user.password : '••••••'}
                      </span>
                      <button onClick={() => setShowPassword((p) => ({ ...p, [user.id]: !p[user.id] }))} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                        {showPassword[user.id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                      </button>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-sm">
                    <Badge variant="primary">{typeof user.profile === 'string' ? user.profile : user.profile?.name || '-'}</Badge>
                  </td>
                  <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{user.macAddress || '-'}</td>
                  <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{user.uptime || '-'}</td>
                  <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{formatBytes(user.bytesIn || 0)}</td>
                  <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{formatBytes(user.bytesOut || 0)}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{user.comment || '-'}</td>
                  <td className="px-4 py-3 text-center">
                    <div className="flex items-center justify-center gap-1">
                      <button onClick={() => openModal(user)} className="p-1.5 rounded-lg text-warning-600 hover:bg-warning-50 dark:hover:bg-warning-900/20">
                        <Edit2 className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => { if (confirm('Delete this user?')) deleteMutation.mutate(user.id) }}
                        className="p-1.5 rounded-lg text-danger-600 hover:bg-danger-50 dark:hover:bg-danger-900/20"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {!isLoading && filteredUsers.length === 0 && (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">No users found</div>
        )}
      </div>

      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={pageCount}
        pageSize={pageSize}
        totalRows={filteredUsers.length}
        onPageChange={setPageIndex}
        onPageSizeChange={(size) => { setPageSize(size); setPageIndex(0) }}
      />

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingUser ? 'Edit User' : 'Add User'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Name" {...register('name')} error={errors.name?.message} />
          <Input label="Password" type="password" {...register('password')} error={errors.password?.message} />
          <Select label="Profile" {...register('profile')} options={profiles?.map((p) => ({ value: p.name, label: p.name })) || []} error={errors.profile?.message} />
          <Input label="MAC Address" {...register('macAddress')} error={errors.macAddress?.message} />
          <Input label="Comment" {...register('comment')} error={errors.comment?.message} />
          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>Cancel</Button>
            <Button type="submit" isLoading={createMutation.isPending || updateMutation.isPending}>
              {editingUser ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
