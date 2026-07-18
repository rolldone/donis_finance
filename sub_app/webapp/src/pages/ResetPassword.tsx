import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import AuthLayout from '../components/AuthLayout'
import { resetPassword } from '../api'

export default function ResetPassword() {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') || ''

  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const { t } = useTranslation()

  useEffect(() => {
    if (!token) {
      setError(t('auth.reset_link_sent'))
    }
  }, [token, t])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (password.length < 6) {
      setError('Password minimal 6 karakter')
      return
    }
    if (password !== confirm) {
      setError('Password dan konfirmasi tidak cocok')
      return
    }

    setLoading(true)
    try {
      await resetPassword(token, password)
      setSuccess(true)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t('common.error_generic'))
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <AuthLayout>
        <div className="text-center">
          <div className="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">✅</span>
          </div>
          <h2 className="text-lg font-semibold text-title mb-2">{t('auth.reset_success')}</h2>
          <p className="text-secondary text-sm mb-6">
            {t('auth.reset_success')}
          </p>
          <Link
            to="/member/auth/login"
            className="inline-block w-full py-2.5 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition text-center"
          >
            ← {t('auth.login')}
          </Link>
        </div>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <h2 className="text-lg font-semibold text-title mb-1">{t('auth.reset_title')}</h2>
      <p className="text-secondary text-sm mb-6">{t('auth.reset_desc')}</p>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Password */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.new_password')}</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm"
            placeholder={t('auth.new_password')}
            required
            minLength={6}
          />
        </div>

        {/* Confirm Password */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.confirm_password')}</label>
          <input
            type="password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm"
            placeholder={t('auth.confirm_password')}
            required
          />
        </div>

        {error && (
          <p className="text-red-500 text-sm">{error}</p>
        )}

        {/* Submit */}
        <button
          type="submit"
          disabled={loading || !token}
          className="w-full py-2.5 bg-gray-900 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium rounded-lg transition duration-200 active:scale-[0.98]"
        >
          {loading ? t('common.loading') : t('auth.reset_password')}
        </button>
      </form>

      {/* Back to login */}
      <p className="text-center text-sm text-muted mt-6">
        <Link to="/member/auth/login" className="text-title hover:text-body font-medium transition">
          ← {t('auth.login')}
        </Link>
      </p>
    </AuthLayout>
  )
}
