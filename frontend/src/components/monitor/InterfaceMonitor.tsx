import { useState, useEffect, useRef, useCallback } from 'react'
import { Network, RefreshCw, ArrowDown, ArrowUp } from 'lucide-react'
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

interface InterfaceStats {
  name: string
  txBitsPerSecond: number
  rxBitsPerSecond: number
  txPacketsPerSecond: number
  rxPacketsPerSecond: number
  timestamp: number
}

interface ChartPoint {
  time: string
  rx: number
  tx: number
}

const MAX_HISTORY = 30

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

export function InterfaceMonitor() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  const [interfaces, setInterfaces] = useState<string[]>([])
  const [selectedInterface, setSelectedInterface] = useState('')
  const [stats, setStats] = useState<InterfaceStats | null>(null)
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
    const fetchInterfaces = async () => {
      try {
        if (!routerId) return
        const interfacesData = await mikrotikApi.getInterfaces(routerId.toString())
        const names = interfacesData.map((i) => i.name).filter((n): n is string => !!n)
        setInterfaces(names)
        if (names.length > 0 && names[0]) setSelectedInterface(names[0])
      } catch (err) {
        console.error('Failed to fetch interfaces:', err)
      }
    }
    fetchInterfaces()
  }, [routerId])

  const connect = useCallback(() => {
    if (!selectedInterface) return

    const wsKey = localStorage.getItem('ws_key') || 'mikhmon-ws-internal-key'
    const wsUrl = `ws://${window.location.host}/api/v1/ws/mikrotik/monitor/interface/${routerId}?key=${wsKey}`

    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      setError('')
      ws.send(JSON.stringify({ action: 'start', name: selectedInterface, interval: 1 }))
      setIsMonitoring(true)
    }

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      if (data.type === 'error') {
        setError(data.message)
        return
      }
      if (data.type === 'status') {
        console.log('Interface monitor status:', data.status)
        return
      }
      const point: InterfaceStats = data
      setStats(point)
      setHistory((prev) => {
        const next = [
          ...prev,
          {
            time: new Date(point.timestamp).toLocaleTimeString(),
            rx: point.rxBitsPerSecond,
            tx: point.txBitsPerSecond,
          },
        ]
        return next.slice(-MAX_HISTORY)
      })
    }

    ws.onerror = (err) => {
      console.error('Interface monitor error:', err)
      setError('Connection error')
    }

    ws.onclose = () => {
      if (!mountedRef.current) return
      setIsConnected(false)
      setIsMonitoring(false)
      reconnectTimeoutRef.current = setTimeout(() => {
        if (mountedRef.current && selectedInterface) connect()
      }, 3000)
    }
  }, [routerId, selectedInterface])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    if (wsRef.current) {
      const ws = wsRef.current
      wsRef.current = null
      if (ws.readyState === WebSocket.OPEN) {
        try { ws.send(JSON.stringify({ action: 'stop' })) } catch (_) { }
      }
      ws.onclose = null
      ws.close()
    }
    setIsMonitoring(false)
    setIsConnected(false)
  }, [])

  useEffect(() => () => { disconnect() }, [disconnect])

  useEffect(() => {
    if (selectedInterface) {
      setHistory([])
      disconnect()
      connect()
    }
  }, [selectedInterface, connect, disconnect])

  return (
    <Card>
      <Card.Header>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Network className="w-5 h-5 text-secondary-500" />
            <h3 className="font-semibold text-gray-900 dark:text-white">Interface Monitor</h3>
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
        {/* Interface Selector */}
        <div className="flex gap-2">
          <Select
            options={[
              { value: '', label: 'Select Interface' },
              ...interfaces.map((iface) => ({ value: iface, label: iface })),
            ]}
            value={selectedInterface}
            onChange={(e) => setSelectedInterface(e.target.value)}
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

        {selectedInterface && stats ? (
          <div className="space-y-4">
            {/* Interface Name */}
            <div className="text-center">
              <span className="text-lg font-semibold text-gray-900 dark:text-white">{stats.name}</span>
            </div>

            {/* Current Stats */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-success-50 dark:bg-success-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-success-600 mb-2">
                  <ArrowDown className="w-4 h-4" />
                  <span className="text-xs font-medium">RX (Download)</span>
                </div>
                <div className="text-lg font-bold text-success-700">
                  {formatBitsPerSecond(stats.rxBitsPerSecond)}
                </div>
                <div className="text-xs text-success-600 mt-1">{stats.rxPacketsPerSecond} pps</div>
              </div>

              <div className="bg-primary-50 dark:bg-primary-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-primary-600 mb-2">
                  <ArrowUp className="w-4 h-4" />
                  <span className="text-xs font-medium">TX (Upload)</span>
                </div>
                <div className="text-lg font-bold text-primary-700">
                  {formatBitsPerSecond(stats.txBitsPerSecond)}
                </div>
                <div className="text-xs text-primary-600 mt-1">{stats.txPacketsPerSecond} pps</div>
              </div>
            </div>

            {/* Chart */}
            {history.length > 1 && (
              <div className="h-48">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={history} margin={{ top: 4, right: 4, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="rxGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3} />
                        <stop offset="95%" stopColor="#22c55e" stopOpacity={0} />
                      </linearGradient>
                      <linearGradient id="txGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3} />
                        <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <XAxis dataKey="time" tick={{ fontSize: 10 }} interval="preserveStartEnd" />
                    <YAxis tickFormatter={yAxisFormatter} tick={{ fontSize: 10 }} width={48} />
                    <Tooltip
                      formatter={((v: number, name: string) => [
                        formatBitsPerSecond(v),
                        name === 'rx' ? 'RX' : 'TX',
                      ]) as any}
                      labelStyle={{ fontSize: 11 }}
                      contentStyle={{ fontSize: 11 }}
                    />
                    <Legend formatter={(v) => (v === 'rx' ? 'RX (Download)' : 'TX (Upload)')} wrapperStyle={{ fontSize: 11 }} />
                    <Area type="monotone" dataKey="rx" stroke="#22c55e" fill="url(#rxGrad)" strokeWidth={2} dot={false} isAnimationActive={false} />
                    <Area type="monotone" dataKey="tx" stroke="#6366f1" fill="url(#txGrad)" strokeWidth={2} dot={false} isAnimationActive={false} />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            )}

            {/* Total Throughput */}
            <div className="bg-gray-50 dark:bg-dark-700 p-4 rounded-xl">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600 dark:text-gray-400">Total Throughput</span>
                <span className="text-lg font-bold text-gray-900 dark:text-white">
                  {formatBitsPerSecond(stats.rxBitsPerSecond + stats.txBitsPerSecond)}
                </span>
              </div>
            </div>

            {/* Last Update */}
            <div className="text-center text-xs text-gray-400">
              Last update: {new Date(stats.timestamp).toLocaleTimeString()}
            </div>
          </div>
        ) : selectedInterface ? (
          <div className="text-center py-8 text-gray-400">
            <Network className="w-8 h-8 mx-auto mb-2 animate-pulse" />
            <p>Waiting for data...</p>
          </div>
        ) : (
          <div className="text-center py-8 text-gray-400">
            <p>Select an interface to monitor</p>
          </div>
        )}
      </Card.Body>
    </Card>
  )
}
