import { useEffect, useRef, useState, useCallback } from 'react'

export interface TrafficStats {
    rxBitsPerSecond: number
    txBitsPerSecond: number
    timestamp: string | Date
}

export function useTrafficWebSocket(routerId: string | number | undefined, interfaceName?: string) {
    const [stats, setStats] = useState<TrafficStats | null>(null)
    const [isConnected, setIsConnected] = useState(false)

    const wsRef = useRef<WebSocket | null>(null)
    const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const connect = useCallback(() => {
        if (!routerId || !interfaceName) {
            setStats(null)
            setIsConnected(false)
            return
        }

        const internalKey = import.meta.env.VITE_WS_KEY || 'mikhmon-ws-internal-key'
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsHost = import.meta.env.DEV
            ? `${window.location.host}`
            : (import.meta.env.VITE_API_URL
                ? import.meta.env.VITE_API_URL.replace(/^https?:\/\//, '')
                : `${window.location.hostname}:8080`)

        // Using the path specified in implementation plan & Golang router
        let wsUrl = `${wsProtocol}//${wsHost}/api/v1/ws/mikrotik/monitor/interface/${routerId}?key=${internalKey}`
        if (interfaceName) {
            wsUrl += `&interface=${encodeURIComponent(interfaceName)}`
        }

        console.log('[WS Traffic] Connecting to:', wsUrl)

        const ws = new WebSocket(wsUrl)

        ws.onopen = () => {
            console.log('[WS Traffic] Connected!')
            setIsConnected(true)
        }

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                if (data && typeof data.rxBitsPerSecond === 'number') {
                    setStats(data as TrafficStats)
                }
            } catch {
                // Ignore parse errors
            }
        }

        ws.onclose = (event) => {
            console.log('[WS Traffic] Closed:', event.code, event.reason)
            setIsConnected(false)
            reconnectTimeoutRef.current = setTimeout(connect, 3000)
        }

        ws.onerror = (error) => {
            console.error('[WS Traffic] Error:', error)
            ws.close()
        }

        wsRef.current = ws
    }, [routerId, interfaceName])

    const disconnect = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current)
        }
        if (wsRef.current) {
            wsRef.current.close()
        }
        wsRef.current = null
        setIsConnected(false)
    }, [])

    useEffect(() => {
        connect()
        return disconnect
    }, [connect, disconnect])

    return {
        stats,
        isConnected,
    }
}
