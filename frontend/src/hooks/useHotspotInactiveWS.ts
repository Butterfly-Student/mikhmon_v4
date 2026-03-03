import { useEffect, useRef, useState, useCallback } from 'react'
import type { HotspotUser } from '../types'

function getWsUrl(endpoint: string, routerId: string | number): string {
    const key = import.meta.env.VITE_WS_KEY || 'mikhmon-ws-internal-key'
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.DEV
        ? window.location.host
        : (import.meta.env.VITE_API_URL
            ? import.meta.env.VITE_API_URL.replace(/^https?:\/\//, '')
            : `${window.location.hostname}:8080`)
    return `${proto}//${host}/api/v1/ws/mikrotik/monitor/${endpoint}/${routerId}?key=${key}`
}

export function useHotspotInactiveWS(routerId: string | number | undefined) {
    const [users, setUsers] = useState<HotspotUser[]>([])
    const [isConnected, setIsConnected] = useState(false)
    const [lastUpdate, setLastUpdate] = useState<Date | null>(null)
    const wsRef = useRef<WebSocket | null>(null)
    const reconnectRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const connect = useCallback(() => {
        if (!routerId) return
        const url = getWsUrl('hotspot-inactive', routerId)
        console.log('[WS HotspotInactive] Connecting to:', url)
        const ws = new WebSocket(url)

        ws.onopen = () => {
            console.log('[WS HotspotInactive] Connected')
            setIsConnected(true)
            if (reconnectRef.current) clearTimeout(reconnectRef.current)
        }

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                if (Array.isArray(data)) {
                    setUsers(data as HotspotUser[])
                    setLastUpdate(new Date())
                }
            } catch { /* ignore */ }
        }

        ws.onclose = (ev) => {
            console.log('[WS HotspotInactive] Closed:', ev.code)
            setIsConnected(false)
            reconnectRef.current = setTimeout(connect, 3000)
        }

        ws.onerror = () => ws.close()
        wsRef.current = ws
    }, [routerId])

    const disconnect = useCallback(() => {
        if (reconnectRef.current) clearTimeout(reconnectRef.current)
        wsRef.current?.close()
        wsRef.current = null
    }, [])

    useEffect(() => {
        if (!routerId) { setUsers([]); setIsConnected(false) }
    }, [routerId])

    useEffect(() => {
        connect()
        return disconnect
    }, [connect, disconnect])

    return { users, isConnected, lastUpdate }
}
