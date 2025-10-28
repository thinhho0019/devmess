// src/contexts/SocketContext.tsx
import React, { createContext, useContext, useEffect, useState, type ReactNode, useRef, useCallback } from 'react';
import WebSocketService, { wsUrl } from '../services/websocketService';

type SocketContextType = {
    readyState: number;
    lastMessage: MessageEvent | null;
    sendMessage: (data: string) => void;
};

const SocketContext = createContext<SocketContextType | undefined>(undefined);

/* eslint-disable react-refresh/only-export-components */
export const useSocket = (): SocketContextType => {
    const context = useContext(SocketContext);
    if (context === undefined) {
        throw new Error('useSocket must be used within a SocketProvider');
    }
    return context;
};

interface SocketProviderProps {
    children: ReactNode;
}

export const SocketProvider: React.FC<SocketProviderProps> = ({ children }) => {
    const [readyState, setReadyState] = useState<number>(WebSocket.CLOSED);
    const [lastMessage, setLastMessage] = useState<MessageEvent | null>(null);
    const webSocketService = useRef<WebSocketService | null>(null);

    useEffect(() => {
        // Create a new instance for each SocketProvider
        const ws = new WebSocketService(wsUrl);
        webSocketService.current = ws;

        const handleOpen = () => setReadyState(ws.getReadyState());
        const handleClose = () => setReadyState(ws.getReadyState());
        const handleError = () => setReadyState(ws.getReadyState());
        const handleMessage = (event: MessageEvent) => setLastMessage(event);

        // Assign handlers
        ws.onOpen = handleOpen;
        ws.onClose = handleClose;
        ws.onError = handleError;
        ws.addMessageListener(handleMessage);

        // Connect
        ws.connect();

        // Cleanup: close connection and remove listeners
        return () => {
            ws.removeMessageListener(handleMessage);
            ws.close();
        };
    }, []); // Empty dependency array ensures this runs once per provider instance

    const sendMessage = useCallback((data: string) => {
        if (webSocketService.current) {
            webSocketService.current.sendMessage(data);
            try {
                // Also set lastMessage locally so outgoing messages can be processed by listeners
                // Create a synthetic MessageEvent with the same shape as incoming events
                const evt = new MessageEvent('message', { data: String(data) });
                setLastMessage(evt);
            } catch   {
                // In some environments MessageEvent constructor may not be available; ignore silently
            }
        }
    }, []);

    const contextValue: SocketContextType = {
        readyState,
        lastMessage,
        sendMessage,
    };

    return (
        <SocketContext.Provider value={contextValue}>
            {children}
        </SocketContext.Provider>
    );
};
