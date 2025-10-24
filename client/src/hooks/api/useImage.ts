import { useState, useEffect } from "react";

export function useImage(url?: string) {
  const [src, setSrc] = useState<string | null>(null);
  const [loadingImage, setLoadingImage] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    if (!url) return;

    const img = new Image();
    img.src = `${import.meta.env.VITE_API_URL}/protected?filename=${encodeURIComponent(url)}&token=${localStorage.getItem("access_token") || ""}`;
    console.log("Loading image from URL:", img.src);
    img.onload = () => {
      setSrc(img.src);
      setLoadingImage(false);
    };

    img.onerror = () => {
      setError(true);
      setLoadingImage(false);
    };

    return () => {
      // cleanup khi component unmount
      img.onload = null;
      img.onerror = null;
    };
  }, [url]);

  return { loadingImage, src, error };
}
