import { useState, useEffect } from "react";

export function useImage(url?: string) {
  const [src, setSrc] = useState<string | null>(null);
  const [loadingImage, setLoadingImage] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    if (!url) return;

    const img = new Image();
    img.src = url;

    img.onload = () => {
      setSrc(url);
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
