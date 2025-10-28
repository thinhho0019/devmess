import { useState, useRef, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { FaBell } from "react-icons/fa";

interface Notification {
  id: number;
  title: string;
  message: string;
  time: string;
  read?: boolean;
}

export default function NotificationBell() {
  const [showPopup, setShowPopup] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>([
    {
      id: 1,
      title: "Thông báo mới từ DevTH",
      message: "Chào mừng bạn đến với web nhắn tin realtime",
      time: "2 min ago",
      read: false,
    },
  ]);

  const unreadCount = notifications.filter((n) => !n.read).length;
  const bellRef = useRef<HTMLDivElement>(null);

  const handleToggle = () => setShowPopup((prev) => !prev);
  const handleCheckNotification = (id: number) => {
    setNotifications((prev) =>
      prev.map((n) => 
        n.id === id ? { ...n, read: true } : n
      )
    );
  }
  const markAllAsRead = () => {
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
  };

  // ✅ Detect click outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (bellRef.current && !bellRef.current.contains(event.target as Node)) {
        setShowPopup(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div ref={bellRef} className="relative inline-block">
      {/* Bell button */}
      <button
        onClick={handleToggle}
        className="relative p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition"
      >
        <motion.div
          animate={
            unreadCount > 0
              ? { rotate: [0, -10, 10, -10, 10, 0] }
              : {}
          }
          transition={{ duration: 0.4 }}
        >
          <FaBell className="w-5 h-5 text-gray-600 dark:text-gray-300" />
        </motion.div>

        {/* Red dot */}
        {unreadCount > 0 && (
          <span className="absolute top-1 right-1 w-2.5 h-2.5 bg-red-500 rounded-full border border-white dark:border-gray-800"></span>
        )}
      </button>

      {/* Popup notifications */}
      <AnimatePresence>
        {showPopup && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="absolute right-0 mt-2 w-80 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden z-50"
          >
            <div className="flex justify-between items-center px-4 py-2 border-b border-gray-200 dark:border-gray-700">
              <h3 className="font-semibold text-gray-800 dark:text-white">
                Notifications
              </h3>
              {notifications.some((n) => !n.read) && (
                <button
                  onClick={markAllAsRead}
                  className="text-xs text-indigo-500 hover:underline"
                >
                  Mark all as read
                </button>
              )}
            </div>

            <div className="max-h-60 overflow-y-auto">
              {notifications.length === 0 ? (
                <p className="text-center text-gray-400 py-6">
                  No notifications
                </p>
              ) : (
                notifications.map((n) => (
                  <div onClick={() => handleCheckNotification(n.id)}
                    key={n.id}
                    className={`px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-700 transition ${
                      !n.read ? "bg-indigo-50 dark:bg-indigo-900/20" : ""
                    }`}
                  >
                    <h4 className="font-medium text-sm text-gray-800 dark:text-gray-100">
                      {n.title}
                    </h4>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {n.message}
                    </p>
                    <span className="text-xs text-gray-400">{n.time}</span>
                  </div>
                ))
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
