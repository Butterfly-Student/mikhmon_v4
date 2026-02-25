import { motion } from 'framer-motion'
import {
  Palette,
  Bell,
  Shield,
} from 'lucide-react'

import { Card, Badge } from '../../components/ui'
import { ThemeToggle } from '../../components/common/ThemeToggle'
import { useThemeStore } from '../../stores/themeStore'

const colorSchemes = [
  { value: 'indigo', label: 'Indigo', color: 'bg-indigo-500' },
  { value: 'pink', label: 'Pink', color: 'bg-pink-500' },
  { value: 'emerald', label: 'Emerald', color: 'bg-emerald-500' },
  { value: 'amber', label: 'Amber', color: 'bg-amber-500' },
  { value: 'cyan', label: 'Cyan', color: 'bg-cyan-500' },
  { value: 'purple', label: 'Purple', color: 'bg-purple-500' },
]

export function SettingsPage() {
  const { colorScheme, setColorScheme } = useThemeStore()

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Settings</h1>
        <p className="text-gray-500 dark:text-gray-400">
          Customize your Mikhmon experience
        </p>
      </div>

      {/* Appearance */}
      <Card>
        <Card.Header>
          <div className="flex items-center gap-2">
            <Palette className="w-5 h-5 text-primary-500" />
            <h3 className="font-semibold text-gray-900 dark:text-white">Appearance</h3>
          </div>
        </Card.Header>
        <Card.Body className="space-y-6">
          {/* Theme Toggle */}
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900 dark:text-white">Dark Mode</p>
              <p className="text-sm text-gray-500">Toggle between light and dark theme</p>
            </div>
            <ThemeToggle />
          </div>

          {/* Color Scheme */}
          <div>
            <p className="font-medium text-gray-900 dark:text-white mb-3">Accent Color</p>
            <div className="grid grid-cols-3 sm:grid-cols-6 gap-3">
              {colorSchemes.map((scheme) => (
                <button
                  key={scheme.value}
                  onClick={() => setColorScheme(scheme.value as any)}
                  className={`flex flex-col items-center gap-2 p-3 rounded-xl transition-all ${
                    colorScheme === scheme.value
                      ? 'bg-gray-100 dark:bg-dark-700 ring-2 ring-primary-500'
                      : 'hover:bg-gray-50 dark:hover:bg-dark-700'
                  }`}
                >
                  <div className={`w-8 h-8 rounded-full ${scheme.color}`} />
                  <span className="text-xs font-medium">{scheme.label}</span>
                </button>
              ))}
            </div>
          </div>
        </Card.Body>
      </Card>

      {/* Notifications */}
      <Card>
        <Card.Header>
          <div className="flex items-center gap-2">
            <Bell className="w-5 h-5 text-warning-500" />
            <h3 className="font-semibold text-gray-900 dark:text-white">Notifications</h3>
          </div>
        </Card.Header>
        <Card.Body className="space-y-4">
          {[
            { label: 'User Login', description: 'Get notified when a user logs in' },
            { label: 'Voucher Expired', description: 'Get notified when vouchers expire' },
            { label: 'Low Balance', description: 'Get notified when router balance is low' },
          ].map((item) => (
            <div key={item.label} className="flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900 dark:text-white">{item.label}</p>
                <p className="text-sm text-gray-500">{item.description}</p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 dark:peer-focus:ring-primary-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary-600" />
              </label>
            </div>
          ))}
        </Card.Body>
      </Card>

      {/* About */}
      <Card>
        <Card.Header>
          <div className="flex items-center gap-2">
            <Shield className="w-5 h-5 text-success-500" />
            <h3 className="font-semibold text-gray-900 dark:text-white">About</h3>
          </div>
        </Card.Header>
        <Card.Body className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900 dark:text-white">Version</p>
              <p className="text-sm text-gray-500">Mikhmon v4.0.0</p>
            </div>
            <Badge variant="success">Latest</Badge>
          </div>
          <div className="pt-4 border-t border-gray-100 dark:border-dark-700">
            <p className="text-sm text-gray-500 text-center">
              © 2024 Mikhmon. Hotspot Management System.
            </p>
          </div>
        </Card.Body>
      </Card>
    </motion.div>
  )
}
