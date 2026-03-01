import { api } from './axios'
import type { NetworkInterface, LogEntry, ApiResponse } from '../types'

const DEBUG = import.meta.env.DEV || localStorage.getItem('debug_mikrotik') === 'true'

export const mikrotikApi = {
    // GET /mikrotik/:id/interfaces
    getInterfaces: async (routerId: string): Promise<NetworkInterface[]> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching interfaces for router:', routerId)
        const { data } = await api.get<ApiResponse<NetworkInterface[]>>(`/mikrotik/${routerId}/interfaces`)
        if (DEBUG) console.log('[Mikrotik API] Interfaces response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get interfaces'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },

    // GET /mikrotik/:id/queues
    getAllQueues: async (routerId: string): Promise<string[]> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching all queues for router:', routerId)
        const { data } = await api.get<ApiResponse<string[]>>(`/mikrotik/${routerId}/queues/parent`)
        if (DEBUG) console.log('[Mikrotik API] Queues response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get queues'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },

    // GET /mikrotik/:id/nat
    getNATRules: async (routerId: string): Promise<any[]> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching NAT rules for router:', routerId)
        const { data } = await api.get<ApiResponse<any[]>>(`/mikrotik/${routerId}/nat`)
        if (DEBUG) console.log('[Mikrotik API] NAT rules response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get NAT rules'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },

    // GET /mikrotik/:id/logs?limit=N
    getLogs: async (routerId: string, limit?: number): Promise<LogEntry[]> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching logs for router:', routerId)
        const qs = limit ? `?limit=${limit}` : ''
        const { data } = await api.get<ApiResponse<LogEntry[]>>(`/mikrotik/${routerId}/logs${qs}`)
        if (DEBUG) console.log('[Mikrotik API] Logs response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get logs'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },

    // GET /mikrotik/:id/interfaces/:name/traffic
    getInterfaceTraffic: async (routerId: string, interfaceName: string): Promise<any> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching interface traffic for router:', routerId, 'interface:', interfaceName)
        const { data } = await api.get<ApiResponse<any>>(`/mikrotik/${routerId}/interfaces/${encodeURIComponent(interfaceName)}/traffic`)
        if (DEBUG) console.log('[Mikrotik API] Interface traffic response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get interface traffic'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },
}
