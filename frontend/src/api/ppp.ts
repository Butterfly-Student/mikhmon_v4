import { api } from './axios'
import type { PPPSecret, PPPProfile, ApiResponse } from '../types'

const DEBUG = import.meta.env.DEV

export const pppApi = {
    // ─── Secrets ────────────────────────────────────────────────────────────────
    getSecrets: async (routerId: string): Promise<PPPSecret[]> => {
        const { data } = await api.get<ApiResponse<PPPSecret[]>>(`/mikrotik/${routerId}/ppp/secrets`)
        if (DEBUG) console.log('[PPP API] Secrets:', data)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to get PPP secrets')
        return data.data
    },

    createSecret: async (routerId: string, payload: Partial<PPPSecret>): Promise<PPPSecret> => {
        const { data } = await api.post<ApiResponse<PPPSecret>>(`/mikrotik/${routerId}/ppp/secrets`, payload)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to create PPP secret')
        return data.data
    },

    updateSecret: async (routerId: string, id: string, payload: Partial<PPPSecret>): Promise<PPPSecret> => {
        const { data } = await api.put<ApiResponse<PPPSecret>>(`/mikrotik/${routerId}/ppp/secrets/${id}`, payload)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to update PPP secret')
        return data.data
    },

    deleteSecret: async (routerId: string, id: string): Promise<void> => {
        const { data } = await api.delete<ApiResponse<void>>(`/mikrotik/${routerId}/ppp/secrets/${id}`)
        if (!data.success) throw new Error(data.error || 'Failed to delete PPP secret')
    },

    enableSecret: async (routerId: string, id: string): Promise<void> => {
        const { data } = await api.patch<ApiResponse<void>>(`/mikrotik/${routerId}/ppp/secrets/${id}/enable`)
        if (!data.success) throw new Error(data.error || 'Failed to enable PPP secret')
    },

    disableSecret: async (routerId: string, id: string): Promise<void> => {
        const { data } = await api.patch<ApiResponse<void>>(`/mikrotik/${routerId}/ppp/secrets/${id}/disable`)
        if (!data.success) throw new Error(data.error || 'Failed to disable PPP secret')
    },

    // ─── Profiles ───────────────────────────────────────────────────────────────
    getProfiles: async (routerId: string): Promise<PPPProfile[]> => {
        const { data } = await api.get<ApiResponse<PPPProfile[]>>(`/mikrotik/${routerId}/ppp/profiles`)
        if (DEBUG) console.log('[PPP API] Profiles:', data)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to get PPP profiles')
        return data.data
    },

    createProfile: async (routerId: string, payload: Partial<PPPProfile>): Promise<PPPProfile> => {
        const { data } = await api.post<ApiResponse<PPPProfile>>(`/mikrotik/${routerId}/ppp/profiles`, payload)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to create PPP profile')
        return data.data
    },

    updateProfile: async (routerId: string, id: string, payload: Partial<PPPProfile>): Promise<PPPProfile> => {
        const { data } = await api.put<ApiResponse<PPPProfile>>(`/mikrotik/${routerId}/ppp/profiles/${id}`, payload)
        if (!data.success || !data.data) throw new Error(data.error || 'Failed to update PPP profile')
        return data.data
    },

    deleteProfile: async (routerId: string, id: string): Promise<void> => {
        const { data } = await api.delete<ApiResponse<void>>(`/mikrotik/${routerId}/ppp/profiles/${id}`)
        if (!data.success) throw new Error(data.error || 'Failed to delete PPP profile')
    },

    // ─── Active Action (kick) ────────────────────────────────────────────────────
    // NOTE: Active/Inactive data is fetched via WebSocket, not HTTP.
    // This endpoint is only for the "kick/disconnect" action.
    disconnectActive: async (routerId: string, id: string): Promise<void> => {
        const { data } = await api.delete<ApiResponse<void>>(`/mikrotik/${routerId}/ppp/active/${id}`)
        if (!data.success) throw new Error(data.error || 'Failed to disconnect PPP session')
    },
}
