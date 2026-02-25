import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import {
  Ticket,
  Users,
  Settings,
  Clock,
  Type,
  Printer,
  Trash2,
  Copy,
  Check,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { Card, Button, Input, Select } from '../../components/ui'
import { vouchersApi } from '../../api/vouchers'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'
import type { Voucher } from '../../types'

const generateSchema = z.object({
  quantity: z.number().min(1).max(500),
  server: z.string().optional(),
  mode: z.enum(['vc', 'up']),
  nameLength: z.number().min(4).max(8),
  prefix: z.string().optional(),
  characterSet: z.string(),
  profile: z.string().min(1, 'Profile is required'),
  timeLimit: z.string().optional(),
  dataLimit: z.string().optional(),
  comment: z.string().optional(),
})

type GenerateForm = z.infer<typeof generateSchema>

const characterSets = [
  { value: 'lower', label: 'abcd (lowercase)' },
  { value: 'upper', label: 'ABCD (uppercase)' },
  { value: 'upplow', label: 'aBcD (mixed case)' },
  { value: 'mix', label: '5ab2c34d (alphanumeric)' },
  { value: 'mix1', label: '5AB2C34D (upper alphanumeric)' },
  { value: 'num', label: '1234 (numbers only)' },
]

const nameLengths = [
  { value: '4', label: '4 characters' },
  { value: '5', label: '5 characters' },
  { value: '6', label: '6 characters' },
  { value: '7', label: '7 characters' },
  { value: '8', label: '8 characters' },
]

export function GeneratePage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  const [generatedVouchers, setGeneratedVouchers] = useState<Voucher[]>([])
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null)

  const { data: profiles } = useQuery({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
  })

  const generateMutation = useMutation({
    mutationFn: (data: GenerateForm) => vouchersApi.generate(routerId, data),
    onSuccess: (data) => {
      setGeneratedVouchers(data)
      toast.success(`${data.length} vouchers generated successfully`)
      queryClient.invalidateQueries({ queryKey: ['users', routerId] })
    },
    onError: (error: any) => toast.error(error.message),
  })

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<GenerateForm>({
    resolver: zodResolver(generateSchema),
    defaultValues: {
      quantity: 10,
      mode: 'vc',
      nameLength: 6,
      characterSet: 'mix',
    },
  })

  const mode = watch('mode')

  const onSubmit = (data: GenerateForm) => {
    generateMutation.mutate(data)
  }

  const copyToClipboard = (text: string, index: number) => {
    navigator.clipboard.writeText(text)
    setCopiedIndex(index)
    toast.success('Copied to clipboard')
    setTimeout(() => setCopiedIndex(null), 2000)
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Generate Vouchers</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Create hotspot vouchers in batch
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Form */}
        <Card>
          <Card.Body>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <Input
                  label="Quantity"
                  type="number"
                  min={1}
                  max={500}
                  {...register('quantity', { valueAsNumber: true })}
                  leftIcon={<Users className="w-4 h-4" />}
                />
                <Select
                  label="User Mode"
                  options={[
                    { value: 'vc', label: 'Username = Password' },
                    { value: 'up', label: 'Username & Password' },
                  ]}
                  {...register('mode')}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Select
                  label="Name Length"
                  options={nameLengths}
                  {...register('nameLength', { valueAsNumber: true })}
                />
                <Input
                  label="Prefix"
                  placeholder="e.g. WIFI-"
                  {...register('prefix')}
                  leftIcon={<Type className="w-4 h-4" />}
                />
              </div>

              <Select
                label="Character Set"
                options={characterSets}
                {...register('characterSet')}
              />

              <Select
                label="Profile"
                options={profiles?.map((p) => ({ value: p.name, label: p.name })) || []}
                {...register('profile')}
                error={errors.profile?.message}
              />

              <div className="grid grid-cols-2 gap-4">
                <Input
                  label="Time Limit"
                  placeholder="e.g. 3h, 30m"
                  {...register('timeLimit')}
                  leftIcon={<Clock className="w-4 h-4" />}
                />
                <Input
                  label="Data Limit"
                  placeholder="e.g. 1GB"
                  {...register('dataLimit')}
                  leftIcon={<Settings className="w-4 h-4" />}
                />
              </div>

              <Input
                label="Comment"
                placeholder="e.g. Promo Jan 2024"
                {...register('comment')}
              />

              <Button
                type="submit"
                variant="gradient"
                className="w-full"
                isLoading={generateMutation.isPending}
                leftIcon={<Ticket className="w-5 h-5" />}
              >
                Generate Vouchers
              </Button>
            </form>
          </Card.Body>
        </Card>

        {/* Preview */}
        <Card>
          <Card.Header>
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-gray-900 dark:text-white">Generated Vouchers</h3>
              {generatedVouchers.length > 0 && (
                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    leftIcon={<Printer className="w-4 h-4" />}
                  >
                    Print
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    leftIcon={<Trash2 className="w-4 h-4 text-danger-500" />}
                    onClick={() => setGeneratedVouchers([])}
                  >
                    Clear
                  </Button>
                </div>
              )}
            </div>
          </Card.Header>
          <Card.Body>
            {generatedVouchers.length === 0 ? (
              <div className="text-center py-12">
                <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mx-auto mb-4">
                  <Ticket className="w-8 h-8 text-gray-400" />
                </div>
                <p className="text-gray-500">No vouchers generated yet</p>
                <p className="text-sm text-gray-400 mt-1">
                  Fill the form and click Generate
                </p>
              </div>
            ) : (
              <div className="space-y-2 max-h-96 overflow-y-auto">
                {generatedVouchers.map((voucher, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between p-3 bg-gray-50 dark:bg-dark-700 rounded-lg"
                  >
                    <div className="font-mono text-sm">
                      <span className="font-semibold text-gray-900 dark:text-white">
                        {voucher.username}
                      </span>
                      {mode === 'up' && voucher.password && (
                        <span className="text-gray-500 ml-2">/ {voucher.password}</span>
                      )}
                    </div>
                    <button
                      onClick={() =>
                        copyToClipboard(
                          mode === 'up'
                            ? `${voucher.username} / ${voucher.password}`
                            : voucher.username,
                          index
                        )
                      }
                      className="p-2 text-gray-400 hover:text-primary-500 transition-colors"
                    >
                      {copiedIndex === index ? (
                        <Check className="w-4 h-4 text-success-500" />
                      ) : (
                        <Copy className="w-4 h-4" />
                      )}
                    </button>
                  </div>
                ))}
              </div>
            )}
          </Card.Body>
        </Card>
      </div>
    </motion.div>
  )
}
