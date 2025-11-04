import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import type { UserResponse } from "../../types/UserResponse";
import api from "../../api/api";

export function useAuth(requireAuth = false) {
    const [user, setUser] = useState<UserResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();
    const location = window.location.pathname;
    
    useEffect(() => {
        
        const token = localStorage.getItem("access_token");
        const isPageLogin = location === "/l" || location === "/r";
        // 🔹 Nếu không có token và trang yêu cầu login
        if (!token) {
            if (requireAuth) navigate("/l");
            setLoading(false);
            return;
        }

        // 🔹 Gọi API /auth-me để xác thực (sử dụng axios với interceptors)
        const checkAuth = async () => {
            try {
                 
                const response = await api.get('/auth-me');
                
                
                
                // Kiểm tra nếu response không phải JSON
                if (typeof response.data === 'string' && response.data.includes('<html>')) {
                    throw new Error("Server returned HTML instead of JSON - check API endpoint");
                }
                
                const data: UserResponse = response.data;
                console.log("✅ Parsed user data:", data);
                // if current page login and authenticated, redirect to home
                if (requireAuth && data && isPageLogin) {
                    navigate("/t");
                }
                // Lưu thông tin user vào localStorage
                if (data) {
                    localStorage.setItem("email", data.email || "");
                    localStorage.setItem("name", data.name || "");
                    localStorage.setItem("avatar", data.avatar || "");
                    localStorage.setItem("createdAt", data.createdAt || "");
                    localStorage.setItem("updatedAt", data.updatedAt || "");
                    // ID có thể hữu ích cho các API calls khác
                    if (data.id) localStorage.setItem("user_id", data.id.toString());
                }
                
                setUser(data);
            // eslint-disable-next-line
            } catch (err: any) {
                console.error("❌ Auth check failed:", err);
                console.error("❌ Error response:", err.response);
                console.error("❌ Error config:", err.config);
                
                // Cleanup tất cả thông tin user khỏi localStorage
                localStorage.removeItem("access_token");
                localStorage.removeItem("email");
                localStorage.removeItem("name");
                localStorage.removeItem("avatar");
                localStorage.removeItem("createdAt");
                localStorage.removeItem("updatedAt");
                localStorage.removeItem("user_id");
                
                // Interceptor đã xử lý refresh token và redirect nếu thất bại
                if (requireAuth && err.response?.status === 401) {
                    navigate("/l");
                }
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, [requireAuth, navigate,location]);
    
    return { user, loading, isAuthenticated: !!user };
}
