import { useState, useEffect, useRef, useCallback } from 'react'
import { Network, RefreshCw, ArrowDown, ArrowUp } from 'lucide-react'
import { Card, Button, Select } from '../ui'
import { useRouterStore } from '../../stores/routerStore'

interface InterfaceStats {
  name: string
  txBitsPerSecond: number
  rxBitsPerSecond: number
  txPacketsPerSecond: number
  rxPacketsPerSecond: number
  timestamp: number
}

export function InterfaceMonitor() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'
  
  const [interfaces, setInterfaces] = useState<string[]>([])
  const [selectedInterface, setSelectedInterface] = useState('')
  const [stats, setStats] = useState<InterfaceStats | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [isMonitoring, setIsMonitoring] = useState(false)
  const [error, setError] = useState('')
  
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // Fetch available interfaces
  useEffect(() => {
    // TODO: Fetch interfaces from API
    // For now, using common interface names - replace with actual API call
    setInterfaces(['ether1', 'ether2', 'ether3', 'wlan1', 'wlan2', 'bridge1'])
  }, [routerId])

  const connect = useCallback(() => {
    if (!selectedInterface) return
    
    const wsKey = localStorage.getItem('ws_key') || 'mikhmon-ws-internal-key'
    const wsUrl = `ws://localhost:8080/api/v1/ws/monitor/interface/${routerId}?key=${wsKey}`
    
    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      setError('')
      // Start monitoring
      ws.send(JSON.stringify({ 
        action: 'start', 
        name: selectedInterface,
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
        console.log('Interface monitor status:', data.status)
        return
      }
      setStats(data)
    }

    ws.onerror = (err) => {
      console.error('Interface monitor error:', err)
      setError('Connection error')
    }

    ws.onclose = () => {
      setIsConnected(false)
      setIsMonitoring(false)
      // Auto reconnect after 3s
      reconnectTimeoutRef.current = setTimeout(() => {
        if (selectedInterface) connect()
      }, 3000)
    }
  }, [routerId, selectedInterface])

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
    if (selectedInterface) {
      disconnect()
      connect()
    }
  }, [selectedInterface, connect, disconnect])

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
              ...interfaces.map(iface => ({ value: iface, label: iface }))
            ]}
            value={selectedInterface}
            onChange={(e) => setSelectedInterface(e.target.value)}
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

        {selectedInterface && stats ? (
          <div className="space-y-4">
            {/* Interface Name */}
            <div className="text-center">
              <span className="text-lg font-semibold text-gray-900 dark:text-white">
                {stats.name}
              </span>
            </div>

            {/* Traffic Stats */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-success-50 dark:bg-success-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-success-600 mb-2">
                  <ArrowDown className="w-4 h-4" />
                  <span className="text-xs font-medium">RX (Download)</span>
                </div>
                <div className="text-lg font-bold text-success-700">
                  {formatBitsPerSecond(stats.rxBitsPerSecond)}
                </div>
                <div className="text-xs text-success-600 mt-1">
                  {stats.rxPacketsPerSecond} pps
                </div>
              </div>

              <div className="bg-primary-50 dark:bg-primary-900/20 p-4 rounded-xl">
                <div className="flex items-center gap-2 text-primary-600 mb-2">
                  <ArrowUp className="w-4 h-4" />
                  <span className="text-xs font-medium">TX (Upload)</span>
                </div>
                <div className="text-lg font-bold text-primary-700">
                  {formatBitsPerSecond(stats.txBitsPerSecond)}
                </div>
                <div className="text-xs text-primary-600 mt-1">
                  {stats.txPacketsPerSecond} pps
                </div>
              </div>
            </div>

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
