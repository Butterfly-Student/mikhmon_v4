import { api } from './axios'
import type {
  HotspotUser,
  UserProfile,
  HotspotActive,
  HotspotHost,
  UserFilter,
  ApiResponse
} from '../types'

const DEBUG = import.meta.env.DEV || localStorage.getItem('debug_hotspot') === 'true'

export const hotspotApi = {
  // Users
  getUsers: async (routerId: string, filter?: UserFilter): Promise<HotspotUser[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching users for router:', routerId, 'filter:', filter)
    const params = new URLSearchParams()
    if (filter?.profile) params.append('profile', filter.profile)
    if (filter?.comment) params.append('comment', filter.comment)

    const { data } = await api.get<ApiResponse<HotspotUser[]>>(`/hotspot/${routerId}/users?${params}`)
    if (DEBUG) console.log('[Hotspot API] Users response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get users'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getUser: async (routerId: string, userId: string): Promise<HotspotUser> => {
    if (DEBUG) console.log('[Hotspot API] Fetching user:', userId, 'for router:', routerId)
    const { data } = await api.get<ApiResponse<HotspotUser>>(`/hotspot/${routerId}/users/${userId}`)
    if (DEBUG) console.log('[Hotspot API] User response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get user'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  createUser: async (routerId: string, user: Partial<HotspotUser>): Promise<HotspotUser> => {
    if (DEBUG) console.log('[Hotspot API] Creating user for router:', routerId, 'user:', user)
    const { data } = await api.post<ApiResponse<HotspotUser>>(`/hotspot/${routerId}/users`, user)
    if (DEBUG) console.log('[Hotspot API] Create user response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to create user'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  updateUser: async (routerId: string, userId: string, user: Partial<HotspotUser>): Promise<HotspotUser> => {
    if (DEBUG) console.log('[Hotspot API] Updating user:', userId, 'for router:', routerId, 'user:', user)
    const { data } = await api.put<ApiResponse<HotspotUser>>(`/hotspot/${routerId}/users/${userId}`, user)
    if (DEBUG) console.log('[Hotspot API] Update user response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to update user'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  deleteUser: async (routerId: string, userId: string): Promise<void> => {
    if (DEBUG) console.log('[Hotspot API] Deleting user:', userId, 'for router:', routerId)
    const { data } = await api.delete<ApiResponse<void>>(`/hotspot/${routerId}/users/${userId}`)
    if (DEBUG) console.log('[Hotspot API] Delete user response:', data)

    if (!data.success) {
      const error = data.error || 'Failed to delete user'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
  },

  // Profiles
  getProfiles: async (routerId: string): Promise<UserProfile[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching profiles for router:', routerId)
    const { data } = await api.get<ApiResponse<UserProfile[]>>(`/hotspot/${routerId}/profiles`)
    if (DEBUG) console.log('[Hotspot API] Profiles response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get profiles'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  createProfile: async (routerId: string, profile: Partial<UserProfile>): Promise<UserProfile> => {
    if (DEBUG) console.log('[Hotspot API] Creating profile for router:', routerId, 'profile:', profile)
    const { data } = await api.post<ApiResponse<UserProfile>>(`/hotspot/${routerId}/profiles`, profile)
    if (DEBUG) console.log('[Hotspot API] Create profile response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to create profile'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  updateProfile: async (routerId: string, profileId: string, profile: Partial<UserProfile>): Promise<UserProfile> => {
    if (DEBUG) console.log('[Hotspot API] Updating profile:', profileId, 'for router:', routerId, 'profile:', profile)
    const { data } = await api.put<ApiResponse<UserProfile>>(`/hotspot/${routerId}/profiles/${profileId}`, profile)
    if (DEBUG) console.log('[Hotspot API] Update profile response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to update profile'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  deleteProfile: async (routerId: string, profileId: string): Promise<void> => {
    if (DEBUG) console.log('[Hotspot API] Deleting profile:', profileId, 'for router:', routerId)
    const { data } = await api.delete<ApiResponse<void>>(`/hotspot/${routerId}/profiles/${profileId}`)
    if (DEBUG) console.log('[Hotspot API] Delete profile response:', data)

    if (!data.success) {
      const error = data.error || 'Failed to delete profile'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
  },

  // Active Sessions
  getActive: async (routerId: string): Promise<HotspotActive[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching active sessions for router:', routerId)
    const { data } = await api.get<ApiResponse<HotspotActive[]>>(`/hotspot/${routerId}/active`)
    if (DEBUG) console.log('[Hotspot API] Active sessions response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get active sessions'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  // Hosts
  getHosts: async (routerId: string): Promise<HotspotHost[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching hosts for router:', routerId)
    const { data } = await api.get<ApiResponse<HotspotHost[]>>(`/hotspot/${routerId}/hosts`)
    if (DEBUG) console.log('[Hotspot API] Hosts response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get hosts'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },
}
