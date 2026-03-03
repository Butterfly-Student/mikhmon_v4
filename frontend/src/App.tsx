import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'

import { useAuthStore } from './stores/authStore'
import { useThemeStore } from './stores/themeStore'
import { LoginPage } from './pages/LoginPage'
import { DashboardPage } from './pages/DashboardPage'
import { AppLayout } from './components/layout/AppLayout'

// Hotspot Pages
import { HotspotPage, UsersPage, ProfilesPage, ActivePage, HostsPage, InactivePage } from './pages/hotspot'

// Voucher Pages
import { GeneratePage, PrintPage } from './pages/vouchers'

// Reports Pages
import { SalesPage } from './pages/reports'

// Routers Page
import { RoutersPage } from './pages/routers/RoutersPage'

// Settings Page
import { SettingsPage } from './pages/settings/SettingsPage'

// PPPoE Pages
import { PPPoEPage } from './pages/pppoe'

// Logs Pages
import { LogsPage } from './pages/logs'

// Protected Route wrapper
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

// Public Route wrapper (redirect to dashboard if authenticated)
function PublicRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />
  }

  return <>{children}</>
}

function App() {
  const { mode } = useThemeStore()

  // Initialize theme on mount
  useEffect(() => {
    const isDark = mode === 'system'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
      : mode === 'dark'

    if (isDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [mode])

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            <PublicRoute>
              <LoginPage />
            </PublicRoute>
          }
        />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <AppLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<DashboardPage />} />

          {/* Hotspot Routes */}
          <Route path="hotspot" element={<HotspotPage />} />
          <Route path="hotspot/users" element={<UsersPage />} />
          <Route path="hotspot/profiles" element={<ProfilesPage />} />
          <Route path="hotspot/active" element={<ActivePage />} />
          <Route path="hotspot/inactive" element={<InactivePage />} />
          <Route path="hotspot/hosts" element={<HostsPage />} />

          {/* PPPoE Routes */}
          <Route path="pppoe" element={<PPPoEPage />} />

          {/* Logs Routes */}
          <Route path="logs" element={<LogsPage />} />

          {/* Voucher Routes */}
          <Route path="vouchers" element={<Navigate to="/vouchers/generate" replace />} />
          <Route path="vouchers/generate" element={<GeneratePage />} />
          <Route path="vouchers/print" element={<PrintPage />} />

          {/* Reports Routes */}
          <Route path="reports" element={<Navigate to="/reports/sales" replace />} />
          <Route path="reports/sales" element={<SalesPage />} />

          {/* Routers */}
          <Route path="routers" element={<RoutersPage />} />

          {/* Settings */}
          <Route path="settings" element={<SettingsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
