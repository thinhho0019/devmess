import { useEffect, useState } from "react";
import type { ChatMessage } from "../../pages/HomeChat";
import { fetchConversations } from "../../api/conversation";


export const useConversation = () => {
    // Add your hook logic here
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [chats, setConversations] = useState<ChatMessage[]>([]);
    useEffect(() => {
        const load = async () => {
            setLoading(true);
            try {
                const data: ChatMessage[] = await fetchConversations();
                setConversations(data);
                setError(null);
            // eslint-disable-next-line
            } catch (err: any) {
                setError(err?.message ?? String(err));
            } finally {
                setLoading(false);
            }
        };
        void load();
    }, []);
    return { loading, error, chats };
}