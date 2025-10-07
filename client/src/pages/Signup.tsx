import React, { useState } from "react";
import { Link } from "react-router-dom";

const Signup: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");

  const handleSignup = (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirm) {
      alert("Mật khẩu không khớp!");
      return;
    }
    console.log("Signup with", email, password);
    // gọi API: POST /api/signup
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white shadow-xl rounded-xl p-8 w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6">Đăng ký</h2>
        <form onSubmit={handleSignup} className="space-y-4">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full border rounded-lg p-2"
            required
          />
          <input
            type="password"
            placeholder="Mật khẩu"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full border rounded-lg p-2"
            required
          />
          <input
            type="password"
            placeholder="Nhập lại mật khẩu"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            className="w-full border rounded-lg p-2"
            required
          />
          <button
            type="submit"
            className="w-full bg-green-500 hover:bg-green-600 text-white py-2 rounded-lg font-medium"
          >
            Đăng ký
          </button>
        </form>
        <p className="text-center text-sm text-gray-600 mt-4">
          Đã có tài khoản?{" "}
          <Link to="/login" className="text-blue-600 font-medium">
            Đăng nhập
          </Link>
        </p>
      </div>
    </div>
  );
};

export default Signup;
