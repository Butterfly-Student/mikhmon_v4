import { useState } from 'react'
import { RefreshCw, Activity, Ticket, Printer, UserPlus } from 'lucide-react'

import { Button } from '../../../components/ui'
import { toggleApiDebug } from '../../../api/axios'
import type { HotspotUser } from '../../../types'

interface UserActionsBarProps {
  users: HotspotUser[]
  onRefresh: () => void
  onAddUser: () => void
}

export function UserActionsBar({ users, onRefresh, onAddUser }: UserActionsBarProps) {
  const [showPrintModal, setShowPrintModal] = useState(false)

  const handlePrint = (size: 'small' | 'default') => {
    localStorage.setItem('printUsers', JSON.stringify(users || []))
    localStorage.setItem('printSize', size)
    window.open('/vouchers/print', '_blank')
  }

  return (
    <div className="flex items-center gap-2">
      <button
        onClick={onRefresh}
        className="p-2 rounded-xl bg-gray-100 dark:bg-dark-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-dark-600 transition-colors"
        title="Refresh"
      >
        <RefreshCw className="w-4 h-4" />
      </button>
      <Button variant="ghost" size="sm" onClick={() => toggleApiDebug()} title="Toggle Debug">
        <Activity className="w-4 h-4" />
      </Button>
      <Button variant="secondary" size="sm" onClick={() => window.location.href = '/vouchers/generate'}>
        <Ticket className="w-4 h-4 mr-1" /> Generate
      </Button>
      <div className="relative">
        <button
          onClick={() => setShowPrintModal(!showPrintModal)}
          disabled={!users?.length}
          className="px-3 py-1.5 rounded-xl text-sm font-medium bg-success-50 dark:bg-success-900/20 text-success-600 dark:text-success-400 hover:bg-success-100 dark:hover:bg-success-900/30 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center gap-1.5"
        >
          <Printer className="w-4 h-4" /> Print
        </button>
        {showPrintModal && users?.length && (
          <>
            <div className="fixed inset-0 z-10" onClick={() => setShowPrintModal(false)} />
            <div className="absolute right-0 mt-2 w-44 bg-white dark:bg-dark-800 rounded-xl shadow-lg border border-gray-200 dark:border-dark-700 z-20 py-1 overflow-hidden">
              {(['small', 'default'] as const).map((s) => (
                <button
                  key={s}
                  onClick={() => { handlePrint(s); setShowPrintModal(false) }}
                  className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-dark-700 flex items-center gap-2"
                >
                  <Printer className="w-4 h-4" /> Print {s === 'small' ? 'Small' : 'Default'}
                </button>
              ))}
            </div>
          </>
        )}
      </div>
      <Button onClick={onAddUser} leftIcon={<UserPlus className="w-4 h-4" />}>Add</Button>
    </div>
  )
}
