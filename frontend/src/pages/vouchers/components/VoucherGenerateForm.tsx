import { Ticket, Users, Type, Clock, Settings } from 'lucide-react'
import { Button, Input, Select } from '../../../components/ui'
import type { UserProfile } from '../../../types'

const nameLengths = [
  { value: '4', label: '4 characters' },
  { value: '5', label: '5 characters' },
  { value: '6', label: '6 characters' },
  { value: '7', label: '7 characters' },
  { value: '8', label: '8 characters' },
]

interface VoucherGenerateFormProps {
  register: any
  errors: any
  profiles?: UserProfile[]
  servers?: string[]
  characterOptions: { value: string; label: string }[]
  isLoading: boolean
  onSubmit: (e?: React.BaseSyntheticEvent) => Promise<void>
}

export function VoucherGenerateForm({
  register,
  errors,
  profiles,
  servers,
  characterOptions,
  isLoading,
  onSubmit,
}: VoucherGenerateFormProps) {
  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <Input
          label="Quantity"
          type="number"
          min={1}
          max={500}
          {...register('quantity')}
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

      <Select
        label="Server"
        options={(servers || ['all']).map((s) => ({ value: s, label: s }))}
        {...register('server')}
      />

      <div className="grid grid-cols-2 gap-4">
        <Select label="Name Length" options={nameLengths} {...register('nameLength')} />
        <Input
          label="Prefix"
          placeholder="e.g. WIFI-"
          {...register('prefix')}
          leftIcon={<Type className="w-4 h-4" />}
        />
      </div>

      <Select label="Character Set" options={characterOptions} {...register('characterSet')} />

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

      <Input label="Comment" placeholder="e.g. Promo Jan 2024" {...register('comment')} />

      <Button
        type="submit"
        variant="primary"
        className="w-full"
        isLoading={isLoading}
        leftIcon={<Ticket className="w-5 h-5" />}
      >
        Generate Vouchers
      </Button>
    </form>
  )
}
