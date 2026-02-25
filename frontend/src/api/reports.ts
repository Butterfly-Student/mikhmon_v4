import { api } from './axios'
import type { SalesReport, ApiResponse } from '../types'

export const reportsApi = {
  getSalesReport: async (
    routerId: string,
    filters: { date?: string; month?: string; year?: string }
  ): Promise<SalesReport[]> => {
    const params = new URLSearchParams()
    if (filters.date) params.append('date', filters.date)
    if (filters.month) params.append('month', filters.month)
    if (filters.year) params.append('year', filters.year)

    const { data } = await api.get<ApiResponse<SalesReport[]>>(
      `/reports/${routerId}/sales?${params}`
    )
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to get sales report')
    }
    return data.data
  },

  getSummary: async (
    routerId: string,
    month?: string,
    year?: string
  ): Promise<{ totalSales: number; totalTransactions: number; averageTicket: number }> => {
    const params = new URLSearchParams()
    if (month) params.append('month', month)
    if (year) params.append('year', year)

    const { data } = await api.get<ApiResponse<any>>(`/reports/${routerId}/summary?${params}`)
    if (!data.success || !data.data) {
      throw new Error(data.error || 'Failed to get summary')
    }
    return data.data
  },

  exportToCSV: async (
    routerId: string,
    filters: { date?: string; month?: string; year?: string }
  ): Promise<string> => {
    const params = new URLSearchParams()
    if (filters.date) params.append('date', filters.date)
    if (filters.month) params.append('month', filters.month)
    if (filters.year) params.append('year', filters.year)

    const response = await api.get(`/reports/${routerId}/export?${params}`, {
      responseType: 'text',
    })
    return response.data
  },
}
