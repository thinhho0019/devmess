import { useEffect, useState } from "react";
import type { ChatMessage } from "../../pages/HomeChat";
import { fetchConversations } from "../../api/conversation";

const conversationAdminDefault: ChatMessage = {
  id: "conversation_admin_default",
  name: "Admin Chat",
  type: "conversation", // hoặc 'text' nếu bạn muốn giữ kiểu cũ
  last_message_id: "msg-001",
  last_message: {
    id: "msg-001",
    conversation_id: "conversation_admin_default",
    content: "Xin chào các bạn đã đến với hệ thống chat của chúng tôi! Đây là tin nhắn thông báo bạn không được trả lời.Hiện tại hệ thống vẫn còn trong giai đoạn phát triển, các chức năng nhắn tin giữa các người dùng hầu như hoạt động ổn định rồi! chúc các bạn 1 ngày tốt lành",
    type: "system",      // text | image | file | video | system
    status: "sent",      // sent | delivered | read
    is_edited: false,
    deleted: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    deleted_at: null,
  },
  messages: [], // mảng rỗng tạm thời
  participants: [
    {
      id: "a464d5b6-ef24-4d2b-a690-4786220cf9e4",
      user_id: "a34c1302-b4c2-4506-a88b-3b732eaf92d7",
      conversation_id: "conversation_admin_default",
      role: "member",
      joined_at: new Date().toISOString(),
      deleted_at: null,
      user: {
          id: "a34c1302-b4c2-4506-a88b-3b732eaf92d7",
          name: "Admin",
          email: "admin@gmail.com",
          avatar: "avatar/a34c1302-b4c2-4506-a88b-3b732eaf92d7.jpg",
          provider: "google",
          created_at: new Date().toISOString(),
          status: "offline",
          last_seen: new Date(0).toISOString(),
          updated_at: new Date().toISOString(),
          devices: null,
          participants: null
      }
    },
    
  ],
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  deleted_at: null,
};

export const useConversation = () => {
    // Add your hook logic here
    const [loadingConversation, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [chats, setConversations] = useState<ChatMessage[]>([]);
    const load = async () => {
        setLoading(true);
        try {
            const data: ChatMessage[] = await fetchConversations();
            if (!data) {
                setConversations([conversationAdminDefault]);
                throw new Error("No data received");
            }
            setConversations([...data, conversationAdminDefault]);
            setError(null);
            // eslint-disable-next-line
        } catch (err: any) {
            setError(err?.message ?? String(err));
        } finally {
            setLoading(false);
        }
    };
    useEffect(() => {
        void load();
    }, []);
    return { loadingConversation, error, chats, setConversations, reFetchConversation: load };
}