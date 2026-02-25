import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import {
  Printer,
  Search,
  FileText,
  Grid,
  List,
  QrCode,
} from 'lucide-react'

import { Card, Button, Input, Select, Badge } from '../../components/ui'
import { hotspotApi } from '../../api/hotspot'
import { useRouterStore } from '../../stores/routerStore'

const templates = [
  { value: 'default', label: 'Default (6 per page)' },
  { value: 'small', label: 'Small (10 per page)' },
  { value: 'thermal-58', label: 'Thermal 58mm' },
  { value: 'thermal-80', label: 'Thermal 80mm' },
]

export function PrintPage() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  const [searchComment, setSearchComment] = useState('')
  const [selectedTemplate, setSelectedTemplate] = useState('default')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')

  const { data: vouchers } = useQuery({
    queryKey: ['vouchers', routerId, searchComment],
    queryFn: () =>
      searchComment
        ? hotspotApi.getUsers(routerId, { comment: searchComment })
        : Promise.resolve([]),
    enabled: !!searchComment,
  })

  const handlePrint = () => {
    window.print()
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Print Vouchers</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Print generated vouchers with templates
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="ghost"
            leftIcon={<Grid className="w-4 h-4" />}
            onClick={() => setViewMode('grid')}
            className={viewMode === 'grid' ? 'bg-gray-100 dark:bg-dark-700' : ''}
          >
            Grid
          </Button>
          <Button
            variant="ghost"
            leftIcon={<List className="w-4 h-4" />}
            onClick={() => setViewMode('list')}
            className={viewMode === 'list' ? 'bg-gray-100 dark:bg-dark-700' : ''}
          >
            List
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <Card.Body className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <Input
              placeholder="Search by comment..."
              leftIcon={<Search className="w-4 h-4" />}
              value={searchComment}
              onChange={(e) => setSearchComment(e.target.value)}
            />
          </div>
          <div className="w-full sm:w-64">
            <Select
              options={templates}
              value={selectedTemplate}
              onChange={(e) => setSelectedTemplate(e.target.value)}
            />
          </div>
          <Button
            variant="gradient"
            leftIcon={<Printer className="w-4 h-4" />}
            onClick={handlePrint}
            disabled={!vouchers?.length}
          >
            Print
          </Button>
        </Card.Body>
      </Card>

      {/* Preview */}
      <Card>
        <Card.Header>
          <div className="flex items-center justify-between">
            <h3 className="font-semibold text-gray-900 dark:text-white">Preview</h3>
            <Badge variant="primary">{vouchers?.length || 0} vouchers</Badge>
          </div>
        </Card.Header>
        <Card.Body>
          {!vouchers?.length ? (
            <div className="text-center py-12">
              <div className="w-16 h-16 rounded-full bg-gray-100 dark:bg-dark-700 flex items-center justify-center mx-auto mb-4">
                <FileText className="w-8 h-8 text-gray-400" />
              </div>
              <p className="text-gray-500">Search for vouchers to print</p>
            </div>
          ) : viewMode === 'grid' ? (
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              {vouchers.map((voucher) => (
                <div
                  key={voucher.id}
                  className="p-4 border-2 border-dashed border-gray-300 dark:border-dark-600 rounded-xl text-center"
                >
                  <div className="text-xs text-gray-500 mb-2">{voucher.comment}</div>
                  <div className="font-mono text-lg font-bold text-gray-900 dark:text-white">
                    {voucher.name}
                  </div>
                  {voucher.password && (
                    <div className="font-mono text-sm text-gray-600 dark:text-gray-400 mt-1">
                      Pass: {voucher.password}
                    </div>
                  )}
                  <div className="mt-2">
                    <Badge variant="primary" size="sm">
                      {typeof voucher.profile === 'string' ? voucher.profile : voucher.profile?.name || '-'}
                    </Badge>
                  </div>
                  <div className="mt-3 flex justify-center">
                    <div className="w-16 h-16 bg-gray-100 dark:bg-dark-700 rounded flex items-center justify-center">
                      <QrCode className="w-8 h-8 text-gray-400" />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="space-y-2">
              {vouchers.map((voucher) => (
                <div
                  key={voucher.id}
                  className="flex items-center justify-between p-3 bg-gray-50 dark:bg-dark-700 rounded-lg"
                >
                  <div className="flex items-center gap-4">
                    <div className="w-10 h-10 bg-gray-100 dark:bg-dark-600 rounded flex items-center justify-center">
                      <QrCode className="w-5 h-5 text-gray-400" />
                    </div>
                    <div>
                      <div className="font-mono font-semibold">{voucher.name}</div>
                      <div className="text-sm text-gray-500">
                        {typeof voucher.profile === 'string' ? voucher.profile : voucher.profile?.name || '-'} • {voucher.comment}
                      </div>
                    </div>
                  </div>
                  {voucher.password && (
                    <div className="font-mono text-sm text-gray-600 dark:text-gray-400">
                      Pass: {voucher.password}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </Card.Body>
      </Card>
    </motion.div>
  )
}
