import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getMembers, createMember, updateMember, deleteMember, approveMember, rejectMember } from '../../api'
import Drawer from '../../components/Drawer'
import Skeleton from '../../components/Skeleton'

export default function Members() {
  const { t } = useTranslation()
  const [members, setMembers] = useState<any[]>([]) 
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<any>(null)
  const [form, setForm] = useState({ name: '', username: '', password: '' })
  const [saving, setSaving] = useState(false)

  const fetchMembers = async () => {
    try {
      setLoading(true)
      const res = await getMembers()
      setMembers(res.members || [])
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchMembers() }, [])

  const openAddForm = () => {
    setEditing(null)
    setForm({ name: '', username: '', password: '' })
    setFormOpen(true)
  }

  const openEditForm = (member: any) => {
    setEditing(member)
    setForm({ name: member.name || '', username: member.username || '', password: '' })
    setFormOpen(true)
  }

  const closeForm = () => {
    setFormOpen(false)
    setEditing(null)
    setForm({ name: '', username: '', password: '' })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      if (editing) {
        const payload: Record<string, string> = {}
        if (form.name) payload.name = form.name
        if (form.username) payload.username = form.username
        await updateMember(editing.id, payload)
      } else {
        await createMember({
          name: form.name,
          username: form.username,
          password: form.password,
        })
      }
      closeForm()
      await fetchMembers()
    } catch (err: any) {
      alert(t('common.error_save') + ': ' + err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm(t('members.confirm_delete'))) return
    try {
      await deleteMember(id)
      await fetchMembers()
    } catch (err: any) {
      alert(t('common.error_delete') + ': ' + err.message)
    }
  }

  const handleApprove = async (id: string) => {
    try {
      await approveMember(id)
      await fetchMembers()
    } catch (err: any) {
      alert(t('common.error_generic') + ': ' + err.message)
    }
  }

  const handleReject = async (id: string) => {
    try {
      await rejectMember(id)
      await fetchMembers()
    } catch (err: any) {
      alert(t('common.error_generic') + ': ' + err.message)
    }
  }

  const statusBadge = (status: string) => {
    switch (status) {
      case 'active':
        return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700">Active</span>
      case 'pending':
        return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-700">Pending</span>
      case 'rejected':
        return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-700">Rejected</span>
      default:
        return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600">{status}</span>
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-semibold text-gray-900">{t('members.title')}</h2>
          <p className="text-gray-400 text-sm mt-1">{t('members.subtitle')}</p>
        </div>
        <button
          onClick={openAddForm}
          className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition"
        >
          + {t('members.add')}
        </button>
      </div>

      {/* Error */}
      {error && (
        <div className="px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{error}</div>
      )}

      {/* Table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-100">
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-400 uppercase tracking-wider">{t('members.member')}</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-400 uppercase tracking-wider hidden sm:table-cell">{t('auth.username')}</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-400 uppercase tracking-wider hidden md:table-cell">{t('members.status')}</th>
                <th className="text-center px-5 py-3 text-xs font-medium text-gray-400 uppercase tracking-wider">{t('members.actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {loading ? (
                Array.from({ length: 6 }).map((_, i) => (
                  <tr key={i}>
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-3">
                        <Skeleton.Base className="w-8 h-8 rounded-full" />
                        <Skeleton.Base className="h-4 w-32" />
                      </div>
                    </td>
                    <td className="px-5 py-4 hidden sm:table-cell"><Skeleton.Base className="h-4 w-24" /></td>
                    <td className="px-5 py-4 hidden md:table-cell"><Skeleton.Base className="h-4 w-16" /></td>
                    <td className="px-5 py-4 text-center"><Skeleton.Base className="h-4 w-16 mx-auto" /></td>
                  </tr>
                ))
              ) : members.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-5 py-12 text-center text-sm text-gray-400">{t('members.no_members')}</td>
                </tr>
              ) : (
                members.map((member) => (
                  <tr key={member.id} className="hover:bg-gray-50 transition">
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center text-sm font-medium text-gray-600">
                          {(member.name || member.username || 'M').charAt(0).toUpperCase()}
                        </div>
                        <span className="text-sm font-medium text-gray-900">{member.name || member.username}</span>
                      </div>
                    </td>
                    <td className="px-5 py-4 text-sm text-gray-500 hidden sm:table-cell">{member.username}</td>
                    <td className="px-5 py-4 hidden md:table-cell">{statusBadge(member.status || 'active')}</td>
                    <td className="px-5 py-4 text-center">
                      <div className="flex items-center justify-center gap-2">
                        {member.status === 'pending' && (
                          <>
                            <button
                              onClick={() => handleApprove(member.id)}
                              className="text-sm text-green-500 hover:text-green-700 transition"
                              title={t('members.approve')}
                            >
                              ✅
                            </button>
                            <button
                              onClick={() => handleReject(member.id)}
                              className="text-sm text-red-400 hover:text-red-600 transition"
                              title={t('members.reject')}
                            >
                              ❌
                            </button>
                          </>
                        )}
                        <button
                          onClick={() => openEditForm(member)}
                          className="text-sm text-gray-400 hover:text-gray-600 transition"
                          title={t('common.edit')}
                        >
                          ✏️
                        </button>
                        <button
                          onClick={() => handleDelete(member.id)}
                          className="text-sm text-gray-400 hover:text-red-600 transition"
                          title={t('common.delete')}
                        >
                          🗑️
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Drawer tambah/edit member */}
      <Drawer open={formOpen} onClose={closeForm} title={editing ? t('members.edit_title') : t('members.add_title')}>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('members.name')}</label>
            <input
              type="text" required
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder={t('members.name_placeholder')}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.username')}</label>
            <input
              type="text" required
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder={t('auth.username')}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.password')} {!editing && <span className="text-red-400">*</span>}</label>
            <input
              type="password" required={!editing}
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder={editing ? t('members.password_empty_hint') : t('members.password_min_hint')}
              minLength={6}
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
