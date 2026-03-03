import { Input, Select } from '../../../components/ui'

const expireModes = [
  { value: '0', label: 'None' },
  { value: 'rem', label: 'Remove' },
  { value: 'ntf', label: 'Notice' },
  { value: 'remc', label: 'Remove & Record' },
  { value: 'ntfc', label: 'Notice & Record' },
]

const lockOptions = [
  { value: 'Disable', label: 'Disable' },
  { value: 'Enable', label: 'Enable' },
]

interface ProfileFormProps {
  register: any
  errors: any
  expireMode: string
  addressPools?: string[]
  parentQueues?: string[]
}

export function ProfileForm({ register, errors, expireMode, addressPools, parentQueues }: ProfileFormProps) {
  return (
    <div className="space-y-6">
      {/* General Settings */}
      <div>
        <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">General</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Input label="Name" {...register('name')} error={errors.name?.message} />
          <Input
            label="Shared Users"
            type="number"
            {...register('sharedUsers', { valueAsNumber: true })}
          />
          <Input label="Rate Limit" {...register('rateLimit')} placeholder="e.g. 1M/1M" />
          <Select
            label="Address Pool"
            options={[
              { value: 'none', label: 'none' },
              ...((addressPools || []).map((p) => ({ value: p, label: p }))),
            ]}
            {...register('addressPool')}
          />
          <Select
            label="Parent Queue"
            options={[
              { value: 'none', label: 'none' },
              ...((parentQueues || []).map((q) => ({ value: q, label: q }))),
            ]}
            {...register('parentQueue')}
          />
        </div>
      </div>

      {/* Expiration Settings */}
      <div>
        <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Expiration</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Select label="Expire Mode" options={expireModes} {...register('expireMode')} />
          {expireMode !== '0' && (
            <Input label="Validity" {...register('validity')} placeholder="e.g. 30d, 12h" />
          )}
        </div>
      </div>

      {/* Pricing */}
      <div>
        <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Pricing</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Input label="Price" type="number" {...register('price', { valueAsNumber: true })} />
          <Input label="Selling Price" type="number" {...register('sellingPrice', { valueAsNumber: true })} />
        </div>
      </div>

      {/* Lock Settings */}
      <div>
        <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Lock Settings</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Select label="Lock User (MAC)" options={lockOptions} {...register('lockUser')} />
          <Select label="Lock Server" options={lockOptions} {...register('lockServer')} />
        </div>
      </div>
    </div>
  )
}
