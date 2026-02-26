import type { ApiRequest } from './types'

const requests: ApiRequest[] = []

// Helper to dispatch API request events
export function logApiRequest(config: { method?: string; url?: string; headers?: unknown; data?: unknown }) {
  const request: ApiRequest = {
    id: `${Date.now()}-${Math.random()}`,
    method: config.method?.toUpperCase() || 'GET',
    url: config.url || '',
    requestTime: new Date(),
    requestHeaders: (config.headers as Record<string, unknown>) || undefined,
    requestData: config.data,
  }
  requests.push(request)
  window.dispatchEvent(new CustomEvent('api-request', { detail: request }))
  return request.id
}

// Helper to dispatch API response events
export function logApiResponse(
  requestId: string | undefined,
  response: { status?: number; statusText?: string; headers?: unknown; data?: unknown } | undefined,
  error?: Error
) {
  if (!requestId) {
    return
  }

  const responseTime = new Date()
  const existingRequestIndex = requests.findIndex(r => r.id === requestId)

  if (existingRequestIndex >= 0) {
    const existingRequest = requests[existingRequestIndex]
    const apiResponse: ApiRequest = {
      ...existingRequest,
      responseTime,
      duration: responseTime.getTime() - existingRequest.requestTime.getTime(),
      status: error ? undefined : response?.status,
      statusText: error ? undefined : response?.statusText,
      responseHeaders: error ? undefined : (response?.headers as Record<string, unknown>),
      responseData: error ? undefined : response?.data,
      error: error?.message,
    }
    requests[existingRequestIndex] = apiResponse
    window.dispatchEvent(new CustomEvent('api-response', { detail: apiResponse }))
  }
}

// Export requests for initial state - use different name to avoid conflict
export const initialRequests = [...requests]

