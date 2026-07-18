import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getBudgetStatus, getCategories, setBudget } from '../../api'
import Drawer from '../../components/Drawer'
import Skeleton from '../../components/Skeleton'

function formatCurrency(amount: number | string) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(Math.abs(Number(amount)))
}

export default function Budget() {
  const { t } = useTranslation()
  const now = new Date()
  const month = now.getMonth() + 1
  const year = now.getFullYear()

  const [budgets, setBudgets] = useState<any[]>([]) 
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Form
  const [formOpen, setFormOpen] = useState(false)
  const [categories, setCategories] = useState<any[]>([])
  const [form, setForm] = useState({ category_id: '', amount: '' })
  const [saving, setSaving] = useState(false)

  const fetchBudgets = async () => {
    try {
      setLoading(true)
      const res = await getBudgetStatus({ month, year })
      setBudgets(res.budgets || [])
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchBudgets() }, []) // eslint-disable-line react-hooks/exhaustive-deps

  const openForm = async () => {
    setFormOpen(true)
    try {
      const res = await getCategories('expense')
      setCategories(res.categories || [])
    } catch (err: any) {
      console.error('Failed to load categories:', err)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await setBudget({
        category_id: form.category_id,
        amount: parseInt(form.amount),
        month,
        year,
      })
      setFormOpen(false)
      setForm({ category_id: '', amount: '' })
      await fetchBudgets()
    } catch (err: any) {
      alert(t('common.error_save') + ': ' + err.message)
    } finally {
      setSaving(false)
    }
  }

  const totalBudget = budgets.reduce((sum, b) => sum + parseInt(b.amount || 0), 0)
  const totalSpent = budgets.reduce((sum, b) => sum + parseInt(b.spent || 0), 0)
  const totalRemaining = totalBudget - totalSpent

  if (loading) {
    return (
      <div className="space-y-6">
        {/* Header */}
        <div>
          <Skeleton.Base className="h-7 w-32 mb-2" />
          <Skeleton.Base className="h-4 w-56" />
        </div>
        {/* Budget cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {Array.from({ length: 4 }).map((_, i) => <Skeleton.BudgetCard key={i} />)}
        </div>
        {/* Summary */}
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <Skeleton.Base className="h-4 w-36 mb-4" />
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div><Skeleton.Base className="h-3 w-20 mb-2" /><Skeleton.Base className="h-6 w-28" /></div>
            <div><Skeleton.Base className="h-3 w-20 mb-2" /><Skeleton.Base className="h-6 w-28" /></div>
            <div><Skeleton.Base className="h-3 w-12 mb-2" /><Skeleton.Base className="h-6 w-28" /></div>
          </div>
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
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-semibold text-gray-900">{t('budget.title')}</h2>
          <p className="text-gray-400 text-sm mt-1">{t('budget.subtitle')}</p>
        </div>
        <button
          onClick={openForm}
          className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition"
        >
          + {t('budget.add')}
        </button>
      </div>

      {/* Budget list */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {budgets.length === 0 && (
          <div className="bg-white rounded-xl border border-gray-200 p-5 col-span-full">
            <p className="text-sm text-gray-400 text-center">{t('budget.no_budget')}</p>
          </div>
        )}
        {budgets.map((budget) => {
          const percent = budget.percentage || 0
          const isOver = percent >= 100
          const remaining = parseInt(budget.amount || 0) - parseInt(budget.spent || 0)

          return (
            <div key={budget.id} className="bg-white rounded-xl border border-gray-200 p-5">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <span className="text-xl">🎯</span>
                  <h3 className="text-sm font-semibold text-gray-900">{budget.category_name || t('common.without_category')}</h3>
                </div>
                {isOver && (
                  <span className="px-2 py-0.5 bg-red-50 text-red-600 text-xs font-medium rounded-full">
                    {t('budget.over')}
                  </span>
                )}
              </div>

              {/* Progress bar */}
              <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden mb-3">
                <div
                  className={`h-full rounded-full transition-all ${
                    isOver ? 'bg-red-500' : percent > 80 ? 'bg-yellow-500' : 'bg-gray-900'
                  }`}
                  style={{ width: `${Math.min(percent, 100)}%` }}
                />
              </div>

              {/* Stats */}
              <div className="flex items-center justify-between text-sm gap-2">
                <span className="text-gray-500 truncate">
                  {formatCurrency(budget.spent || 0)} / {formatCurrency(budget.amount || 0)}
                </span>
                <span className={`font-medium whitespace-nowrap ${isOver ? 'text-red-600' : 'text-gray-900'}`}>
                  {percent}%
                </span>
              </div>

              {/* Remaining */}
              <p className={`text-xs mt-2 ${isOver ? 'text-red-500' : 'text-gray-400'}`}>
                {isOver
                  ? `⚠️ ${t('budget.exceeded')} ${formatCurrency(Math.abs(remaining))}`
                  : `${t('budget.remaining')} ${formatCurrency(Math.max(remaining, 0))}`
                }
              </p>
            </div>
          )
        })}
      </div>

      {/* Summary */}
      <div className="bg-white rounded-xl border border-gray-200 p-5">
        <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('budget.summary')}</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <p className="text-xs text-gray-400">{t('budget.total_budget')}</p>
            <p className="text-lg font-bold text-gray-900">{formatCurrency(totalBudget)}</p>
          </div>
          <div>
            <p className="text-xs text-gray-400">{t('budget.total_spent')}</p>
            <p className="text-lg font-bold text-gray-900">{formatCurrency(totalSpent)}</p>
          </div>
          <div>
            <p className="text-xs text-gray-400">{t('budget.remaining_label')}</p>
            <p className="text-lg font-bold text-green-600">{formatCurrency(Math.max(totalRemaining, 0))}</p>
          </div>
        </div>
      </div>

      {/* Drawer tambah budget */}
      <Drawer open={formOpen} onClose={() => setFormOpen(false)} title={t('budget.add_title')}>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('budget.category')}</label>
            <select
              value={form.category_id}
              onChange={(e) => setForm({ ...form, category_id: e.target.value })}
              required
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            >
              <option value="">{t('budget.select_category')}</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('budget.amount_label')}</label>
            <input
              type="number" required min="1"
              value={form.amount}
              onChange={(e) => setForm({ ...form, amount: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder="200000"
            />
          </div>
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => setFormOpen(false)} className="flex-1 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium rounded-lg transition">
              {t('common.cancel')}
            </button>
            <button type="submit" disabled={saving} className="flex-1 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition disabled:opacity-50">
              {saving ? t('common.saving') : t('common.save')}
            </button>
          </div>
        </form>
      </Drawer>
    </div>
  )
}
