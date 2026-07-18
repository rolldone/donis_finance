import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { getMembers, getAdminTransactions, getAdminTransactionSummaryWithBreakdown, getAdminMonthlySeries } from '../../api'
import Skeleton from '../../components/Skeleton'
import ExpenseChart from '../../components/ExpenseChart'
import IncomeVsExpenseChart from '../../components/IncomeVsExpenseChart'

function formatCurrency(amount: number | string) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(Math.abs(Number(amount)))
}

export default function AdminDashboard() {
  const { t } = useTranslation()
  const now = new Date()
  const month = now.getMonth() + 1
  const year = now.getFullYear()

  const [members, setMembers] = useState<any[]>([])
  const [transactions, setTransactions] = useState<any[]>([])
  const [categoryBreakdown, setCategoryBreakdown] = useState<any[]>([])
  const [monthlySeries, setMonthlySeries] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function fetchData() {
      try {
        const [memRes, txRes, summRes, seriesRes] = await Promise.all([
          getMembers(),
          getAdminTransactions({ month, year, limit: 100 }),
          getAdminTransactionSummaryWithBreakdown({ month, year }),
          getAdminMonthlySeries({ months: 6 }),
        ])
        setMembers(memRes.members || [])
        setTransactions(txRes.transactions || [])
        setCategoryBreakdown(summRes.category_breakdown || [])
        setMonthlySeries(seriesRes.series || [])
      } catch (err: any) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [month, year])

  const totalTx = transactions.length
  const totalIncome = transactions
    .filter((tx) => tx.type === 'income')
    .reduce((sum, tx) => sum + (parseInt(tx.amount) || 0), 0)
  const totalExpense = transactions
    .filter((tx) => tx.type === 'expense')
    .reduce((sum, tx) => sum + (parseInt(tx.amount) || 0), 0)

  if (loading) {
    return (
      <div className="space-y-6">
        {/* Header */}
        <div>
          <Skeleton.Base className="h-7 w-48 mb-2" />
          <Skeleton.Base className="h-4 w-64" />
        </div>
        {/* Stat cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <Skeleton.Card />
          <Skeleton.Card />
          <Skeleton.Card />
          <Skeleton.Card />
        </div>
        {/* Charts row */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Skeleton.ChartCard />
          <Skeleton.ChartCard />
        </div>
        {/* Recent transactions */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <Skeleton.Base className="h-4 w-36 mb-4" />
          {Array.from({ length: 5 }).map((_, i) => <Skeleton.TxItem key={i} />)}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="px-4 py-3 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">
        {t('common.error_load')}: {error}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-semibold text-gray-900">{t('dashboard.admin_title')}</h2>
        <p className="text-gray-400 text-sm mt-1">{t('dashboard.admin_subtitle')}</p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.total_members')}</p>
          <p className="text-2xl font-bold text-gray-900">{members.length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.total_transactions')}</p>
          <p className="text-2xl font-bold text-gray-900">{totalTx}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.income')}</p>
          <p className="text-2xl font-bold text-green-600">{formatCurrency(totalIncome)}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.expense')}</p>
          <p className="text-2xl font-bold text-red-600">{formatCurrency(totalExpense)}</p>
        </div>
      </div>

      {/* Charts row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('dashboard.expense_by_category')}</h3>
          <ExpenseChart data={categoryBreakdown} />
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('dashboard.income_vs_expense')}</h3>
          <IncomeVsExpenseChart data={monthlySeries} />
        </div>
      </div>

      {/* Recent transactions */}
      <div className="bg-white rounded-xl border border-gray-200 p-5">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-semibold text-gray-900">{t('dashboard.recent_transactions')}</h3>
          <Link to="/admin/transactions" className="text-xs text-gray-400 hover:text-gray-900 transition">
            {t('common.view_all')} →
          </Link>
        </div>
        <div className="space-y-3">
          {transactions.length === 0 && (
            <p className="text-sm text-gray-400">{t('dashboard.no_transactions_this_month')}</p>
          )}
          {transactions.slice(0, 5).map((tx) => {
            const amount = parseInt(tx.amount) || 0
            const isIncome = tx.type === 'income'
            const isTransfer = tx.type === 'transfer'
            // Get member name from member_id (last 8 chars as fallback)
            const member = members.find((m) => m.id === tx.member_id)
            const memberName = member?.name || member?.username || tx.member_id?.slice(-8) || ''
            return (
              <div key={tx.id} className="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
                <div className="flex items-center gap-3">
                  <div className={`w-8 h-8 rounded-lg flex items-center justify-center text-sm ${
                    isIncome ? 'bg-green-50 text-green-600'
                    : isTransfer ? 'bg-blue-50 text-blue-600'
                    : 'bg-red-50 text-red-600'
                  }`}>
                    {isIncome ? '📈' : isTransfer ? '🔄' : '📉'}
                  </div>
                  <div>
                    <p className="text-sm text-gray-900">{tx.description}</p>
                    <p className="text-xs text-gray-400">{memberName} • {tx.date}{tx.category_name ? ` • ${tx.category_name}` : ''}</p>
                  </div>
                </div>
                <span className={`text-sm font-medium ${isIncome || isTransfer ? 'text-green-600' : 'text-red-600'}`}>
                  {isIncome || isTransfer ? '+' : '-'}{formatCurrency(amount)}
                </span>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
