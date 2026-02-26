import { useState, useMemo } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import {
  Search,
  Edit2,
  Trash2,
  RefreshCw,
  Eye,
  EyeOff,
  Printer,
  Ticket,
  Filter,
  X,
  UserPlus,
  AlertTriangle,
  Activity,
} from 'lucide-react'
import toast from 'react-hot-toast'
import { clsx } from 'clsx'

import { Card, Button, Input, Badge, Modal, Select } from '../../components/ui'
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

  if (import.meta.env.DEV) {
    console.log('[UsersPage] Rendered - routerId:', routerId, 'selectedRouter:', selectedRouter)
  }

  const [searchQuery, setSearchQuery] = useState('')
  const [selectedProfile, setSelectedProfile] = useState('')
  const [selectedComment, setSelectedComment] = useState('')
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

  // Show error toast if there's a connection error
  if (usersError && !isLoading) {
    const errorMsg = usersError instanceof Error ? usersError.message : 'Failed to load users'
    if (errorMsg.includes('connection') || errorMsg.includes('Network Error') || errorMsg.includes('timeout')) {
      toast.error('Cannot connect to router. Please check your network settings.', { id: 'users-error', duration: 5000 })
    }
  }

  // Get unique comments for filter
  const uniqueComments = useMemo(() => {
    if (!users) return []
    const comments = new Set<string>()
    users.forEach(user => {
      if (user.comment) comments.add(user.comment)
    })
    return Array.from(comments).sort()
  }, [users])

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

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<UserForm>({
    resolver: zodResolver(userSchema),
  })

  const onSubmit = (data: UserForm) => {
    if (editingUser) {
      updateMutation.mutate({ id: editingUser.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const openModal = (user?: HotspotUser) => {
    if (user) {
      setEditingUser(user)
      reset({
        name: user.name,
        profile: typeof user.profile === 'string' ? user.profile : user.profile?.name || '',
        server: user.server,
        macAddress: user.macAddress,
        timeLimit: user.limitUptime || '',
        dataLimit: '',
        comment: user.comment,
      })
    } else {
      setEditingUser(null)
      reset({})
    }
    setIsModalOpen(true)
  }

  const togglePassword = (userId: string) => {
    setShowPassword((prev) => ({ ...prev, [userId]: !prev[userId] }))
  }

  // Filter users
  const filteredUsers = users?.filter((user) => {
    const profileStr = typeof user.profile === 'string' ? user.profile : user.profile?.name || ''
    const matchesSearch =
      user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      profileStr.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.comment?.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesProfile = !selectedProfile || profileStr === selectedProfile
    const matchesComment = !selectedComment || user.comment === selectedComment
    return matchesSearch && matchesProfile && matchesComment
  })

  // Users to print (filtered)
  const usersToPrint = filteredUsers || []

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const clearFilters = () => {
    setSearchQuery('')
    setSelectedProfile('')
    setSelectedComment('')
  }

  const hasFilters = searchQuery || selectedProfile || selectedComment

  const handlePrint = (size: 'small' | 'default') => {
    // Store filtered users in localStorage for print page
    localStorage.setItem('printUsers', JSON.stringify(usersToPrint))
    localStorage.setItem('printSize', size)
    window.open('/vouchers/print', '_blank')
  }

  if (!selectedRouter) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
          <UserPlus className="w-8 h-8 text-primary-500" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Router Selected</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 max-w-sm">
          Silahkan pilih router terlebih dahulu untuk melihat dan mengelola hotspot users.
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
      className="space-y-4"
    >
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Hotspot Users</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Manage hotspot users and their credentials
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
        </div>
      </div>

      {/* Error Banner */}
      {usersError && !isLoading && (
        <Card className="bg-danger-50 dark:bg-danger-900/20 border-danger-200 dark:border-danger-800">
          <Card.Body className="flex items-start gap-3">
            <AlertTriangle className="w-5 h-5 text-danger-600 dark:text-danger-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <h3 className="font-semibold text-danger-900 dark:text-danger-100 mb-1">
                Failed to Load Users
              </h3>
              <p className="text-sm text-danger-700 dark:text-danger-300">
                {usersError instanceof Error ? usersError.message : 'An unknown error occurred'}
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

      {/* Toolbar - Compact Design like PHP version */}
      <Card>
        <Card.Body className="p-3">
          <div className="flex flex-col lg:flex-row gap-3">
            {/* Left Group: Total & Refresh */}
            <div className="flex items-center gap-2">
              <span className="px-3 py-2 bg-gray-100 dark:bg-dark-700 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300">
                Total: {users?.length || 0}
              </span>
              <button
                onClick={() => refetch()}
                className="p-2 rounded-lg bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600 transition-colors"
                title="Refresh"
              >
                <RefreshCw className="w-4 h-4" />
              </button>
            </div>

            {/* Middle Group: Search & Filters */}
            <div className="flex-1 flex flex-col sm:flex-row gap-2">
              {/* Search Input */}
              <div className="relative flex-1 min-w-[200px]">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Filter users..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full pl-9 pr-8 py-2 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
                {searchQuery && (
                  <button
                    onClick={() => setSearchQuery('')}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  >
                    <X className="w-4 h-4" />
                  </button>
                )}
              </div>

              {/* Profile Filter */}
              <Select
                value={selectedProfile}
                onChange={(e) => setSelectedProfile(e.target.value)}
                options={[
                  { value: '', label: 'All Profiles' },
                  ...(profiles?.map((p) => ({ value: p.name, label: p.name })) || []),
                ]}
                className="w-full sm:w-40"
              />

              {/* Comment Filter */}
              <Select
                value={selectedComment}
                onChange={(e) => setSelectedComment(e.target.value)}
                options={[
                  { value: '', label: 'All Comments' },
                  ...uniqueComments.map((c) => ({ value: c, label: c })),
                ]}
                className="w-full sm:w-40"
              />

              {/* Clear Filter Button */}
              <button
                onClick={clearFilters}
                disabled={!hasFilters}
                className={clsx(
                  "px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-1",
                  hasFilters
                    ? "bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600"
                    : "bg-gray-50 dark:bg-dark-800 text-gray-400 cursor-not-allowed"
                )}
                title="Clear all filters"
              >
                <Filter className="w-4 h-4" />
                <span className="hidden sm:inline">Clear</span>
              </button>
            </div>

            {/* Right Group: Action Buttons */}
            <div className="flex items-center gap-2">
              {/* Add User Button */}
              <Button
                onClick={() => openModal()}
                size="sm"
                className="whitespace-nowrap"
              >
                <UserPlus className="w-4 h-4 mr-1" />
                Add
              </Button>

              {/* Generate Button */}
              <Button
                variant="secondary"
                size="sm"
                onClick={() => window.location.href = '/vouchers/generate'}
                className="whitespace-nowrap"
              >
                <Ticket className="w-4 h-4 mr-1" />
                Generate
              </Button>

              {/* Print Dropdown */}
              <div className="relative">
                <button
                  onClick={() => setShowPrintModal(!showPrintModal)}
                  className={clsx(
                    "px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-1",
                    usersToPrint.length > 0
                      ? "bg-success-50 dark:bg-success-900/20 text-success-600 dark:text-success-400 hover:bg-success-100 dark:hover:bg-success-900/30"
                      : "bg-gray-100 dark:bg-dark-700 text-gray-400 cursor-not-allowed"
                  )}
                  disabled={usersToPrint.length === 0}
                >
                  <Printer className="w-4 h-4" />
                  <span className="hidden sm:inline">Print</span>
                </button>

                {/* Print Dropdown Menu */}
                {showPrintModal && usersToPrint.length > 0 && (
                  <>
                    <div
                      className="fixed inset-0 z-10"
                      onClick={() => setShowPrintModal(false)}
                    />
                    <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-dark-800 rounded-lg shadow-lg border border-gray-200 dark:border-dark-700 z-20 py-1">
                      <button
                        onClick={() => {
                          handlePrint('small')
                          setShowPrintModal(false)
                        }}
                        className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 flex items-center gap-2"
                      >
                        <Printer className="w-4 h-4" />
                        Print Small
                      </button>
                      <button
                        onClick={() => {
                          handlePrint('default')
                          setShowPrintModal(false)
                        }}
                        className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 flex items-center gap-2"
                      >
                        <Printer className="w-4 h-4" />
                        Print Default
                      </button>
                    </div>
                  </>
                )}
              </div>
            </div>
          </div>
        </Card.Body>
      </Card>

      {/* Table */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-dark-700">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">User</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Profile</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">MAC Address</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Uptime</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Bytes</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Comment</th>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="text-center py-8">
                    <RefreshCw className="w-6 h-6 animate-spin mx-auto text-gray-400" />
                  </td>
                </tr>
              ) : filteredUsers?.length === 0 ? (
                <tr>
                  <td colSpan={7} className="text-center py-8 text-gray-500">
                    No users found
                  </td>
                </tr>
              ) : (
                filteredUsers?.map((user, index) => (
                  <tr key={user.id || `user-${index}`} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                    <td className="px-4 py-3">
                      <div>
                        <p className="font-medium text-gray-900 dark:text-white">{user.name}</p>
                        {user.password && (
                          <div className="flex items-center gap-2 text-sm text-gray-500">
                            <span className="font-mono">
                              {showPassword[user.id] ? user.password : '••••••'}
                            </span>
                            <button
                              onClick={() => togglePassword(user.id)}
                              className="text-gray-400 hover:text-gray-600"
                            >
                              {showPassword[user.id] ? (
                                <EyeOff className="w-3 h-3" />
                              ) : (
                                <Eye className="w-3 h-3" />
                              )}
                            </button>
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant="primary">{typeof user.profile === 'string' ? user.profile : user.profile?.name || '-'}</Badge>
                    </td>
                    <td className="px-4 py-3 font-mono text-sm text-gray-600 dark:text-gray-400">{user.macAddress || '-'}</td>
                    <td className="px-4 py-3 text-gray-600 dark:text-gray-400">{user.uptime || '-'}</td>
                    <td className="px-4 py-3">
                      <div className="text-sm">
                        <span className="text-success-600 dark:text-success-400">↓ {formatBytes(user.bytesIn)}</span>
                        <span className="mx-1 text-gray-400">|</span>
                        <span className="text-primary-600 dark:text-primary-400">↑ {formatBytes(user.bytesOut)}</span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-gray-600 dark:text-gray-400">{user.comment || '-'}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center justify-center gap-1">
                        <button
                          onClick={() => openModal(user)}
                          className="p-1.5 rounded-lg text-warning-600 hover:bg-warning-50 dark:hover:bg-warning-900/20"
                          title="Edit"
                        >
                          <Edit2 className="w-4 h-4" />
                        </button>
                        <button
                          onClick={() => {
                            if (confirm('Are you sure you want to delete this user?')) {
                              deleteMutation.mutate(user.id)
                            }
                          }}
                          className="p-1.5 rounded-lg text-danger-600 hover:bg-danger-50 dark:hover:bg-danger-900/20"
                          title="Delete"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>

      {/* Add/Edit Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={editingUser ? 'Edit User' : 'Add User'}
        footer={
          <div className="flex justify-end gap-3">
            <Button variant="ghost" onClick={() => setIsModalOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmit(onSubmit)}
              isLoading={createMutation.isPending || updateMutation.isPending}
            >
              {editingUser ? 'Update' : 'Create'}
            </Button>
          </div>
        }
      >
        <form className="space-y-4">
          <Input
            label="Name"
            {...register('name')}
            error={errors.name?.message}
          />
          <Input
            label="Password"
            type="password"
            {...register('password')}
            error={errors.password?.message}
          />
          <Select
            label="Profile"
            options={profiles?.map((p) => ({ value: p.name, label: p.name })) || []}
            {...register('profile')}
          />
          <Input
            label="MAC Address"
            placeholder="AA:BB:CC:DD:EE:FF"
            {...register('macAddress')}
          />
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Time Limit"
              placeholder="mis: 1h, 30m, 1d"
              {...register('timeLimit')}
            />
            <Input
              label="Data Limit"
              placeholder="mis: 100M, 1G, 500K"
              {...register('dataLimit')}
            />
          </div>
          <Input
            label="Comment"
            {...register('comment')}
          />
        </form>
      </Modal>
    </motion.div>
  )
}
