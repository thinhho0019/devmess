export const getAuthToken = () => {
    return localStorage.getItem("access_token") || localStorage.getItem("auth_token") || null;
};