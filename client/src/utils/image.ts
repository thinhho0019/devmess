

export const defaultProxyImageUrl = (url: string) => {
  return `${import.meta.env.VITE_API_URL}/protected?filename=${encodeURIComponent(url)}&token=${localStorage.getItem("access_token") || ""}`;
}