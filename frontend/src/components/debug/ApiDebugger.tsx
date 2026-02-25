/* eslint-disable react-refresh/only-export-components */
import { useState, useEffect } from 'react'
import { Copy, Check, Activity } from 'lucide-react'
import { Modal } from '../ui/Modal'
import { Button } from '../ui/Button'
import type { ApiRequest } from './types'
import { initialRequests } from './helpers'

// Hook to track API requests
export function useApiDebugger() {
  const [requests, setRequests] = useState<ApiRequest[]>([...initialRequests])
  const [selectedRequest, setSelectedRequest] = useState<ApiRequest | null>(null)
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    // Set up listener for API requests from window events
    const handleRequest = (e: CustomEvent) => {
      const request = e.detail as ApiRequest
      setRequests(prev => [...prev, request])
    }

    const handleResponse = (e: CustomEvent) => {
      const response = e.detail as ApiRequest
      setRequests(prev => prev.map(r => r.id === response.id ? response : r))
    }

    window.addEventListener('api-request', handleRequest as EventListener)
    window.addEventListener('api-response', handleResponse as EventListener)

    return () => {
      window.removeEventListener('api-request', handleRequest as EventListener)
      window.removeEventListener('api-response', handleResponse as EventListener)
    }
  }, [])

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const clearRequests = () => {
    setRequests([])
  }

  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms.toFixed(0)}ms`
    return `${(ms / 1000).toFixed(2)}s`
  }

  const getStatusColor = (status?: number) => {
    if (!status) return 'text-gray-500'
    if (status >= 200 && status < 300) return 'text-success-600'
    if (status >= 300 && status < 400) return 'text-warning-600'
    if (status >= 400 && status < 500) return 'text-danger-600'
    return 'text-gray-500'
  }

  return {
    apiRequests: requests,
    selectedRequest,
    setSelectedRequest,
    copyToClipboard,
    clearRequests,
    formatDuration,
    getStatusColor,
    copied,
  }
}

/* eslint-enable react-refresh/only-export-components */

// Component to display API debugger
export function ApiDebugger() {
  const {
    apiRequests,
    selectedRequest,
    setSelectedRequest,
    copyToClipboard,
    clearRequests,
    formatDuration,
    getStatusColor,
    copied,
  } = useApiDebugger()
  const [isOpen, setIsOpen] = useState(false)

  // Store in window for external access
  useEffect(() => {
    (window as { apiDebugger?: { toggle: () => void; open: () => void; close: () => void } }).apiDebugger = {
      toggle: () => setIsOpen(!isOpen),
      open: () => setIsOpen(true),
      close: () => setIsOpen(false),
    }

    // Log instructions
    console.log('[API Debugger] Available via window.apiDebugger')
    console.log('[API Debugger] Use window.apiDebugger.toggle() to open/close')
  }, [isOpen])

  return (
    <>
      {/* Floating button */}
      {!isOpen && (
        <button
          onClick={() => setIsOpen(true)}
          className="fixed bottom-4 right-4 z-50 p-3 bg-primary-600 hover:bg-primary-700 text-white rounded-full shadow-lg transition-all"
          title="Open API Debugger"
        >
          <Activity className="w-5 h-5" />
        </button>
      )}

      {/* Modal */}
      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="API Request Debugger"
        size="full"
        footer={
          <div className="flex justify-end gap-2">
            <Button variant="ghost" onClick={clearRequests}>
              Clear All
            </Button>
            <Button onClick={() => setIsOpen(false)}>
              Close
            </Button>
          </div>
        }
      >
        <div className="flex gap-4 h-[600px]">
          {/* Request list */}
          <div className="w-1/3 overflow-y-auto border-r border-gray-200 dark:border-dark-700">
            <div className="space-y-2">
              {apiRequests.length === 0 && (
                <div className="p-4 text-center text-gray-500">
                  No API requests yet
                </div>
              )}
              {apiRequests.map((req) => (
                <button
                  key={req.id}
                  onClick={() => setSelectedRequest(req)}
                  className={`w-full p-3 text-left rounded-lg transition-colors ${
                    selectedRequest?.id === req.id
                      ? 'bg-primary-100 dark:bg-primary-900/30 border-2 border-primary-500'
                      : 'bg-gray-50 dark:bg-dark-800 hover:bg-gray-100 dark:hover:bg-dark-700'
                  }`}
                >
                  <div className="flex items-center justify-between mb-1">
                    <span className={`text-xs font-medium ${
                      req.method === 'GET' ? 'text-success-600' :
                      req.method === 'POST' ? 'text-primary-600' :
                      req.method === 'PUT' ? 'text-warning-600' :
                      req.method === 'DELETE' ? 'text-danger-600' :
                      'text-gray-600'
                    }`}>
                      {req.method}
                    </span>
                    <span className={`text-xs font-medium ${getStatusColor(req.status)}`}>
                      {req.status || '...'}
                    </span>
                  </div>
                  <div className="text-xs text-gray-600 dark:text-gray-300 truncate">
                    {req.url}
                  </div>
                  <div className="text-xs text-gray-400 mt-1">
                    {req.duration ? formatDuration(req.duration) : '...'}
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Request details */}
          <div className="w-2/3 overflow-y-auto">
            {!selectedRequest && (
              <div className="p-8 text-center text-gray-500">
                Select a request to view details
              </div>
            )}
            {selectedRequest && (
              <div className="space-y-4 p-4">
                {/* Request Info */}
                <div>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">
                    Request Information
                  </h3>
                  <div className="space-y-2 text-sm">
                    <div>
                      <span className="text-gray-500">Method:</span>{' '}
                      <span className="font-mono">{selectedRequest.method}</span>
                    </div>
                    <div>
                      <span className="text-gray-500">URL:</span>{' '}
                      <span className="font-mono break-all">{selectedRequest.url}</span>
                    </div>
                    <div>
                      <span className="text-gray-500">Time:</span>{' '}
                      {selectedRequest.requestTime.toLocaleString()}
                    </div>
                    {selectedRequest.duration !== undefined && (
                      <div>
                        <span className="text-gray-500">Duration:</span>{' '}
                        <span className={selectedRequest.duration > 1000 ? 'text-warning-600' : 'text-success-600'}>
                          {formatDuration(selectedRequest.duration)}
                        </span>
                      </div>
                    )}
                    {selectedRequest.status && (
                      <div>
                        <span className="text-gray-500">Status:</span>{' '}
                        <span className={getStatusColor(selectedRequest.status)}>
                          {selectedRequest.status} {selectedRequest.statusText}
                        </span>
                      </div>
                    )}
                  </div>
                </div>

                {/* Error */}
                {selectedRequest.error && (
                  <div className="p-3 bg-danger-50 dark:bg-danger-900/20 border border-danger-200 dark:border-danger-800 rounded-lg">
                    <h3 className="text-sm font-semibold text-danger-900 dark:text-danger-100 mb-1">
                      Error
                    </h3>
                    <pre className="text-xs text-danger-700 dark:text-danger-300 whitespace-pre-wrap">
                      {selectedRequest.error}
                    </pre>
                  </div>
                )}

                {/* Request Headers */}
                {selectedRequest.requestHeaders && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">
                      Request Headers
                    </h3>
                    <pre className="text-xs bg-gray-100 dark:bg-dark-800 p-3 rounded-lg overflow-x-auto">
                      {JSON.stringify(selectedRequest.requestHeaders, null, 2)}
                    </pre>
                  </div>
                )}

                {/* Request Body */}
                {selectedRequest.requestData && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">
                      Request Body
                    </h3>
                    <pre className="text-xs bg-gray-100 dark:bg-dark-800 p-3 rounded-lg overflow-x-auto max-h-48 overflow-y-auto">
                      {JSON.stringify(selectedRequest.requestData, null, 2)}
                    </pre>
                  </div>
                )}

                {/* Response Headers */}
                {selectedRequest.responseHeaders && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">
                      Response Headers
                    </h3>
                    <pre className="text-xs bg-gray-100 dark:bg-dark-800 p-3 rounded-lg overflow-x-auto">
                      {JSON.stringify(selectedRequest.responseHeaders, null, 2)}
                    </pre>
                  </div>
                )}

                {/* Response Body */}
                {selectedRequest.responseData && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">
                      Response Body
                    </h3>
                    <pre className="text-xs bg-gray-100 dark:bg-dark-800 p-3 rounded-lg overflow-x-auto max-h-64 overflow-y-auto">
                      {JSON.stringify(selectedRequest.responseData, null, 2)}
                    </pre>
                  </div>
                )}

                {/* Copy JSON button */}
                <div className="flex gap-2 pt-4 border-t border-gray-200 dark:border-dark-700">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => copyToClipboard(JSON.stringify(selectedRequest, null, 2))}
                  >
                    {copied ? <Check className="w-4 h-4 mr-1" /> : <Copy className="w-4 h-4 mr-1" />}
                    {copied ? 'Copied!' : 'Copy JSON'}
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      </Modal>
    </>
  )
}

