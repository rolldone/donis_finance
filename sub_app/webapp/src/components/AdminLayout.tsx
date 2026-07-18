import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTheme } from '../context/ThemeContext'
import { useTranslation } from 'react-i18next'
import LanguageSwitcher from './LanguageSwitcher'

const NAV_ITEMS = [
  { path: '/admin', labelKey: 'nav.dashboard', icon: '📊' },
  { path: '/admin/members', labelKey: 'nav.members', icon: '👥' },
  { path: '/admin/categories', labelKey: 'nav.categories', icon: '🏷️' },
  { path: '/admin/transactions', labelKey: 'nav.transactions', icon: '💳' },
  { path: '/admin/budget', labelKey: 'nav.budget', icon: '🎯' },
  { path: '/admin/settings', labelKey: 'nav.settings', icon: '⚙️' },
]

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const location = useLocation()
  const navigate = useNavigate()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { user, logout } = useAuth()
  const { theme, toggle: toggleTheme } = useTheme()
  const { t } = useTranslation()

  const handleLogout = () => {
    logout()
    navigate('/admin/auth/login')
  }

  return (
    <div className="min-h-screen bg-base flex lg:h-screen lg:overflow-hidden">
      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-64 bg-sidebar text-title
        transform transition-transform duration-200 ease-in-out
        lg:sticky lg:top-0 lg:h-screen lg:translate-x-0 lg:z-auto
        flex flex-col
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
      `}>
        {/* Logo */}
        <div className="h-20 flex items-center px-6 border-b border-light shrink-0">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gray-900 dark:bg-white rounded-lg flex items-center justify-center">
              <span className="text-white dark:text-gray-900 text-sm font-bold">D</span>
            </div>
            <div>
              <span className="font-semibold text-sm text-title">donis_finance</span>
              <p className="text-xs text-muted">Admin Panel</p>
            </div>
          </div>
        </div>

        {/* Nav */}
        <nav className="p-4 space-y-1 overflow-y-auto flex-1">
          {NAV_ITEMS.map((item) => {
            const isActive = location.pathname === item.path ||
              (item.path !== '/admin' && location.pathname.startsWith(item.path))
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
        <div className="p-4 border-t border-light shrink-0">
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-sidebar-text hover:bg-red-50 dark:hover:bg-red-500/10 hover:text-red-500 dark:hover:text-red-400 transition w-full"
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
      <div className="flex-1 flex flex-col min-w-0 overflow-y-auto">
        {/* Topbar */}
        <header className="h-20 bg-topbar border-b border-default flex items-center justify-between px-4 sm:px-6 sticky top-0 z-30 shrink-0">
          <button
            onClick={() => setSidebarOpen(true)}
            className="lg:hidden p-2 -ml-2 text-secondary hover:text-title"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>

          <div className="hidden lg:block">
            <h1 className="text-xl font-semibold text-title">
              {t(NAV_ITEMS.find(n => location.pathname === n.path)?.labelKey || 'nav.dashboard')}
            </h1>
          </div>

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

            <div className="text-right hidden sm:block">
              <p className="text-sm font-medium text-title">{user?.username || 'Admin'}</p>
              <p className="text-xs text-muted">{user?.role || 'admin'}</p>
            </div>
            <div className="w-9 h-9 bg-gray-900 dark:bg-white/10 rounded-full flex items-center justify-center text-white dark:text-title font-semibold text-sm">
              {(user?.username || 'A').charAt(0).toUpperCase()}
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
