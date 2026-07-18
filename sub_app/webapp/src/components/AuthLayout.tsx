import { useTranslation } from 'react-i18next'

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation()
  return (
    <div className="min-h-screen bg-base flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        {/* Logo */}
        <div className="text-center mb-10">
          <div className="w-12 h-12 bg-gray-900 dark:bg-white rounded-xl flex items-center justify-center mx-auto mb-4">
            <span className="text-white dark:text-gray-900 text-xl font-bold">D</span>
          </div>
          <h1 className="text-2xl font-semibold text-title tracking-tight">
            {t('app.name')}
          </h1>
          <p className="text-secondary text-sm mt-1">{t('app.subtitle')}</p>
        </div>

        {/* Card */}
        <div className="bg-surface rounded-xl border border-default p-8">
          {children}
        </div>

        {/* Footer */}
        <p className="text-center text-muted text-xs mt-8">
          {t('footer.copyright')}
        </p>
      </div>
    </div>
  )
}
