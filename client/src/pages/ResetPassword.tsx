// src/pages/ResetPasswordConfirm.tsx
import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import { Key, Eye, EyeOff } from "lucide-react";
import { confirmResetPassword } from "../services"; // gọi API reset password

const ResetPasswordConfirm: React.FC = () => {
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [notify, setNotify] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const navigate = useNavigate();

  // Lấy token từ query param (link trong email)
  const token = new URLSearchParams(window.location.search).get("token") || "";

  useEffect(() => {
    if (!token) {
      setNotify("Token không hợp lệ hoặc bị thiếu.");
    }
  }, [token]);

  const validate = () => {
    if (!token) {
      setNotify("Thiếu token. Vui lòng dùng link đúng trong email.");
      return false;
    }
    if (password.length < 8) {
      setNotify("Mật khẩu phải có ít nhất 8 ký tự.");
      return false;
    }
    if (password !== confirm) {
      setNotify("Mật khẩu xác nhận không khớp.");
      return false;
    }
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setNotify("");
    if (!validate()) return;

    setLoading(true);
    try {
      await confirmResetPassword(token, password);
      setSuccess(true);
      setNotify("Đặt lại mật khẩu thành công! Chuyển tới trang đăng nhập...");
      setTimeout(() => {
        navigate("/l");
      }, 1800);
    } catch (err) {
      if (err instanceof Error) setNotify(err.message);
      else setNotify("Đặt lại mật khẩu thất bại. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-black relative overflow-hidden text-gray-100">
      {/* Hiệu ứng nền */}
      <motion.div
        className="absolute w-[28rem] h-[28rem] bg-purple-500/20 rounded-full blur-3xl"
        animate={{ x: [0, 100, -100, 0], y: [0, 60, -60, 0], scale: [1, 1.3, 0.9, 1] }}
        transition={{ duration: 16, repeat: Infinity, ease: "easeInOut" }}
      />
      <motion.div
        className="absolute w-[22rem] h-[22rem] bg-blue-400/20 rounded-full blur-3xl"
        animate={{ x: [100, -50, 80, 100], y: [50, -60, 40, 50], scale: [1, 1.1, 1, 1] }}
        transition={{ duration: 14, repeat: Infinity, ease: "easeInOut" }}
      />

      {/* Form card */}
      <motion.div
        initial={{ opacity: 0, y: 40, scale: 0.96 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        className="relative z-10 bg-gray-900/70 backdrop-blur-xl border border-gray-700/40 shadow-2xl rounded-2xl p-8 w-full max-w-md"
      >
        <div className="text-center mb-6">
          <div className="flex justify-center mb-3">
            <motion.div whileHover={{ rotate: 8 }} transition={{ type: "spring", stiffness: 300 }}>
              <Key className="w-10 h-10 text-blue-400" />
            </motion.div>
          </div>
          <h2 className="text-2xl font-extrabold text-white tracking-tight">Đặt mật khẩu mới</h2>
          <p className="text-gray-400 mt-1">Nhập mật khẩu mới cho tài khoản của bạn</p>
        </div>

        {success ? (
          <div className="text-center space-y-4">
            <p className="text-green-400 font-semibold">{notify}</p>
            <Link to="/l" className="text-blue-400 font-medium hover:underline hover:text-blue-300 transition">
              Về trang đăng nhập
            </Link>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Mật khẩu mới */}
            <div className="relative">
              <label className="block text-sm font-semibold text-gray-300 mb-1">Mật khẩu mới</label>
              <div className="relative">
                <input
                  type={showPassword ? "text" : "password"}
                  placeholder="Nhập mật khẩu mới (ít nhất 8 ký tự)"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 pr-10 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
                  required
                  minLength={8}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-200"
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            {/* Xác nhận mật khẩu */}
            <div className="relative">
              <label className="block text-sm font-semibold text-gray-300 mb-1">Xác nhận mật khẩu mới</label>
              <div className="relative">
                <input
                  type={showConfirm ? "text" : "password"}
                  placeholder="Nhập lại mật khẩu để xác nhận"
                  value={confirm}
                  onChange={(e) => setConfirm(e.target.value)}
                  className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 pr-10 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
                  required
                  minLength={8}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirm(!showConfirm)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-200"
                >
                  {showConfirm ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            {/* Thông báo lỗi */}
            {notify && (<div className="relative w-full mt-2 px-2 py-2 text-red-500 font-medium z-50 h-6">
              {notify}
            </div>)}

            {/* Nút submit */}
            <motion.button
              whileTap={{ scale: 0.97 }}
              whileHover={{ scale: 1.03 }}
              type="submit"
              disabled={loading}
              className={`w-full ${loading
                  ? "bg-gray-700 text-gray-300"
                  : "bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500"
                } text-white py-3 rounded-lg font-semibold shadow-lg transition-all`}
            >
              {loading ? "Đang xử lý..." : "Đặt lại mật khẩu"}
            </motion.button>

            <p className="text-center text-sm text-gray-400 mt-2">
              Hoặc quay lại{" "}
              <Link to="/l" className="text-blue-400 font-medium hover:underline hover:text-blue-300 transition">
                đăng nhập
              </Link>
            </p>
          </form>
        )}
      </motion.div>
    </div>
  );
};

export default ResetPasswordConfirm;
