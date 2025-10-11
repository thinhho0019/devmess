import "./LoadingFullScreen.css";

export default function LoadingFullScreen() {
  return (
    <div className="loading-overlay" role="status" aria-live="polite" aria-busy="true">
      <div className="loading-center">
        <div className="spinner" />
        <span className="sr-only">Loading…</span>
      </div>
    </div>
  );
}
