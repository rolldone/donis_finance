import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import AuthLayout from '../components/AuthLayout'
import { useAuth } from '../context/AuthContext'

export default function Login() {
  const [form, setForm] = useState({ username: '', password: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { loginMember } = useAuth()
  const navigate = useNavigate()
  const { t } = useTranslation()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await loginMember(form.username, form.password)
      navigate('/member')
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout>
      <h2 className="text-lg font-semibold text-title mb-1">{t('auth.login')}</h2>
      <p className="text-secondary text-sm mb-6">{t('auth.login_title')}</p>

      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">
            {error}
          </div>
        )}

        {/* Username */}
        <div>
          <label className="block text-sm font-medium text-secondary mb-1.5">{t('auth.username')}</label>
          <input
            type="text"
            value={form.username}
            onChange={(e) => setForm({ ...form, username: e.target.value })}
            className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm disabled:opacity-50"
            placeholder={t('auth.username')}
            disabled={loading}
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
              className="w-full px-3 py-2 bg-input border border-default rounded-lg text-title placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm pr-10 disabled:opacity-50"
              placeholder={t('auth.password')}
              disabled={loading}
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

        {/* Forgot */}
        <div className="flex items-center justify-end text-sm">
          <Link to="/member/auth/forgot-password" className="text-muted hover:text-title transition">
            {t('auth.forgot_password')}?
          </Link>
        </div>

        {/* Submit */}
        <button
          type="submit"
          disabled={loading}
          className="w-full py-2.5 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition duration-200 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? t('common.loading') : t('auth.login')}
        </button>
      </form>

      {/* Register link */}
      <p className="text-center text-sm text-muted mt-6">
        {t('auth.no_account')}{' '}
        <Link to="/member/auth/register" className="text-gray-900 hover:text-gray-700 font-medium transition">
          {t('auth.register')}
        </Link>
      </p>
    </AuthLayout>
  )
}
