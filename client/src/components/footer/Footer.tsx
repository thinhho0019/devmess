import { FiMail, FiGithub, FiHeart } from "react-icons/fi";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";

export default function Footer() {
  const { t, ready } = useTranslation();

  // Debug log
  console.log('Footer i18n ready:', ready);
  console.log('Footer translation test:', t("footer_description", "Default description"));

  return (
    <footer className="relative z-10 border-t border-gray-800/50 bg-gradient-to-b from-gray-900/50 to-black/80 backdrop-blur-sm">
      <div className="container mx-auto px-6 py-12">
        {/* Main Footer Content */}
        <div className="grid lg:grid-cols-3 gap-8 mb-8">
          {/* Brand Section */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="text-center lg:text-left"
          >
            <h3 className="text-2xl font-bold mb-4 bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              {t("header_title")}
            </h3>
            <p className="text-gray-400 text-sm leading-relaxed max-w-md">
              {t("footer_description", "Modern, secure, and fast chat application. Connect with friends easily and safely.")}
            </p>
          </motion.div>

          {/* Quick Links */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-center"
          >
            <h4 className="text-lg font-semibold text-white mb-4">{t("footer_quick_links", "Quick Links")}</h4>
            <nav className="space-y-2">
              <a href="#features" className="block text-gray-400 hover:text-blue-400 transition-colors">
                {t("features", "Features")}
              </a>
              <a href="#how-it-works" className="block text-gray-400 hover:text-blue-400 transition-colors">
                {t("how_it_works", "How It Works")}
              </a>
              <a href="#testimonials" className="block text-gray-400 hover:text-blue-400 transition-colors">
                {t("testimonials", "Testimonials")}
              </a>
              <Link to="/l" className="block text-gray-400 hover:text-blue-400 transition-colors">
                {t("launch_app", "Launch App")}
              </Link>
            </nav>
          </motion.div>

          {/* Contact & Social */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="text-center lg:text-right"
          >
            <h4 className="text-lg font-semibold text-white mb-4">{t("footer_contact", "Contact")}</h4>
            <div className="space-y-3">
              <motion.a
                href="mailto:hoxuanthinh68@gmail.com"
                whileHover={{ scale: 1.05 }}
                className="flex items-center justify-center lg:justify-end space-x-2 text-gray-400 hover:text-blue-400 transition-colors group"
              >
                <FiMail className="w-4 h-4 group-hover:animate-bounce" />
                <span className="text-sm">hoxuanthinh68@gmail.com</span>
              </motion.a>
              <motion.a
                href="https://github.com/thinhho0019"
                target="_blank"
                rel="noopener noreferrer"
                whileHover={{ scale: 1.05 }}
                className="flex items-center justify-center lg:justify-end space-x-2 text-gray-400 hover:text-blue-400 transition-colors group"
              >
                <FiGithub className="w-4 h-4 group-hover:animate-pulse" />
                <span className="text-sm">github.com/thinhho0019</span>
              </motion.a>
            </div>
          </motion.div>
        </div>

        {/* Divider */}
        <div className="border-t border-gray-800/50 pt-8">
          {/* Bottom Footer */}
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0"
          >
            {/* Copyright */}
            <div className="flex items-center space-x-2 text-gray-400 text-sm">
              <span>{t("footer_copyright_prefix", "© 2025 Made with")}</span>
              <motion.span
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ duration: 1, repeat: Infinity }}
              >
                <FiHeart className="w-4 h-4 text-red-500" />
              </motion.span>
              <span>{t("footer_copyright_by", "by")}</span>
              <motion.span
                whileHover={{ scale: 1.1 }}
                className="font-semibold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent"
              >
                DevTh
              </motion.span>
            </div>

            {/* Tech Stack */}
            <div className="flex items-center space-x-4 text-xs text-gray-500">
              <span>{t("footer_built_with", "Built with React + TypeScript + Vite")}</span>
            </div>
          </motion.div>
        </div>
      </div>

      {/* Floating Particles Effect */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        {[...Array(6)].map((_, i) => (
          <motion.div
            key={i}
            className="absolute w-2 h-2 bg-blue-400/20 rounded-full"
            animate={{
              x: [0, 100, 0],
              y: [0, -100, 0],
              opacity: [0, 1, 0],
            }}
            transition={{
              duration: 8 + i * 2,
              repeat: Infinity,
              delay: i * 1.5,
            }}
            style={{
              left: `${10 + i * 15}%`,
              bottom: `${10 + (i % 3) * 10}%`,
            }}
          />
        ))}
      </div>
    </footer>
  );
}