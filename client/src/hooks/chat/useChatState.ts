import { useState, useRef, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import type { Chat, EmojiData } from '../../components/chat/types';
import { getTimeIsoCurrent } from '../../utils/date';

export const useChatState = (initialMessages: Chat[]) => {
  const [messages, setMessages] = useState<Chat[]>(initialMessages);
  const chatContainerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom on new message
  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [messages]);

  const sendMessage = (text: string, currentUserId: string = "user_002") => {
    const newMessage: Chat = {
      id: uuidv4(),
      message: text,
      created_at: getTimeIsoCurrent(),
      user_id: currentUserId,
      type: 'text',
      status: 'sent',
    };
    setMessages(prev => [...prev, newMessage]);
  };

  const addReaction = (messageId: string, emoji: EmojiData, currentUserId: string = "user_001") => {
    setMessages(prevMessages =>
      prevMessages.map(msg => {
        if (msg.id === messageId) {
          const newReactions = [...(msg.reactions || [])];
          const reactionIndex = newReactions.findIndex(r => r.type === emoji.type);

          if (reactionIndex > -1) {
            // Reaction exists, update it
            const existingReaction = { ...newReactions[reactionIndex] };
            if (!existingReaction.user_ids.includes(currentUserId)) {
              existingReaction.count++;
              existingReaction.user_ids.push(currentUserId);
              newReactions[reactionIndex] = existingReaction;
            }
          } else {
            // New reaction
            newReactions.push({
              emoji: emoji.emoji,
              count: 1,
              user_ids: [currentUserId],
              type: emoji.type,
            });
          }
          return { ...msg, reactions: newReactions };
        }
        return msg;
      })
    );
  };

  return {
    messages,
    sendMessage,
    addReaction,
    chatContainerRef,
  };
};
