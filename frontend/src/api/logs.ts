import { api } from './axios'
import type { LogEntry, ApiResponse } from '../types'

const DEBUG = import.meta.env.DEV

export const logsApi = {
    getSystemLogs: async (routerId: string): Promise<LogEntry[]> => {
        const { data } = await api.get<ApiResponse<LogEntry[]>>(`/mikrotik/${routerId}/logs`)
        if (DEBUG) console.log('[Logs API] System logs:', data)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to get system logs')
        return data.data
    },

    getHotspotLogs: async (routerId: string): Promise<LogEntry[]> => {
        const { data } = await api.get<ApiResponse<LogEntry[]>>(`/mikrotik/${routerId}/logs/hotspot`)
        if (DEBUG) console.log('[Logs API] Hotspot logs:', data)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to get hotspot logs')
        return data.data
    },

    getPPPLogs: async (routerId: string): Promise<LogEntry[]> => {
        const { data } = await api.get<ApiResponse<LogEntry[]>>(`/mikrotik/${routerId}/logs/ppp`)
        if (DEBUG) console.log('[Logs API] PPP logs:', data)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to get PPP logs')
        return data.data
    },
}
