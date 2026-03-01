import { useState, useEffect, useRef, useCallback } from 'react'
import { Activity, TrendingDown, TrendingUp, RefreshCw } from 'lucide-react'
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
} from 'recharts'
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

interface ChartPoint {
  time: string
  rateIn: number
  rateOut: number
}

const MAX_HISTORY = 30

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

const yAxisFormatter = (v: number) => {
  if (v === 0) return '0'
  const k = 1000
  const sizes = ['', 'K', 'M', 'G']
  const i = Math.floor(Math.log(v) / Math.log(k))
  return parseFloat((v / Math.pow(k, i)).toFixed(1)) + sizes[i]
}

export function QueueMonitor() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  // API returns string[] directly
  const [queues, setQueues] = useState<string[]>([])
  const [selectedQueue, setSelectedQueue] = useState('')
  const [stats, setStats] = useState<QueueStats | null>(null)
  const [history, setHistory] = useState<ChartPoint[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const [isMonitoring, setIsMonitoring] = useState(false)
  const [error, setError] = useState('')

  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const mountedRef = useRef(true)

  useEffect(() => {
    mountedRef.current = true
    return () => { mountedRef.current = false }
  }, [])

  useEffect(() => {
    const fetchQueues = async () => {
      try {
        if (!routerId) return
        // getAllQueues already returns string[]
        const queueNames = await mikrotikApi.getAllQueues(routerId.toString())
        setQueues(queueNames)
        if (queueNames.length > 0) setSelectedQueue(queueNames[0])
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
      ws.send(JSON.stringify({ action: 'start', name: selectedQueue, interval: 1 }))
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
      const point: QueueStats = data
      setStats(point)
      setHistory((prev) => {
        const next = [
          ...prev,
          {
            time: new Date(point.timestamp).toLocaleTimeString(),
            rateIn: point.rateIn,
            rateOut: point.rateOut,
          },
        ]
        return next.slice(-MAX_HISTORY)
      })
    }

    ws.onerror = (err) => {
      console.error('Queue monitor error:', err)
      setError('Connection error')
    }

    ws.onclose = () => {
      if (!mountedRef.current) return
      setIsConnected(false)
      setIsMonitoring(false)
      reconnectTimeoutRef.current = setTimeout(() => {
        if (mountedRef.current && selectedQueue) connect()
      }, 3000)
    }
  }, [routerId, selectedQueue])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    if (wsRef.current) {
      const ws = wsRef.current
      wsRef.current = null
      if (ws.readyState === WebSocket.OPEN) {
        try { ws.send(JSON.stringify({ action: 'stop' })) } catch (_) {}
      }
      ws.onclose = null
      ws.close()
    }
    setIsMonitoring(false)
    setIsConnected(false)
  }, [])

  useEffect(() => () => { disconnect() }, [disconnect])

  useEffect(() => {
    if (selectedQueue) {
      setHistory([])
      disconnect()
      connect()
    }
  }, [selectedQueue, connect, disconnect])

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
              ...queues.map((q) => ({ value: q, label: q })),
            ]}
            value={selectedQueue}
            onChange={(e) => setSelectedQueue(e.target.value)}
            className="flex-1"
          />
          <Button
            variant="ghost"
            size="sm"
            leftIcon={<RefreshCw className={`w-4 h-4 ${isMonitoring ? 'animate-spin' : ''}`} />}
            onClick={() => { disconnect(); connect() }}
          >
            Refresh
          </Button>
        </div>

        {error && (
          <div className="p-3 bg-danger-50 text-danger-700 rounded-lg text-sm">{error}</div>
        )}

        {selectedQueue && stats ? (
          <div className="space-y-4">
            {/* Queue Name */}
            <div className="text-center">
              <span className="text-lg font-semibold text-gray-900 dark:text-white">{stats.name}</span>
            </div>

            {/* Current Stats */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-success-50 dark:bg-success-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-success-600 mb-1">
                  <TrendingDown className="w-4 h-4" />
                  <span className="text-xs font-medium">Download (RX)</span>
                </div>
                <div className="text-lg font-bold text-success-700">
                  {formatBitsPerSecond(stats.rateIn)}
                </div>
                <div className="text-xs text-success-600 mt-1">Total: {formatBytes(stats.bytesIn)}</div>
              </div>

              <div className="bg-primary-50 dark:bg-primary-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-primary-600 mb-1">
                  <TrendingUp className="w-4 h-4" />
                  <span className="text-xs font-medium">Upload (TX)</span>
                </div>
                <div className="text-lg font-bold text-primary-700">
                  {formatBitsPerSecond(stats.rateOut)}
                </div>
                <div className="text-xs text-primary-600 mt-1">Total: {formatBytes(stats.bytesOut)}</div>
              </div>
            </div>

            {/* Chart */}
            {history.length > 1 && (
              <div className="h-48">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={history} margin={{ top: 4, right: 4, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="qRxGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3} />
                        <stop offset="95%" stopColor="#22c55e" stopOpacity={0} />
                      </linearGradient>
                      <linearGradient id="qTxGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3} />
                        <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <XAxis dataKey="time" tick={{ fontSize: 10 }} interval="preserveStartEnd" />
                    <YAxis tickFormatter={yAxisFormatter} tick={{ fontSize: 10 }} width={48} />
                    <Tooltip
                      formatter={((v: number, name: string) => [
                        formatBitsPerSecond(v),
                        name === 'rateIn' ? 'Download' : 'Upload',
                      ]) as any}
                      labelStyle={{ fontSize: 11 }}
                      contentStyle={{ fontSize: 11 }}
                    />
                    <Legend
                      formatter={(v) => (v === 'rateIn' ? 'Download (RX)' : 'Upload (TX)')}
                      wrapperStyle={{ fontSize: 11 }}
                    />
                    <Area type="monotone" dataKey="rateIn" stroke="#22c55e" fill="url(#qRxGrad)" strokeWidth={2} dot={false} isAnimationActive={false} />
                    <Area type="monotone" dataKey="rateOut" stroke="#6366f1" fill="url(#qTxGrad)" strokeWidth={2} dot={false} isAnimationActive={false} />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            )}

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
