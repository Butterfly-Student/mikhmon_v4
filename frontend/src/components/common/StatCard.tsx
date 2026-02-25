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
}

const gradients = {
  cyan: 'from-cyan-500 to-blue-600',
  indigo: 'from-indigo-500 to-purple-600',
  emerald: 'from-emerald-500 to-teal-600',
  amber: 'from-amber-500 to-orange-600',
  pink: 'from-pink-500 to-rose-600',
}

export function StatCard({ title, value, icon: Icon, subtitle, gradient, onClick, className }: StatCardProps) {
  return (
    <div
      onClick={onClick}
      className={cn(
        'relative overflow-hidden rounded-2xl p-6 text-white transition-all duration-300',
        'bg-gradient-to-br shadow-lg hover:shadow-xl hover:-translate-y-1',
        gradients[gradient],
        onClick && 'cursor-pointer',
        className
      )}
    >
      {/* Background decoration */}
      <div className="absolute -right-6 -top-6 h-24 w-24 rounded-full bg-white/10" />
      <div className="absolute -bottom-8 -left-8 h-32 w-32 rounded-full bg-white/5" />
      
      <div className="relative flex items-start justify-between">
        <div>
          <p className="text-sm font-medium text-white/80">{title}</p>
          <p className="mt-2 text-3xl font-bold tracking-tight">{value}</p>
          {subtitle && <p className="mt-1 text-sm text-white/70">{subtitle}</p>}
        </div>
        <div className="rounded-xl bg-white/20 p-3">
          <Icon className="h-6 w-6 text-white" />
        </div>
      </div>
    </div>
  )
}
