// Hotspot Page with Tabs - Users, Profiles, Active, Hosts
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
  Users,
  PieChart,
  Wifi,
  Laptop,
  Eye,
  EyeOff,
  X,
  Ticket,
  Printer,
  Filter,
  UserPlus,
} from 'lucide-react'
import toast from 'react-hot-toast'
import { clsx } from 'clsx'

import { Card, Button, Input, Badge, Modal, Select } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'
import type { HotspotUser, UserProfile, HotspotActive, HotspotHost } from '../../types'

const tabs = [
  { id: 'users', name: 'Users', icon: Users },
  { id: 'profiles', name: 'User Profile', icon: PieChart },
  { id: 'active', name: 'Active', icon: Wifi },
  { id: 'hosts', name: 'Hosts', icon: Laptop },
]

// User Schema
const userSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  password: z.string().optional(),
  profile: z.string().min(1, 'Profile is required'),
  server: z.string().optional(),
  macAddress: z.string().optional(),
  comment: z.string().optional(),
})

type UserForm = z.infer<typeof userSchema>

// Profile Schema
const profileSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  sharedUsers: z.string().optional(),
  rateLimit: z.string().optional(),
})

interface ProfileForm {
  name: string
  sharedUsers?: string
  rateLimit?: string
}

export function HotspotPage() {
  const [activeTab, setActiveTab] = useState('users')
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  // Data queries
  const { data: users, isLoading: usersLoading, refetch: refetchUsers } = useQuery({
    queryKey: ['users', routerId],
    queryFn: () => hotspotApi.getUsers(routerId),
  })

  const { data: profiles, isLoading: profilesLoading, refetch: refetchProfiles } = useQuery({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
  })

  const { data: activeUsers, isLoading: activeLoading, refetch: refetchActive } = useQuery({
    queryKey: ['active', routerId],
    queryFn: () => hotspotApi.getActive(routerId),
  })

  const { data: hosts, isLoading: hostsLoading, refetch: refetchHosts } = useQuery({
    queryKey: ['hosts', routerId],
    queryFn: () => hotspotApi.getHosts(routerId),
  })

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      {/* Header with Tabs */}
      <Card>
        <Card.Header className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <Wifi className="w-5 h-5 text-primary-500" />
            <h1 className="text-xl font-bold text-gray-900 dark:text-white">Hotspot</h1>
          </div>
          
          {/* Tabs */}
          <div className="flex flex-wrap gap-1">
            {tabs.map((tab) => {
              const Icon = tab.icon
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={clsx(
                    'flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
                    activeTab === tab.id
                      ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                      : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700 hover:text-gray-900 dark:hover:text-gray-100'
                  )}
                >
                  <Icon className="w-4 h-4" />
                  <span className="hidden sm:inline">{tab.name}</span>
                </button>
              )
            })}
          </div>
        </Card.Header>
      </Card>

      {/* Tab Content */}
      <div className="mt-6">
        {activeTab === 'users' && (
          <UsersTab
            users={users || []}
            profiles={profiles || []}
            isLoading={usersLoading}
            routerId={routerId}
            queryClient={queryClient}
            refetch={refetchUsers}
          />
        )}
        {activeTab === 'profiles' && (
          <ProfilesTab
            profiles={profiles || []}
            isLoading={profilesLoading}
            routerId={routerId}
            queryClient={queryClient}
            refetch={refetchProfiles}
          />
        )}
        {activeTab === 'active' && (
          <ActiveTab
            activeUsers={activeUsers || []}
            isLoading={activeLoading}
            refetch={refetchActive}
          />
        )}
        {activeTab === 'hosts' && (
          <HostsTab
            hosts={hosts || []}
            isLoading={hostsLoading}
            refetch={refetchHosts}
          />
        )}
      </div>
    </motion.div>
  )
}

