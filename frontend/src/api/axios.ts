import axios from 'axios'
import { useAuthStore } from '../stores/authStore'
import { logApiRequest, logApiResponse } from '../components/debug'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const DEBUG_API = import.meta.env.DEV || localStorage.getItem('debug_api') === 'true'

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 15000, // 15 seconds - more reasonable for debugging
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authState = useAuthStore.getState()
    const token = authState.token

    if (DEBUG_API) {
      console.log('[Auth] Token present:', !!token)
      console.log('[Auth] Authenticated:', authState.isAuthenticated)
    }

    if (token) {
      config.headers.Authorization = `Bearer ${token}`
      if (DEBUG_API) {
        console.log('[Auth] Authorization header set: Bearer ***')
      }
    } else if (DEBUG_API) {
      console.warn('[Auth] No token found in auth store!')
    }

    // Log to API debugger
    const requestId = logApiRequest(config)

    // Attach requestId to config for response interceptor
    ;(config as { metadata?: { requestId: string } }).metadata = { requestId }

    if (DEBUG_API) {
      console.group(`[API Request] ${config.method?.toUpperCase()} ${config.url}`)
      console.log('Base URL:', config.baseURL)
      console.log('Full URL:', API_BASE_URL + config.url)
      console.log('Headers:', {
        'Content-Type': config.headers['Content-Type'],
        'Authorization': config.headers.Authorization ? 'Bearer ***' : 'No token',
      })
      console.log('Data:', config.data)
      console.log('Params:', config.params)
      console.log('Request ID:', requestId)
      console.groupEnd()
    }

    return config
  },
  (error) => {
    if (DEBUG_API) {
      console.error('[API Request Error]', error)
    }
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    const requestId = (response.config as { metadata?: { requestId: string } })?.metadata?.requestId
    logApiResponse(requestId, response)

    if (DEBUG_API) {
      console.group(`[API Response] ${response.config.method?.toUpperCase()} ${response.config.url}`)
      console.log('Status:', response.status)
      console.log('Headers:', response.headers)
      console.log('Data:', response.data)
      console.log('Request ID:', requestId)
      console.groupEnd()
    }
    return response
  },
  (error) => {
    const requestId = (error.config as { metadata?: { requestId: string } })?.metadata?.requestId
    logApiResponse(requestId, error.response, error)

    if (DEBUG_API) {
      console.group(`[API Error] ${error.config?.method?.toUpperCase()} ${error.config?.url}`)
      console.error('Error:', error)
      if (error.response) {
        console.log('Response Status:', error.response.status)
        console.log('Response Data:', error.response.data)
        console.log('Response Headers:', error.response.headers)
      } else if (error.request) {
        console.log('No Response Received - Request:', error.request)
      } else {
        console.log('Error Message:', error.message)
      }
      console.log('Request ID:', requestId)
      console.groupEnd()
    }

    if (error.response?.status === 401) {
      console.warn('[Auth] Unauthorized - logging out...')
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Helper function to toggle debug mode
export function toggleApiDebug(enabled?: boolean) {
  if (enabled === undefined) {
    const current = localStorage.getItem('debug_api') === 'true'
    localStorage.setItem('debug_api', (!current).toString())
    console.log(`[Debug Mode] API logging ${!current ? 'ENABLED' : 'DISABLED'}`)
  } else {
    localStorage.setItem('debug_api', enabled.toString())
    console.log(`[Debug Mode] API logging ${enabled ? 'ENABLED' : 'DISABLED'}`)
  }
}
