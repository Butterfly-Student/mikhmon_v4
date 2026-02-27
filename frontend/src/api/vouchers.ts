import { api } from './axios'
import type { VoucherBatchResult, GenerateVoucherRequest, ApiResponse, HotspotUser } from '../types'

export const vouchersApi = {
  generate: async (routerId: string, request: GenerateVoucherRequest): Promise<VoucherBatchResult> => {
    const { data } = await api.post<ApiResponse<VoucherBatchResult>>(`/mikrotik/${routerId}/vouchers/generate`, request)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to generate vouchers')
    }
    return data.data
  },

  getByComment: async (routerId: string, comment: string): Promise<HotspotUser[]> => {
    const { data } = await api.get<ApiResponse<HotspotUser[]>>(`/mikrotik/${routerId}/vouchers?comment=${encodeURIComponent(comment)}`)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to get vouchers')
    }
    return data.data
  },

  deleteByComment: async (routerId: string, comment: string): Promise<void> => {
    const { data } = await api.delete<ApiResponse<void>>(`/mikrotik/${routerId}/vouchers?comment=${encodeURIComponent(comment)}`)
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete vouchers')
    }
  },
}