// Users Tab Component
function UsersTab({
  users,
  profiles,
  isLoading,
  routerId,
  queryClient,
  refetch,
}: {
  users: HotspotUser[]
  profiles: UserProfile[]
  isLoading: boolean
  routerId: string
  queryClient: any
  refetch: () => void
}) {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedProfile, setSelectedProfile] = useState('')
  const [selectedComment, setSelectedComment] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<HotspotUser | null>(null)
  const [showPassword, setShowPassword] = useState<Record<string, boolean>>({})
  const [showPrintMenu, setShowPrintMenu] = useState(false)

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

  const togglePassword = (userId: string) => {
    setShowPassword((prev) => ({ ...prev, [userId]: !prev[userId] }))
  }

  // Get unique comments for filter
  const uniqueComments = useMemo(() => {
    const comments = new Set<string>()
    users.forEach(user => {
      if (user.comment) comments.add(user.comment)
    })
    return Array.from(comments).sort()
  }, [users])

  const filteredUsers = users.filter((user) => {
    const profileStr = typeof user.profile === 'string' ? user.profile : user.profile?.name || ''
    const matchesSearch = 
      user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      profileStr.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.comment?.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesProfile = !selectedProfile || profileStr === selectedProfile
    const matchesComment = !selectedComment || user.comment === selectedComment
    return matchesSearch && matchesProfile && matchesComment
  })

  const clearFilters = () => {
    setSearchQuery('')
    setSelectedProfile('')
    setSelectedComment('')
  }

  const hasFilters = searchQuery || selectedProfile || selectedComment

  const handlePrint = (size: 'small' | 'default') => {
    localStorage.setItem('printUsers', JSON.stringify(filteredUsers))
    localStorage.setItem('printSize', size)
    window.open('/vouchers/print', '_blank')
  }

  return (
    <Card>
      <Card.Body className="p-3">
        {/* Toolbar - Compact Design */}
        <div className="flex flex-col xl:flex-row gap-3 mb-4">
          {/* Left Group: Total & Refresh */}
          <div className="flex items-center gap-2">
            <span className="px-3 py-2 bg-gray-100 dark:bg-dark-700 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300">
              Total: {users.length}
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
            <Select
              value={selectedProfile}
              onChange={(e) => setSelectedProfile(e.target.value)}
              options={[
                { value: '', label: 'All Profiles' },
                ...profiles.map((p) => ({ value: p.name, label: p.name })),
              ]}
              className="w-full sm:w-40"
            />
            <Select
              value={selectedComment}
              onChange={(e) => setSelectedComment(e.target.value)}
              options={[
                { value: '', label: 'All Comments' },
                ...uniqueComments.map((c) => ({ value: c, label: c })),
              ]}
              className="w-full sm:w-40"
            />
            <button
              onClick={clearFilters}
              disabled={!hasFilters}
              className={clsx(
                "px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-1",
                hasFilters
                  ? "bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600"
                  : "bg-gray-50 dark:bg-dark-800 text-gray-400 cursor-not-allowed"
              )}
            >
              <Filter className="w-4 h-4" />
              <span className="hidden sm:inline">Clear</span>
            </button>
          </div>

          {/* Right Group: Action Buttons */}
          <div className="flex items-center gap-2">
            <Button onClick={() => openModal()} size="sm" className="whitespace-nowrap">
              <UserPlus className="w-4 h-4 mr-1" />
              Add
            </Button>
            <Button 
              variant="secondary" 
              size="sm" 
              onClick={() => window.location.href = '/vouchers/generate'}
              className="whitespace-nowrap"
            >
              <Ticket className="w-4 h-4 mr-1" />
              Generate
            </Button>
            <div className="relative">
              <button
                onClick={() => setShowPrintMenu(!showPrintMenu)}
                disabled={filteredUsers.length === 0}
                className={clsx(
                  "px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-1",
                  filteredUsers.length > 0
                    ? "bg-success-50 dark:bg-success-900/20 text-success-600 dark:text-success-400 hover:bg-success-100 dark:hover:bg-success-900/30"
                    : "bg-gray-100 dark:bg-dark-700 text-gray-400 cursor-not-allowed"
                )}
              >
                <Printer className="w-4 h-4" />
                <span className="hidden sm:inline">Print</span>
              </button>
              {showPrintMenu && filteredUsers.length > 0 && (
                <>
                  <div className="fixed inset-0 z-10" onClick={() => setShowPrintMenu(false)} />
                  <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-dark-800 rounded-lg shadow-lg border border-gray-200 dark:border-dark-700 z-20 py-1">
                    <button
                      onClick={() => { handlePrint('small'); setShowPrintMenu(false); }}
                      className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 flex items-center gap-2"
                    >
                      <Printer className="w-4 h-4" />
                      Print Small
                    </button>
                    <button
                      onClick={() => { handlePrint('default'); setShowPrintMenu(false); }}
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

        {/* Table */}
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-dark-700">
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
                {filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                    <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{user.server || 'all'}</td>
                    <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{user.name}</td>
                    <td className="px-4 py-3 text-sm">
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-gray-600 dark:text-gray-400">
                          {showPassword[user.id] ? user.password : '••••••'}
                        </span>
                        <button
                          onClick={() => togglePassword(user.id)}
                          className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                        >
                          {showPassword[user.id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                        </button>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <Badge variant="primary">
                        {typeof user.profile === 'string' ? user.profile : user.profile?.name || '-'}
                      </Badge>
                    </td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{user.macAddress || '-'}</td>
                    <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{user.uptime || '-'}</td>
                    <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{user.bytesIn || '-'}</td>
                    <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{user.bytesOut || '-'}</td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{user.comment || '-'}</td>
                    <td className="px-4 py-3 text-center">
                      <div className="flex items-center justify-center gap-1">
                        <button
                          onClick={() => openModal(user)}
                          className="p-1.5 rounded-lg text-warning-600 hover:bg-warning-50 dark:hover:bg-warning-900/20"
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
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filteredUsers.length === 0 && (
              <div className="p-8 text-center text-gray-500 dark:text-gray-400">No users found</div>
            )}
          </div>
        )}
      </Card.Body>

      {/* Add/Edit Modal */}
      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingUser ? 'Edit User' : 'Add User'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
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
            {...register('profile')}
            options={profiles?.map((p) => ({ value: p.name, label: p.name })) || []}
            error={errors.profile?.message}
          />
          <Input
            label="MAC Address"
            {...register('macAddress')}
            error={errors.macAddress?.message}
          />
          <Input
            label="Comment"
            {...register('comment')}
            error={errors.comment?.message}
          />
          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" isLoading={createMutation.isPending || updateMutation.isPending}>
              {editingUser ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Modal>
    </Card>
  )
}

// Profiles Tab Component
function ProfilesTab({
  profiles,
  isLoading,
  routerId,
  queryClient,
  refetch,
}: {
  profiles: UserProfile[]
  isLoading: boolean
  routerId: string
  queryClient: any
  refetch: () => void
}) {
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProfile, setEditingProfile] = useState<UserProfile | null>(null)

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

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
  })

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
        sharedUsers: profile.sharedUsers?.toString(),
        rateLimit: profile.rateLimit,
      })
    } else {
      setEditingProfile(null)
      reset({ name: '', sharedUsers: undefined, rateLimit: '' })
    }
    setIsModalOpen(true)
  }

  const filteredProfiles = profiles.filter((profile) =>
    profile.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <Card>
      <Card.Header>
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Total: <span className="font-medium text-gray-900 dark:text-white">{profiles.length}</span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <button
              onClick={() => refetch()}
              className="p-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700"
              title="Refresh"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Filter..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 pr-8 py-2 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
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
            <Button onClick={() => openModal()} size="sm">
              <UserPlus className="w-4 h-4 mr-1" />
              Add Profile
            </Button>
            <Button variant="secondary" size="sm" onClick={() => window.location.href = '/hotspot/users'}>
              <UserPlus className="w-4 h-4 mr-1" />
              Add User
            </Button>
            <Button variant="secondary" size="sm" onClick={() => window.location.href = '/vouchers/generate'}>
              <Ticket className="w-4 h-4 mr-1" />
              Generate
            </Button>
          </div>
        </div>
      </Card.Header>

      <Card.Body className="p-0">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-dark-700">
                <tr>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase"></th>
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
                {filteredProfiles.map((profile) => (
                  <tr key={profile.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                    <td className="px-4 py-3 text-center">
                      <div className="flex items-center justify-center gap-1">
                        <button
                          onClick={() => openModal(profile)}
                          className="p-1.5 rounded-lg text-warning-600 hover:bg-warning-50 dark:hover:bg-warning-900/20"
                        >
                          <Edit2 className="w-4 h-4" />
                        </button>
                        <button
                          onClick={() => {
                            if (confirm('Are you sure you want to delete this profile?')) {
                              deleteMutation.mutate(profile.id)
                            }
                          }}
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
            {filteredProfiles.length === 0 && (
              <div className="p-8 text-center text-gray-500 dark:text-gray-400">No profiles found</div>
            )}
          </div>
        )}
      </Card.Body>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingProfile ? 'Edit Profile' : 'Add Profile'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Name" {...register('name')} error={errors.name?.message} />
          <Input label="Shared Users" {...register('sharedUsers')} />
          <Input label="Rate Limit" {...register('rateLimit')} placeholder="e.g. 1M/1M" />
          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" isLoading={createMutation.isPending || updateMutation.isPending}>
              {editingProfile ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Modal>
    </Card>
  )
}

// Active Tab Component
function ActiveTab({
  activeUsers,
  isLoading,
  refetch,
}: {
  activeUsers: HotspotActive[]
  isLoading: boolean
  refetch: () => void
}) {
  const [searchQuery, setSearchQuery] = useState('')

  const filteredActive = activeUsers.filter((active) =>
    active.user.toLowerCase().includes(searchQuery.toLowerCase()) ||
    active.address.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <Card>
      <Card.Header>
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Total Active: <span className="font-medium text-gray-900 dark:text-white">{activeUsers.length}</span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <button
              onClick={() => refetch()}
              className="p-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700"
              title="Refresh"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Filter..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 pr-8 py-2 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
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
          </div>
        </div>
      </Card.Header>

      <Card.Body className="p-0">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-dark-700">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Server</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">User</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Address</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">MAC Address</th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Uptime</th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Bytes In</th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Bytes Out</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Time Left</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Login By</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Comment</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
                {filteredActive.map((active) => (
                  <tr key={active.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                    <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{active.server}</td>
                    <td className="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{active.user}</td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{active.address}</td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{active.macAddress}</td>
                    <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{active.uptime}</td>
                    <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{active.bytesIn}</td>
                    <td className="px-4 py-3 text-sm text-right text-gray-600 dark:text-gray-400">{active.bytesOut}</td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{active.sessionTimeLeft || '-'}</td>
                    <td className="px-4 py-3 text-sm">
                      <Badge variant="success">{active.loginBy}</Badge>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{active.comment || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filteredActive.length === 0 && (
              <div className="p-8 text-center text-gray-500 dark:text-gray-400">No active users</div>
            )}
          </div>
        )}
      </Card.Body>
    </Card>
  )
}

// Hosts Tab Component
function HostsTab({
  hosts,
  isLoading,
  refetch,
}: {
  hosts: HotspotHost[]
  isLoading: boolean
  refetch: () => void
}) {
  const [searchQuery, setSearchQuery] = useState('')
  const [filterType, setFilterType] = useState<'all' | 'authorized' | 'bypassed'>('all')

  const filteredHosts = hosts.filter((host) => {
    const matchesSearch = host.macAddress.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (host.address || '').toLowerCase().includes(searchQuery.toLowerCase())
    const matchesFilter = filterType === 'all' || (filterType === 'authorized' && host.authorized) || (filterType === 'bypassed' && host.bypassed)
    return matchesSearch && matchesFilter
  })

  return (
    <Card>
      <Card.Header>
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            Total Hosts: <span className="font-medium text-gray-900 dark:text-white">{hosts.length}</span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <button
              onClick={() => refetch()}
              className="p-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700"
              title="Refresh"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Filter..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 pr-8 py-2 text-sm rounded-lg border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
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
            <div className="flex items-center gap-1">
              <button
                onClick={() => setFilterType('all')}
                className={clsx(
                  'px-3 py-2 text-xs font-medium rounded-lg transition-colors',
                  filterType === 'all'
                    ? 'bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400'
                    : 'bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600'
                )}
              >
                All
              </button>
              <button
                onClick={() => setFilterType('authorized')}
                className={clsx(
                  'px-3 py-2 text-xs font-medium rounded-lg transition-colors',
                  filterType === 'authorized'
                    ? 'bg-success-100 dark:bg-success-900/30 text-success-600 dark:text-success-400'
                    : 'bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600'
                )}
              >
                A
              </button>
              <button
                onClick={() => setFilterType('bypassed')}
                className={clsx(
                  'px-3 py-2 text-xs font-medium rounded-lg transition-colors',
                  filterType === 'bypassed'
                    ? 'bg-warning-100 dark:bg-warning-900/30 text-warning-600 dark:text-warning-400'
                    : 'bg-gray-100 dark:bg-dark-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600'
                )}
              >
                P
              </button>
            </div>
          </div>
        </div>
      </Card.Header>

      <Card.Body className="p-0">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-dark-700">
                <tr>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase"></th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase"></th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">MAC Address</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Address</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">To Address</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Server</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase">Comment</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-dark-700">
                {filteredHosts.map((host) => (
                  <tr key={host.id} className="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                    <td className="px-4 py-3 text-center">
                      <button className="text-gray-400 hover:text-gray-600">
                        <Edit2 className="w-4 h-4" />
                      </button>
                    </td>
                    <td className="px-4 py-3 text-center">
                      <button className="text-gray-400 hover:text-gray-600">
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-900 dark:text-white">{host.macAddress}</td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{host.address}</td>
                    <td className="px-4 py-3 text-sm font-mono text-gray-600 dark:text-gray-400">{host.toAddress || '-'}</td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{host.server || '-'}</td>
                    <td className="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">{host.comment || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filteredHosts.length === 0 && (
              <div className="p-8 text-center text-gray-500 dark:text-gray-400">No hosts found</div>
            )}
          </div>
        )}
      </Card.Body>
    </Card>
  )
}
