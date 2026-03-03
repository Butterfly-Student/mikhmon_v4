import { useEffect, useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { motion } from 'framer-motion'
import toast from 'react-hot-toast'

import { Card } from '../../components/ui'
import { vouchersApi } from '../../api/vouchers'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'
import { VoucherGenerateForm } from './components/VoucherGenerateForm'
import { VoucherPreview } from './components/VoucherPreview'
import type { Voucher } from '../../types'

const generateSchema = z.object({
  quantity: z.coerce.number().min(1).max(500),
  server: z.string().optional(),
  mode: z.enum(['vc', 'up']),
  gencode: z.string().optional(),
  nameLength: z.coerce.number().min(4).max(8),
  prefix: z.string().optional(),
  characterSet: z.string(),
  profile: z.string().min(1, 'Profile is required'),
  timeLimit: z.string().optional(),
  dataLimit: z.string().optional(),
  comment: z.string().optional(),
})

type GenerateForm = z.infer<typeof generateSchema>

const upCharacterSets = [
  { value: 'lower', label: 'abcd (lowercase)' },
  { value: 'upper', label: 'ABCD (uppercase)' },
  { value: 'upplow', label: 'aBcD (mixed case)' },
  { value: 'mix', label: '5ab2c34d (alphanumeric)' },
  { value: 'mix1', label: '5AB2C34D (upper alphanumeric)' },
  { value: 'mix2', label: '5aB2c34D (mixed alphanumeric)' },
]

const vcCharacterSets = [
  { value: 'lower1', label: 'abcd2345 (lower+num)' },
  { value: 'upper1', label: 'ABCD2345 (upper+num)' },
  { value: 'upplow1', label: 'aBcD2345 (mix+num)' },
  { value: 'mix', label: '5ab2c34d (alphanumeric)' },
  { value: 'mix1', label: '5AB2C34D (upper alphanumeric)' },
  { value: 'mix2', label: '5aB2c34D (mixed alphanumeric)' },
  { value: 'num', label: '1234 (numbers only)' },
]

export function GeneratePage() {
  const queryClient = useQueryClient()
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '1')

  const [generatedVouchers, setGeneratedVouchers] = useState<Voucher[]>([])
  const [generatedComment, setGeneratedComment] = useState('')

  const { data: profiles } = useQuery({
    queryKey: ['profiles', routerId],
    queryFn: () => hotspotApi.getProfiles(routerId),
  })

  const { data: servers } = useQuery({
    queryKey: ['hotspot-servers', routerId],
    queryFn: () => hotspotApi.getServers(routerId),
  })

  const generateMutation = useMutation({
    mutationFn: (data: GenerateForm) => vouchersApi.generate(routerId, data),
    onSuccess: (result) => {
      setGeneratedVouchers(result.vouchers)
      setGeneratedComment(result.comment)
      toast.success(`${result.count} vouchers generated successfully`)
      queryClient.invalidateQueries({ queryKey: ['users', routerId] })
    },
    onError: (error: any) => toast.error(error.message),
  })

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(generateSchema),
    defaultValues: {
      quantity: 10,
      server: 'all',
      mode: 'vc',
      nameLength: 8,
      characterSet: 'lower1',
      timeLimit: '3h',
    },
  })

  const mode = watch('mode')
  const characterSet = watch('characterSet')

  const characterOptions = useMemo(
    () => (mode === 'up' ? upCharacterSets : vcCharacterSets),
    [mode]
  )

  useEffect(() => {
    const valid = new Set(characterOptions.map((c) => c.value))
    if (!valid.has(characterSet)) {
      setValue('characterSet', mode === 'up' ? 'lower' : 'lower1')
    }
    setValue('nameLength', mode === 'up' ? 4 : 8)
  }, [mode, characterSet, characterOptions, setValue])

  const onSubmit = (data: any) => {
    const parsed = generateSchema.parse(data)
    const gencode = `${Math.floor(Math.random() * 899) + 101}`
    generateMutation.mutate({ ...parsed, gencode })
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
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Generate Vouchers</h1>
          <p className="text-gray-500 dark:text-gray-400">Create hotspot vouchers in batch</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <Card.Body>
            <VoucherGenerateForm
              register={register}
              errors={errors}
              profiles={profiles}
              servers={servers}
              characterOptions={characterOptions}
              isLoading={generateMutation.isPending}
              onSubmit={handleSubmit(onSubmit)}
            />
          </Card.Body>
        </Card>

        <VoucherPreview
          vouchers={generatedVouchers}
          mode={mode}
          comment={generatedComment}
          onClear={() => setGeneratedVouchers([])}
        />
      </div>
    </motion.div>
  )
}
