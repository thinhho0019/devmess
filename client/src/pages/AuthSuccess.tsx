import { useEffect } from "react";
import {  useNavigate } from "react-router-dom";
import type { UserResponse } from "../types/UserResponse";
 

export default function AuthSuccess() {
    const navigate = useNavigate();

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const token = params.get("token");
        const getUserInfor = async (): Promise<UserResponse> => {
            const response = await fetch(import.meta.env.VITE_API_URL + "/auth-me", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`, // ✅ chèn Bearer token vào header
                },
            });
            if (!response.ok) {
                throw new Error("Failed to fetch user infor")
            }
            const data: UserResponse = await response.json();
            return data

        }
        //request check auth me infor user
        getUserInfor().then((user) => {
            if (user) {

                localStorage.setItem("access_token", token || "");
                localStorage.setItem("email", user.email || "");
                localStorage.setItem("name", user.name || "");
                localStorage.setItem("avatar", user.avatar || "");
                localStorage.setItem("createdAt", user.createdAt || "");
                localStorage.setItem("updatedAt", user.updatedAt || "");

            }
            navigate("/t", { replace: true });
        }).catch((error) => {
            console.error("❌ Error fetching user info:", error);
            navigate("/l");
        })
    }, [navigate]);
    return <p>Đang đăng nhập, vui lòng chờ...</p>;
}
