import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import AuthLayout from '../components/AuthLayout'
import { forgotPassword } from '../api'

export default function ForgotPassword() {
  const [email, setEmail] = useState('')
  const [sent, setSent] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const { t } = useTranslation()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      await forgotPassword(email)
      setSent(true)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t('common.error_generic'))
    } finally {
      setLoading(false)
    }
  }

  if (sent) {
    return (
      <AuthLayout>
        <div className="text-center">
          <div className="w-12 h-12 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">📧</span>
          </div>
          <h2 className="text-lg font-semibold text-title mb-2">{t('auth.forgot_success')}</h2>
          <p className="text-secondary text-sm mb-4">{email}</p>
          <p className="text-muted text-xs mb-6">{t('auth.reset_link_sent')}</p>
          <Link
            to="/member/auth/login"
            className="inline-block w-full py-2.5 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-body text-sm font-medium rounded-lg transition text-center"
          >
            ← {t('auth.login')}
          </Link>
        </div>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <h2 className="text-lg font-semibold text-title mb-1">{t('auth.forgot_title')}</h2>
      <p className="text-secondary text-sm mb-6">{t('auth.forgot_desc')}</p>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Email */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.email')}</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm"
            placeholder="email@domain.com"
            required
          />
        </div>

        {error && (
          <p className="text-red-500 text-sm">{error}</p>
        )}

        {/* Submit */}
        <button
          type="submit"
          disabled={loading}
          className="w-full py-2.5 bg-gray-900 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium rounded-lg transition duration-200 active:scale-[0.98]"
        >
          {loading ? t('common.loading') : t('auth.send_reset_link')}
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
