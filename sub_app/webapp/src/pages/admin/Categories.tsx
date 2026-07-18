import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getAdminCategories, createAdminCategory, updateAdminCategory, deleteAdminCategory, type CategoryData } from '../../api'
import Drawer from '../../components/Drawer'
import Skeleton from '../../components/Skeleton'
import CategoryIcon from '../../components/CategoryIcon'

const ICON_OPTIONS = [
  'briefcase', 'gift', 'building', 'trending-up', 'plus-circle',
  'shopping-cart', 'truck', 'shopping-bag', 'file-text', 'music',
  'heart', 'book', 'home', 'piggy-bank', 'alert-triangle',
  'wallet', 'car', 'plane', 'graduation', 'stethoscope',
  'gamepad', 'shirt', 'paw', 'lightbulb', 'smartphone',
  'dumbbell', 'baby', 'wrench', 'package', 'credit-card',
]

export default function Categories() {
  const { t } = useTranslation()
  const [categories, setCategories] = useState<CategoryData[]>([]) 
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<CategoryData | null>(null)
  const [form, setForm] = useState({ name: '', type: 'expense', icon: 'briefcase' })
  const [saving, setSaving] = useState(false)

  const fetchCategories = async () => {
    try {
      setLoading(true)
      const res = await getAdminCategories()
      setCategories(res.categories || [])
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchCategories() }, [])

  const openAddForm = (catType: string) => {
    setEditing(null)
    setForm({ name: '', type: catType, icon: 'briefcase' })
    setFormOpen(true)
  }

  const openEditForm = (cat: CategoryData) => {
    setEditing(cat)
    setForm({ name: cat.name, type: cat.type, icon: cat.icon || 'briefcase' })
    setFormOpen(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim()) return
    setSaving(true)
    try {
      if (editing) {
        await updateAdminCategory(editing.id, form)
      } else {
        await createAdminCategory(form)
      }
      setFormOpen(false)
      setEditing(null)
      setForm({ name: '', type: 'expense', icon: 'briefcase' })
      await fetchCategories()
    } catch (err: any) {
      alert(t('common.error_save') + ': ' + err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (cat: CategoryData) => {
    if (!confirm(t('categories.confirm_delete', { name: cat.name }))) return
    try {
      await deleteAdminCategory(cat.id)
      await fetchCategories()
    } catch (err: any) {
      alert(t('common.error_delete') + ': ' + err.message)
    }
  }

  const income = categories.filter(c => c.type === 'income')
  const expense = categories.filter(c => c.type === 'expense')

  const renderCategoryGrid = (items: CategoryData[]) => (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
      {items.length === 0 && (
        <p className="text-sm text-gray-400 col-span-full">{t('categories.no_categories')}</p>
      )}
      {items.map((cat) => (
        <div key={cat.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg group">
          <div className="flex items-center gap-3">
            <span className="text-lg text-gray-600"><CategoryIcon name={cat.icon} /></span>
            <span className="text-sm text-gray-900">{cat.name}</span>
          </div>
          <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition">
            <button onClick={() => openEditForm(cat)} className="p-1 text-gray-400 hover:text-gray-600 transition text-xs" title={t('common.edit')}>✏️</button>
            <button onClick={() => handleDelete(cat)} className="p-1 text-gray-400 hover:text-red-500 transition text-xs" title={t('common.delete')}>🗑️</button>
          </div>
        </div>
      ))}
    </div>
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-semibold text-gray-900">{t('categories.title')}</h2>
          <p className="text-gray-400 text-sm mt-1">{t('categories.subtitle')}</p>
        </div>
      </div>

      {error && (
        <div className="px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{error}</div>
      )}

      {loading ? (
        <div className="space-y-6">
          <div className="bg-white rounded-xl border border-gray-200 p-5">
            <Skeleton.Base className="h-4 w-28 mb-4" />
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {Array.from({ length: 6 }).map((_, i) => (
                <div key={i} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                  <Skeleton.Base className="w-6 h-6 rounded" />
                  <Skeleton.Base className="h-4 w-24" />
                </div>
              ))}
            </div>
          </div>
          <div className="bg-white rounded-xl border border-gray-200 p-5">
            <Skeleton.Base className="h-4 w-28 mb-4" />
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {Array.from({ length: 6 }).map((_, i) => (
                <div key={i} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                  <Skeleton.Base className="w-6 h-6 rounded" />
                  <Skeleton.Base className="h-4 w-24" />
                </div>
              ))}
            </div>
          </div>
        </div>
      ) : (
        <>
          {/* Income */}
          <div className="bg-white rounded-xl border border-gray-200 p-5">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-semibold text-gray-900">📈 {t('transactions.income')}</h3>
              <button onClick={() => openAddForm('income')} className="text-xs text-gray-400 hover:text-gray-600 transition">+ {t('common.add')}</button>
            </div>
            {renderCategoryGrid(income)}
          </div>

          {/* Expense */}
          <div className="bg-white rounded-xl border border-gray-200 p-5">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-semibold text-gray-900">📉 {t('transactions.expense')}</h3>
              <button onClick={() => openAddForm('expense')} className="text-xs text-gray-400 hover:text-gray-600 transition">+ {t('common.add')}</button>
            </div>
            {renderCategoryGrid(expense)}
          </div>
        </>
      )}

      {/* Add/Edit Drawer */}
      <Drawer
        open={formOpen}
        onClose={() => { setFormOpen(false); setEditing(null) }}
        title={editing ? t('categories.edit_title') : t('categories.add_title')}
      >
        <form onSubmit={handleSubmit} className="space-y-5">
          {/* Type */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('transactions.type')}</label>
            <div className="flex gap-2">
              {['expense', 'income'].map(type => (
                <button key={type} type="button" onClick={() => setForm({ ...form, type })}
                  className={`px-3 py-1.5 text-sm rounded-lg border transition ${
                    form.type === type
                      ? 'bg-gray-900 text-white border-gray-900'
                      : 'bg-white text-gray-600 border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  {type === 'income' ? '📈 ' + t('transactions.income') : '📉 ' + t('transactions.expense')}
                </button>
              ))}
            </div>
          </div>

          {/* Name */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('categories.name')}</label>
            <input
              type="text" value={form.name}
              onChange={e => setForm({ ...form, name: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm"
              placeholder={t('categories.name_placeholder')} required
            />
          </div>

          {/* Icon */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('categories.icon')}</label>
            <div className="flex flex-wrap gap-2">
              {ICON_OPTIONS.map(icon => (
                <button key={icon} type="button" onClick={() => setForm({ ...form, icon })}
                  className={`w-9 h-9 flex items-center justify-center rounded-lg border transition ${
                    form.icon === icon
                      ? 'border-gray-900 bg-gray-50 shadow-sm'
                      : 'border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  <CategoryIcon name={icon} className="w-4 h-4" />
                </button>
              ))}
            </div>
          </div>

          {/* Submit */}
          <button type="submit" disabled={saving || !form.name.trim()}
            className="w-full py-2.5 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition duration-200 active:scale-[0.98] disabled:opacity-50"
          >
            {saving ? t('common.saving') : editing ? t('common.save') : t('common.add')}
          </button>
        </form>
      </Drawer>
    </div>
  )
}
