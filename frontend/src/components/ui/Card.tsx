import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  hover?: boolean
  glass?: boolean
}

export function Card({ className, hover = false, glass = false, children, ...props }: CardProps) {
  return (
    <div
      className={cn(
        'rounded-xl sm:rounded-2xl border transition-all duration-300',
        glass
          ? 'glass-strong'
          : 'bg-white dark:bg-dark-800 border-gray-100 dark:border-dark-700 shadow-card',
        hover && 'hover:shadow-card-hover hover:-translate-y-1 cursor-pointer',
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}

Card.Header = function CardHeader({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'px-3 py-3 sm:px-5 sm:py-4 border-b border-gray-100 dark:border-dark-700 flex items-center justify-between',
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}

Card.Body = function CardBody({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={cn('p-3 sm:p-5', className)} {...props}>
      {children}
    </div>
  )
}

Card.Footer = function CardFooter({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'px-3 py-3 sm:px-5 sm:py-4 border-t border-gray-100 dark:border-dark-700',
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}
