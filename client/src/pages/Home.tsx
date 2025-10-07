import { FiArrowRight, FiLock, FiZap, FiSmartphone, FiGlobe, FiChevronDown } from "react-icons/fi";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useState } from "react";

export default function HomePage() {
  const { t, i18n } = useTranslation();
  const [langDropdownOpen, setLangDropdownOpen] = useState(false);

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    setLangDropdownOpen(false);
  };

  return (
    <div className="chat-scroll bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-200 min-h-screen">
      {/* Header */}
      <header className="container mx-auto px-6 py-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold text-blue-600 dark:text-blue-400">{t('header_title')}</h1>
        <nav className="flex items-center space-x-6">
          <a href="#features" className="hover:text-blue-500">{t('features')}</a>
          <a href="#about" className="hover:text-blue-500">{t('about')}</a>
          
          {/* Language Selector */}
          <div className="relative">
            <button 
              onClick={() => setLangDropdownOpen(!langDropdownOpen)}
              className="flex items-center hover:text-blue-500"
            >
              <FiGlobe className="mr-1" />
              <span>{i18n.language.toUpperCase()}</span>
              <FiChevronDown className={`ml-1 transition-transform ${langDropdownOpen ? 'rotate-180' : ''}`} />
            </button>
            {langDropdownOpen && (
              <div className="absolute right-0 mt-2 py-2 w-28 bg-white dark:bg-gray-800 rounded-md shadow-xl z-20">
                <button onClick={() => changeLanguage('en')} className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-blue-500 hover:text-white dark:hover:bg-blue-500 w-full text-left">
                  English
                </button>
                <button onClick={() => changeLanguage('vi')} className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-blue-500 hover:text-white dark:hover:bg-blue-500 w-full text-left">
                  Tiếng Việt
                </button>
              </div>
            )}
          </div>

          <Link to="/t" className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full transition duration-300">
            {t('launch_app')}
          </Link>
        </nav>
      </header>

      {/* Hero Section */}
      <main className="container mx-auto px-6 text-center pt-24 pb-16">
        <h2 className="text-5xl md:text-6xl font-extrabold leading-tight mb-4">
          {t('hero_title')}
        </h2>
        <p className="text-lg md:text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto mb-8">
          {t('hero_subtitle')}
        </p>
        <Link to="/t">
          <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-4 px-8 rounded-full text-lg transition duration-300 inline-flex items-center">
            {t('start_chatting')} <FiArrowRight className="ml-2" />
          </button>
        </Link>
      </main>

      {/* Features Section */}
      <section id="features" className="py-20 bg-white dark:bg-gray-800">
        <div className="container mx-auto px-6">
          <div className="text-center mb-12">
            <h3 className="text-4xl font-bold">{t('why_choose')}</h3>
            <p className="text-gray-600 dark:text-gray-400 mt-2">{t('features_subtitle')}</p>
          </div>
          <div className="grid md:grid-cols-3 gap-12">
            <div className="text-center p-6">
              <div className="inline-block p-4 bg-blue-100 dark:bg-blue-900 rounded-full mb-4">
                <FiLock className="w-8 h-8 text-blue-600 dark:text-blue-400" />
              </div>
              <h4 className="text-xl font-semibold mb-2">{t('feature_encryption_title')}</h4>
              <p className="text-gray-600 dark:text-gray-400">
                {t('feature_encryption_desc')}
              </p>
            </div>
            <div className="text-center p-6">
              <div className="inline-block p-4 bg-green-100 dark:bg-green-900 rounded-full mb-4">
                <FiZap className="w-8 h-8 text-green-600 dark:text-green-400" />
              </div>
              <h4 className="text-xl font-semibold mb-2">{t('feature_fast_title')}</h4>
              <p className="text-gray-600 dark:text-gray-400">
                {t('feature_fast_desc')}
              </p>
            </div>
            <div className="text-center p-6">
              <div className="inline-block p-4 bg-purple-100 dark:bg-purple-900 rounded-full mb-4">
                <FiSmartphone className="w-8 h-8 text-purple-600 dark:text-purple-400" />
              </div>
              <h4 className="text-xl font-semibold mb-2">{t('feature_modern_title')}</h4>
              <p className="text-gray-600 dark:text-gray-400">
                {t('feature_modern_desc')}
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="container mx-auto px-6 py-8 text-center text-gray-500">
        <p dangerouslySetInnerHTML={{ __html: t('footer_copyright') }} />
        <p className="mt-1">{t('footer_tagline')}</p>
      </footer>
    </div>
  );
}
