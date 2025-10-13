import { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { FiGlobe, FiChevronDown } from "react-icons/fi";
import { useTranslation } from "react-i18next";

export const LanguageDropdown = ({ openOnHover = false }: { openOnHover?: boolean }) => {
  const { i18n } = useTranslation();
  const [open, setOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement | null>(null);

  const changeLanguage = (lng: string) => {
    console.log("Changing language to:", lng);
    i18n.changeLanguage(lng);
    setOpen(false);
    // focus back to button for accessibility
    buttonRef.current?.focus();
  };

  // Đóng dropdown khi click ra ngoài
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleKey);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKey);
    };
  }, []);

  // Optional: open on hover
  const onMouseEnter = () => {
    if (openOnHover) setOpen(true);
  };
  const onMouseLeave = () => {
    if (openOnHover) setOpen(false);
  };

  return (
    <div
      ref={dropdownRef}
      className="relative inline-block"
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      {/* Enlarged, easier-to-hit button */}
      <button
        ref={buttonRef}
        onClick={(e) => {
          e.stopPropagation();
          setOpen((v) => !v);
        }}
        aria-haspopup="menu"
        aria-expanded={open}
        className="inline-flex items-center gap-2 px-3 py-2 rounded-md hover:bg-gray-800/30 transition-colors cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <FiGlobe className="mr-0 w-5 h-5 text-gray-200" />
        <span className="text-sm font-medium text-gray-200 tracking-wide">
          {i18n.language.toUpperCase()}
        </span>
        <FiChevronDown
          className={`ml-1 w-4 h-4 text-gray-300 transition-transform ${open ? "rotate-180" : ""}`}
        />
      </button>

      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: -6 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: -6 }}
            transition={{ duration: 0.15, ease: "easeOut" }}
            className="absolute right-0 mt-3 w-36 bg-gray-800 border border-gray-700 rounded-xl shadow-lg overflow-hidden z-[9999] pointer-events-auto"
            role="menu"
            aria-label="Language selector"
            onClick={(e) => e.stopPropagation()} // keep clicks inside dropdown
          >
            <button
              onClick={(e) => {
                e.stopPropagation();
                changeLanguage("en");
              }}
              className="flex items-center gap-2 w-full px-4 py-2 text-sm text-left text-gray-200 hover:bg-blue-600 transition-colors cursor-pointer"
              role="menuitem"
            >
              English
            </button>

            <button
              onClick={(e) => {
                e.stopPropagation();
                changeLanguage("vi");
              }}
              className="flex items-center gap-2 w-full px-4 py-2 text-sm text-left text-gray-200 hover:bg-blue-600 transition-colors cursor-pointer"
              role="menuitem"
            >
              Tiếng Việt
            </button>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
