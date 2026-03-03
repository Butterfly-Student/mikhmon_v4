import { Ticket, Printer, Trash2, Copy, Check } from 'lucide-react'
import { useState } from 'react'
import toast from 'react-hot-toast'

import { Card, Button } from '../../../components/ui'
import type { Voucher } from '../../../types'

interface VoucherPreviewProps {
  vouchers: Voucher[]
  mode: string
  comment: string
  onClear: () => void
}

export function VoucherPreview({ vouchers, mode, comment, onClear }: VoucherPreviewProps) {
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null)

  const copyToClipboard = (text: string, index: number) => {
    navigator.clipboard.writeText(text)
    setCopiedIndex(index)
    toast.success('Copied to clipboard')
    setTimeout(() => setCopiedIndex(null), 2000)
  }

  return (
    <Card>
      <Card.Header>
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-gray-900 dark:text-white">Generated Vouchers</h3>
          {vouchers.length > 0 && (
            <div className="flex gap-2">
              <Button variant="ghost" size="sm" leftIcon={<Printer className="w-4 h-4" />}>
                Print
              </Button>
              <Button
                variant="ghost"
                size="sm"
                leftIcon={<Trash2 className="w-4 h-4 text-danger-500" />}
                onClick={onClear}
              >
                Clear
              </Button>
            </div>
          )}
        </div>
      </Card.Header>
      <Card.Body>
        {vouchers.length === 0 ? (
          <div className="text-center py-12">
            <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mx-auto mb-4">
              <Ticket className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-500">No vouchers generated yet</p>
            <p className="text-sm text-gray-400 mt-1">Fill the form and click Generate</p>
          </div>
        ) : (
          <div className="space-y-2 max-h-96 overflow-y-auto">
            {comment && (
              <div className="p-2 text-xs rounded bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300">
                Batch Comment: {comment}
              </div>
            )}
            {vouchers.map((voucher, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-3 bg-gray-50 dark:bg-dark-700 rounded-lg"
              >
                <div className="font-mono text-sm">
                  <span className="font-semibold text-gray-900 dark:text-white">
                    {voucher.username}
                  </span>
                  {mode === 'up' && voucher.password && (
                    <span className="text-gray-500 ml-2">/ {voucher.password}</span>
                  )}
                </div>
                <button
                  onClick={() =>
                    copyToClipboard(
                      mode === 'up'
                        ? `${voucher.username} / ${voucher.password}`
                        : voucher.username,
                      index
                    )
                  }
                  className="p-2 text-gray-400 hover:text-primary-500 transition-colors"
                >
                  {copiedIndex === index ? (
                    <Check className="w-4 h-4 text-success-500" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            ))}
          </div>
        )}
      </Card.Body>
    </Card>
  )
}
