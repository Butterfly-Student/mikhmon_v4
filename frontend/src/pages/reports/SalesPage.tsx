import { useState, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Download, DollarSign, ShoppingCart, TrendingUp } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'

import { Card, Button, Badge, DataTable } from '../../components/ui'
import { reportsApi } from '../../api/reports'
import { useRouterStore } from '../../stores/routerStore'

const months = [
  { value: '1', label: 'January' }, { value: '2', label: 'February' },
  { value: '3', label: 'March' }, { value: '4', label: 'April' },
  { value: '5', label: 'May' }, { value: '6', label: 'June' },
  { value: '7', label: 'July' }, { value: '8', label: 'August' },
  { value: '9', label: 'September' }, { value: '10', label: 'October' },
  { value: '11', label: 'November' }, { value: '12', label: 'December' },
]

type SaleReport = {
  date: string
  time: string
  username: string
  price: number
  ipAddress: string
  macAddress: string
  validity: string
  profile: string
}

const formatCurrency = (value: number) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value)

export function SalesPage() {
  const selectedRouter = useRouterStore((state) => state.selectedRouter)
  const routerId = String(selectedRouter?.id ?? '1')

  const currentMonth = new Date().getMonth() + 1
  const currentYear = new Date().getFullYear()

  const [selectedMonth, setSelectedMonth] = useState(currentMonth.toString())
  const [selectedYear, setSelectedYear] = useState(currentYear.toString())

  const { data: reports, isLoading } = useQuery({
    queryKey: ['reports', routerId, selectedMonth, selectedYear],
    queryFn: () => reportsApi.getSalesReport(routerId, { month: selectedMonth, year: selectedYear }),
  })

  const { data: summary } = useQuery({
    queryKey: ['summary', routerId, selectedMonth, selectedYear],
    queryFn: () => reportsApi.getSummary(routerId, selectedMonth, selectedYear),
  })

  const handleExport = async () => {
    try {
      const csv = await reportsApi.exportToCSV(routerId, { month: selectedMonth, year: selectedYear })
      const blob = new Blob([csv], { type: 'text/csv' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `sales-report-${selectedMonth}-${selectedYear}.csv`
      a.click()
    } catch (err) { console.error('Export failed:', err) }
  }

  const columns = useMemo<ColumnDef<SaleReport, any>[]>(() => [
    {
      accessorKey: 'date',
      header: 'Date',
      cell: ({ row }) => (
        <div className="text-xs">
          <p className="font-medium text-gray-800 dark:text-gray-200">{row.original.date}</p>
          <p className="text-gray-400">{row.original.time}</p>
        </div>
      ),
    },
    {
      accessorKey: 'username',
      header: 'User',
      cell: ({ getValue }) => (
        <span className="font-medium text-gray-900 dark:text-white text-sm">{getValue()}</span>
      ),
    },
    {
      accessorKey: 'price',
      header: 'Price',
      cell: ({ getValue }) => (
        <span className="font-semibold text-success-600 dark:text-success-400 text-sm">
          {formatCurrency(getValue())}
        </span>
      ),
    },
    {
      accessorKey: 'ipAddress',
      header: 'IP Address',
      cell: ({ getValue }) => (
        <span className="font-mono text-xs text-gray-600 dark:text-gray-300">{getValue()}</span>
      ),
    },
    {
      accessorKey: 'macAddress',
      header: 'MAC Address',
      cell: ({ getValue }) => (
        <span className="font-mono text-xs text-gray-500 dark:text-gray-400">{getValue()}</span>
      ),
    },
    {
      accessorKey: 'validity',
      header: 'Validity',
      cell: ({ getValue }) => <Badge variant="info">{getValue()}</Badge>,
    },
    {
      accessorKey: 'profile',
      header: 'Profile',
      cell: ({ getValue }) => <Badge variant="primary">{getValue()}</Badge>,
    },
  ], [])

  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="space-y-4">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">Sales Report</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">Track your hotspot sales and revenue</p>
        </div>
        <Button variant="primary" size="sm" leftIcon={<Download className="w-4 h-4" />} onClick={handleExport}>
          Export CSV
        </Button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-3 gap-3">
        {[
          { label: 'Total Sales', value: formatCurrency(summary?.totalSales || 0), icon: DollarSign, color: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600' },
          { label: 'Transactions', value: String(summary?.totalTransactions || 0), icon: ShoppingCart, color: 'bg-primary-100 dark:bg-primary-900/30 text-primary-600' },
          { label: 'Avg. Ticket', value: formatCurrency(summary?.averageTicket || 0), icon: TrendingUp, color: 'bg-cyan-100 dark:bg-cyan-900/30 text-cyan-600' },
        ].map(({ label, value, icon: Icon, color }) => (
          <Card key={label}>
            <Card.Body className="flex items-center gap-3 py-4">
              <div className={`w-10 h-10 rounded-xl flex items-center justify-center shrink-0 ${color}`}>
                <Icon className="w-5 h-5" />
              </div>
              <div className="min-w-0">
                <p className="text-xs text-gray-500 dark:text-gray-400 truncate">{label}</p>
                <p className="text-base font-bold text-gray-900 dark:text-white truncate">{value}</p>
              </div>
            </Card.Body>
          </Card>
        ))}
      </div>

      {/* Filters + Table */}
      <Card>
        <Card.Header>
          <div className="flex flex-wrap items-center gap-2">
            <select
              value={selectedMonth}
              onChange={(e) => setSelectedMonth(e.target.value)}
              className="text-sm border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 rounded-xl px-3 py-1.5 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              {months.map((m) => <option key={m.value} value={m.value}>{m.label}</option>)}
            </select>
            <select
              value={selectedYear}
              onChange={(e) => setSelectedYear(e.target.value)}
              className="text-sm border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 rounded-xl px-3 py-1.5 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              {[2022, 2023, 2024, 2025, 2026].map((y) => <option key={y} value={y}>{y}</option>)}
            </select>
          </div>
        </Card.Header>
        <Card.Body>
          <DataTable
            data={(reports as SaleReport[]) || []}
            columns={columns}
            isLoading={isLoading}
            searchPlaceholder="Search user, IP, profile..."
            emptyMessage="No sales data found"
          />
        </Card.Body>
      </Card>
    </motion.div>
  )
}

