import { getAuthToken } from "../../utils/Auth";



export const findUserByEmail = async (email: string) => {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/search?email=${encodeURIComponent(email)}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${getAuthToken()}`
        }
    });

    if (!response.ok) {
        throw new Error("Failed to fetch user");
    }

    const data = await response.json();

    return data;
};