import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import {
  Router,
  Plus,
  Edit2,
  Trash2,
  RefreshCw,
  Globe,
  Shield,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { Card, Button, Input, Modal, Badge } from '../../components/ui'
import { routersApi } from '../../api/routers'
import { useRouterStore } from '../../stores/routerStore'
import type { Router as RouterType } from '../../types'

const routerSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  host: z.string().min(1, 'Host is required'),
  port: z.coerce.number().min(1).max(65535).default(8728),
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
  useSsl: z.boolean().default(false),
  description: z.string().optional(),
})

type RouterForm = z.infer<typeof routerSchema>

export function RoutersPage() {
  const queryClient = useQueryClient()
  const setSelectedRouter = useRouterStore((state) => state.setSelectedRouter)
  const selectedRouter = useRouterStore((state) => state.selectedRouter)

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingRouter, setEditingRouter] = useState<RouterType | null>(null)
  const [testingId, setTestingId] = useState<string | number | null>(null)

  const { data: routers, isLoading } = useQuery({
    queryKey: ['routers'],
    queryFn: () => routersApi.getAll(),
  })

  const createMutation = useMutation({
    mutationFn: (data: RouterForm) => routersApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routers'] })
      toast.success('Router added successfully')
      setIsModalOpen(false)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string | number; data: RouterForm }) =>
      routersApi.update(id.toString(), data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routers'] })
      toast.success('Router updated successfully')
      setIsModalOpen(false)
      setEditingRouter(null)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string | number) => routersApi.delete(id.toString()),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routers'] })
      toast.success('Router deleted successfully')
    },
    onError: (error: any) => toast.error(error.message),
  })

  const testMutation = useMutation({
    mutationFn: async (id: string | number) => {
      setTestingId(id)
      await routersApi.testConnection(id.toString())
    },
    onSuccess: () => {
      toast.success('Connection successful')
      setTestingId(null)
    },
    onError: (error: any) => {
      toast.error(error.message)
      setTestingId(null)
    },
  })

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(routerSchema),
    defaultValues: {
      port: 8728,
      useSsl: false,
    },
  })

  const onSubmit = (data: RouterForm) => {
    if (editingRouter) {
      updateMutation.mutate({ id: editingRouter.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const openModal = (router?: RouterType) => {
    if (router) {
      setEditingRouter(router)
      reset({
        name: router.name,
        host: router.host,
        port: router.port,
        // username not shown in edit form
        password: '', // Don't show password
        useSsl: router.useSsl,
        description: router.description,
      })
    } else {
      setEditingRouter(null)
      reset({
        port: 8728,
        useSsl: false,
      })
    }
    setIsModalOpen(true)
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
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Routers</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Manage your MikroTik routers
          </p>
        </div>
        <Button onClick={() => openModal()} leftIcon={<Plus className="w-4 h-4" />}>
          Add Router
        </Button>
      </div>

      {/* Routers Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {[1, 2].map((i) => (
            <Card key={i} className="h-64 animate-pulse" />
          ))}
        </div>
      ) : routers?.length === 0 ? (
        <Card>
          <Card.Body className="text-center py-12">
            <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mx-auto mb-4">
              <Router className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-500">No routers found</p>
          </Card.Body>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {routers?.map((router) => (
            <Card
              key={router.id}
              className={selectedRouter?.id === router.id ? 'ring-2 ring-primary-500' : ''}
            >
              <Card.Body>
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div
                      className={`w-12 h-12 rounded-xl flex items-center justify-center ${
                        router.isActive
                          ? 'bg-success-100 dark:bg-success-900/30'
                          : 'bg-gray-100 dark:bg-dark-700'
                      }`}
                    >
                      <Router
                        className={`w-6 h-6 ${
                          router.isActive ? 'text-success-600' : 'text-gray-400'
                        }`}
                      />
                    </div>
                    <div>
                      <h3 className="font-semibold text-lg text-gray-900 dark:text-white">
                        {router.name}
                      </h3>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge variant={router.isActive ? 'success' : 'default'}>
                          {router.isActive ? 'Online' : 'Offline'}
                        </Badge>
                        {router.useSsl && (
                          <Badge variant="primary">
                            <Shield className="w-3 h-3 mr-1" />
                            SSL
                          </Badge>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => openModal(router)}
                    >
                      <Edit2 className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        if (confirm('Are you sure?')) {
                          deleteMutation.mutate(router.id)
                        }
                      }}
                    >
                      <Trash2 className="w-4 h-4 text-danger-500" />
                    </Button>
                  </div>
                </div>

                <div className="space-y-2 text-sm">
                  <div className="flex items-center gap-2 text-gray-600 dark:text-gray-300">
                    <Globe className="w-4 h-4 text-gray-400" />
                    <span className="font-mono">
                      {router.host}:{router.port}
                    </span>
                  </div>
                  {router.description && (
                    <p className="text-gray-500">{router.description}</p>
                  )}
                  {router.lastConnected && (
                    <p className="text-xs text-gray-400">
                      Last connected: {new Date(router.lastConnected).toLocaleString()}
                    </p>
                  )}
                </div>

                <div className="mt-4 flex gap-2">
                  <Button
                    variant={selectedRouter?.id === router.id ? 'primary' : 'ghost'}
                    size="sm"
                    className="flex-1"
                    onClick={() => setSelectedRouter(router)}
                  >
                    {selectedRouter?.id === router.id ? 'Selected' : 'Select'}
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    isLoading={testingId === router.id}
                    onClick={() => testMutation.mutate(router.id)}
                  >
                    <RefreshCw className="w-4 h-4" />
                  </Button>
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
        title={editingRouter ? 'Edit Router' : 'Add Router'}
        footer={
          <div className="flex justify-end gap-3">
            <Button variant="ghost" onClick={() => setIsModalOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmit(onSubmit)}
              isLoading={createMutation.isPending || updateMutation.isPending}
            >
              {editingRouter ? 'Update' : 'Create'}
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
            label="Host"
            placeholder="192.168.1.1 or router.local"
            {...register('host')}
            error={errors.host?.message}
          />
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Port"
              type="number"
              {...register('port', { valueAsNumber: true })}
            />
            <div className="pt-6">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  {...register('useSsl')}
                  className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm text-gray-700 dark:text-gray-300">Use SSL</span>
              </label>
            </div>
          </div>
          <Input
            label="Username"
            {...register('username')}
            error={errors.username?.message}
          />
          <Input
            label="Password"
            type="password"
            {...register('password')}
            error={errors.password?.message}
          />
          <Input
            label="Description"
            {...register('description')}
          />
        </form>
      </Modal>
    </motion.div>
  )
}

