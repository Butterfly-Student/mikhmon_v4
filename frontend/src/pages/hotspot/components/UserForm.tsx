import { Input, Select } from '../../../components/ui'
import type { UserProfile } from '../../../types'

interface UserFormProps {
  register: any
  errors: any
  profiles?: UserProfile[]
}

export function UserForm({ register, errors, profiles }: UserFormProps) {
  return (
    <div className="space-y-4">
      <Input label="Name" {...register('name')} error={errors.name?.message} />
      <Input label="Password" type="password" {...register('password')} />
      <Select
        label="Profile"
        options={profiles?.map(p => ({ value: p.name, label: p.name })) || []}
        {...register('profile')}
      />
      <Input label="MAC Address" placeholder="AA:BB:CC:DD:EE:FF" {...register('macAddress')} />
      <div className="grid grid-cols-2 gap-4">
        <Input label="Time Limit" placeholder="1h, 30m, 1d" {...register('timeLimit')} />
        <Input label="Data Limit" placeholder="100M, 1G" {...register('dataLimit')} />
      </div>
      <Input label="Comment" {...register('comment')} />
    </div>
  )
}
