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

        // 🔹 Gọi API /auth-me để xác thực
        const checkAuth = async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/auth-me`, {
                    headers: {
                        "Authorization": `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) throw new Error("Unauthorized");

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
