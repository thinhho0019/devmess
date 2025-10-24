import { FiArrowRight, FiLock, FiZap, FiSmartphone, FiMessageCircle, FiUsers, FiShield } from "react-icons/fi";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";
import type { Variants } from "framer-motion";
 
import imgAvatar from "../assets/img.jpg";
import { LanguageDropdown } from "../components/dropdown/LanguageDropdown";

const fadeInUp: Variants = {
  hidden: { opacity: 0, y: 40 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: i * 0.15, duration: 0.6, ease: [0.22, 1, 0.36, 1] },
  }),
};

interface FeatureCardProps {
  icon: React.ReactNode;
  title: React.ReactNode;
  children?: React.ReactNode;
  index?: number;
}

const FeatureCard = ({ icon, title, children, index }: FeatureCardProps) => (
  <motion.div
    variants={fadeInUp}
    initial="hidden"
    whileInView="visible"
    viewport={{ once: true, amount: 0.3 }}
    custom={index}
    className="bg-white/5 dark:bg-gray-800/40 p-8 rounded-2xl border border-white/10 shadow-lg backdrop-blur-lg text-center hover:border-blue-400/40 transition-all"
  >
    <motion.div
      whileHover={{ rotate: 10, scale: 1.1 }}
      transition={{ type: "spring", stiffness: 200 }}
      className="inline-block p-4 bg-blue-500/10 rounded-full mb-4 ring-2 ring-blue-500/30"
    >
      {icon}
    </motion.div>
    <h4 className="text-xl font-semibold mb-2 text-white">{title}</h4>
    <p className="text-gray-400">{children}</p>
  </motion.div>
);

interface TestimonialCardProps {
  quote: string;
  name: string;
  role: string;
  avatar: string;
  index?: number;
}

const TestimonialCard = ({ quote, name, role, avatar, index }: TestimonialCardProps) => (
  <motion.div
    variants={fadeInUp}
    initial="hidden"
    whileInView="visible"
    viewport={{ once: true }}
    custom={index}
    className="bg-gray-800/40 p-6 rounded-xl border border-gray-700/60 hover:border-blue-500/40 transition-all"
  >
    <motion.p
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.3 }}
      className="text-gray-300 italic"
    >
      "{quote}"
    </motion.p>
    <div className="flex items-center mt-4">
      <img src={avatar} alt={name} className="w-12 h-12 rounded-full object-cover" />
      <div className="ml-4">
        <p className="font-semibold text-white">{name}</p>
        <p className="text-sm text-gray-400">{role}</p>
      </div>
    </div>
  </motion.div>
);

interface StepCardProps {
  step: number;
  icon: React.ReactNode;
  title: React.ReactNode;
  description: React.ReactNode;
  index?: number;
}

const StepCard = ({ step, icon, title, description, index }: StepCardProps) => (
  <motion.div
    variants={fadeInUp}
    initial="hidden"
    whileInView="visible"
    viewport={{ once: true, amount: 0.3 }}
    custom={index}
    className="relative bg-gradient-to-br from-gray-800/60 to-gray-900/60 p-8 rounded-2xl border border-gray-700/50 hover:border-blue-500/40 transition-all duration-300 backdrop-blur-sm group"
  >
    {/* Step number */}
    <motion.div
      whileHover={{ scale: 1.1, rotate: 5 }}
      transition={{ type: "spring", stiffness: 200 }}
      className="absolute -top-4 -left-4 w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold text-lg shadow-lg"
    >
      {step}
    </motion.div>

    {/* Icon */}
    <motion.div
      whileHover={{ scale: 1.1, y: -5 }}
      transition={{ type: "spring", stiffness: 200 }}
      className="inline-block p-4 bg-blue-500/10 rounded-full mb-6 ring-2 ring-blue-500/20 group-hover:ring-blue-400/40"
    >
      {icon}
    </motion.div>

    {/* Content */}
    <h4 className="text-xl font-semibold mb-3 text-white group-hover:text-blue-300 transition-colors">
      {title}
    </h4>
    <p className="text-gray-400 leading-relaxed">
      {description}
    </p>

    {/* Connecting line for desktop */}
    <div className="hidden lg:block absolute top-1/2 -right-8 w-16 h-px bg-gradient-to-r from-blue-500/50 to-transparent"></div>
  </motion.div>
);

