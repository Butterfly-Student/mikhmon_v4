import { useState, useEffect, useRef, useCallback } from 'react'
import { Activity, TrendingDown, TrendingUp, RefreshCw } from 'lucide-react'
import { Card, Button, Select } from '../ui'
import { useRouterStore } from '../../stores/routerStore'
import { mikrotikApi } from '../../api/mikrotik'

interface QueueStats {
  name: string
  bytesIn: number
  bytesOut: number
  rateIn: number
  rateOut: number
  timestamp: number
}

export function QueueMonitor() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  const [queues, setQueues] = useState<string[]>([])
  const [selectedQueue, setSelectedQueue] = useState('')
  const [stats, setStats] = useState<QueueStats | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [isMonitoring, setIsMonitoring] = useState(false)
  const [error, setError] = useState('')

  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // Fetch available queues
  useEffect(() => {
    const fetchQueues = async () => {
      try {
        if (!routerId) return
        const queuesData = await mikrotikApi.getAllQueues(routerId.toString())
        // Adjust depending on the structure of queuesData, assuming objects with a 'name' field
        const queueNames = queuesData.map((q: any) => q.name || q['.id'])
        setQueues(queueNames)
      } catch (err) {
        console.error('Failed to fetch queues:', err)
      }
    }
    fetchQueues()
  }, [routerId])

  const connect = useCallback(() => {
    if (!selectedQueue) return

    const wsKey = localStorage.getItem('ws_key') || 'mikhmon-ws-internal-key'
    const wsUrl = `ws://${window.location.host}/api/v1/ws/mikrotik/monitor/queue/${routerId}?key=${wsKey}`

    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      setError('')
      // Start monitoring
      ws.send(JSON.stringify({
        action: 'start',
        name: selectedQueue,
        interval: 1
      }))
      setIsMonitoring(true)
    }

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      if (data.type === 'error') {
        setError(data.message)
        return
      }
      if (data.type === 'status') {
        console.log('Queue monitor status:', data.status)
        return
      }
      setStats(data)
    }

    ws.onerror = (err) => {
      console.error('Queue monitor error:', err)
      setError('Connection error')
    }

    ws.onclose = () => {
      setIsConnected(false)
      setIsMonitoring(false)
      // Auto reconnect after 3s
      reconnectTimeoutRef.current = setTimeout(() => {
        if (selectedQueue) connect()
      }, 3000)
    }
  }, [routerId, selectedQueue])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    if (wsRef.current) {
      wsRef.current.send(JSON.stringify({ action: 'stop' }))
      wsRef.current.close()
      wsRef.current = null
    }
    setIsMonitoring(false)
    setIsConnected(false)
  }, [])

  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [disconnect])

  useEffect(() => {
    if (selectedQueue) {
      disconnect()
      connect()
    }
  }, [selectedQueue, connect, disconnect])

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const formatBitsPerSecond = (bps: number) => {
    if (bps === 0) return '0 bps'
    const k = 1000
    const sizes = ['bps', 'Kbps', 'Mbps', 'Gbps']
    const i = Math.floor(Math.log(bps) / Math.log(k))
    return parseFloat((bps / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  return (
    <Card>
      <Card.Header>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Activity className="w-5 h-5 text-primary-500" />
            <h3 className="font-semibold text-gray-900 dark:text-white">Queue Monitor</h3>
          </div>
          <div className="flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-success-500 animate-pulse' : 'bg-gray-300'}`} />
            <span className="text-xs text-gray-500">
              {isConnected ? 'Connected' : 'Disconnected'}
            </span>
          </div>
        </div>
      </Card.Header>
      <Card.Body className="space-y-4">
        {/* Queue Selector */}
        <div className="flex gap-2">
          <Select
            options={[
              { value: '', label: 'Select Queue' },
              ...queues.map(q => ({ value: q, label: q }))
            ]}
            value={selectedQueue}
            onChange={(e) => setSelectedQueue(e.target.value)}
            className="flex-1"
          />
          <Button
            variant="ghost"
            size="sm"
            leftIcon={<RefreshCw className={`w-4 h-4 ${isMonitoring ? 'animate-spin' : ''}`} />}
            onClick={() => {
              disconnect()
              connect()
            }}
          >
            Refresh
          </Button>
        </div>

        {error && (
          <div className="p-3 bg-danger-50 text-danger-700 rounded-lg text-sm">
            {error}
          </div>
        )}

        {selectedQueue && stats ? (
          <div className="space-y-4">
            {/* Queue Name */}
            <div className="text-center">
              <span className="text-lg font-semibold text-gray-900 dark:text-white">
                {stats.name}
              </span>
            </div>

            {/* Traffic Stats */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-success-50 dark:bg-success-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-success-600 mb-1">
                  <TrendingDown className="w-4 h-4" />
                  <span className="text-xs font-medium">Download (RX)</span>
                </div>
                <div className="text-lg font-bold text-success-700">
                  {formatBitsPerSecond(stats.rateIn)}
                </div>
                <div className="text-xs text-success-600 mt-1">
                  Total: {formatBytes(stats.bytesIn)}
                </div>
              </div>

              <div className="bg-primary-50 dark:bg-primary-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-primary-600 mb-1">
                  <TrendingUp className="w-4 h-4" />
                  <span className="text-xs font-medium">Upload (TX)</span>
                </div>
                <div className="text-lg font-bold text-primary-700">
                  {formatBitsPerSecond(stats.rateOut)}
                </div>
                <div className="text-xs text-primary-600 mt-1">
                  Total: {formatBytes(stats.bytesOut)}
                </div>
              </div>
            </div>

            {/* Last Update */}
            <div className="text-center text-xs text-gray-400">
              Last update: {new Date(stats.timestamp).toLocaleTimeString()}
            </div>
          </div>
        ) : selectedQueue ? (
          <div className="text-center py-8 text-gray-400">
            <Activity className="w-8 h-8 mx-auto mb-2 animate-pulse" />
            <p>Waiting for data...</p>
          </div>
        ) : (
          <div className="text-center py-8 text-gray-400">
            <p>Select a queue to monitor</p>
          </div>
        )}
      </Card.Body>
    </Card>
  )
}
