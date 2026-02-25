import { api } from './axios'
import type { Router, ApiResponse } from '../types'

export const routersApi = {
  getAll: async (): Promise<Router[]> => {
    const { data } = await api.get<ApiResponse<Router[]>>('/routers')
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to get routers')
    }
    return data.data
  },

  getById: async (id: string): Promise<Router> => {
    const { data } = await api.get<ApiResponse<Router>>(`/routers/${id}`)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to get router')
    }
    return data.data
  },

  create: async (router: Partial<Router>): Promise<Router> => {
    const { data } = await api.post<ApiResponse<Router>>('/routers', router)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to create router')
    }
    return data.data
  },

  update: async (id: string, router: Partial<Router>): Promise<Router> => {
    const { data } = await api.put<ApiResponse<Router>>(`/routers/${id}`, router)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to update router')
    }
    return data.data
  },

  delete: async (id: string): Promise<void> => {
    const { data } = await api.delete<ApiResponse<void>>(`/routers/${id}`)
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete router')
    }
  },

  testConnection: async (id: string): Promise<void> => {
    const { data } = await api.post<ApiResponse<void>>(`/routers/${id}/test`)
    if (!data.success) {
      throw new Error(data.error || 'Connection failed')
    }
  },
}
