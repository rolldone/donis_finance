import { useTranslation } from 'react-i18next'

export default function LanguageSwitcher() {
  const { i18n } = useTranslation()

  const toggleLang = () => {
    const next = i18n.language === 'id' ? 'en' : 'id'
    i18n.changeLanguage(next)
  }

  return (
    <button
      onClick={toggleLang}
      className="p-2 rounded-lg text-secondary hover-bg hover:text-title transition text-xs font-medium"
      title={i18n.language === 'id' ? 'Switch to English' : 'Ganti ke Indonesia'}
    >
      {i18n.language === 'id' ? 'EN' : 'ID'}
    </button>
  )
}
