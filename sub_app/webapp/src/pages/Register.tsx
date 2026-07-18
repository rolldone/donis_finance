import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import AuthLayout from '../components/AuthLayout'
import { registerMember } from '../api'

export default function Register() {
  const [form, setForm] = useState({ name: '', email: '', password: '', confirmPassword: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const { t } = useTranslation()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (form.password !== form.confirmPassword) {
      setError('Password tidak cocok!')
      return
    }

    setLoading(true)
    try {
      await registerMember(form.name, form.email, form.password)
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
          <h2 className="text-lg font-semibold text-title mb-2">{t('auth.register_success')}</h2>
          <p className="text-secondary text-sm mb-6">
            {t('auth.register')}
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
      <h2 className="text-lg font-semibold text-title mb-1">{t('auth.register')}</h2>
      <p className="text-secondary text-sm mb-6">{t('auth.register_title')}</p>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Name */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('members.name')}</label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm"
            placeholder={t('members.name')}
            required
          />
        </div>

        {/* Email */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.email')}</label>
          <input
            type="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm"
            placeholder="email@domain.com"
            required
          />
        </div>

        {/* Password */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.password')}</label>
          <div className="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm pr-10"
              placeholder={t('auth.password')}
              minLength={6}
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-secondary hover:text-title transition text-sm"
            >
              {showPassword ? '🙈' : '👁️'}
            </button>
          </div>
        </div>

        {/* Confirm Password */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.confirm_password')}</label>
          <input
            type={showPassword ? 'text' : 'password'}
            value={form.confirmPassword}
            onChange={(e) => setForm({ ...form, confirmPassword: e.target.value })}
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
          disabled={loading}
          className="w-full py-2.5 bg-gray-900 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium rounded-lg transition duration-200 active:scale-[0.98]"
        >
          {loading ? t('common.loading') : t('auth.register')}
        </button>
      </form>

      {/* Login link */}
      <p className="text-center text-sm text-muted mt-6">
        {t('auth.have_account')}{' '}
        <Link to="/member/auth/login" className="text-title hover:text-body font-medium transition">
          {t('auth.login')}
        </Link>
      </p>
    </AuthLayout>
  )
}
