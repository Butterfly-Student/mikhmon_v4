import type { LucideIcon } from 'lucide-react'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

interface StatCardProps {
  title: string
  value: string | number
  icon: LucideIcon
  subtitle?: string
  gradient: 'cyan' | 'indigo' | 'emerald' | 'amber' | 'pink'
  onClick?: () => void
  className?: string
  trend?: { value: number; label: string }
}

const gradients = {
  cyan: 'from-cyan-500 to-blue-600',
  indigo: 'from-indigo-500 to-violet-600',
  emerald: 'from-emerald-500 to-teal-600',
  amber: 'from-amber-400 to-orange-500',
  pink: 'from-pink-500 to-rose-600',
}

export function StatCard({ title, value, icon: Icon, subtitle, gradient, onClick, className, trend }: StatCardProps) {
  return (
    <div
      onClick={onClick}
      className={cn(
        'relative overflow-hidden rounded-xl sm:rounded-2xl p-3 sm:p-5 text-white',
        'bg-gradient-to-br shadow-lg transition-all duration-300',
        'hover:shadow-xl hover:-translate-y-0.5 sm:hover:-translate-y-1',
        gradients[gradient],
        onClick && 'cursor-pointer',
        className
      )}
    >
      {/* Decorative circles */}
      <div className="absolute -right-4 -top-4 h-14 w-14 sm:h-20 sm:w-20 rounded-full bg-white/10" />
      <div className="absolute -bottom-6 -left-4 h-20 w-20 sm:h-28 sm:w-28 rounded-full bg-white/5" />
      <div className="absolute top-1/2 right-6 h-8 w-8 sm:h-10 sm:w-10 rounded-full bg-white/5" />

      <div className="relative flex items-start justify-between gap-2">
        <div className="min-w-0 flex-1">
          <p className="text-[10px] sm:text-xs font-semibold text-white/75 uppercase tracking-wider mb-0.5 sm:mb-1 truncate">{title}</p>
          <p className="text-2xl sm:text-3xl font-extrabold tracking-tight leading-none truncate">{value}</p>
          {subtitle && (
            <p className="mt-1 text-[10px] sm:text-xs text-white/65 truncate">{subtitle}</p>
          )}
          {trend && (
            <p className="mt-1 sm:mt-2 text-[10px] sm:text-xs text-white/80 font-medium">
              <span className={trend.value >= 0 ? 'text-white' : 'text-white/60'}>
                {trend.value >= 0 ? '↑' : '↓'} {Math.abs(trend.value)}%
              </span>
              {' '}<span className="text-white/60">{trend.label}</span>
            </p>
          )}
        </div>
        <div className="rounded-lg sm:rounded-xl bg-white/20 backdrop-blur-sm p-2 sm:p-2.5 shrink-0">
          <Icon className="h-4 w-4 sm:h-6 sm:w-6 text-white drop-shadow" />
        </div>
      </div>
    </div>
  )
}
