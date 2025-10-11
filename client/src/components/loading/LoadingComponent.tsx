import "./LoadingComponent.css"

export default function LoadingComponent() {
    return (
        <div className="loading-overlay" role="status" aria-live="polite" aria-busy="true">
            <div className="loading-center">
                <div className="spinner" />
                <span className="sr-only">Loading…</span>
            </div>
        </div>
    );
}
