import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import {
  Download,
  DollarSign,
  ShoppingCart,
  TrendingUp,
} from 'lucide-react'

import { Card, Button, Badge } from '../../components/ui'
import { reportsApi } from '../../api/reports'
import { useRouterStore } from '../../stores/routerStore'

const months = [
  { value: '1', label: 'January' },
  { value: '2', label: 'February' },
  { value: '3', label: 'March' },
  { value: '4', label: 'April' },
  { value: '5', label: 'May' },
  { value: '6', label: 'June' },
  { value: '7', label: 'July' },
  { value: '8', label: 'August' },
  { value: '9', label: 'September' },
  { value: '10', label: 'October' },
  { value: '11', label: 'November' },
  { value: '12', label: 'December' },
]

export function SalesPage() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = selectedRouter?.id || '1'

  const currentMonth = new Date().getMonth() + 1
  const currentYear = new Date().getFullYear()

  const [selectedMonth, setSelectedMonth] = useState(currentMonth.toString())
  const [selectedYear, setSelectedYear] = useState(currentYear.toString())

  const { data: reports, isLoading } = useQuery({
    queryKey: ['reports', routerId, selectedMonth, selectedYear],
    queryFn: () =>
      reportsApi.getSalesReport(routerId, {
        month: selectedMonth,
        year: selectedYear,
      }),
  })

  const { data: summary } = useQuery({
    queryKey: ['summary', routerId, selectedMonth, selectedYear],
    queryFn: () => reportsApi.getSummary(routerId, selectedMonth, selectedYear),
  })

  const handleExport = async () => {
    try {
      const csv = await reportsApi.exportToCSV(routerId, {
        month: selectedMonth,
        year: selectedYear,
      })
      const blob = new Blob([csv], { type: 'text/csv' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `sales-report-${selectedMonth}-${selectedYear}.csv`
      a.click()
    } catch (error) {
      console.error('Export failed:', error)
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(value)
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Sales Report</h1>
          <p className="text-gray-500 dark:text-gray-400">
            Track your hotspot sales and revenue
          </p>
        </div>
        <Button
          variant="gradient"
          leftIcon={<Download className="w-4 h-4" />}
          onClick={handleExport}
        >
          Export CSV
        </Button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <Card.Body className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-emerald-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Total Sales</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {formatCurrency(summary?.totalSales || 0)}
              </p>
            </div>
          </Card.Body>
        </Card>

        <Card>
          <Card.Body className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center">
              <ShoppingCart className="w-6 h-6 text-primary-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Transactions</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {summary?.totalTransactions || 0}
              </p>
            </div>
          </Card.Body>
        </Card>

        <Card>
          <Card.Body className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-cyan-100 dark:bg-cyan-900/30 flex items-center justify-center">
              <TrendingUp className="w-6 h-6 text-cyan-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Average Ticket</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {formatCurrency(summary?.averageTicket || 0)}
              </p>
            </div>
          </Card.Body>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <Card.Body className="flex flex-col sm:flex-row gap-4">
          <div className="flex gap-4">
            <select
              value={selectedMonth}
              onChange={(e) => setSelectedMonth(e.target.value)}
              className="input"
            >
              {months.map((m) => (
                <option key={m.value} value={m.value}>
                  {m.label}
                </option>
              ))}
            </select>
            <select
              value={selectedYear}
              onChange={(e) => setSelectedYear(e.target.value)}
              className="input"
            >
              {[2022, 2023, 2024, 2025].map((y) => (
                <option key={y} value={y}>
                  {y}
                </option>
              ))}
            </select>
          </div>
        </Card.Body>
      </Card>

      {/* Table */}
      <Card>
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Time</th>
                <th>User</th>
                <th>Price</th>
                <th>IP Address</th>
                <th>MAC Address</th>
                <th>Validity</th>
                <th>Profile</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr>
                  <td colSpan={8} className="text-center py-8">
                    <div className="animate-spin w-6 h-6 border-2 border-primary-500 border-t-transparent rounded-full mx-auto" />
                  </td>
                </tr>
              ) : reports?.length === 0 ? (
                <tr>
                  <td colSpan={8} className="text-center py-8 text-gray-500">
                    No sales data found
                  </td>
                </tr>
              ) : (
                reports?.map((report, index) => (
                  <tr key={index}>
                    <td>{report.date}</td>
                    <td>{report.time}</td>
                    <td className="font-medium">{report.username}</td>
                    <td className="text-emerald-600 font-medium">
                      {formatCurrency(report.price)}
                    </td>
                    <td className="font-mono text-sm">{report.ipAddress}</td>
                    <td className="font-mono text-sm">{report.macAddress}</td>
                    <td>
                      <Badge variant="info">{report.validity}</Badge>
                    </td>
                    <td>
                      <Badge variant="primary">{report.profile}</Badge>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </motion.div>
  )
}
