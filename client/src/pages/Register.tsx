import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import { UserPlus, Eye, EyeOff } from "lucide-react";
import GoogleLoginButton from "../components/button/GoogleLoginButton";
import { registerUser } from "../services";

const Register: React.FC = () => {
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [notify, setNotify] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const navigate = useNavigate();

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        setNotify("");

        if (password !== confirmPassword) {
            setNotify("Mật khẩu nhập lại không khớp.");
            return;
        }

        try {
            // Giả sử API trả về token sau khi đăng ký thành công
            // Cần cập nhật interface trong `auth.ts` nếu cần
            // check email đã tồn tại
            const data = await registerUser(name, email, password) as { access_token?: string, message: string };

            if (data.access_token) {
                localStorage.setItem("access_token", data.access_token);
                navigate(`/auth/success?token=${data.access_token}`);
            } else {
                // Nếu API chỉ trả về tin nhắn thành công mà không có token
                setNotify(data.message || "Đăng ký thành công! Vui lòng đăng nhập.");
                setTimeout(() => navigate("/login"), 2000); // Chuyển hướng sau 2s
            }
        } catch (error) {
            if (error instanceof Error) {
                setNotify(error.message);
            } else {
                setNotify("Đã xảy ra lỗi không xác định. Vui lòng thử lại.");
            }
            console.error("Registration failed:", error);
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

            {/* Form đăng ký */}
            <motion.div
                initial={{ opacity: 0, y: 40, scale: 0.96 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ duration: 0.6, ease: "easeOut" }}
                className="relative z-10 bg-gray-900/70 backdrop-blur-xl border border-gray-700/40 shadow-2xl rounded-2xl p-8 w-full max-w-md"
            >
                <div className="text-center mb-8">
                    <div className="flex justify-center mb-3">
                        <motion.div whileHover={{ rotate: 10 }} transition={{ type: "spring", stiffness: 300 }}>
                            <UserPlus className="w-10 h-10 text-blue-400" />
                        </motion.div>
                    </div>
                    <h2 className="text-3xl font-extrabold text-white tracking-tight">
                        Tạo tài khoản mới 👋
                    </h2>
                    <p className="text-gray-400 mt-1">
                        Bắt đầu hành trình trò chuyện của bạn
                    </p>
                </div>

                <form onSubmit={handleRegister} className="space-y-4">
                    <div>
                        <label className="block text-sm font-semibold text-gray-300 mb-1">Tên của bạn</label>
                        <input
                            type="text"
                            placeholder="John Doe"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-semibold text-gray-300 mb-1">Email</label>
                        <input
                            type="email"
                            placeholder="you@example.com"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
                            required
                        />
                    </div>

                    <div className="relative">
                        <label className="block text-sm font-semibold text-gray-300 mb-1">Mật khẩu</label>
                        <input
                            type={showPassword ? "text" : "password"}
                            placeholder="••••••••"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 pr-10 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
                            required
                        />
                        <button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute inset-y-0 right-0 top-6 flex items-center px-3 text-gray-400 hover:text-white">
                            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                        </button>
                    </div>
                    <div className="relative">
                        <label className="block text-sm font-semibold text-gray-300 mb-1">Nhập lại mật khẩu</label>
                        <input
                            type={showConfirmPassword ? "text" : "password"}
                            placeholder="••••••••"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            className="w-full border border-gray-700 bg-gray-800 rounded-lg p-3 pr-10 focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-white placeholder-gray-400"
                            required
                        />
                        <button type="button" onClick={() => setShowConfirmPassword(!showConfirmPassword)} className="absolute inset-y-0 right-0 top-6 flex items-center px-3 text-gray-400 hover:text-white">
                            {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                        </button>
                    </div>

                    {notify && (<div className="w-full text-center py-1 text-red-400 font-medium h-6">
                        {notify}
                    </div>)}

                    <motion.button
                        whileTap={{ scale: 0.97 }}
                        whileHover={{ scale: 1.03 }}
                        type="submit"
                        className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 text-white py-3 rounded-lg font-semibold shadow-lg transition-all"
                    >
                        Đăng ký
                    </motion.button>
                </form>

                <div className="mt-5">
                    <GoogleLoginButton onLogin={(token) => console.log("Google Token:", token)} />
                </div>

                <p className="text-center text-sm text-gray-400 mt-6">
                    Đã có tài khoản?{" "}
                    <Link to="/login" className="text-blue-400 font-medium hover:underline hover:text-blue-300 transition">
                        Đăng nhập ngay
                    </Link>
                </p>
            </motion.div>
        </div>
    );
};

export default Register;
