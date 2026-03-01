import { api } from './axios'
import type { SystemResources, SystemInfo, DashboardData, SystemHealth, ApiResponse } from '../types'

const DEBUG = import.meta.env.DEV || localStorage.getItem('debug_dashboard') === 'true'

export const dashboardApi = {
  // getDashboardData has been removed to fetch data explicitly from individual endpoints

  getSystemResources: async (routerId: string): Promise<SystemResources> => {
    if (DEBUG) console.log('[Dashboard API] Fetching system resources for router:', routerId)
    const { data } = await api.get<ApiResponse<SystemResources>>(`/mikrotik/${routerId}/system/resources`)
    if (DEBUG) console.log('[Dashboard API] System resources response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get system resources'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getSystemInfo: async (routerId: string): Promise<SystemInfo> => {
    if (DEBUG) console.log('[Dashboard API] Fetching system info for router:', routerId)
    const { data } = await api.get<ApiResponse<SystemInfo>>(`/mikrotik/${routerId}/status`)
    if (DEBUG) console.log('[Dashboard API] System info response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get system info'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getIdentity: async (routerId: string): Promise<{ name: string }> => {
    if (DEBUG) console.log('[Dashboard API] Fetching identity for router:', routerId)
    const { data } = await api.get<ApiResponse<{ name: string }>>(`/mikrotik/${routerId}/system/identity`)
    if (DEBUG) console.log('[Dashboard API] Identity response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get identity'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getRouterBoard: async (routerId: string): Promise<{ model: string, boardName: string, version: string }> => {
    if (DEBUG) console.log('[Dashboard API] Fetching routerboard for router:', routerId)
    const { data } = await api.get<ApiResponse<{ model: string, boardName: string, version: string }>>(`/mikrotik/${routerId}/system/routerboard`)
    if (DEBUG) console.log('[Dashboard API] Routerboard response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get routerboard'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  // GET /mikrotik/:id/system/health
  getSystemHealth: async (routerId: string): Promise<SystemHealth> => {
    if (DEBUG) console.log('[Dashboard API] Fetching system health for router:', routerId)
    const { data } = await api.get<ApiResponse<SystemHealth>>(`/mikrotik/${routerId}/system/health`)
    if (DEBUG) console.log('[Dashboard API] System health response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get system health'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  // GET /mikrotik/:id/system/clock
  getSystemClock: async (routerId: string): Promise<{ time: string; date: string; timeZoneName: string }> => {
    if (DEBUG) console.log('[Dashboard API] Fetching system clock for router:', routerId)
    const { data } = await api.get<ApiResponse<{ time: string; date: string; timeZoneName: string }>>(`/mikrotik/${routerId}/system/clock`)
    if (DEBUG) console.log('[Dashboard API] System clock response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get system clock'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  // GET /mikrotik/:id/dashboard — data lengkap dalam satu request
  getDashboard: async (routerId: string): Promise<DashboardData> => {
    if (DEBUG) console.log('[Dashboard API] Fetching dashboard for router:', routerId)
    const { data } = await api.get<ApiResponse<DashboardData>>(`/mikrotik/${routerId}/dashboard`)
    if (DEBUG) console.log('[Dashboard API] Dashboard response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get dashboard data'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },
}
