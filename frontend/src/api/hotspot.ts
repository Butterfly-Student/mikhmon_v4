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

  getUsersCount: async (routerId: string): Promise<number> => {
    if (DEBUG) console.log('[Hotspot API] Fetching users count for router:', routerId)
    const { data } = await api.get<ApiResponse<{ totalUsers: number }>>(`/hotspot/${routerId}/users/count`)
    if (DEBUG) console.log('[Hotspot API] Users count response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get users count'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data.totalUsers
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

  createUser: async (routerId: string, user: Partial<HotspotUser> & {
    timeLimit?: string;
    dataLimit?: string;
  }): Promise<any> => {
    if (DEBUG) console.log('[Hotspot API] Creating user for router:', routerId, 'user:', user)
    const payload = {
      server: user.server,
      name: user.name,
      password: user.password,
      profile: typeof user.profile === 'string' ? user.profile : user.profile?.name,
      macAddress: user.macAddress,
      timeLimit: user.timeLimit || user.limitUptime,
      dataLimit: user.dataLimit,
      comment: user.comment,
      disabled: user.disabled,
    }
    const { data } = await api.post<ApiResponse<any>>(`/hotspot/${routerId}/users`, payload)
    if (DEBUG) console.log('[Hotspot API] Create user response:', data)

    if (!data.success) {
      const error = data.error || 'Failed to create user'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  updateUser: async (routerId: string, userId: string, user: Partial<HotspotUser> & {
    timeLimit?: string;
    dataLimit?: string;
    reset?: boolean;
  }): Promise<any> => {
    if (DEBUG) console.log('[Hotspot API] Updating user:', userId, 'for router:', routerId, 'user:', user)
    const payload = {
      server: user.server,
      name: user.name,
      password: user.password,
      profile: typeof user.profile === 'string' ? user.profile : user.profile?.name,
      macAddress: user.macAddress,
      timeLimit: user.timeLimit || user.limitUptime,
      dataLimit: user.dataLimit,
      comment: user.comment,
      disabled: user.disabled ?? false,
      reset: user.reset ?? false,
    }
    const { data } = await api.put<ApiResponse<any>>(`/hotspot/${routerId}/users/${userId}`, payload)
    if (DEBUG) console.log('[Hotspot API] Update user response:', data)

    if (!data.success) {
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

  deleteActiveSession: async (routerId: string, activeId: string): Promise<void> => {
    if (DEBUG) console.log('[Hotspot API] Deleting active session:', activeId, 'for router:', routerId)
    const { data } = await api.delete<ApiResponse<void>>(`/hotspot/${routerId}/active/${activeId}`)
    if (DEBUG) console.log('[Hotspot API] Delete active session response:', data)

    if (!data.success) {
      const error = data.error || 'Failed to remove active session'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
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

  deleteHost: async (routerId: string, hostId: string): Promise<void> => {
    if (DEBUG) console.log('[Hotspot API] Deleting host:', hostId, 'for router:', routerId)
    const { data } = await api.delete<ApiResponse<void>>(`/hotspot/${routerId}/hosts/${hostId}`)
    if (DEBUG) console.log('[Hotspot API] Delete host response:', data)

    if (!data.success) {
      const error = data.error || 'Failed to remove host'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
  },

  // Metadata
  getServers: async (routerId: string): Promise<string[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching servers for router:', routerId)
    const { data } = await api.get<ApiResponse<string[]>>(`/hotspot/${routerId}/servers`)
    if (DEBUG) console.log('[Hotspot API] Servers response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get hotspot servers'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getAddressPools: async (routerId: string): Promise<string[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching address pools for router:', routerId)
    const { data } = await api.get<ApiResponse<string[]>>(`/hotspot/${routerId}/address-pools`)
    if (DEBUG) console.log('[Hotspot API] Address pools response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get address pools'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getParentQueues: async (routerId: string): Promise<string[]> => {
    if (DEBUG) console.log('[Hotspot API] Fetching parent queues for router:', routerId)
    const { data } = await api.get<ApiResponse<string[]>>(`/hotspot/${routerId}/parent-queues`)
    if (DEBUG) console.log('[Hotspot API] Parent queues response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to get parent queues'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  // Expire Monitor
  setupExpireMonitor: async (routerId: string, script?: string): Promise<{ message: string; status: string }> => {
    if (DEBUG) console.log('[Hotspot API] Setup expire monitor for router:', routerId)
    const payload = script ? { script } : {}
    const { data } = await api.post<ApiResponse<{ message: string; status: string }>>(
      `/hotspot/${routerId}/expire-monitor`,
      payload
    )
    if (DEBUG) console.log('[Hotspot API] Setup expire monitor response:', data)

    if (!data.success || !data.data) {
      const error = data.error || 'Failed to setup expire monitor'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data
  },

  getExpireMonitorScript: async (routerId: string): Promise<string> => {
    if (DEBUG) console.log('[Hotspot API] Fetching expire monitor script for router:', routerId)
    const { data } = await api.get<ApiResponse<{ script: string }>>(`/hotspot/${routerId}/expire-monitor/script`)
    if (DEBUG) console.log('[Hotspot API] Expire monitor script response:', data)

    if (!data.success || !data.data?.script) {
      const error = data.error || 'Failed to get expire monitor script'
      if (DEBUG) console.error('[Hotspot API] Error:', error)
      throw new Error(error)
    }
    return data.data.script
  },
}
