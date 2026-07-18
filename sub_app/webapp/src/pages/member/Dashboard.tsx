import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../../context/AuthContext'
import { getAccounts, getTransactionSummaryWithBreakdown, getTransactions, getBudgetStatus, getMonthlySeries } from '../../api'
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

const MONTH_NAMES = ['', 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

export default function Dashboard() {
  const { user } = useAuth()
  const { t } = useTranslation()
  const now = new Date()
  const month = now.getMonth() + 1
  const year = now.getFullYear()

  const [accounts, setAccounts] = useState<any[]>([])
  const [summary, setSummary] = useState<{ total_income: number | string; total_expense: number | string }>({ total_income: 0, total_expense: 0 })
  const [recentTxs, setRecentTxs] = useState<any[]>([])
  const [budgets, setBudgets] = useState<any[]>([])
  const [categoryBreakdown, setCategoryBreakdown] = useState<any[]>([])
  const [monthlySeries, setMonthlySeries] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function fetchData() {
      try {
        const [acctsRes, summRes, txsRes, budgRes, seriesRes] = await Promise.all([
          getAccounts(),
          getTransactionSummaryWithBreakdown({ month, year }),
          getTransactions({ month, year, limit: 5 }),
          getBudgetStatus({ month, year }),
          getMonthlySeries({ months: 6 }),
        ])
        setAccounts(acctsRes.accounts || [])
        setSummary(summRes)
        setCategoryBreakdown(summRes.category_breakdown || [])
        setRecentTxs(txsRes.transactions || [])
        setBudgets(budgRes.budgets || [])
        setMonthlySeries(seriesRes.series || [])
      } catch (err: any) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [month, year])

  if (loading) {
    return (
      <div className="space-y-6">
        {/* Greeting skeleton */}
        <div>
          <Skeleton.Base className="h-7 w-72 mb-2" />
          <Skeleton.Base className="h-4 w-56" />
        </div>
        {/* Stat cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
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
          {Array.from({ length: 3 }).map((_, i) => <Skeleton.TxItem key={i} />)}
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

  const totalBalance = accounts.reduce((sum, a) => sum + (a.balance || 0), 0)

  return (
    <div className="space-y-6">
      {/* Greeting */}
      <div>
        <h2 className="text-2xl font-semibold text-gray-900">{t('dashboard.greeting', { username: user?.username || 'User' })} 👋</h2>
        <p className="text-gray-400 text-sm mt-1">{t('dashboard.summary', { month: MONTH_NAMES[month], year })}</p>
      </div>

      {/* Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* Total Saldo */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.total_balance')}</p>
          <p className="text-2xl font-bold text-gray-900">{formatCurrency(totalBalance)}</p>
          <p className="text-xs text-gray-300 mt-2">{accounts.length} {t('dashboard.accounts')}</p>
        </div>

        {/* Income */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.income')}</p>
          <p className="text-2xl font-bold text-green-600">{formatCurrency(summary.total_income || 0)}</p>
          <p className="text-xs text-gray-300 mt-2">{MONTH_NAMES[month]} {year}</p>
        </div>

        {/* Expense */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <p className="text-sm text-gray-400 mb-1">{t('dashboard.expense')}</p>
          <p className="text-2xl font-bold text-red-600">{formatCurrency(summary.total_expense || 0)}</p>
          <p className="text-xs text-gray-300 mt-2">{MONTH_NAMES[month]} {year}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Chart: Expense by Category */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('dashboard.expense_by_category')}</h3>
          <ExpenseChart data={categoryBreakdown} />
        </div>

        {/* Chart: Income vs Expense */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('dashboard.income_vs_expense')}</h3>
          <IncomeVsExpenseChart data={monthlySeries} />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Budget Status */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('dashboard.budget_status')}</h3>
          <div className="space-y-4">
            {budgets.length === 0 && (
              <p className="text-sm text-gray-400">{t('dashboard.no_budget')}</p>
            )}
            {budgets.map((budget) => {
              const percent = budget.percentage || 0
              const isOver = budget.remaining <= 0
              return (
                <div key={budget.id}>
                  <div className="flex items-center justify-between mb-1.5 gap-2">
                    <span className="text-sm text-gray-600 truncate">{budget.category_name || t('common.without_category')}</span>
                    <span className={`text-sm font-medium whitespace-nowrap ${isOver ? 'text-red-600' : 'text-gray-900'}`}>
                      {formatCurrency(budget.spent)} / {formatCurrency(budget.amount)}
                    </span>
                  </div>
                  <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all ${isOver ? 'bg-red-500' : percent > 80 ? 'bg-yellow-500' : 'bg-gray-900'}`}
                      style={{ width: `${Math.min(percent, 100)}%` }}
                    />
                  </div>
                  {isOver && (
                    <p className="text-xs text-red-500 mt-1">⚠️ {t('dashboard.over_budget')}</p>
                  )}
                </div>
              )
            })}
          </div>
        </div>

        {/* Accounts */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('dashboard.accounts_title')}</h3>
          <div className="space-y-3">
            {accounts.length === 0 && (
              <p className="text-sm text-gray-400">{t('dashboard.no_accounts')}</p>
            )}
            {accounts.map((account) => (
              <div key={account.id} className="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 bg-gray-100 rounded-lg flex items-center justify-center text-sm">
                    💰
                  </div>
                  <span className="text-sm text-gray-600">{account.name}</span>
                </div>
                <span className="text-sm font-medium text-gray-900">{formatCurrency(account.balance || 0)}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Recent Transactions */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 overflow-x-auto">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-semibold text-gray-900">{t('dashboard.recent_transactions')}</h3>
          <Link to="/member/transactions" className="text-xs text-gray-400 hover:text-gray-900 transition">
            {t('common.view_all')} →
          </Link>
        </div>
        <div className="space-y-3">
          {recentTxs.length === 0 && (
            <p className="text-sm text-gray-400">{t('dashboard.no_transactions_this_month')}</p>
          )}
          {recentTxs.map((tx) => {
            const amount = parseInt(tx.amount) || 0
            const isIncome = tx.type === 'income'
            const isTransfer = tx.type === 'transfer'
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
                    <p className="text-xs text-gray-400">{tx.date}{tx.category_name ? ` • ${tx.category_name}` : ''}</p>
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
