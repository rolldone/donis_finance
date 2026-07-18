import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTheme } from '../context/ThemeContext'
import { useTranslation } from 'react-i18next'
import LanguageSwitcher from './LanguageSwitcher'

const NAV_ITEMS = [
  { path: '/member', labelKey: 'nav.dashboard', icon: '📊' },
  { path: '/member/transactions', labelKey: 'nav.transactions', icon: '💳' },
  { path: '/member/budget', labelKey: 'nav.budget', icon: '🎯' },
  { path: '/member/profile', labelKey: 'nav.profile', icon: '👤' },
]

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const location = useLocation()
  const navigate = useNavigate()
  const { user, logout } = useAuth()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const { theme, toggle: toggleTheme } = useTheme()
  const { t } = useTranslation()

  const handleLogout = () => {
    logout()
    navigate('/member/auth/login')
  }

  return (
    <div className="min-h-screen bg-base flex">
      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-64 bg-sidebar border-r border-default
        transform transition-transform duration-200 ease-in-out
        lg:translate-x-0 lg:static lg:z-auto
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
      `}>
        {/* Logo */}
        <div className="h-16 flex items-center px-6 border-b border-light">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gray-900 dark:bg-white rounded-lg flex items-center justify-center">
              <span className="text-white dark:text-gray-900 text-sm font-bold">D</span>
            </div>
            <span className="font-semibold text-title">donis_finance</span>
          </div>
        </div>

        {/* Nav */}
        <nav className="p-4 space-y-1">
          {NAV_ITEMS.map((item) => {
            const isActive = location.pathname === item.path ||
              (item.path !== '/member' && location.pathname.startsWith(item.path))
            return (
              <Link
                key={item.path}
                to={item.path}
                onClick={() => setSidebarOpen(false)}
                className={`
                  flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition
                  ${isActive
                    ? 'bg-sidebar-active-bg text-sidebar-active-text'
                    : 'text-sidebar-text hover:bg-sidebar-active-bg/50 hover:text-sidebar-active-text'
                  }
                `}
              >
                <span className="text-base">{item.icon}</span>
                {t(item.labelKey)}
              </Link>
            )
          })}
        </nav>

        {/* Logout */}
        <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-light">
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-sidebar-text hover:bg-red-50 dark:hover:bg-red-500/10 hover:text-red-600 transition w-full"
          >
            <span className="text-base">🚪</span>
            {t('auth.logout')}
          </button>
        </div>
      </aside>

      {/* Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 lg:hidden"
          style={{ backgroundColor: 'var(--overlay)' }}
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Topbar */}
        <header className="h-16 bg-topbar border-b border-default flex items-center justify-between px-4 sm:px-6 sticky top-0 z-30">
          {/* Mobile menu button */}
          <button
            onClick={() => setSidebarOpen(true)}
            className="lg:hidden p-2 -ml-2 text-secondary hover:text-title"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>

          {/* Page title (hidden on mobile) */}
          <div className="hidden lg:block">
            <h1 className="text-lg font-semibold text-title">
              {t(NAV_ITEMS.find(n => location.pathname === n.path)?.labelKey || 'nav.dashboard')}
            </h1>
          </div>

          {/* Right side */}
          <div className="flex items-center gap-2 sm:gap-3">
            {/* Language switcher */}
            <LanguageSwitcher />
            {/* Theme toggle */}
            <button
              onClick={toggleTheme}
              className="p-2 rounded-lg text-secondary hover-bg hover:text-title transition"
              title={theme === 'dark' ? 'Mode Terang' : 'Mode Gelap'}
            >
              {theme === 'dark' ? (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m8.66-14.66l-.71.71M6.05 6.05l-.71-.71M21 12h1M3 12H2m15.66 7.66l-.71-.71M6.05 17.95l-.71.71M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
              )}
            </button>

            {/* User info */}
            <div className="text-right hidden sm:block">
              <p className="text-sm font-medium text-title">{user?.username || 'User'}</p>
              <p className="text-xs text-muted">{user?.role || 'member'}</p>
            </div>
            <div className="w-9 h-9 bg-gray-100 dark:bg-white/10 rounded-full flex items-center justify-center text-title font-semibold text-sm">
              {(user?.username || 'U').charAt(0).toUpperCase()}
            </div>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 p-4 sm:p-6">
          {children}
        </main>
      </div>
    </div>
  )
}
