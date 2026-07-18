import { Navigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../context/AuthContext'

export function MemberRoute({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation()
  const { isAuthenticated, isMember, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen bg-base flex items-center justify-center">
        <div className="text-muted text-sm">{t('common.loading')}</div>
      </div>
    )
  }

  if (!isAuthenticated || !isMember) {
    return <Navigate to="/member/auth/login" replace />
  }

  return children
}

export function AdminRoute({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation()
  const { isAuthenticated, isAdmin, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-gray-400 text-sm">{t('common.loading')}</div>
      </div>
    )
  }

  if (!isAuthenticated || !isAdmin) {
    return <Navigate to="/admin/auth/login" replace />
  }

  return children
}

export function GuestRoute({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation()
  const { isAuthenticated, isAdmin, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen bg-base flex items-center justify-center">
        <div className="text-muted text-sm">{t('common.loading')}</div>
      </div>
    )
  }

  // Redirect if already logged in
  if (isAuthenticated) {
    if (isAdmin) return <Navigate to="/admin" replace />
    return <Navigate to="/member" replace />
  }

  return children
}