export default function HomePage() {
  const { t } = useTranslation();
  

  return (
    <div className="chat-scroll bg-gray-900 text-gray-200 min-h-screen overflow-x-hidden relative">
      {/* Floating Gradient Lights */}
      <motion.div
        className="absolute w-[30rem] h-[30rem] bg-blue-500/20 rounded-full blur-3xl top-10 -left-20"
        animate={{ x: [0, 30, -30, 0], y: [0, 20, -20, 0] }}
        transition={{ duration: 10, repeat: Infinity }}
      />
      <motion.div
        className="absolute w-[40rem] h-[40rem] bg-purple-600/20 rounded-full blur-3xl bottom-10 -right-20"
        animate={{ x: [0, -30, 30, 0], y: [0, -20, 20, 0] }}
        transition={{ duration: 10, repeat: Infinity }}
      />

      {/* Header */}
      <header className="relative container mx-auto px-6 py-4 flex justify-between items-center z-10">
        <motion.h1
          whileHover={{ scale: 1.05 }}
          className="text-2xl font-bold text-white"
        >
          {t("header_title")}
        </motion.h1>

        <nav className="hidden md:flex items-center space-x-6 bg-white/5 backdrop-blur-md px-4 py-2 rounded-full border border-white/10">
          <a href="#features" className="hover:text-blue-400">{t("features")}</a>
          <a href="#how-it-works" className="hover:text-blue-400">{t("how_it_works")}</a>
          <a href="#testimonials" className="hover:text-blue-400">{t("testimonials")}</a>
          <LanguageDropdown />
        </nav>

        <Link to="/l" className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-5 rounded-full transition duration-300 shadow-lg shadow-blue-500/20">
          {t("launch_app")}
        </Link>
      </header>

      {/* Hero Section */}
      <main className="relative container mx-auto px-6 text-center pt-28 pb-20 z-0">
        <motion.h2
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-5xl md:text-7xl font-extrabold leading-tight mb-4 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-400"
        >
          {t("hero_title")}
        </motion.h2>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.6 }}
          className="text-lg md:text-xl text-gray-400 max-w-3xl mx-auto mb-10"
        >
          {t("hero_subtitle")}
        </motion.p>
        <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
          <Link to="/r">
            <button className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-bold py-4 px-8 rounded-full text-lg transition-all duration-300 inline-flex items-center shadow-2xl shadow-blue-500/30">
              {t("start_chatting")} <FiArrowRight className="ml-3" />
            </button>
          </Link>
        </motion.div>
      </main>

      {/* Features */}
      <section id="features" className="relative py-20 z-10">
        <div className="container mx-auto px-6 text-center mb-16">
          <h3 className="text-4xl font-bold text-white">{t("why_choose")}</h3>
          <p className="text-gray-400 mt-2 max-w-2xl mx-auto">{t("features_subtitle")}</p>
        </div>
        <div className="grid md:grid-cols-3 gap-8 px-6">
          <FeatureCard index={0} icon={<FiLock className="w-8 h-8 text-blue-400" />} title={t("feature_encryption_title")}>{t("feature_encryption_desc")}</FeatureCard>
          <FeatureCard index={1} icon={<FiZap className="w-8 h-8 text-blue-400" />} title={t("feature_fast_title")}>{t("feature_fast_desc")}</FeatureCard>
          <FeatureCard index={2} icon={<FiSmartphone className="w-8 h-8 text-blue-400" />} title={t("feature_modern_title")}>{t("feature_modern_desc")}</FeatureCard>
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="relative py-20 z-10">
        <div className="container mx-auto px-6">
          {/* Section Header */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="text-center mb-20"
          >
            <h3 className="text-4xl md:text-5xl font-bold text-white mb-4">
              {t("how_it_works")}
            </h3>
            <p className="text-gray-400 text-lg max-w-3xl mx-auto">
              {t("how_it_works_subtitle")}
            </p>
          </motion.div>

          {/* Steps Grid */}
          <div className="grid lg:grid-cols-3 gap-8 lg:gap-12 max-w-6xl mx-auto">
            <StepCard
              index={0}
              step={1}
              icon={<FiUsers className="w-8 h-8 text-blue-400" />}
              title={t("step1_title")}
              description={t("step1_desc")}
            />
            <StepCard
              index={1}
              step={2}
              icon={<FiMessageCircle className="w-8 h-8 text-blue-400" />}
              title={t("step2_title")}
              description={t("step2_desc")}
            />
            <StepCard
              index={2}
              step={3}
              icon={<FiShield className="w-8 h-8 text-blue-400" />}
              title={t("step3_title")}
              description={t("step3_desc")}
            />
          </div>

          {/* Demo Video/Animation Area */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.8, delay: 0.3 }}
            className="mt-20 text-center"
          >
            <div className="relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 rounded-3xl p-8 border border-gray-700/50 backdrop-blur-sm max-w-4xl mx-auto">
              {/* Mockup Chat Interface */}
              <div className="bg-gray-900/80 rounded-2xl p-6 border border-gray-700/30">
                <div className="flex items-center justify-between mb-4 pb-4 border-b border-gray-700/50">
                  <div className="flex items-center space-x-3">
                    <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                  <h4 className="text-white font-semibold">{t("chat_demo_title")}</h4>
                  <div className="w-16"></div>
                </div>
                
                {/* Chat Messages */}
                <div className="space-y-4">
                  <motion.div
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 1, duration: 0.5 }}
                    className="flex items-start space-x-3"
                  >
                    <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-bold">
                      A
                    </div>
                    <div className="bg-gray-700/50 rounded-lg px-4 py-2 max-w-xs">
                      <p className="text-white text-sm">{t("demo_message1")}</p>
                    </div>
                  </motion.div>
                  
                  <motion.div
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 1.5, duration: 0.5 }}
                    className="flex items-start space-x-3 justify-end"
                  >
                    <div className="bg-blue-600 rounded-lg px-4 py-2 max-w-xs">
                      <p className="text-white text-sm">{t("demo_message2")}</p>
                    </div>
                    <div className="w-8 h-8 bg-purple-500 rounded-full flex items-center justify-center text-white text-sm font-bold">
                      B
                    </div>
                  </motion.div>

                  <motion.div
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 2, duration: 0.5 }}
                    className="flex items-start space-x-3"
                  >
                    <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-bold">
                      A
                    </div>
                    <div className="bg-gray-700/50 rounded-lg px-4 py-2 max-w-xs">
                      <p className="text-white text-sm">{t("demo_message3")}</p>
                    </div>
                  </motion.div>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Testimonials */}
      <section id="testimonials" className="relative py-20 z-10">
        <div className="container mx-auto px-6 text-center mb-16">
          <h3 className="text-4xl font-bold text-white">{t("testimonials_title")}</h3>
          <p className="text-gray-400 mt-2">{t("testimonials_subtitle")}</p>
        </div>
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8 px-6">
          <TestimonialCard index={0} name="Alex Johnson" role="Developer" avatar={imgAvatar} quote={t("testimonial1_quote")} />
          <TestimonialCard index={1} name="Maria Garcia" role="Designer" avatar={imgAvatar} quote={t("testimonial2_quote")} />
          <TestimonialCard index={2} name="Sam Lee" role="Product Manager" avatar={imgAvatar} quote={t("testimonial3_quote")} />
        </div>
      </section>

      {/* Call To Action */}
      <section className="relative py-24 z-10 text-center">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true, amount: 0.5 }}
          transition={{ duration: 0.7 }}
          className="max-w-2xl mx-auto"
        >
          <h3 className="text-4xl font-bold text-white mb-6">{t("cta_title")}</h3>
          <p className="text-gray-400 mb-8">{t("cta_subtitle")}</p>
          <Link to="/r">
            <motion.button
              whileHover={{ scale: 1.1 }}
              className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-4 rounded-full font-bold text-lg shadow-lg shadow-blue-500/30"
            >
              {t("cta_button")}
            </motion.button>
          </Link>
        </motion.div>
      </section>
    </div>
  );
}
