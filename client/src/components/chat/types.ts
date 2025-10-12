// Type constants thay vì enum
export const MessageType = {
  TEXT: 'text',
  IMAGE: 'image',
  FILE: 'file',
  AUDIO: 'audio',
  VIDEO: 'video',
  SYSTEM: 'system'
} as const;

export const MessageStatus = {
  SENDING: 'sending',
  SENT: 'sent',
  DELIVERED: 'delivered',
  READ: 'read',
  FAILED: 'failed'
} as const;

export const UserRole = {
  USER: 'user',
  ADMIN: 'admin',
  BOT: 'bot',
  MODERATOR: 'moderator'
} as const;

// Type definitions
export type MessageTypeValue = typeof MessageType[keyof typeof MessageType];
export type MessageStatusValue = typeof MessageStatus[keyof typeof MessageStatus];
export type UserRoleValue = typeof UserRole[keyof typeof UserRole];

// Interface cho file attachment
export interface ChatAttachment {
  id: string;
  filename: string;
  url: string;
  type: string; // mime type
  size: number; // bytes
  thumbnail?: string; // for images/videos
}

// Interface cho reaction
export interface MessageReaction {
  id?: string,
  emoji: string;
  user_ids: string[];
  count: number;
  type:string;
  description?:string
}

// Interface cho reply/thread
export interface MessageReply {
  reply_to_id: string;
  reply_to_message?: string; // preview of original message
  reply_to_user?: string; // username of original sender
}

// Enhanced Chat interface
export interface Chat {
  id: string; // Required, không nên optional
  message: string;
  created_at: string; // Sửa typo: create_at → created_at
  updated_at?: string; // Thời gian edit message
  user_id: string; // Sửa: id_user → user_id (convention)
  
  // Message metadata
  type?: MessageTypeValue;
  status?: MessageStatusValue;
  
  // Content enhancements
  attachments?: ChatAttachment[];
  reactions?: MessageReaction[];
  reply?: MessageReply;
  
  // Message features
  is_edited?: boolean;
  is_deleted?: boolean;
  deleted_at?: string;
  
  // Rich content
  mentions?: string[]; // user IDs mentioned in message
  hashtags?: string[];
  links?: string[]; // extracted URLs
  
  // Moderation
  is_flagged?: boolean;
  flagged_reason?: string;
  
  // Thread/conversation
  thread_id?: string;
  parent_message_id?: string;
}

// Enhanced User interface
export interface ChatUser {
  id: string;
  name: string;
  username?: string;
  email?: string;
  img?: string; // avatar URL
  alt?: string; // alt text for avatar
  
  // Status
  is_online?: boolean;
  last_seen?: string;
  status?: string; // custom status message
  
  // Role & permissions
  role?: UserRoleValue;
  permissions?: string[];
  
  // Profile info
  bio?: string;
  location?: string;
  timezone?: string;
  
  // Preferences
  notification_settings?: {
    mute_until?: string;
    sound_enabled?: boolean;
    desktop_notifications?: boolean;
  };
}

// Enhanced ChatProps interface
export interface ChatProps {
  id: string; // Required - conversation ID
  name?: string; // Chat room/conversation name
  description?: string; // Chat description
  is_mobile?: boolean; // Mobile view
  // Visual
  img?: string; // Chat avatar/cover image
  alt?: string; // Alt text for image
  
  // Messages & Users
  chats?: Chat[];
  users?: ChatUser[];
  current_user_id?: string;
  
  // Callbacks
  onBack?: () => void;

  // Chat metadata
  created_at?: string;
  updated_at?: string;
  created_by?: string; // user ID who created chat
  
  // Chat type & settings
  type?: 'direct' | 'group' | 'channel' | 'public' | 'private';
  is_archived?: boolean;
  is_muted?: boolean;
  is_pinned?: boolean;
  
  // Group chat specific
  admin_ids?: string[];
  max_members?: number;
  invite_link?: string;
  
  // Features
  features?: {
    file_sharing?: boolean;
    voice_messages?: boolean;
    video_calls?: boolean;
    screen_sharing?: boolean;
    message_reactions?: boolean;
    message_replies?: boolean;
    message_forwarding?: boolean;
  };
  
  // Statistics
  stats?: {
    total_messages?: number;
    total_members?: number;
    unread_count?: number;
    last_message?: Chat;
    last_activity?: string;
  };
}

// Interface cho chat list/sidebar
export interface ChatListItem {
  id: string;
  name: string;
  img?: string;
  last_message?: {
    text: string;
    timestamp: string;
    sender_name: string;
  };
  unread_count: number;
  is_online?: boolean;
  is_muted?: boolean;
  is_pinned?: boolean;
  type: 'direct' | 'group';
}

// Interface cho typing indicator
export interface TypingUser {
  user_id: string;
  username: string;
  started_at: string;
}

// Interface cho chat events (WebSocket)
export interface ChatEvent {
  type: 'message' | 'typing' | 'user_join' | 'user_leave' | 'reaction' | 'delete' | 'edit';
  chat_id: string;
  user_id: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any;
  timestamp: string;
}

// Interface cho search/filter
export interface ChatSearchParams {
  query?: string;
  user_id?: string;
  message_type?: MessageTypeValue;
  date_from?: string;
  date_to?: string;
  has_attachments?: boolean;
}

// Interface cho pagination
export interface ChatPagination {
  page: number;
  limit: number;
  total: number;
  has_more: boolean;
}

// Response interfaces
export interface ChatResponse {
  chats: Chat[];
  pagination?: ChatPagination;
  users?: ChatUser[];
}

export interface SendMessageRequest {
  message: string;
  chat_id: string;
  type?: MessageTypeValue;
  attachments?: File[];
  reply_to?: string;
  mentions?: string[];
}

export interface SendMessageResponse {
  success: boolean;
  chat: Chat;
  error?: string;
}
export type EmojiData = {
  emoji: string;
  type: string;
  description: string;
};