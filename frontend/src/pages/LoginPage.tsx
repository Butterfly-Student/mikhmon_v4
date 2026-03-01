import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useNavigate } from 'react-router-dom'
import { WifiIcon, Eye, EyeOff } from 'lucide-react'
import { motion } from 'framer-motion'
import toast from 'react-hot-toast'

import { authApi } from '../api/auth'
import { useAuthStore } from '../stores/authStore'
import { ThemeToggle } from '../components/common/ThemeToggle'
import { Button } from '../components/ui/Button'

const loginSchema = z.object({
  username: z.string().min(3, 'Username minimal 3 karakter'),
  password: z.string().min(4, 'Password minimal 4 karakter'),
})

type LoginForm = z.infer<typeof loginSchema>

export function LoginPage() {
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const setAuth = useAuthStore((state) => state.setAuth)

  const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>({ resolver: zodResolver(loginSchema) })

  const onSubmit = async (data: LoginForm) => {
    setIsLoading(true)
    try {
      const response = await authApi.login(data)
      setAuth(response.token, response.user)
      toast.success('Login berhasil!')
      navigate('/dashboard')
    } catch (error: any) {
      toast.error(error.message || 'Login gagal')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 via-primary-50/30 to-secondary-50/20 dark:from-dark-900 dark:via-dark-800 dark:to-dark-900 p-4 relative overflow-hidden">
      {/* Animated background blobs */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden">
        <div className="absolute -top-40 -left-40 w-96 h-96 bg-primary-400/20 dark:bg-primary-500/10 rounded-full blur-3xl animate-pulse-slow" />
        <div className="absolute -bottom-40 -right-40 w-96 h-96 bg-secondary-400/20 dark:bg-secondary-500/10 rounded-full blur-3xl animate-pulse-slow" style={{ animationDelay: '1.5s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-72 h-72 bg-cyan-400/10 rounded-full blur-3xl" />
      </div>

      {/* Theme toggle */}
      <div className="fixed top-4 right-4 z-50">
        <ThemeToggle />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 24, scale: 0.97 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.5, ease: 'easeOut' }}
        className="relative w-full max-w-md"
      >
        {/* Card */}
        <div className="bg-white/80 dark:bg-dark-800/80 backdrop-blur-xl rounded-3xl p-8 shadow-2xl border border-white/50 dark:border-dark-700/50">
          {/* Logo */}
          <div className="flex flex-col items-center mb-8">
            <div className="relative mb-4">
              <div className="absolute inset-0 rounded-2xl bg-gradient-to-br from-primary-500 to-secondary-500 blur-lg opacity-40" />
              <div className="relative w-20 h-20 rounded-2xl bg-gradient-to-br from-primary-500 via-violet-500 to-secondary-500 flex items-center justify-center shadow-lg">
                <WifiIcon className="w-10 h-10 text-white drop-shadow-lg" />
              </div>
            </div>
            <h1 className="text-2xl font-extrabold text-gray-900 dark:text-white tracking-tight">
              Mikhmon <span className="text-primary-500">v4</span>
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Hotspot Management System</p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {/* Username */}
            <div>
              <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Username</label>
              <input
                type="text"
                {...register('username')}
                placeholder="Masukkan username"
                autoComplete="username"
                className="w-full px-4 py-3 rounded-xl border bg-white dark:bg-dark-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all border-gray-200 dark:border-dark-600"
              />
              {errors.username && <p className="mt-1 text-xs text-danger-500">{errors.username.message}</p>}
            </div>

            {/* Password */}
            <div>
              <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1.5">Password</label>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  {...register('password')}
                  placeholder="Masukkan password"
                  autoComplete="current-password"
                  className="w-full px-4 py-3 pr-11 rounded-xl border bg-white dark:bg-dark-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all border-gray-200 dark:border-dark-600"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                >
                  {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                </button>
              </div>
              {errors.password && <p className="mt-1 text-xs text-danger-500">{errors.password.message}</p>}
            </div>

            {/* Submit */}
            <Button type="submit" variant="primary" size="lg" className="w-full mt-2" isLoading={isLoading}>
              Login
            </Button>
          </form>

          <p className="mt-6 text-center text-xs text-gray-400 dark:text-gray-600">
            © {new Date().getFullYear()} Mikhmon v4 · Hotspot Manager
          </p>
        </div>
      </motion.div>
    </div>
  )
}

