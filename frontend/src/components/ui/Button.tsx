import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { Loader2 } from 'lucide-react'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

type ButtonVariant = 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'ghost' | 'outline'
type ButtonSize = 'xs' | 'sm' | 'md' | 'lg'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
  isLoading?: boolean
  leftIcon?: React.ReactNode
  rightIcon?: React.ReactNode
}

const variantStyles: Record<ButtonVariant, string> = {
  primary: 'bg-gradient-to-r from-primary-500 to-primary-600 hover:from-primary-600 hover:to-primary-700 text-white shadow-sm hover:shadow-glow',
  secondary: 'bg-gradient-to-r from-secondary-500 to-secondary-600 hover:from-secondary-600 hover:to-secondary-700 text-white shadow-sm hover:shadow-glow-pink',
  success: 'bg-gradient-to-r from-success-500 to-success-600 hover:from-success-600 hover:to-success-700 text-white shadow-sm hover:shadow-glow-success',
  danger: 'bg-gradient-to-r from-danger-500 to-danger-600 hover:from-danger-600 hover:to-danger-700 text-white shadow-sm',
  warning: 'bg-gradient-to-r from-warning-400 to-warning-500 hover:from-warning-500 hover:to-warning-600 text-white shadow-sm',
  ghost: 'bg-transparent hover:bg-gray-100 dark:hover:bg-dark-700 text-gray-700 dark:text-gray-300',
  outline: 'border border-gray-300 dark:border-dark-600 bg-transparent hover:bg-gray-50 dark:hover:bg-dark-700 text-gray-700 dark:text-gray-200',
}

/* Mobile-first: xs & sm use tighter sizes, md/lg are "desktop" defaults */
const sizeStyles: Record<ButtonSize, string> = {
  xs: 'px-2 py-1 text-[11px] rounded-lg gap-1',
  sm: 'px-2.5 py-1.5 text-xs rounded-lg gap-1 sm:px-3 sm:gap-1.5',
  md: 'px-3 py-1.5 text-xs rounded-xl gap-1.5 sm:px-4 sm:py-2 sm:text-sm sm:gap-2',
  lg: 'px-4 py-2 text-sm rounded-xl gap-2 sm:px-5 sm:py-2.5 sm:text-base',
}

export function Button({
  variant = 'primary',
  size = 'md',
  isLoading,
  leftIcon,
  rightIcon,
  children,
  className,
  disabled,
  ...props
}: ButtonProps) {
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center font-semibold transition-all duration-200',
        'focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 dark:focus:ring-offset-dark-800',
        'disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none',
        variantStyles[variant],
        sizeStyles[size],
        className
      )}
      disabled={disabled || isLoading}
      {...props}
    >
      {isLoading ? (
        <Loader2 className="w-3.5 h-3.5 sm:w-4 sm:h-4 animate-spin shrink-0" />
      ) : (
        leftIcon && <span className="shrink-0">{leftIcon}</span>
      )}
      {children}
      {rightIcon && !isLoading && <span className="shrink-0">{rightIcon}</span>}
    </button>
  )
}
