import { Input } from '../../../components/ui'

interface RouterFormProps {
  register: any
  errors: any
}

export function RouterForm({ register, errors }: RouterFormProps) {
  return (
    <div className="space-y-4">
      <Input label="Name" {...register('name')} error={errors.name?.message} />
      <Input
        label="Host"
        placeholder="192.168.1.1 or router.local"
        {...register('host')}
        error={errors.host?.message}
      />
      <div className="grid grid-cols-2 gap-4">
        <Input label="Port" type="number" {...register('port', { valueAsNumber: true })} />
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
      <Input label="Username" {...register('username')} error={errors.username?.message} />
      <Input label="Password" type="password" {...register('password')} error={errors.password?.message} />
      <Input label="Description" {...register('description')} />
    </div>
  )
}
