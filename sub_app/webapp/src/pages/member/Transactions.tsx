import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getTransactions, getCategories, getAccounts, createTransaction, updateTransaction, deleteMemberTransaction } from '../../api'
import Drawer from '../../components/Drawer'
import Pagination from '../../components/Pagination'
import Skeleton from '../../components/Skeleton'

function formatCurrency(amount: number | string) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(Math.abs(Number(amount)))
}

const PAGE_SIZE = 20

export default function Transactions() {
  const { t } = useTranslation()
  const now = new Date()
  const [transactions, setTransactions] = useState<any[]>([]) 
  const [filter, setFilter] = useState('all')
  const [search, setSearch] = useState('')
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [year, setYear] = useState(now.getFullYear())
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [sortBy, setSortBy] = useState('date')
  const [sortOrder, setSortOrder] = useState('desc')

  const toggleSort = (column: string) => {
    if (sortBy === column) {
      setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc')
    } else {
      setSortBy(column)
      setSortOrder(column === 'date' ? 'desc' : 'asc')
    }
  }

  const SortArrow = ({ column }: { column: string }) => {
    if (sortBy !== column) return <span className="text-gray-300 ml-1">↕</span>
    return <span className="text-gray-600 ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
  }

  // Form state
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<string | null>(null)
  const [categories, setCategories] = useState<any[]>([])
  const [accounts, setAccounts] = useState<any[]>([])
  const [form, setForm] = useState({
    type: 'expense', account_id: '', to_account_id: '', category_id: '',
    amount: '', description: '', date: new Date().toISOString().slice(0, 10),
  })
  const [saving, setSaving] = useState(false)

  const fetchTransactions = async (type = 'all', q = '', p = 1) => {
    try {
      setLoading(true)
      const params: Record<string, any> = { month, year, limit: PAGE_SIZE, offset: (p - 1) * PAGE_SIZE, sort_by: sortBy, sort_order: sortOrder }
      if (type !== 'all') params.type = type
      if (q) params.q = q
      const res = await getTransactions(params)
      setTransactions(res.transactions || [])
      setTotal(res.total || 0)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    setPage(1)
  }, [filter, month, year, search, sortBy, sortOrder])

  useEffect(() => {
    fetchTransactions(filter === 'all' ? '' : filter, search, page)
  }, [filter, month, year, search, page, sortBy, sortOrder]) // eslint-disable-line react-hooks/exhaustive-deps

  const openForm = async (tx?: any) => {
    setEditing(tx?.id || null)
    setFormOpen(true)
    try {
      const [catsRes, acctsRes] = await Promise.all([
        getCategories(),
        getAccounts(),
      ])
      setCategories(catsRes.categories || [])
      setAccounts(acctsRes.accounts || [])
      if (tx) {
        setForm({
          type: tx.type || 'expense',
          account_id: tx.account_id || '',
          to_account_id: tx.to_account_id || '',
          category_id: tx.category_id || '',
          amount: String(tx.amount || ''),
          description: tx.description || '',
          date: tx.date || new Date().toISOString().slice(0, 10),
        })
      }
    } catch (err: any) {
      console.error('Failed to load form data:', err)
    }
  }

  const closeForm = () => {
    setFormOpen(false)
    setEditing(null)
    setForm({
      type: 'expense', account_id: '', to_account_id: '', category_id: '',
      amount: '', description: '', date: new Date().toISOString().slice(0, 10),
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      const payload = {
        type: form.type,
        account_id: form.account_id,
        to_account_id: form.type === 'transfer' ? form.to_account_id : undefined,
        category_id: form.category_id,
        amount: parseInt(form.amount),
        description: form.description,
        date: form.date,
      }
      if (editing) {
        await updateTransaction(editing, payload)
      } else {
        await createTransaction(payload)
      }
      closeForm()
      await fetchTransactions(filter === 'all' ? '' : filter, search, page)
    } catch (err: any) {
      alert('Gagal menyimpan: ' + err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus transaksi ini?')) return
    try {
      await deleteMemberTransaction(id)
      await fetchTransactions(filter === 'all' ? '' : filter, search, page)
    } catch (err: any) {
      alert('Gagal menghapus: ' + err.message)
    }
  }

  const filtered = filter === 'all' ? transactions : transactions.filter(tx => tx.type === filter)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-semibold text-gray-900">{t('transactions.title')}</h2>
          <p className="text-gray-400 text-sm mt-1">{t('transactions.subtitle')}</p>
        </div>
        <button
          onClick={openForm}
          className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition"
        >
          + {t('transactions.add')}
        </button>
      </div>

      {/* Filters */}
      <div className="flex gap-2 flex-wrap items-center">
        {[
          { value: 'all', label: t('common.all') },
          { value: 'income', label: '📈 ' + t('transactions.income') },
          { value: 'expense', label: '📉 ' + t('transactions.expense') },
          { value: 'transfer', label: '🔄 ' + t('transactions.transfer') },
        ].map((f) => (
          <button
            key={f.value}
            onClick={() => setFilter(f.value)}
            className={`px-3 py-1.5 text-sm rounded-lg border transition ${
              filter === f.value
                ? 'bg-gray-900 text-white border-gray-900'
                : 'bg-white text-gray-600 border-gray-200 hover:bg-gray-50'
            }`}
          >
            {f.label}
          </button>
        ))}

        {/* Search */}
        <input
          type="text" value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder={'🔍 ' + t('transactions.search')}
          className="px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-900/10 w-full sm:w-48"
        />

        {/* Month picker */}
        <div className="ml-auto flex gap-2">
          <select
            value={month}
            onChange={e => setMonth(Number(e.target.value))}
            className="px-2 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-900/10"
          >
            {Array.from({ length: 12 }, (_, i) => i + 1).map(m => (
              <option key={m} value={m}>
                {new Date(2000, m - 1).toLocaleString('id', { month: 'long' })}
              </option>
            ))}
          </select>
          <select
            value={year}
            onChange={e => setYear(Number(e.target.value))}
            className="px-2 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-900/10"
          >
            {Array.from({ length: 5 }, (_, i) => year - 2 + i).map(y => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{error}</div>
      )}

      {/* Transaction list */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-x-auto">
        {loading ? (
          <div>
            <Skeleton.Table rows={8} cols={3} />
          </div>
        ) : (
          <div>
            {/* Sortable column headers */}
            <div className="flex items-center px-5 py-2.5 border-b border-gray-100 bg-gray-50/50 text-xs font-medium text-gray-400 uppercase tracking-wider">
              <div className="flex-1 flex items-center gap-6">
                <button onClick={() => toggleSort('description')} className="flex items-center hover:text-gray-700 transition">
                  {t('transactions.description')} <SortArrow column="description" />
                </button>
                <button onClick={() => toggleSort('date')} className="flex items-center hover:text-gray-700 transition">
                  {t('transactions.date')} <SortArrow column="date" />
                </button>
                <button onClick={() => toggleSort('amount')} className="flex items-center hover:text-gray-700 transition">
                  {t('transactions.amount')} <SortArrow column="amount" />
                </button>
              </div>
              <div className="w-16" />
            </div>
            <div className="divide-y divide-gray-100 min-w-[600px]">
              {filtered.map((tx) => {
                const amount = parseInt(tx.amount) || 0
                const isIncome = tx.type === 'income'
                const isTransfer = tx.type === 'transfer'
                return (
                  <div key={tx.id} className="flex items-center justify-between px-5 py-4 hover:bg-gray-50 transition">
                    <div className="flex items-center gap-4">
                      <div className={`w-10 h-10 rounded-lg flex items-center justify-center text-base ${
                        isIncome ? 'bg-green-50 text-green-600'
                        : isTransfer ? 'bg-blue-50 text-blue-600'
                        : 'bg-red-50 text-red-600'
                      }`}>
                        {isIncome ? '📈' : isTransfer ? '🔄' : '📉'}
                      </div>
                      <div>
                        <p className="text-sm font-medium text-gray-900">{tx.description}</p>
                        <p className="text-xs text-gray-400 mt-0.5">
                          {tx.date} • {tx.account_name || ''}
                          {tx.type === 'transfer' && tx.to_account_name ? ` → ${tx.to_account_name}` : ''}
                          {tx.category_name ? ` • ${tx.category_name}` : ''}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className={`text-sm font-semibold ${isIncome || isTransfer ? 'text-green-600' : 'text-red-600'}`}>
                        {isIncome || isTransfer ? '+' : '-'}{formatCurrency(amount)}
                      </span>
                      <button
                        onClick={() => openForm(tx)}
                        className="text-gray-300 hover:text-gray-600 transition ml-2"
                        title={t('common.edit')}
                      >
                        ✏️
                      </button>
                      <button
                        onClick={() => handleDelete(tx.id)}
                        className="text-gray-300 hover:text-red-500 transition ml-2"
                        title={t('common.delete')}
                      >
                        🗑️
                      </button>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {!loading && filtered.length === 0 && (
          <div className="px-5 py-12 text-center">
            <p className="text-gray-400 text-sm">{t('transactions.no_transactions')}</p>
          </div>
        )}

        {!loading && (
          <Pagination
            page={page}
            totalPages={Math.max(1, Math.ceil(total / PAGE_SIZE))}
            onPageChange={setPage}
          />
        )}
      </div>

      {/* Drawer tambah/edit transaksi */}
      <Drawer open={formOpen} onClose={closeForm} title={editing ? t('transactions.edit_title') : t('transactions.add_title')}>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Tipe */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.type')}</label>
            <select
              value={form.type}
              onChange={(e) => setForm({ ...form, type: e.target.value, category_id: '' })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            >
              <option value="expense">{t('transactions.expense')}</option>
              <option value="income">{t('transactions.income')}</option>
              <option value="transfer">{t('transactions.transfer')}</option>
            </select>
          </div>

          {/* Akun Sumber */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.account')} {form.type === 'transfer' ? t('transactions.source') : ''}</label>
            <select
              value={form.account_id}
              onChange={(e) => setForm({ ...form, account_id: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            >
              <option value="">{t('transactions.select_account')}</option>
              {accounts.map((a) => (
                <option key={a.id} value={a.id}>{a.name}</option>
              ))}
            </select>
          </div>

          {/* Akun Tujuan (khusus transfer) */}
          {form.type === 'transfer' && (
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.destination_account')}</label>
              <select
                value={form.to_account_id}
                onChange={(e) => setForm({ ...form, to_account_id: e.target.value })}
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              >
                <option value="">{t('transactions.select_destination')}</option>
                {accounts
                  .filter((a: any) => a.id !== form.account_id)
                  .map((a: any) => (
                    <option key={a.id} value={a.id}>{a.name}</option>
                  ))}
              </select>
            </div>
          )}

          {/* Kategori */}
          {form.type !== 'transfer' && (
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.category')}</label>
              <select
                value={form.category_id}
                onChange={(e) => setForm({ ...form, category_id: e.target.value })}
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              >
                <option value="">{t('transactions.select_category')}</option>
                {categories
                  .filter((c) => c.type === form.type)
                  .map((c) => (
                    <option key={c.id} value={c.id}>{c.name}</option>
                  ))}
              </select>
            </div>
          )}

          {/* Jumlah */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.amount_label')}</label>
            <input
              type="number" required min="1"
              value={form.amount}
              onChange={(e) => setForm({ ...form, amount: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder="100000"
            />
          </div>

          {/* Deskripsi */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.description_label')}</label>
            <input
              type="text" required
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder={t('transactions.description_placeholder')}
            />
          </div>

          {/* Tanggal */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.date_label')}</label>
            <input
              type="date" required
              value={form.date}
              onChange={(e) => setForm({ ...form, date: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            />
          </div>

          <div className="flex gap-3 pt-2">
            <button type="button" onClick={closeForm} className="flex-1 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium rounded-lg transition">
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
