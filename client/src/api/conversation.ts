import type { ChatMessage } from '../pages/HomeChat';
import api from './api';

export const fetchConversations = async (): Promise<ChatMessage[]> => {
  try {
    const res = await api.get('/conversations/');
    console.log("Fetched conversations:", res.data);
    if (!res.data) {
      throw new Error("No data received from conversations endpoint");
    }
    if (Object.hasOwn(res.data, 'conversations')) {
      return res.data.conversations;
    }
    return res.data ?? [];
  } catch (error) {
    console.error("Error fetching conversations:", error);
    return [];
  }

};  

 

