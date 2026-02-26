export interface ApiRequest {
  id: string
  method: string
  url: string
  status?: number
  statusText?: string
  requestTime: Date
  responseTime?: Date
  duration?: number
  requestHeaders?: Record<string, unknown>
  requestData?: unknown
  responseHeaders?: Record<string, unknown>
  responseData?: unknown
  error?: string
}
