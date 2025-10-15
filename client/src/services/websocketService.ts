type MessageListener = (event: MessageEvent) => void;

class WebSocketService {
    private socket: WebSocket | null = null;
    private listeners: MessageListener[] = [];
    private url: string;
    private reconnectInterval: number = 5000; // 5 seconds
    private reconnectAttempts: number = 5;
    private currentAttempts: number = 0;

    public onOpen: (() => void) | null = null;
    public onClose: (() => void) | null = null;
    public onError: ((event: Event) => void) | null = null;

    constructor(url: string) {
        this.url = url;
    }
    
    public connect() {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.log("WebSocket is already connected.");
            return;
        }

        this.socket = new WebSocket(this.url);
        this.currentAttempts++;

        this.socket.onopen = () => {
            console.log("WebSocket connected");
            this.currentAttempts = 0;
            if (this.onOpen) this.onOpen();
        };

        this.socket.onmessage = (event) => {
            this.listeners.forEach(listener => listener(event));
        };

        this.socket.onerror = (event) => {
            console.error("WebSocket error:", event);
            if (this.onError) this.onError(event);
        };

        this.socket.onclose = () => {
            console.log("WebSocket disconnected.");
            if (this.onClose) this.onClose();

            if (this.currentAttempts < this.reconnectAttempts) {
                console.log(`Attempting to reconnect in ${this.reconnectInterval / 1000}s...`);
                setTimeout(() => this.connect(), this.reconnectInterval);
            } else {
                console.error("Max reconnect attempts reached.");
            }
        };
    }

    public sendMessage(data: string | ArrayBufferLike | Blob | ArrayBufferView) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(data);
        } else {
            console.error("WebSocket is not connected.");
        }
    }

    public addMessageListener(listener: MessageListener) {
        this.listeners.push(listener);
    }

    public removeMessageListener(listener: MessageListener) {
        this.listeners = this.listeners.filter(l => l !== listener);
    }

    public close() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    public getReadyState(): number {
        return this.socket?.readyState ?? WebSocket.CLOSED;
    }
}

const wsUrl = import.meta.env.VITE_WS_URL || "ws://localhost:8080/ws";
// const wsUrlLocal = "ws://192.168.100.249:8080/ws";

export { wsUrl };
export default WebSocketService;
