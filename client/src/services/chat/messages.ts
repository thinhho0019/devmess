// src/services/chat/messages.ts

// Ví dụ: Hàm lấy lịch sử tin nhắn
export const getMessageHistory = async (channelId: string) => {
  // const url = `${import.meta.env.VITE_API_URL}/messages/${channelId}`;
  // const response = await fetch(url);
  // ...
  console.log(`Fetching messages for ${channelId}`);
  return [];
};

// Ví dụ: Hàm gửi tin nhắn
export const sendMessage = async (channelId: string, content: string) => {
  console.log(`Sending message to ${channelId}: ${content}`);
  return { success: true };
};
