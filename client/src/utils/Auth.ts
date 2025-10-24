export const getAuthToken = () => {
    return localStorage.getItem("access_token") || localStorage.getItem("auth_token") || null;
};

/**
 * Clear known auth items from storage without triggering a navigation.
 * Keep this small and explicit so callers can decide how to redirect.
 */
export const clearAuth = () => {
    try {
        localStorage.removeItem("access_token");
        localStorage.removeItem("auth_token");
        localStorage.removeItem("refresh_token");
        localStorage.removeItem("user");
    } catch{
        // ignore
        // In environments without localStorage this would fail silently.
    }
};

/**
 * Logout helper: clear auth and reload the page to reset app state.
 * Callers can pass `false` to avoid an automatic reload.
 */
export const logout = (reload = true) => {
    clearAuth();
    if (reload && typeof window !== "undefined") {
        // reload the page to reset client state and force auth flow
        window.location.reload();
    }
};