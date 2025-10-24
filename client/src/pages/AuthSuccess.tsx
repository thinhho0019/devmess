import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import type { UserResponse } from "../types/UserResponse";
import api from "../api/api";

export default function AuthSuccess() {
    const navigate = useNavigate();

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const token = params.get("token");

        if (!token) {
            console.error("❌ No token found in URL");
            navigate("/l");
            return;
        }

        const getUserInfo = async (): Promise<UserResponse> => {
            // Tạm thời set token vào localStorage để axios interceptor có thể sử dụng
            localStorage.setItem("access_token", token);
            
            const response = await api.get("/auth-me");
            return response.data;
        };

        // Request để lấy thông tin user
        getUserInfo()
            .then((user) => {
                if (user) {
                    // Lưu thông tin user vào localStorage
                    localStorage.setItem("access_token", token);
                    localStorage.setItem("email", user.email || "");
                    localStorage.setItem("name", user.name || "");
                    localStorage.setItem("avatar", user.avatar || "");
                    localStorage.setItem("createdAt", user.createdAt || "");
                    localStorage.setItem("updatedAt", user.updatedAt || "");
                }
                navigate("/t", { replace: true });
            })
            .catch((error) => {
                console.error("❌ Error fetching user info:", error);
                localStorage.removeItem("access_token"); // Cleanup invalid token
                navigate("/l");
            });
    }, [navigate]);

    return <p>Đang đăng nhập, vui lòng chờ...</p>;
}
