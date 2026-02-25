import { api } from './axios'
import type { LoginCredentials, LoginResponse, ApiResponse, UserInfo } from '../types'

export const authApi = {
  login: async (credentials: LoginCredentials): Promise<LoginResponse> => {
    const { data } = await api.post<ApiResponse<LoginResponse>>('/auth/login', credentials)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Login failed')
    }
    return data.data
  },

  getMe: async (): Promise<UserInfo> => {
    const { data } = await api.get<ApiResponse<UserInfo>>('/auth/me')
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to get user info')
    }
    return data.data
  },
}
