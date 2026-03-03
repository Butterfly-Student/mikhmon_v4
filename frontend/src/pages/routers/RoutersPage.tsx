import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import { Router, Plus } from 'lucide-react'
import toast from 'react-hot-toast'

import { Card, Button, Modal } from '../../components/ui'
import { routersApi } from '../../api/routers'
import { useRouterStore } from '../../stores/routerStore'
import { RouterCard } from './components/RouterCard'
import { RouterForm } from './components/RouterForm'
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

type RouterFormType = z.infer<typeof routerSchema>

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
    mutationFn: (data: RouterFormType) => routersApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routers'] })
      toast.success('Router added successfully')
      setIsModalOpen(false)
    },
    onError: (error: any) => toast.error(error.message),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string | number; data: RouterFormType }) =>
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
    defaultValues: { port: 8728, useSsl: false },
  })

  const onSubmit = (data: any) => {
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
        password: '',
        useSsl: router.useSsl,
        description: router.description,
      })
    } else {
      setEditingRouter(null)
      reset({ port: 8728, useSsl: false })
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
          <p className="text-gray-500 dark:text-gray-400">Manage your MikroTik routers</p>
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
            <RouterCard
              key={router.id}
              router={router}
              isSelected={selectedRouter?.id === router.id}
              testingId={testingId}
              onEdit={openModal}
              onDelete={(id) => deleteMutation.mutate(id)}
              onSelect={setSelectedRouter}
              onTest={(id) => testMutation.mutate(id)}
            />
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
        <RouterForm register={register} errors={errors} />
      </Modal>
    </motion.div>
  )
}
