import { useState, useMemo } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import {
  Edit2, Trash2, RefreshCw, Eye, EyeOff, Printer, Ticket, UserPlus, AlertTriangle, Activity,
} from 'lucide-react'
import toast from 'react-hot-toast'
import type { ColumnDef } from '@tanstack/react-table'

import { Card, Button, Input, Badge, Modal, Select, DataTable } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'
import type { HotspotUser } from '../../types'
import { Link } from 'react-router-dom'
import { toggleApiDebug } from '../../api/axios'

const userSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  password: z.string().optional(),
  profile: z.string().min(1, 'Profile is required'),
  server: z.string().optional(),
  macAddress: z.string().optional(),
  timeLimit: z.string().optional(),
  dataLimit: z.string().optional(),
  comment: z.string().optional(),
})

type UserForm = z.infer<typeof userSchema>

export function UsersPage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '')

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<HotspotUser | null>(null)
  const [showPassword, setShowPassword] = useState<Record<string, boolean>>({})
  const [showPrintModal, setShowPrintModal] = useState(false)

  const { data: users, isLoading, error: usersError, refetch } = useQuery({
    queryKey: ['users', routerId],
    queryFn: () => hotspotApi.getUsers(routerId),
    enabled: !!selectedRouter,
    retry: 2,
  })

  const { data: profiles } = useQuery({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
    enabled: !!selectedRouter,
    retry: 2,
  })

  const createMutation = useMutation({
    mutationFn: (data: UserForm) => hotspotApi.createUser(routerId, data),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['users', routerId] }); toast.success('User created'); setIsModalOpen(false) },
    onError: (e: any) => toast.error(e.message),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UserForm }) => hotspotApi.updateUser(routerId, id, data),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['users', routerId] }); toast.success('User updated'); setIsModalOpen(false); setEditingUser(null) },
    onError: (e: any) => toast.error(e.message),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => hotspotApi.deleteUser(routerId, id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['users', routerId] }); toast.success('User deleted') },
    onError: (e: any) => toast.error(e.message),
  })

  const { register, handleSubmit, reset, formState: { errors } } = useForm<UserForm>({ resolver: zodResolver(userSchema) })

  const onSubmit = (data: UserForm) => {
    if (editingUser) { updateMutation.mutate({ id: editingUser.id, data }) }
    else { createMutation.mutate(data) }
  }

  const openModal = (user?: HotspotUser) => {
    if (user) {
      setEditingUser(user)
      reset({ name: user.name, profile: typeof user.profile === 'string' ? user.profile : user.profile?.name || '', server: user.server, macAddress: user.macAddress, timeLimit: user.limitUptime || '', dataLimit: '', comment: user.comment })
    } else {
      setEditingUser(null); reset({})
    }
    setIsModalOpen(true)
  }

  const togglePassword = (userId: string) =>
    setShowPassword(prev => ({ ...prev, [userId]: !prev[userId] }))

  const formatBytes = (bytes: number) => {
    if (!bytes) return '0 B'
    const k = 1024, sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const handlePrint = (size: 'small' | 'default') => {
    localStorage.setItem('printUsers', JSON.stringify(users || []))
    localStorage.setItem('printSize', size)
    window.open('/vouchers/print', '_blank')
  }

  // --- TanStack Table columns ---
  const columns = useMemo<ColumnDef<HotspotUser, any>[]>(() => [
    {
      accessorKey: 'name',
      header: 'User',
      cell: ({ row }) => {
        const u = row.original
        return (
          <div>
            <p className="font-semibold text-gray-900 dark:text-white">{u.name}</p>
            {u.password && (
              <div className="flex items-center gap-1.5 text-xs text-gray-400 mt-0.5">
                <span className="font-mono">{showPassword[u.id] ? u.password : '••••••'}</span>
                <button onClick={() => togglePassword(u.id)} className="hover:text-gray-600 dark:hover:text-gray-200">
                  {showPassword[u.id] ? <EyeOff className="w-3 h-3" /> : <Eye className="w-3 h-3" />}
                </button>
              </div>
            )}
          </div>
        )
      },
    },
    {
      accessorKey: 'profile',
      header: 'Profile',
      cell: ({ getValue }) => {
        const p = getValue()
        const name = typeof p === 'string' ? p : p?.name || '-'
        return <Badge variant="primary">{name}</Badge>
      },
    },
    {
      accessorKey: 'macAddress',
      header: 'MAC Address',
      cell: ({ getValue }) => (
        <span className="font-mono text-xs text-gray-500 dark:text-gray-400">{getValue() || '-'}</span>
      ),
    },
    {
      accessorKey: 'uptime',
      header: 'Uptime',
      cell: ({ getValue }) => <span className="text-gray-600 dark:text-gray-400 text-xs">{getValue() || '-'}</span>,
    },
    {
      id: 'bytes',
      header: 'Bytes',
      accessorFn: (row) => row.bytesIn + row.bytesOut,
      cell: ({ row }) => (
        <div className="text-xs space-x-1">
          <span className="text-success-600 dark:text-success-400">↓ {formatBytes(row.original.bytesIn)}</span>
          <span className="text-gray-300">|</span>
          <span className="text-primary-600 dark:text-primary-400">↑ {formatBytes(row.original.bytesOut)}</span>
        </div>
      ),
    },
    {
      accessorKey: 'comment',
      header: 'Comment',
      cell: ({ getValue }) => (
        <span className="text-gray-500 dark:text-gray-400 text-xs">{getValue() || '-'}</span>
      ),
    },
    {
      id: 'actions',
      header: 'Actions',
      enableSorting: false,
      cell: ({ row }) => {
        const u = row.original
        return (
          <div className="flex items-center gap-1">
            <button
              onClick={() => openModal(u)}
              className="p-1.5 rounded-lg text-warning-600 hover:bg-warning-50 dark:hover:bg-warning-900/20 transition-colors"
              title="Edit"
            >
              <Edit2 className="w-3.5 h-3.5" />
            </button>
            <button
              onClick={() => { if (confirm('Delete this user?')) deleteMutation.mutate(u.id) }}
              className="p-1.5 rounded-lg text-danger-600 hover:bg-danger-50 dark:hover:bg-danger-900/20 transition-colors"
              title="Delete"
            >
              <Trash2 className="w-3.5 h-3.5" />
            </button>
          </div>
        )
      },
    },
  ], [showPassword, deleteMutation, profiles])

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <UserPlus className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">Silahkan pilih router terlebih dahulu untuk mengelola hotspot users.</p>
        <Link to="/routers" className="px-5 py-2.5 rounded-xl bg-primary-500 text-white font-medium hover:bg-primary-600 transition-colors">
          Manage Routers
        </Link>
      </div>
    )
  }

  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="space-y-4">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Hotspot Users</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {selectedRouter.name} &mdash; {users?.length || 0} users registered
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => refetch()}
            className="p-2 rounded-xl bg-gray-100 dark:bg-dark-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600 transition-colors"
            title="Refresh"
          >
            <RefreshCw className="w-4 h-4" />
          </button>
          <Button variant="ghost" size="sm" onClick={() => toggleApiDebug()} title="Toggle Debug">
            <Activity className="w-4 h-4" />
          </Button>
          <Button variant="secondary" size="sm" onClick={() => window.location.href = '/vouchers/generate'}>
            <Ticket className="w-4 h-4 mr-1" /> Generate
          </Button>
          {/* Print button */}
          <div className="relative">
            <button
              onClick={() => setShowPrintModal(!showPrintModal)}
              disabled={!users?.length}
              className="px-3 py-1.5 rounded-xl text-sm font-medium bg-success-50 dark:bg-success-900/20 text-success-600 dark:text-success-400 hover:bg-success-100 dark:hover:bg-success-900/30 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center gap-1.5"
            >
              <Printer className="w-4 h-4" /> Print
            </button>
            {showPrintModal && users?.length && (
              <>
                <div className="fixed inset-0 z-10" onClick={() => setShowPrintModal(false)} />
                <div className="absolute right-0 mt-2 w-44 bg-white dark:bg-dark-800 rounded-xl shadow-lg border border-gray-200 dark:border-dark-700 z-20 py-1 overflow-hidden">
                  {['small', 'default'].map((s) => (
                    <button key={s} onClick={() => { handlePrint(s as any); setShowPrintModal(false) }} className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-dark-700 flex items-center gap-2">
                      <Printer className="w-4 h-4" /> Print {s === 'small' ? 'Small' : 'Default'}
                    </button>
                  ))}
                </div>
              </>
            )}
          </div>
          <Button onClick={() => openModal()} leftIcon={<UserPlus className="w-4 h-4" />}>Add</Button>
        </div>
      </div>

      {/* Error Banner */}
      {usersError && !isLoading && (
        <Card className="bg-danger-50 dark:bg-danger-900/20 border-danger-200 dark:border-danger-800">
          <Card.Body className="flex items-start gap-3">
            <AlertTriangle className="w-5 h-5 text-danger-600 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <h3 className="font-semibold text-danger-900 dark:text-danger-100">Failed to Load Users</h3>
              <p className="text-sm text-danger-700 dark:text-danger-300 mt-1">
                {usersError instanceof Error ? usersError.message : 'Unknown error'}
              </p>
              <Button variant="ghost" size="sm" onClick={() => refetch()} className="mt-2">
                <RefreshCw className="w-4 h-4 mr-1" /> Retry
              </Button>
            </div>
          </Card.Body>
        </Card>
      )}

      {/* DataTable */}
      <Card>
        <Card.Body>
          <DataTable
            data={users || []}
            columns={columns}
            isLoading={isLoading}
            searchPlaceholder="Search users by name, profile, comment..."
            emptyMessage="No users found"
          />
        </Card.Body>
      </Card>

      {/* Add/Edit Modal */}
      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingUser ? 'Edit User' : 'Add User'} footer={
        <div className="flex justify-end gap-2">
          <Button variant="ghost" onClick={() => setIsModalOpen(false)}>Cancel</Button>
          <Button onClick={handleSubmit(onSubmit)} isLoading={createMutation.isPending || updateMutation.isPending}>
            {editingUser ? 'Update' : 'Create'}
          </Button>
        </div>
      }>
        <form className="space-y-4">
          <Input label="Name" {...register('name')} error={errors.name?.message} />
          <Input label="Password" type="password" {...register('password')} />
          <Select label="Profile" options={profiles?.map(p => ({ value: p.name, label: p.name })) || []} {...register('profile')} />
          <Input label="MAC Address" placeholder="AA:BB:CC:DD:EE:FF" {...register('macAddress')} />
          <div className="grid grid-cols-2 gap-4">
            <Input label="Time Limit" placeholder="1h, 30m, 1d" {...register('timeLimit')} />
            <Input label="Data Limit" placeholder="100M, 1G" {...register('dataLimit')} />
          </div>
          <Input label="Comment" {...register('comment')} />
        </form>
      </Modal>
    </motion.div>
  )
}

