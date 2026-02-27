import { useEffect, useRef, useState, useCallback } from 'react'

export interface ResourceStats {
    cpuUsed: number
    freeMemory: number
    totalMemory: number
    freeHddSpace: number
    totalHddSpace: number
    writeSectSinceReboot: number
    uptime: string
    timestamp: string | Date
}

export function useResourceWebSocket(routerId: string | number | undefined) {
    const [stats, setStats] = useState<ResourceStats | null>(null)
    const [isConnected, setIsConnected] = useState(false)
    const wsRef = useRef<WebSocket | null>(null)
    const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const connect = useCallback(() => {
        if (!routerId) {
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
        const wsUrl = `${wsProtocol}//${wsHost}/api/v1/ws/mikrotik/monitor/resource/${routerId}?key=${internalKey}`

        console.log('[WS Resource] Connecting to:', wsUrl)

        const ws = new WebSocket(wsUrl)

        ws.onopen = () => {
            console.log('[WS Resource] Connected!')
            setIsConnected(true)
        }

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                if (data && typeof data.cpuUsed === 'number') {
                    setStats(data as ResourceStats)
                }
            } catch {
                // Ignore parse errors
            }
        }

        ws.onclose = (event) => {
            console.log('[WS Resource] Closed:', event.code, event.reason)
            setIsConnected(false)
            reconnectTimeoutRef.current = setTimeout(connect, 3000)
        }

        ws.onerror = (error) => {
            console.error('[WS Resource] Error:', error)
            ws.close()
        }

        wsRef.current = ws
    }, [routerId])

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
