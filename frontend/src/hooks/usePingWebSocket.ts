import { useEffect, useRef, useState, useCallback } from 'react'

interface PingStats {
  current: number
  min: number
  max: number
  avg: number
}

export interface PingConfig {
  /** Target IP/domain (default: '8.8.8.8') */
  address?: string
  /** Interval antar ping dalam detik (default: 1) */
  interval?: number
  /** Jumlah ping — 0 = infinite (default: 0) */
  count?: number
  /** Ukuran packet dalam bytes (default: 64) */
  size?: number
}

const DEFAULT_CONFIG: Required<PingConfig> = {
  address: '8.8.8.8',
  interval: 1,
  count: 0,
  size: 64,
}

export function usePingWebSocket(
  routerId: string | number | undefined,
  config: PingConfig = {}
) {
  // Merge user config dengan defaults
  const resolvedConfig: Required<PingConfig> = {
    ...DEFAULT_CONFIG,
    ...config,
  }

  const [latency, setLatency] = useState<number | null>(null)
  const [stats, setStats] = useState<PingStats>({ current: 0, min: 0, max: 0, avg: 0 })
  const [isConnected, setIsConnected] = useState(false)
  const [isPinging, setIsPinging] = useState(false)

  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const latencyHistoryRef = useRef<number[]>([])
  const connectRef = useRef<(() => void) | null>(null)
  // Simpan resolved config terbaru agar tidak stale di closure
  const configRef = useRef(resolvedConfig)

  const connect = useCallback(() => {
    if (!routerId) return

    // Internal key untuk WebSocket auth (sama dengan backend)
    const internalKey = import.meta.env.VITE_WS_KEY || 'mikhmon-ws-internal-key'

    // Build WebSocket URL (gunakan path relatif untuk lewat Vite proxy)
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsHost = import.meta.env.DEV
      ? `${window.location.host}`  // Vite dev server dengan proxy
      : (import.meta.env.VITE_API_URL
        ? import.meta.env.VITE_API_URL.replace(/^https?:\/\//, '')
        : `${window.location.hostname}:8080`)
    const wsUrl = `${wsProtocol}//${wsHost}/api/v1/ws/mikrotik/monitor/ping/${routerId}?key=${internalKey}`

    console.log('[WS Ping] Connecting to:', wsUrl)

    const ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('[WS Ping] Connected!')
      setIsConnected(true)
      // Kirim konfigurasi ping ke backend saat koneksi terbuka
      const cfg = configRef.current
      ws.send(JSON.stringify({
        action: 'start',
        address: cfg.address,
        interval: cfg.interval,
        count: cfg.count,
        size: cfg.size,
      }))
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        if (data.type === 'status' && data.status === 'started') {
          setIsPinging(true)
        }

        if (data.type === 'status' && data.status === 'stopped') {
          setIsPinging(false)
        }

        if (data.received && typeof data.timeMs === 'number') {
          setLatency(data.timeMs)

          // Update stats
          latencyHistoryRef.current.push(data.timeMs)
          if (latencyHistoryRef.current.length > 20) {
            latencyHistoryRef.current.shift()
          }

          const times = latencyHistoryRef.current
          setStats({
            current: data.timeMs,
            min: Math.min(...times),
            max: Math.max(...times),
            avg: times.reduce((a, b) => a + b, 0) / times.length,
          })
        }
      } catch {
        // Ignore parse errors
      }
    }

    ws.onclose = (event) => {
      console.log('[WS Ping] Closed:', event.code, event.reason)
      setIsConnected(false)
      setIsPinging(false)
      // Auto reconnect after 3 seconds
      reconnectTimeoutRef.current = window.setTimeout(() => {
        if (connectRef.current) {
          connectRef.current()
        }
      }, 3000)
    }

    ws.onerror = (error) => {
      console.error('[WS Ping] Error:', error)
      ws.close()
    }

    wsRef.current = ws
  }, [routerId])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    if (wsRef.current) {
      if (wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ action: 'stop' }))
      }
      wsRef.current.close()
    }
    wsRef.current = null
  }, [])

  /** Ganti address ping (dan opsional config) tanpa reconnect */
  const changeAddress = useCallback((newAddress: string, newConfig: PingConfig = {}) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const cfg = { ...configRef.current, ...newConfig, address: newAddress }
      wsRef.current.send(JSON.stringify({
        action: 'start',
        address: cfg.address,
        interval: cfg.interval,
        count: cfg.count,
        size: cfg.size,
      }))
    }
  }, [])

  // Keep refs in sync after render — avoids render-time ref mutation
  useEffect(() => {
    configRef.current = resolvedConfig
    connectRef.current = connect
  })

  // Reset state when there is no routerId (avoids synchronous setState inside connect/disconnect)
  useEffect(() => {
    if (!routerId) {
      setLatency(null)
      setIsConnected(false)
      setIsPinging(false)
      latencyHistoryRef.current = []
    }
  }, [routerId])

  useEffect(() => {
    connect()
    return disconnect
  }, [connect, disconnect])

  return {
    latency,
    stats,
    isConnected,
    isPinging,
    changeAddress,
  }
}
