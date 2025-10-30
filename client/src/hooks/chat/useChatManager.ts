import { useState, useRef, useEffect, useCallback } from 'react';
import { v4 as uuidv4 } from 'uuid';
import type { MessageReaction, EmojiData } from '../../components/chat/types';
import { convertTimeToOnlineStatus,  getTimeIsoCurrent } from '../../utils/date';
import { MessageType, MessageStatus } from '../../components/chat/types';
import { useSocket } from '../socket/useSocket';
import { fetchMessages, SendMessage } from '../../api/message';
import { useParams } from 'react-router-dom';
import type { Messages } from '../../pages/HomeChat';

export const useChatManager = (initialChats: Messages[], currentUserId: string, userID: string,
    onUpdateLastMessage: (conversationId: string, newMessage: Messages) => void,
    handlerChangeStatus: (user_id: string, newStatus: string) => void
) => {
    const [messages, setMessages] = useState<Messages[]>(initialChats);

    const { conversation_id } = useParams<{ conversation_id: string }>();
    const [loadingMessages, setLoadingMessages] = useState<boolean>(false);
    const chatContainerRef = useRef<HTMLDivElement>(null);
    const { lastMessage, sendMessage: sendSocketMessage } = useSocket();

    // auto sending check online_user

    useEffect(() => {
        const presencePayload = {
            type: 'is_online',
            from: currentUserId,
            to: userID
        };
        sendSocketMessage(JSON.stringify(presencePayload));
        const interval = setInterval(() => {
            console.log("Sending presence payload (interval):", presencePayload);
            sendSocketMessage(JSON.stringify(presencePayload));
        }, 30000);

        return () => clearInterval(interval);
    }, [userID, sendSocketMessage, currentUserId]);
    // Fetch messages when conversation_id changes (initial load only)
    useEffect(() => {
        if (conversation_id) {
            if (conversation_id === "conversation_admin_default") {
                setMessages([{
                    id: "msg_admin_001",
                    content: "This is a notification message. You cannot reply to this message.",
                    conversation_id: "conversation_admin_default",
                    sender_id: "admin_system",
                    type: MessageType.SYSTEM,
                    status: MessageStatus.SENT,
                    is_edited: false,
                    deleted: false,
                    created_at: new Date().toISOString(),
                    updated_at: new Date().toISOString(),
                    deleted_at: null,
                    sender: {
                        id: "admin_system",
                        name: "System Admin",
                        email: "",
                        provider: "system",
                        created_at: new Date().toISOString(),
                        status: "online",
                        last_seen: new Date().toISOString(),
                        updated_at: new Date().toISOString(),
                        devices: null,
                        participants: null,
                        avatar: ''
                    }
                },{
                    id: "msg_admin_001",
                    content: "Xin chào các bạn đã đến với hệ thống chat của chúng tôi! Đây là tin nhắn thông báo bạn không được trả lời.Hiện tại hệ thống vẫn còn trong giai đoạn phát triển, các chức năng nhắn tin giữa các người dùng hầu như hoạt động ổn định rồi! chúc các bạn 1 ngày tốt lành",
                    conversation_id: "conversation_admin_default",
                    sender_id: "admin_system",
                    type: MessageType.TEXT,
                    status: MessageStatus.SENT,
                    is_edited: false,
                    deleted: false,
                    created_at: new Date().toISOString(),
                    updated_at: new Date().toISOString(),
                    deleted_at: null,
                    sender: {
                        id: "admin_system",
                        name: "System Admin",
                        email: "",
                        provider: "system",
                        created_at: new Date().toISOString(),
                        status: "online",
                        last_seen: new Date().toISOString(),
                        updated_at: new Date().toISOString(),
                        devices: null,
                        participants: null,
                        avatar: ''
                    }
                }]);
                return;
            }
            setLoadingMessages(true);
            fetchMessages(conversation_id, 50, new Date())
                .then((data) => {
                    console.log("Messages data for chat", conversation_id, ":", data);
                    if (Array.isArray(data)) {
                        data = data.slice().reverse();
                    }
                    setMessages(data);
                })
                .catch((err) => {
                    console.error("Failed to fetch messages:", err);
                })
                .finally(() => {
                    setLoadingMessages(false);
                });
        }
    }, [conversation_id]); // Only depend on conversation_id, not lastMessage

    // Listen for incoming messages from WebSocket (separate effect)
    useEffect(() => {
        if (lastMessage !== null) {
            try {
                const eventData = JSON.parse(lastMessage.data);
                console.log("Received WebSocket message:", eventData);
                switch (eventData.type) {
                    case 'new_message': {
                        const newMessage: Messages = eventData.payload;
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
                    case 'receive_message': {
                        console.log("Received 'receive_message' event:", eventData);
                        const messsage: Messages = eventData.message;
                        const conversation = eventData.conversation;
                        // Only add to messages list if this message belongs to the currently open conversation
                        if (conversation === conversation_id) {
                            setMessages((prev) => [...prev, messsage]);
                        }
                        // Always update last_message preview for the ACTUAL conversation (not the focused one)
                        onUpdateLastMessage(conversation || '', messsage);
                        break;
                    }
                    case 'is_online_response': {
                        if (conversation_id === "conversation_admin_default")break; 
                        const { user_id, is_online, time_online } = eventData;
                        console.log(`User ${user_id} is ${is_online ? 'online' : 'offline'}`);
                        const message_show = is_online ? 'online' : `offline since ${convertTimeToOnlineStatus(time_online)}`;
                        handlerChangeStatus(user_id, message_show);
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
    }, [lastMessage]); // Only depend on lastMessage for websocket updates


    // Auto-scroll to bottom when messages change
    useEffect(() => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
        }
    }, [messages]);

    // Auto-scroll to bottom when first loaded (after fetch completes)
    useEffect(() => {
        if (messages.length > 0 && chatContainerRef.current) {
            // Small delay to ensure DOM has rendered
            const timer = setTimeout(() => {
                if (chatContainerRef.current) {
                    chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
                }
            }, 100);
            return () => clearTimeout(timer);
        }
    }, [messages.length]);

    const sendMessage = useCallback((text: string, file: File | null = null) => {

        if (!text.trim() && !file) return;
        console.log(currentUserId);
        const newMessage: Messages = {
            id: `msg_${uuidv4()}`,
            content: text,
            created_at: getTimeIsoCurrent(),
            sender_id: currentUserId,
            type: file ? MessageType.FILE : MessageType.TEXT,
            status: MessageStatus.SENDING, // Set status to sending
            attachments: file ? [{
                id: `att_${uuidv4()}`,
                filename: file.name,
                url: URL.createObjectURL(file), // Temporary URL
                type: file.type,
                size: file.size
            }] : undefined,
            conversation_id: '',
            is_edited: false,
            deleted: false,
            updated_at: new Date().toISOString(),
        };

        const messagePayload = {
            type: 'new_message',
            payload: newMessage
        };
        if (!conversation_id) {
            console.error('Conversation ID is missing. Cannot send message.');
            return;
        }
        SendMessage({
            content: text,
            conversation_id: conversation_id, // Replace with actual conversation ID
        }).then((res) => {
            console.log('Message sent successfully:', res);
        }).catch((error) => {
            console.error('Error sending message:', error);
        });
        // Optimistically update the UI
        setMessages(prev => [...prev, newMessage]);
        onUpdateLastMessage(conversation_id, newMessage);
        // Send the message via WebSocket
        sendSocketMessage(JSON.stringify(messagePayload));

    }, [currentUserId, sendSocketMessage, conversation_id]);

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
                            // existingReaction.user_ids = existingReaction.user_ids.filter(sender_id => uid !== currentUserId);
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
        loadingMessages,conversation_id
    };
};
