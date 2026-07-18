import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getAdminProfile, updateAdminProfile, changeAdminPassword, getMembers, exportTransactionsCSV, getAdminSMTPConfig, saveAdminSMTPConfig } from '../../api'
import type { SMTPConfig } from '../../api'
import Skeleton from '../../components/Skeleton'

const INIT_SMTP: SMTPConfig = {
  host: '', port: '', user: '', pass: '', from_email: '', from_name: '',
  use_tls: false, use_starttls: false, skip_verify: false,
}

export default function Settings() {
  const { t } = useTranslation()
  const [profile, setProfile] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  // Edit
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  // SMTP
  const [smtp, setSmtp] = useState<SMTPConfig>(INIT_SMTP)
  const [notifEmail, setNotifEmail] = useState('')
  const [envSmtp, setEnvSmtp] = useState<SMTPConfig | null>(null)
  const [smtpOverridden, setSmtpOverridden] = useState(false)
  const [showEnvSmtp, setShowEnvSmtp] = useState(false)
  const [smtpSaving, setSmtpSaving] = useState(false)
  const [smtpMsg, setSmtpMsg] = useState('')
  const [smtpErr, setSmtpErr] = useState('')

  // Password
  const [showPass, setShowPass] = useState(false)
  const [oldPass, setOldPass] = useState('')
  const [newPass, setNewPass] = useState('')
  const [passMsg, setPassMsg] = useState('')
  const [passErr, setPassErr] = useState('')
  const [changingPass, setChangingPass] = useState(false)

  // Export
  const [members, setMembers] = useState<any[]>([])
  const [exportMemberId, setExportMemberId] = useState('')
  const [exporting, setExporting] = useState(false)
  const [exportErr, setExportErr] = useState('')

  const fetchProfile = async () => {
    try {
      setLoading(true)
      const data = await getAdminProfile()
      setProfile(data)
      setUsername(data.username)
      setEmail(data.email || '')
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProfile()
    getMembers().then(res => setMembers(res.members || [])).catch(() => {})
    getAdminSMTPConfig().then(res => {
      if (res.smtp) setSmtp({ ...INIT_SMTP, ...res.smtp })
      if (res.env_smtp) setEnvSmtp(res.env_smtp)
      setSmtpOverridden(res.override)
      setNotifEmail(res.notif_email || '')
    }).catch(() => {})
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  const now = new Date()
  const handleExportAll = async () => {
    setExporting(true)
    setExportErr('')
    try {
      await exportTransactionsCSV({ month: now.getMonth() + 1, year: now.getFullYear() })
    } catch (err: any) {
      setExportErr(err.message)
    } finally {
      setExporting(false)
    }
  }

  const handleExportMember = async () => {
    if (!exportMemberId) return
    setExporting(true)
    setExportErr('')
    try {
      const member = members.find((m: any) => m.id === exportMemberId)
      await exportTransactionsCSV({
        month: now.getMonth() + 1,
        year: now.getFullYear(),
        member_id: exportMemberId,
        member_name: member?.name || member?.username || '',
      })
    } catch (err: any) {
      setExportErr(err.message)
    } finally {
      setExporting(false)
    }
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError('')
    setSuccess('')
    try {
      const updated = await updateAdminProfile({ username, email })
      setProfile(updated)
      setSuccess(t('settings.save_success'))
    } catch (err: any) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    setChangingPass(true)
    setPassErr('')
    setPassMsg('')
    try {
      const res = await changeAdminPassword(oldPass, newPass)
      setPassMsg(res.message)
      setOldPass('')
      setNewPass('')
      setShowPass(false)
    } catch (err: any) {
      setPassErr(err.message)
    } finally {
      setChangingPass(false)
    }
  }

  const handleSaveSMTP = async (e: React.FormEvent) => {
    e.preventDefault()
    setSmtpSaving(true)
    setSmtpErr('')
    setSmtpMsg('')
    try {
      const res = await saveAdminSMTPConfig({ ...smtp, notif_email: notifEmail })
      setSmtpMsg(res.message)
      if (res.smtp) setSmtp({ ...smtp, ...res.smtp })
      if (res.env_smtp) setEnvSmtp(res.env_smtp)
      setSmtpOverridden(res.override)
      setNotifEmail(res.notif_email || '')
    } catch (err: any) {
      setSmtpErr(err.message)
    } finally {
      setSmtpSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl space-y-6">
        <div>
          <Skeleton.Base className="h-7 w-32 mb-2" />
          <Skeleton.Base className="h-4 w-48" />
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <Skeleton.Base className="h-4 w-28 mb-6" />
          <div className="space-y-4">
            <div><Skeleton.Base className="h-3 w-16 mb-1.5" /><Skeleton.Base className="h-10 w-full rounded-lg" /></div>
            <div><Skeleton.Base className="h-3 w-12 mb-1.5" /><Skeleton.Base className="h-10 w-full rounded-lg" /></div>
            <div><Skeleton.Base className="h-3 w-20 mb-1.5" /><Skeleton.Base className="h-5 w-36" /></div>
            <Skeleton.Base className="h-9 w-28 rounded-lg" />
          </div>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <Skeleton.Base className="h-4 w-24 mb-6" />
          <Skeleton.Base className="h-3 w-40 mb-4" />
          <div className="space-y-4">
            <div><Skeleton.Base className="h-3 w-28 mb-1.5" /><Skeleton.Base className="h-10 w-full rounded-lg" /></div>
            <div><Skeleton.Base className="h-3 w-24 mb-1.5" /><Skeleton.Base className="h-10 w-full rounded-lg" /></div>
            <Skeleton.Base className="h-9 w-32 rounded-lg" />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-semibold text-gray-900">{t('settings.title')}</h2>
        <p className="text-gray-400 text-sm mt-1">{t('settings.subtitle')}</p>
      </div>

      {error && (
        <div className="px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{error}</div>
      )}
      {success && (
        <div className="px-4 py-2 bg-green-50 border border-green-200 rounded-lg text-green-600 text-sm">{success}</div>
      )}

      {/* Profile */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h3 className="text-sm font-semibold text-gray-900 mb-4">{t('settings.admin_profile')}</h3>
        <form onSubmit={handleSave} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.username')}</label>
            <input
              type="text" required
              value={username}
              onChange={e => setUsername(e.target.value)}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.email')}</label>
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              placeholder="admin@donis.finance"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('members.joined')}</label>
            <p className="text-sm text-gray-900">{profile?.created_at || '-'}</p>
          </div>
          <button
            type="submit" disabled={saving}
            className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition disabled:opacity-50"
          >
            {saving ? t('common.saving') : t('settings.save_profile')}
          </button>
        </form>
      </div>

      {/* Password */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h3 className="text-sm font-semibold text-gray-900 mb-4">🔒 {t('profile.security')}</h3>

        {passMsg && (
          <div className="mb-4 px-4 py-2 bg-green-50 border border-green-200 rounded-lg text-green-600 text-sm">{passMsg}</div>
        )}
        {passErr && (
          <div className="mb-4 px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{passErr}</div>
        )}

        {showPass ? (
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
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('auth.new_password')}</label>
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
                type="button" onClick={() => { setShowPass(false); setPassErr(''); setPassMsg('') }}
                className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-600 text-sm font-medium rounded-lg transition"
              >
                {t('common.cancel')}
              </button>
            </div>
          </form>
        ) : (
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-900">{t('auth.password')}</p>
              <p className="text-xs text-gray-400">{t('profile.password_desc')}</p>
            </div>
            <button
              type="button"
              onClick={() => setShowPass(true)}
              className="px-3 py-1.5 text-xs font-medium text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg transition"
            >
              {t('profile.change')}
            </button>
          </div>
        )}
      </div>

      {/* SMTP Settings */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-semibold text-gray-900">📧 {t('settings.smtp_title')}</h3>
          {envSmtp && (
            <button
              type="button"
              onClick={() => setShowEnvSmtp(!showEnvSmtp)}
              className="text-xs text-gray-500 hover:text-gray-700 underline"
            >
              {showEnvSmtp ? 'Hide' : 'View'} .env defaults
            </button>
          )}
        </div>

        {smtpOverridden && (
          <div className="mb-4 px-4 py-3 bg-amber-50 border border-amber-200 rounded-lg text-amber-700 text-sm flex items-start gap-2">
            <span className="mt-0.5">⚡</span>
            <div>
              <strong>SMTP Override Active</strong> — The settings below are being used instead of the .env defaults.
              To revert to .env values, clear the fields and save.
            </div>
          </div>
        )}

        {/* .env fallback viewer */}
        {showEnvSmtp && envSmtp && (
          <div className="mb-4 px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-xs font-mono text-gray-600 space-y-1">
            <div className="font-semibold text-gray-700 mb-1.5">📄 .env SMTP defaults:</div>
            <div><span className="text-gray-500">SMTP_HOST</span> = {envSmtp.host || <span className="italic text-gray-400">(empty)</span>}</div>
            <div><span className="text-gray-500">SMTP_PORT</span> = {envSmtp.port || <span className="italic text-gray-400">(empty)</span>}</div>
            <div><span className="text-gray-500">SMTP_USER</span> = {envSmtp.user || <span className="italic text-gray-400">(empty)</span>}</div>
            <div><span className="text-gray-500">SMTP_PASS</span> = {envSmtp.pass ? '********' : <span className="italic text-gray-400">(empty)</span>}</div>
            <div><span className="text-gray-500">SMTP_FROM_EMAIL</span> = {envSmtp.from_email || <span className="italic text-gray-400">(empty)</span>}</div>
            <div><span className="text-gray-500">SMTP_FROM_NAME</span> = {envSmtp.from_name || <span className="italic text-gray-400">(empty)</span>}</div>
            <div><span className="text-gray-500">SMTP_USE_TLS</span> = {String(envSmtp.use_tls)}</div>
            <div><span className="text-gray-500">SMTP_STARTTLS</span> = {String(envSmtp.use_starttls)}</div>
          </div>
        )}

        {smtpMsg && (
          <div className="mb-4 px-4 py-2 bg-green-50 border border-green-200 rounded-lg text-green-600 text-sm">{smtpMsg}</div>
        )}
        {smtpErr && (
          <div className="mb-4 px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{smtpErr}</div>
        )}

        <form onSubmit={handleSaveSMTP} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('settings.smtp_host')}</label>
            <input
              type="text" required
              value={smtp.host}
              onChange={e => setSmtp({ ...smtp, host: e.target.value })}
              placeholder="smtp.example.com"
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('settings.smtp_port')}</label>
              <input
                type="text" required
                value={smtp.port}
                onChange={e => setSmtp({ ...smtp, port: e.target.value })}
                placeholder="587"
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('settings.smtp_from')}</label>
              <input
                type="email"
                value={smtp.from_email}
                onChange={e => setSmtp({ ...smtp, from_email: e.target.value })}
                placeholder="noreply@donis.finance"
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">Sender Name</label>
              <input
                type="text"
                value={smtp.from_name}
                onChange={e => setSmtp({ ...smtp, from_name: e.target.value })}
                placeholder="donis_finance"
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('settings.smtp_user')}</label>
              <input
                type="text"
                value={smtp.user}
                onChange={e => setSmtp({ ...smtp, user: e.target.value })}
                placeholder="username"
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1.5">{t('settings.smtp_pass')}</label>
              <input
                type="password"
                value={smtp.pass}
                onChange={e => setSmtp({ ...smtp, pass: e.target.value })}
                placeholder="••••••••"
                className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
              />
            </div>
          </div>

          {/* Notification Email */}
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1.5">
              📬 Notification Email <span className="text-gray-400 font-normal">(admin notified on new member registration)</span>
            </label>
            <input
              type="email"
              value={notifEmail}
              onChange={e => setNotifEmail(e.target.value)}
              placeholder="system@rolldev.my.id"
              className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300"
            />
          </div>

          <button
            type="submit" disabled={smtpSaving}
            className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition disabled:opacity-50"
          >
            {smtpSaving ? t('common.saving') : t('settings.save_smtp')}
          </button>
        </form>
      </div>

      {/* Backup CSV */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h3 className="text-sm font-semibold text-gray-900 mb-4">💾 {t('settings.export_csv')}</h3>
        <p className="text-sm text-gray-400 mb-4">{t('settings.export_desc')}</p>

        {exportErr && (
          <div className="mb-4 px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">{exportErr}</div>
        )}

        <div className="flex flex-col gap-4">
          <button
            onClick={handleExportAll}
            disabled={exporting}
            className="px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition disabled:opacity-50"
          >
            {exporting ? t('common.loading') : t('settings.export_all')}
          </button>

          <div className="flex items-center gap-3">
            <select
              value={exportMemberId}
              onChange={e => setExportMemberId(e.target.value)}
              className="flex-1 px-3 py-2 text-sm bg-gray-50 border border-gray-200 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-900/10"
            >
              <option value="">{t('budget.select_member')}</option>
              {members.map((m: any) => (
                <option key={m.id} value={m.id}>{m.name || m.username}</option>
              ))}
            </select>
            <button
              onClick={handleExportMember}
              disabled={exporting || !exportMemberId}
              className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium rounded-lg transition disabled:opacity-50"
            >
              {t('settings.export_per_member')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
