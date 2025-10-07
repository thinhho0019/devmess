import { useState, useRef, useEffect, useCallback } from 'react';
import { v4 as uuidv4 } from 'uuid';
import type { Chat, MessageReaction, EmojiData } from '../../components/chat/types';
import { getTimeIsoCurrent } from '../../utils/date';
import { MessageType, MessageStatus } from '../../components/chat/types';
import { useSocket } from '../socket/useSocket';

export const useChatManager = (initialChats: Chat[], currentUserId: string) => {
    const [messages, setMessages] = useState<Chat[]>(initialChats);
    const chatContainerRef = useRef<HTMLDivElement>(null);
    const { lastMessage, sendMessage: sendSocketMessage } = useSocket();
    // Listen for incoming messages from WebSocket
    useEffect(() => {
        if (lastMessage !== null) {
            try {
                const eventData = JSON.parse(lastMessage.data);

                switch (eventData.type) {
                    case 'new_message': {
                        const newMessage: Chat = eventData.payload;
                        if (newMessage && newMessage.id && !messages.some(msg => msg.id === newMessage.id)) {
                            setMessages(prev => [...prev, newMessage]);
                        }
                        break;
                    }
                    case 'message_status': {
                        const { tempId, newId, status, timestamp } = eventData.payload;
                        setMessages(prevMessages =>
                            prevMessages.map(msg =>
                                msg.id === tempId ? { ...msg, id: newId, status: status, created_at: timestamp } : msg
                            )
                        );
                        break;
                    }
                    case 'reaction_update': {
                        const { messageId, reactions } = eventData.payload;
                        setMessages(prevMessages =>
                            prevMessages.map(msg =>
                                msg.id === messageId ? { ...msg, reactions: reactions } : msg
                            )
                        );
                        break;
                    }
                    default:
                        // Handle legacy format or other message types
                        if (eventData && eventData.id && eventData.message) {
                            if (!messages.some(msg => msg.id === eventData.id)) {
                                setMessages(prev => [...prev, eventData]);
                            }
                        }
                }
            } catch (error) {
                console.log(error);
                console.log("Received non-JSON or malformed data:", lastMessage.data);
            }
        }
    }, [lastMessage, messages]);


    useEffect(() => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
        }
    }, [messages]);

    const sendMessage = useCallback((text: string, file: File | null = null) => {
        if (!text.trim() && !file) return;
        console.log(currentUserId);
        const newMessage: Chat = {
            id: `msg_${uuidv4()}`,
            message: text,
            created_at: getTimeIsoCurrent(),
            user_id: currentUserId,
            type: file ? MessageType.FILE : MessageType.TEXT,
            status: MessageStatus.SENDING, // Set status to sending
            attachments: file ? [{
                id: `att_${uuidv4()}`,
                filename: file.name,
                url: URL.createObjectURL(file), // Temporary URL
                type: file.type,
                size: file.size
            }] : undefined,
        };
        
        const messagePayload = {
            type: 'new_message',
            payload: newMessage
        };

        // Optimistically update the UI
        setMessages(prev => [...prev, newMessage]);

        // Send the message via WebSocket
        sendSocketMessage(JSON.stringify(messagePayload));

    }, [currentUserId, sendSocketMessage]);

    const addReaction = useCallback((messageId: string, emoji: EmojiData | MessageReaction) => {
        const reactionPayload = {
            type: 'add_reaction',
            payload: {
                messageId,
                emoji,
                userId: currentUserId
            }
        };
        sendSocketMessage(JSON.stringify(reactionPayload));

        // Optimistically update UI
        setMessages(prevMessages =>
            prevMessages.map(msg => {
                if (msg.id === messageId) {
                    const newReactions = [...(msg.reactions || [])];
                    const reactionIndex = newReactions.findIndex(r => r.type === emoji.type);

                    if (reactionIndex > -1) {
                        const existingReaction = { ...newReactions[reactionIndex] };
                        if (!existingReaction.user_ids.includes(currentUserId)) {
                            existingReaction.count++;
                            existingReaction.user_ids.push(currentUserId);
                        } else {
                            existingReaction.count--;
                            existingReaction.user_ids = existingReaction.user_ids.filter(uid => uid !== currentUserId);
                        }
                        
                        if (existingReaction.count > 0) {
                            newReactions[reactionIndex] = existingReaction;
                        } else {
                            newReactions.splice(reactionIndex, 1);
                        }

                    } else {
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
    }, [currentUserId, sendSocketMessage]);

    return {
        messages,
        chatContainerRef,
        sendMessage,
        addReaction,
    };
};
