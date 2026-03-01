import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info' | 'default'
  size?: 'sm' | 'md'
}

const variants = {
  default: 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300 ring-1 ring-gray-200 dark:ring-dark-600',
  primary: 'bg-primary-100 text-primary-700 dark:bg-primary-900/40 dark:text-primary-300 ring-1 ring-primary-200 dark:ring-primary-800',
  secondary: 'bg-secondary-100 text-secondary-700 dark:bg-secondary-900/40 dark:text-secondary-300 ring-1 ring-secondary-200 dark:ring-secondary-800',
  success: 'bg-success-100 text-success-700 dark:bg-success-900/40 dark:text-success-300 ring-1 ring-success-200 dark:ring-success-800',
  danger: 'bg-danger-100 text-danger-700 dark:bg-danger-900/40 dark:text-danger-300 ring-1 ring-danger-200 dark:ring-danger-800',
  warning: 'bg-warning-100 text-warning-700 dark:bg-warning-900/40 dark:text-warning-300 ring-1 ring-warning-200 dark:ring-warning-800',
  info: 'bg-info-100 text-info-700 dark:bg-info-900/40 dark:text-info-300 ring-1 ring-info-200 dark:ring-info-800',
}

const sizes = {
  sm: 'px-1.5 py-0.5 text-xs rounded-md',
  md: 'px-2.5 py-0.5 text-xs rounded-full',
}

export function Badge({ className, variant = 'default', size = 'md', children, ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center font-semibold',
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {children}
    </span>
  )
}
