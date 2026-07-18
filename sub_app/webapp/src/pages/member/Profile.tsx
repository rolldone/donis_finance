import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../../context/AuthContext'
import { getMemberProfile, updateMemberProfile, changeMemberPassword } from '../../api'

export default function Profile() {
  const { user } = useAuth()
  const { t } = useTranslation()
  const [profile, setProfile] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Edit mode
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState('')
  const [username, setUsername] = useState('')
  const [saving, setSaving] = useState(false)

  // Password
  const [showPasswordForm, setShowPasswordForm] = useState(false)
  const [oldPass, setOldPass] = useState('')
  const [newPass, setNewPass] = useState('')
  const [passMsg, setPassMsg] = useState('')
  const [passError, setPassError] = useState('')
  const [changingPass, setChangingPass] = useState(false)

  const fetchProfile = async () => {
    try {
      setLoading(true)
      const data = await getMemberProfile()
      setProfile(data)
      setName(data.name || '')
      setUsername(data.username || '')
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProfile()
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError('')
    try {
      const updated = await updateMemberProfile({ name, username })
      setProfile(updated)
      setEditing(false)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    setChangingPass(true)
    setPassError('')
    setPassMsg('')
    try {
      const res = await changeMemberPassword(oldPass, newPass)
      setPassMsg(res.message)
      setOldPass('')
      setNewPass('')
      setShowPasswordForm(false)
    } catch (err: any) {
      setPassError(err.message)
    } finally {
      setChangingPass(false)
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl space-y-6">
        <div>
          <h2 className="text-2xl font-semibold text-gray-900">{t('profile.title')}</h2>
          <p className="text-gray-400 text-sm mt-1">{t('profile.subtitle')}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-12 text-center">
          <p className="text-gray-400 text-sm">{t('common.loading')}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-semibold text-gray-900">{t('profile.title')}</h2>
        <p className="text-gray-400 text-sm mt-1">{t('profile.subtitle')}</p>
      </div>

      {/* Profile card */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        {/* Avatar + info */}
        <div className="flex items-center gap-4 mb-6 pb-6 border-b border-gray-100">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center text-gray-600 font-bold text-xl">
            {(profile?.username || user?.username || 'U').charAt(0).toUpperCase()}
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900">{profile?.name || profile?.username || 'User'}</h3>
            <p className="text-sm text-gray-400">ID: {profile?.id || user?.id || '-'}</p>
            <span className="inline-block mt-1 px-2 py-0.5 bg-gray-100 text-gray-500 text-xs font-medium rounded-full capitalize">
              {user?.role || 'member'}
            </span>
          </div>
        </div>

        {error && (
          <div className="mb-4 px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{error}</div>
        )}

        {editing ? (
          <form onSubmit={handleSaveProfile} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('members.name')}</label>
              <input
                type="text" required
                value={name}
                onChange={e => setName(e.target.value)}
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.username')}</label>
              <input
                type="text" required
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
            <div className="flex gap-3">
              <button
                type="submit" disabled={saving}
                className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition disabled:opacity-50"
              >
                {saving ? t('common.saving') : t('common.save')}
              </button>
              <button
                type="button" onClick={() => setEditing(false)} disabled={saving}
                className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-600 text-sm font-medium rounded-lg transition"
              >
                {t('common.cancel')}
              </button>
            </div>
          </form>
        ) : (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('members.name')}</label>
              <p className="text-sm text-gray-900">{profile?.name || '-'}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.username')}</label>
              <p className="text-sm text-gray-900">{profile?.username || '-'}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('profile.joined')}</label>
              <p className="text-sm text-gray-900">{profile?.created_at || '-'}</p>
            </div>
            <button
              onClick={() => setEditing(true)}
              className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-600 text-sm font-medium rounded-lg transition"
            >
              {t('profile.edit')}
            </button>
          </div>
        )}
      </div>

      {/* Security section */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('profile.security')}</h3>

        {passMsg && (
          <div className="mb-4 px-4 py-2 bg-green-50 border border-green-200 rounded-lg text-green-600 text-sm">{passMsg}</div>
        )}
        {passError && (
          <div className="mb-4 px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{passError}</div>
        )}

        {showPasswordForm ? (
          <form onSubmit={handleChangePassword} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('profile.current_password')}</label>
              <input
                type="password" required
                value={oldPass}
                onChange={e => setOldPass(e.target.value)}
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('profile.new_password')}</label>
              <input
                type="password" required minLength={6}
                value={newPass}
                onChange={e => setNewPass(e.target.value)}
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
                placeholder={t('profile.password_min')}
              />
            </div>
            <div className="flex gap-3">
              <button
                type="submit" disabled={changingPass}
                className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition disabled:opacity-50"
              >
                {changingPass ? t('common.saving') : t('profile.save_password')}
              </button>
              <button
                type="button" onClick={() => { setShowPasswordForm(false); setPassError(''); setPassMsg('') }}
                className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-600 text-sm font-medium rounded-lg transition"
              >
                {t('common.cancel')}
              </button>
            </div>
          </form>
        ) : (
          <div>
            <div className="flex items-center justify-between py-3 border-b border-gray-50 last:border-0">
              <div>
                <p className="text-sm text-gray-900">{t('profile.password')}</p>
                <p className="text-xs text-gray-400">{t('profile.password_desc')}</p>
              </div>
              <button
                onClick={() => setShowPasswordForm(true)}
                className="px-3 py-1.5 text-xs font-medium text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg transition"
              >
                {t('profile.change')}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
