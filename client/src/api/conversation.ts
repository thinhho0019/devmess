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


export const fetchFindConversationTwoUsers = async (userId: string): Promise<string | null> => {
  try {
    const payload = {
      user_id: userId,
    }
    console.log("Payload sent for finding conversation:", payload);
    const res = await api.post(`/conversations/find-conversation`, payload);
    console.log("Payload sent for finding conversation:", payload);
    console.log("Fetched conversation ID between users:", res.data);
    return res.data.conversation_id ?? null;
  } catch (error) {
    console.error("Error fetching conversation between users:", error);
    return null;
  }
};