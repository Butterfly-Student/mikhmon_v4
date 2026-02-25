export interface ApiRequest {
  id: string
  method: string
  url: string
  status?: number
  statusText?: string
  requestTime: Date
  responseTime?: Date
  duration?: number
  requestHeaders?: Record<string, string>
  requestData?: unknown
  responseHeaders?: Record<string, string>
  responseData?: unknown
  error?: string
}
