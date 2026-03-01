import {
    useReactTable,
    getCoreRowModel,
    getSortedRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    flexRender,
    type ColumnDef,
    type SortingState,
    type ColumnFiltersState,
} from '@tanstack/react-table'
import { useState } from 'react'
import {
    ChevronUp,
    ChevronDown,
    ChevronsUpDown,
    ChevronLeft,
    ChevronRight,
    ChevronsLeft,
    ChevronsRight,
    Search,
    Loader2,
    Inbox,
} from 'lucide-react'
import { clsx } from 'clsx'

interface DataTableProps<TData> {
    data: TData[]
    columns: ColumnDef<TData, any>[]
    isLoading?: boolean
    searchPlaceholder?: string
    showSearch?: boolean
    pageSize?: number
    emptyMessage?: string
    emptyIcon?: React.ReactNode
}

export function DataTable<TData>({
    data,
    columns,
    isLoading = false,
    searchPlaceholder = 'Search...',
    showSearch = true,
    pageSize = 25,
    emptyMessage = 'No data found',
    emptyIcon,
}: DataTableProps<TData>) {
    const [sorting, setSorting] = useState<SortingState>([])
    const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
    const [globalFilter, setGlobalFilter] = useState('')

    const table = useReactTable({
        data,
        columns,
        state: { sorting, columnFilters, globalFilter },
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        onGlobalFilterChange: setGlobalFilter,
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        initialState: { pagination: { pageSize } },
    })

    const totalRows = table.getFilteredRowModel().rows.length
    const { pageIndex, pageSize: currentPageSize } = table.getState().pagination
    const fromRow = totalRows === 0 ? 0 : pageIndex * currentPageSize + 1
    const toRow = Math.min((pageIndex + 1) * currentPageSize, totalRows)

    return (
        <div className="flex flex-col gap-2 sm:gap-3">
            {/* Search bar */}
            {showSearch && (
                <div className="relative max-w-full sm:max-w-sm">
                    <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-gray-400" />
                    <input
                        type="text"
                        value={globalFilter}
                        onChange={e => setGlobalFilter(e.target.value)}
                        placeholder={searchPlaceholder}
                        className="w-full pl-8 pr-3 py-1.5 sm:py-2 text-xs sm:text-sm rounded-lg sm:rounded-xl border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
                    />
                </div>
            )}

            {/* Table */}
            <div className="overflow-x-auto rounded-lg sm:rounded-xl border border-gray-200 dark:border-dark-700">
                <table className="w-full text-xs sm:text-sm">
                    <thead>
                        <tr className="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-dark-800 dark:to-dark-700 border-b border-gray-200 dark:border-dark-700">
                            {table.getHeaderGroups()[0]?.headers.map(header => (
                                <th
                                    key={header.id}
                                    colSpan={header.colSpan}
                                    className="px-2 py-2 sm:px-4 sm:py-3 text-left text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider whitespace-nowrap"
                                >
                                    {header.isPlaceholder ? null : (
                                        <div
                                            className={clsx(
                                                'flex items-center gap-1',
                                                header.column.getCanSort() && 'cursor-pointer select-none hover:text-gray-700 dark:hover:text-gray-200 transition-colors'
                                            )}
                                            onClick={header.column.getToggleSortingHandler()}
                                        >
                                            {flexRender(header.column.columnDef.header, header.getContext())}
                                            {header.column.getCanSort() && (
                                                <span className="text-gray-400">
                                                    {header.column.getIsSorted() === 'asc' ? (
                                                        <ChevronUp className="w-3 h-3 text-primary-500" />
                                                    ) : header.column.getIsSorted() === 'desc' ? (
                                                        <ChevronDown className="w-3 h-3 text-primary-500" />
                                                    ) : (
                                                        <ChevronsUpDown className="w-3 h-3 opacity-40" />
                                                    )}
                                                </span>
                                            )}
                                        </div>
                                    )}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100 dark:divide-dark-700 bg-white dark:bg-dark-800">
                        {isLoading ? (
                            <tr>
                                <td colSpan={columns.length} className="py-10 sm:py-16 text-center">
                                    <div className="flex flex-col items-center gap-2 sm:gap-3 text-gray-400">
                                        <Loader2 className="w-6 h-6 sm:w-8 sm:h-8 animate-spin text-primary-500" />
                                        <span className="text-xs sm:text-sm">Loading data...</span>
                                    </div>
                                </td>
                            </tr>
                        ) : table.getRowModel().rows.length === 0 ? (
                            <tr>
                                <td colSpan={columns.length} className="py-10 sm:py-16 text-center">
                                    <div className="flex flex-col items-center gap-2 sm:gap-3 text-gray-400">
                                        {emptyIcon || <Inbox className="w-8 h-8 sm:w-10 sm:h-10 text-gray-300 dark:text-gray-600" />}
                                        <span className="text-xs sm:text-sm">{emptyMessage}</span>
                                    </div>
                                </td>
                            </tr>
                        ) : (
                            table.getRowModel().rows.map(row => (
                                <tr
                                    key={row.id}
                                    className="hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors"
                                >
                                    {row.getVisibleCells().map(cell => (
                                        <td key={cell.id} className="px-2 py-2 sm:px-4 sm:py-3">
                                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                        </td>
                                    ))}
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {/* Pagination */}
            <div className="flex flex-col xs:flex-row items-center justify-between gap-2 px-0.5">
                <div className="text-[10px] sm:text-xs text-gray-500 dark:text-gray-400">
                    {totalRows > 0 ? (
                        <><span className="font-medium text-gray-700 dark:text-gray-300">{fromRow}–{toRow}</span> of <span className="font-medium text-gray-700 dark:text-gray-300">{totalRows}</span></>
                    ) : 'No results'}
                </div>

                <div className="flex items-center gap-1">
                    <select
                        value={currentPageSize}
                        onChange={e => table.setPageSize(Number(e.target.value))}
                        className="mr-1 text-[10px] sm:text-xs border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-800 rounded-lg px-1.5 py-1 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    >
                        {[10, 25, 50, 100].map(n => (
                            <option key={n} value={n}>{n}/pg</option>
                        ))}
                    </select>

                    {[
                        { icon: ChevronsLeft, action: () => table.firstPage(), disabled: !table.getCanPreviousPage(), label: 'First' },
                        { icon: ChevronLeft, action: () => table.previousPage(), disabled: !table.getCanPreviousPage(), label: 'Prev' },
                        { icon: ChevronRight, action: () => table.nextPage(), disabled: !table.getCanNextPage(), label: 'Next' },
                        { icon: ChevronsRight, action: () => table.lastPage(), disabled: !table.getCanNextPage(), label: 'Last' },
                    ].map(({ icon: Icon, action, disabled, label }) => (
                        <button
                            key={label}
                            onClick={action}
                            disabled={disabled}
                            className="p-1 sm:p-1.5 rounded-lg text-gray-500 hover:bg-gray-100 dark:hover:bg-dark-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                            title={label}
                        >
                            <Icon className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
                        </button>
                    ))}

                    <span className="text-[10px] sm:text-xs text-gray-500 dark:text-gray-400 ml-0.5">
                        <span className="font-medium text-gray-700 dark:text-gray-300">{table.getState().pagination.pageIndex + 1}</span>
                        /{table.getPageCount()}
                    </span>
                </div>
            </div>
        </div>
    )
}
