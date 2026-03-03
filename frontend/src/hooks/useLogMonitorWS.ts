import { useEffect, useRef, useState, useCallback } from 'react'
import type { LogEntryWithSeq } from '../types'

export type LogEndpoint = 'logs' | 'hotspot-logs' | 'ppp-logs'

const MAX_ENTRIES = 1000 // Circular buffer size - must match backend

function getWsUrl(endpoint: LogEndpoint, routerId: string | number): string {
    const key = import.meta.env.VITE_WS_KEY || 'mikhmon-ws-internal-key'
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.DEV
        ? window.location.host
        : (import.meta.env.VITE_API_URL
            ? import.meta.env.VITE_API_URL.replace(/^https?:\/\//, '')
            : `${window.location.hostname}:8080`)
    return `${proto}//${host}/api/v1/ws/mikrotik/monitor/${endpoint}/${routerId}?key=${key}`
}

/**
 * Circular Buffer for log entries
 * Maintains max MAX_ENTRIES entries, removing oldest when full
 */
class CircularLogBuffer {
    private buffer: LogEntryWithSeq[] = []
    private maxSize: number = MAX_ENTRIES

    constructor(maxSize?: number) {
        if (maxSize) this.maxSize = maxSize
    }

    /**
     * Initialize buffer with data (replaces existing)
     */
    init(data: LogEntryWithSeq[]): void {
        this.buffer = data.slice(-this.maxSize)
    }

    /**
     * Append entries to buffer
     * If buffer exceeds maxSize, oldest entries are removed
     */
    append(entries: LogEntryWithSeq[]): void {
        if (entries.length === 0) return

        const newBuffer = [...this.buffer, ...entries]
        
        // Keep only the last maxSize entries (circular buffer behavior)
        if (newBuffer.length > this.maxSize) {
            this.buffer = newBuffer.slice(newBuffer.length - this.maxSize)
        } else {
            this.buffer = newBuffer
        }
    }

    /**
     * Get all entries sorted by sequence (ascending)
     */
    getEntries(): LogEntryWithSeq[] {
        return [...this.buffer].sort((a, b) => a.seq - b.seq)
    }

    /**
     * Get current size
     */
    get size(): number {
        return this.buffer.length
    }

    /**
     * Clear buffer
     */
    clear(): void {
        this.buffer = []
    }
}

// Type guards for WebSocket messages
interface WSMessageInit {
    type: 'init'
    data: LogEntryWithSeq[]
    meta: { topics?: string; maxSize?: number; count?: number; routerID?: number }
}

interface WSMessageUpdate {
    type: 'update'
    data: LogEntryWithSeq[]
    meta: { batchSize?: number; totalSeq?: number }
}

interface WSMessageError {
    type: 'error'
    message: string
}

interface WSMessageStatus {
    type: 'status'
    status: string
}

type WSMessage = WSMessageInit | WSMessageUpdate | WSMessageError | WSMessageStatus

function isValidMessage(data: unknown): data is WSMessage {
    if (!data || typeof data !== 'object') return false
    return 'type' in data && typeof (data as WSMessage).type === 'string'
}

function isInitMessage(data: WSMessage): data is WSMessageInit {
    return data.type === 'init' && 'data' in data && Array.isArray((data as WSMessageInit).data)
}

function isUpdateMessage(data: WSMessage): data is WSMessageUpdate {
    return data.type === 'update' && 'data' in data && Array.isArray((data as WSMessageUpdate).data)
}

function isErrorMessage(data: WSMessage): data is WSMessageError {
    return data.type === 'error' && 'message' in data
}

function isStatusMessage(data: WSMessage): data is WSMessageStatus {
    return data.type === 'status' && 'status' in data
}

export function useLogMonitorWS(
    routerId: string | number | undefined,
    endpoint: LogEndpoint,
    paused: boolean = false
) {
    const [entries, setEntries] = useState<LogEntryWithSeq[]>([])
    const [isConnected, setIsConnected] = useState(false)
    const [hasReceivedData, setHasReceivedData] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [meta, setMeta] = useState<{ topics?: string; maxSize?: number }>({})
    const wsRef = useRef<WebSocket | null>(null)
    const reconnectRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const pausedRef = useRef(paused)
    const bufferRef = useRef(new CircularLogBuffer(MAX_ENTRIES))
    const messageCountRef = useRef(0)
    pausedRef.current = paused

    const connect = useCallback(() => {
        if (!routerId) {
            console.log(`[WS Log ${endpoint}] No routerId, skipping connection`)
            return
        }
        
        // Reset state on new connection
        setError(null)
        setHasReceivedData(false)
        messageCountRef.current = 0
        
        const url = getWsUrl(endpoint, routerId)
        console.log(`[WS Log ${endpoint}] Connecting to:`, url)
        
        let ws: WebSocket
        try {
            ws = new WebSocket(url)
        } catch (err) {
            console.error(`[WS Log ${endpoint}] Failed to create WebSocket:`, err)
            setError('Failed to create WebSocket connection')
            reconnectRef.current = setTimeout(connect, 3000)
            return
        }

        ws.onopen = () => {
            console.log(`[WS Log ${endpoint}] Connected successfully`)
            setIsConnected(true)
            setError(null)
            if (reconnectRef.current) clearTimeout(reconnectRef.current)
        }

        ws.onmessage = (event) => {
            messageCountRef.current++
            
            if (pausedRef.current) {
                console.log(`[WS Log ${endpoint}] Message #${messageCountRef.current} received but paused`)
                return
            }
            
            let rawData: unknown
            try {
                rawData = JSON.parse(event.data)
            } catch (err) {
                console.error(`[WS Log ${endpoint}] JSON parse error:`, err)
                console.error(`[WS Log ${endpoint}] Raw data:`, event.data)
                return
            }
            
            console.log(`[WS Log ${endpoint}] Message #${messageCountRef.current}:`, rawData)
            
            // Check if it's a valid message object with type
            if (!isValidMessage(rawData)) {
                console.log(`[WS Log ${endpoint}] Invalid message format (no type field)`)
                // Try legacy format - array of entries
                if (Array.isArray(rawData)) {
                    const entries = rawData as LogEntryWithSeq[]
                    console.log(`[WS Log ${endpoint}] Legacy array format, entries:`, entries.length)
                    if (entries.length > 0) {
                        setHasReceivedData(true)
                        bufferRef.current.append(entries)
                        setEntries(bufferRef.current.getEntries())
                    }
                }
                return
            }
            
            // Handle different message types
            if (isInitMessage(rawData)) {
                console.log(`[WS Log ${endpoint}] Init message with`, rawData.data.length, 'entries')
                setHasReceivedData(true)
                bufferRef.current.init(rawData.data)
                setMeta({
                    topics: rawData.meta?.topics,
                    maxSize: rawData.meta?.maxSize
                })
                setEntries(bufferRef.current.getEntries())
            } else if (isUpdateMessage(rawData)) {
                console.log(`[WS Log ${endpoint}] Update message with`, rawData.data.length, 'entries')
                if (rawData.data.length > 0) {
                    setHasReceivedData(true)
                }
                bufferRef.current.append(rawData.data)
                setEntries(bufferRef.current.getEntries())
            } else if (isErrorMessage(rawData)) {
                console.error(`[WS Log ${endpoint}] Error message:`, rawData.message)
                setError(rawData.message)
            } else if (isStatusMessage(rawData)) {
                console.log(`[WS Log ${endpoint}] Status message:`, rawData.status)
            } else {
                console.log(`[WS Log ${endpoint}] Unknown message type:`, (rawData as WSMessage).type)
            }
        }

        ws.onclose = (ev) => {
            console.log(`[WS Log ${endpoint}] Closed with code:`, ev.code, 'reason:', ev.reason)
            setIsConnected(false)
            reconnectRef.current = setTimeout(connect, 3000)
        }

        ws.onerror = (err) => {
            console.error(`[WS Log ${endpoint}] WebSocket error:`, err)
            setError('WebSocket connection error')
            ws.close()
        }
        
        wsRef.current = ws
    }, [routerId, endpoint])

    const disconnect = useCallback(() => {
        if (reconnectRef.current) clearTimeout(reconnectRef.current)
        wsRef.current?.close()
        wsRef.current = null
    }, [])

    const clearEntries = useCallback(() => {
        bufferRef.current.clear()
        setEntries([])
        setHasReceivedData(false)
    }, [])

    useEffect(() => {
        if (!routerId) { 
            bufferRef.current.clear()
            setEntries([])
            setIsConnected(false)
            setHasReceivedData(false)
            setError(null)
            setMeta({})
        }
    }, [routerId])

    useEffect(() => {
        connect()
        return disconnect
    }, [connect, disconnect])

    return { 
        entries, 
        isConnected, 
        hasReceivedData,
        error,
        meta, 
        clearEntries 
    }
}
