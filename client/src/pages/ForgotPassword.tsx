import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { Mail } from "lucide-react";
import { sendResetPassword } from "../services"; // 🔹 API gửi mail reset

const ForgotPassword: React.FC = () => {
    const [email, setEmail] = useState("");
    const [notify, setNotify] = useState("");
    const [loading, setLoading] = useState(false);
    const [sent, setSent] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setNotify("");
        setLoading(true);

        try {
            await sendResetPassword(email);
            setSent(true);
            setNotify("Liên kết đặt lại mật khẩu đã được gửi đến email của bạn!");
        } catch (error) {
            if (error instanceof Error) {
                setNotify(error.message);
            } else {
                setNotify("Không thể gửi yêu cầu. Vui lòng thử lại sau.");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-black relative overflow-hidden text-gray-100">
            {/* 🌌 Hiệu ứng nền */}
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

            {/* 🧊 Form */}
            <motion.div
                initial={{ opacity: 0, y: 40, scale: 0.96 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ duration: 0.6, ease: "easeOut" }}
                className="relative z-10 bg-gray-900/70 backdrop-blur-xl border border-gray-700/40 shadow-2xl rounded-2xl p-8 w-full max-w-md"
            >
                <div className="text-center mb-8">
                    <div className="flex justify-center mb-3">
                        <motion.div whileHover={{ rotate: 10 }} transition={{ type: "spring", stiffness: 300 }}>
                            <Mail className="w-10 h-10 text-blue-400" />
                        </motion.div>
                    </div>
                    <h2 className="text-3xl font-extrabold text-white tracking-tight">
                        Quên mật khẩu 🔑
                    </h2>
                    <p className="text-gray-400 mt-1">
                        Nhập email để nhận liên kết đặt lại mật khẩu
                    </p>
                </div>

                {!sent ? (
                    <form onSubmit={handleSubmit} className="space-y-5">
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

                        {notify && (
                            <div className="relative w-full mt-2 px-2 py-2 text-red-500 font-medium z-50 h-6">
                                {notify}
                            </div>
                        )}

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
                            {loading ? "Đang gửi..." : "Gửi yêu cầu"}
                        </motion.button>
                    </form>
                ) : (
                    <div className="text-center space-y-4">
                        <p className="text-green-400 font-semibold">{notify}</p>
                        <Link
                            to="/l"
                            className="text-blue-400 font-medium hover:underline hover:text-blue-300 transition"
                        >
                            Quay lại đăng nhập
                        </Link>
                    </div>
                )}
            </motion.div>
        </div>
    );
};

export default ForgotPassword;
