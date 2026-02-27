import { api } from './axios'
import type { ApiResponse } from '../types'

const DEBUG = import.meta.env.DEV || localStorage.getItem('debug_mikrotik') === 'true'

export const mikrotikApi = {
    getInterfaces: async (routerId: string): Promise<any[]> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching interfaces for router:', routerId)
        const { data } = await api.get<ApiResponse<any[]>>(`/mikrotik/${routerId}/interfaces`)
        if (DEBUG) console.log('[Mikrotik API] Interfaces response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get interfaces'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },

    getAllQueues: async (routerId: string): Promise<any[]> => {
        if (DEBUG) console.log('[Mikrotik API] Fetching all queues for router:', routerId)
        const { data } = await api.get<ApiResponse<any[]>>(`/mikrotik/${routerId}/queues`)
        if (DEBUG) console.log('[Mikrotik API] Queues response:', data)

        if (!data.success || !data.data) {
            const error = data.error || 'Failed to get queues'
            if (DEBUG) console.error('[Mikrotik API] Error:', error)
            throw new Error(error)
        }
        return data.data
    },
}
