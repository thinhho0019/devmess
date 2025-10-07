import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { LogIn } from "lucide-react";
import GoogleLoginButton from "../components/button/GoogleLoginButton";

const Login: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Login with", email, password);
    // Gọi API: POST /api/login
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-black relative overflow-hidden text-gray-100">
      {/* 🌌 Hiệu ứng ánh sáng nền di chuyển */}
      <motion.div
        className="absolute w-[28rem] h-[28rem] bg-purple-500/20 rounded-full blur-3xl"
        animate={{
          x: [0, 100, -100, 0],
          y: [0, 60, -60, 0],
          scale: [1, 1.3, 0.9, 1],
        }}
        transition={{ duration: 16, repeat: Infinity, ease: "easeInOut" }}
      />
      <motion.div
        className="absolute w-[22rem] h-[22rem] bg-blue-400/20 rounded-full blur-3xl"
        animate={{
          x: [100, -50, 80, 100],
          y: [50, -60, 40, 50],
          scale: [1, 1.1, 1, 1],
        }}
        transition={{ duration: 14, repeat: Infinity, ease: "easeInOut" }}
      />

      {/* 🧊 Khối form đăng nhập */}
      <motion.div
        initial={{ opacity: 0, y: 40, scale: 0.96 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        className="relative z-10 bg-gray-900/70 backdrop-blur-xl border border-gray-700/40 shadow-2xl rounded-2xl p-8 w-full max-w-md"
      >
        <div className="text-center mb-8">
          <div className="flex justify-center mb-3">
            <motion.div
              whileHover={{ rotate: 10 }}
              transition={{ type: "spring", stiffness: 300 }}
            >
              <LogIn className="w-10 h-10 text-blue-400" />
            </motion.div>
          </div>
          <h2 className="text-3xl font-extrabold text-white tracking-tight">
            Chào mừng trở lại 👋
          </h2>
          <p className="text-gray-400 mt-1">
            Đăng nhập để tiếp tục cuộc trò chuyện
          </p>
        </div>

        <form onSubmit={handleLogin} className="space-y-5">
          <div>
            <label className="block text-sm font-semibold text-gray-300 mb-1">
              Email
            </label>
            <input
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-300 mb-1">
              Mật khẩu
            </label>
            <input
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
              required
            />
          </div>

          <motion.button
            whileTap={{ scale: 0.97 }}
            whileHover={{ scale: 1.03 }}
            type="submit"
            className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white py-3 rounded-lg font-semibold shadow-lg transition-all"
          >
            Đăng nhập
          </motion.button>
        </form>

        <div className="mt-5">
          <GoogleLoginButton
            onLogin={(token) => console.log("Google Token:", token)}
          />
        </div>

        <p className="text-center text-sm text-gray-400 mt-6">
          Chưa có tài khoản?{" "}
          <Link
            to="/signup"
            className="text-blue-400 font-medium hover:underline hover:text-blue-300 transition"
          >
            Đăng ký ngay
          </Link>
        </p>
      </motion.div>
    </div>
  );
};

export default Login;
