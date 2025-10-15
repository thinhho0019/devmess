import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import type { UserResponse } from "../../types/UserResponse";



export function useAuth(requireAuth = false) {
    const [user, setUser] = useState<UserResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem("access_token");

        // 🔹 Nếu không có token và trang yêu cầu login
        if (!token) {
            if (requireAuth) navigate("/l");
            setLoading(false);
            return;
        }
        const fetchRefreshToken = async (token: string) => {
            try {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/refresh`, {
                    method: "POST",
                    headers: {
                        "Authorization": `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });
                if (response.ok) {
                    const data = await response.json();
                    const { access_token } = data;
                    localStorage.setItem("access_token", access_token);
                    return access_token;
                } else {
                    throw new Error("Failed to refresh token");
                }
            } catch (error) {
                console.error("❌ Error refreshing token:", error);
            }
        };

        // 🔹 Gọi API /auth-me để xác thực
        const checkAuth = async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/auth-me`, {
                    headers: {
                        "Authorization": `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    const errorData = await response.json();

                    // Xử lý khi token hết hạn hoặc không hợp lệ
                    if (response.status === 401 && errorData.error === "Access token expired or invalid") {
                        // refresh token
                        const newToken = await fetchRefreshToken(token);
                        if (newToken) {
                            // Thử gọi lại /auth-me với token mới
                            localStorage.setItem("access_token", newToken);
                            return checkAuth(); // Gọi lại hàm checkAuth với token mới
                        }
                    }

                    throw new Error(errorData.error || `HTTP ${response.status}`);
                }

                const data: UserResponse = await response.json();
                console.log("✅ Auth check successful:", data);
                setUser(data);
            } catch (err) {
                console.error("❌ Auth check failed:", err);
                localStorage.removeItem("access_token");
                if (requireAuth) navigate("/l");
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, [requireAuth, navigate]);
    return { user, loading, isAuthenticated: !!user };
}
