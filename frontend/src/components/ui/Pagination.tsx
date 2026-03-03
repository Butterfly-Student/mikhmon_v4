import { ChevronsLeft, ChevronLeft, ChevronRight, ChevronsRight } from 'lucide-react'

interface PaginationProps {
  pageIndex: number
  pageCount: number
  pageSize: number
  totalRows: number
  onPageChange: (index: number) => void
  onPageSizeChange: (size: number) => void
}

export function Pagination({
  pageIndex,
  pageCount,
  pageSize,
  totalRows,
  onPageChange,
  onPageSizeChange,
}: PaginationProps) {
  const fromRow = totalRows === 0 ? 0 : pageIndex * pageSize + 1
  const toRow = Math.min((pageIndex + 1) * pageSize, totalRows)

  return (
    <div className="flex items-center justify-between gap-2 px-4 py-2.5 border-t border-gray-100 dark:border-dark-700 bg-gray-50/50 dark:bg-dark-800/50 flex-shrink-0">
      <div className="text-xs text-gray-500 dark:text-gray-400">
        {totalRows > 0 ? (
          <>
            <span className="font-medium text-gray-700 dark:text-gray-300">{fromRow}–{toRow}</span>
            {' '}of{' '}
            <span className="font-medium text-gray-700 dark:text-gray-300">{totalRows}</span>
          </>
        ) : (
          'No results'
        )}
      </div>

      <div className="flex items-center gap-1">
        <select
          value={pageSize}
          onChange={(e) => onPageSizeChange(Number(e.target.value))}
          className="mr-1 text-xs border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 rounded-lg px-1.5 py-1 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-1 focus:ring-primary-500"
        >
          {[10, 25, 50, 100].map((n) => (
            <option key={n} value={n}>{n}/pg</option>
          ))}
        </select>

        {[
          { Icon: ChevronsLeft, action: () => onPageChange(0), disabled: pageIndex === 0, label: 'First' },
          { Icon: ChevronLeft, action: () => onPageChange(pageIndex - 1), disabled: pageIndex === 0, label: 'Prev' },
          { Icon: ChevronRight, action: () => onPageChange(pageIndex + 1), disabled: pageIndex >= pageCount - 1, label: 'Next' },
          { Icon: ChevronsRight, action: () => onPageChange(pageCount - 1), disabled: pageIndex >= pageCount - 1, label: 'Last' },
        ].map(({ Icon, action, disabled, label }) => (
          <button
            key={label}
            onClick={action}
            disabled={disabled}
            title={label}
            className="p-1.5 rounded-lg text-gray-500 hover:bg-gray-100 dark:hover:bg-dark-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          >
            <Icon className="w-3.5 h-3.5" />
          </button>
        ))}

        <span className="text-xs text-gray-500 dark:text-gray-400 ml-0.5">
          <span className="font-medium text-gray-700 dark:text-gray-300">
            {pageCount === 0 ? 0 : pageIndex + 1}
          </span>
          /{pageCount}
        </span>
      </div>
    </div>
  )
}
