import type { Messages } from "../pages/HomeChat";
import api from "./api";

const sendMessageApiUrl = `/messages/send`;


interface SendMessageRequest {
    content: string;
    conversation_id: string;
}

export const SendMessage = async (data: SendMessageRequest) => {
    try {
        const { data: result } = await api.post(sendMessageApiUrl, data);
        console.log('Message sent successfully:', result);
        return result;
    // eslint-disable-next-line
    } catch (error: any) {
        console.error('Error sending message:', error?.response?.data || error.message);
        throw error;
    }
};

export const fetchMessages = async (conversation_id: string, limit: number, before: Date) => {
    try {
        const res = await api.get(`/messages/`,{
            params: {
                conversation_id,
                limit,
                before: before.toISOString(),
            }
        });
        if(res.status !== 200){
            throw new Error(`Error fetching messages: ${res.statusText}`);
        }
        const result: Messages[] = res.data.messages;
        return result;
    // eslint-disable-next-line
    } catch (error: any) {
        console.error('Error fetching messages:', error?.response?.data || error.message);
        throw error;
    }
};