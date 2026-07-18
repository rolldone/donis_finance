import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getAccounts, createAccount, updateAccount, deleteAccount } from '../../api'
import Drawer from '../../components/Drawer'
import Skeleton from '../../components/Skeleton'

function formatCurrency(amount: number | string) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(Math.abs(Number(amount)))
}

const ACCOUNT_TYPES = ['cash', 'bank', 'e_wallet', 'savings', 'investment'] as const

const TYPE_LABELS: Record<string, string> = {
  cash: 'Cash',
  bank: 'Bank',
  e_wallet: 'E-Wallet',
  savings: 'Tabungan',
  investment: 'Investasi',
}

const TYPE_ICONS: Record<string, string> = {
  cash: '💵',
  bank: '🏦',
  e_wallet: '📱',
  savings: '🏧',
  investment: '📈',
}

export default function Accounts() {
  const { t } = useTranslation()
  const [accounts, setAccounts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Form state
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<string | null>(null)
  const [form, setForm] = useState({ name: '', type: 'bank', initial_balance: '' })
  const [saving, setSaving] = useState(false)

  const fetchAccounts = async () => {
    try {
      setLoading(true)
      const res = await getAccounts()
      setAccounts(res.accounts || [])
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchAccounts()
  }, [])

  const openForm = async (acct?: any) => {
    setEditing(acct?.id || null)
    if (acct) {
      setForm({ name: acct.name, type: acct.type, initial_balance: '' })
    } else {
      setForm({ name: '', type: 'bank', initial_balance: '' })
    }
    setFormOpen(true)
  }

  const handleSave = async () => {
    if (!form.name.trim()) {
      alert(t('accounts.name_required'))
      return
    }
    setSaving(true)
    try {
      if (editing) {
        await updateAccount(editing, form.name.trim(), form.type)
      } else {
        const balance = parseInt(form.initial_balance) || 0
        await createAccount(form.name.trim(), form.type, balance)
      }
      setFormOpen(false)
      setEditing(null)
      await fetchAccounts()
    } catch (err: any) {
      alert(err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(t('accounts.delete_confirm', { name }))) return
    try {
      await deleteAccount(id)
      await fetchAccounts()
    } catch (err: any) {
      alert(err.message)
    }
  }

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <Skeleton.Base className="h-7 w-48" />
          <Skeleton.Base className="h-9 w-32 rounded-lg" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Skeleton.Card />
          <Skeleton.Card />
          <Skeleton.Card />
          <Skeleton.Card />
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
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-semibold text-gray-900">{t('accounts.title')}</h2>
          <p className="text-sm text-gray-400 mt-1">
            {t('accounts.total_balance')}: <span className="font-semibold text-gray-900">{formatCurrency(totalBalance)}</span>
          </p>
        </div>
        <button
          onClick={() => openForm()}
          className="px-4 py-2 bg-gray-900 text-white text-sm font-medium rounded-lg hover:bg-gray-800 transition-colors"
        >
          + {t('accounts.add')}
        </button>
      </div>

      {/* Accounts Grid */}
      {accounts.length === 0 ? (
        <div className="text-center py-16 bg-white rounded-xl border border-gray-200">
          <div className="text-4xl mb-3">🏦</div>
          <p className="text-gray-400 text-sm">{t('accounts.no_accounts')}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {accounts.map((acct) => (
            <div key={acct.id} className="bg-white rounded-xl border border-gray-200 p-5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 hover:shadow-sm transition-shadow">
              <div className="flex items-center gap-4">
                <div className="w-10 h-10 bg-gray-100 rounded-xl flex items-center justify-center text-lg shrink-0">
                  {TYPE_ICONS[acct.type] || '💰'}
                </div>
                <div className="min-w-0">
                  <p className="font-medium text-gray-900 truncate">{acct.name}</p>
                  <p className="text-xs text-gray-400">{TYPE_LABELS[acct.type] || acct.type}</p>
                </div>
              </div>
              <div className="flex items-center gap-3 sm:justify-end">
                <div className="text-right">
                  <p className="font-semibold text-gray-900">{formatCurrency(acct.balance || 0)}</p>
                </div>
                <button
                  onClick={() => openForm(acct)}
                  className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                  title={t('common.edit')}
                >
                  ✏️
                </button>
                <button
                  onClick={() => handleDelete(acct.id, acct.name)}
                  className="p-1.5 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                  title={t('common.delete')}
                >
                  🗑️
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Form Drawer */}
      <Drawer open={formOpen} onClose={() => { setFormOpen(false); setEditing(null) }}>
        <div className="p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-6">
            {editing ? t('accounts.edit_title') : t('accounts.add_title')}
          </h3>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">{t('accounts.name')}</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                placeholder={t('accounts.name_placeholder')}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">{t('accounts.type')}</label>
              <select
                value={form.type}
                onChange={(e) => setForm({ ...form, type: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
              >
                {ACCOUNT_TYPES.map((t) => (
                  <option key={t} value={t}>{TYPE_LABELS[t]}</option>
                ))}
              </select>
            </div>

            {!editing && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">{t('accounts.initial_balance')}</label>
                <input
                  type="number"
                  value={form.initial_balance}
                  onChange={(e) => setForm({ ...form, initial_balance: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                  placeholder="0"
                  min="0"
                />
              </div>
            )}
          </div>

          <div className="flex items-center gap-3 mt-8">
            <button
              onClick={handleSave}
              disabled={saving}
              className="flex-1 px-4 py-2.5 bg-gray-900 text-white text-sm font-medium rounded-lg hover:bg-gray-800 transition-colors disabled:opacity-50"
            >
              {saving ? t('common.saving') : t('common.save')}
            </button>
            <button
              onClick={() => { setFormOpen(false); setEditing(null) }}
              className="px-4 py-2.5 text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors"
            >
              {t('common.cancel')}
            </button>
          </div>
        </div>
      </Drawer>
    </div>
  )
}
