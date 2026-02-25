import { api } from './axios'
import type { DashboardData, SystemResources, SystemInfo, ApiResponse } from '../types'

const DEBUG = import.meta.env.DEV || localStorage.getItem('debug_dashboard') === 'true'

export const dashboardApi = {
  getDashboardData: async (routerId: string): Promise<DashboardData> => {
    if (DEBUG) console.log('[Dashboard API] Fetching dashboard data for router:', routerId)
    const { data } = await api.get<ApiResponse<DashboardData>>(`/dashboard/${routerId}`)
    if (DEBUG) console.log('[Dashboard API] Response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get dashboard data'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getSystemResources: async (routerId: string): Promise<SystemResources> => {
    if (DEBUG) console.log('[Dashboard API] Fetching system resources for router:', routerId)
    const { data } = await api.get<ApiResponse<SystemResources>>(`/dashboard/${routerId}/resources`)
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
    const { data } = await api.get<ApiResponse<SystemInfo>>(`/dashboard/${routerId}/status`)
    if (DEBUG) console.log('[Dashboard API] System info response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get system info'
      if (DEBUG) console.error('[Dashboard API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },
}
