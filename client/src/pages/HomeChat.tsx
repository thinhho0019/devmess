import { useState, useEffect  } from "react";
import { FaSearch } from "react-icons/fa";
import { FiEdit, FiMoreHorizontal, FiSend, FiSmile, FiUsers } from "react-icons/fi";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import imgAvatar from "../assets/img.jpg";
import { Avatar } from "../components/avatar";
import { ChatView } from "../components/chat";

import { useIsMobile } from "../hooks";
import { useAuth } from "../hooks/auth/is_login";
import LoadingFullScreen from "../components/loading/LoadingFullScreen";
import { useImage } from "../hooks/api/useImage";
import LoadingComponent from "../components/loading/LoadingComponent";
import {  useSocket } from "../contexts/SocketContext";
import { PopupProfile } from "../components/modals/Profile";
import { PopupFriendsManager } from "../components/modals/FriendsManager";
import NotificationBell from "../components/notify/NotificationBell";
import { fetchAddFriend } from "../api/friends";
import { convertUtcToDatePart, TypeDate } from "../utils/date";
import { defaultProxyImageUrl } from "../utils/image";

import { useConversation } from "../hooks/chat/useConversation";
import { useNavigate  } from "react-router-dom";


interface LastMessage {
  id: string;
  conversation_id: string;
  content: string;
  type: string;
  status: string;
  is_edited: boolean;
  deleted: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}
interface User {
  id: string;
  name: string;
  email: string;
  avatar: string;
  provider: string;
  created_at: string;
  status: string;
  last_seen: string;
  updated_at: string;
  // eslint-disable-next-line
  devices: any[] | null;
  // eslint-disable-next-line
  participants: any[] | null;
}
interface Participant {
  id: string;
  user_id: string;
  conversation_id: string;
  role: string;
  joined_at: string;
  deleted_at: string | null;
  user: User;
}
export interface Messages {
  id: string;
  conversation_id: string;
  sender_id: string;
  content: string;
  type: 'text' | 'image' | 'file' | 'video' | 'system'; // tuỳ loại message hỗ trợ
  status: 'sending' | 'sent' | 'delivered' | 'read' | 'failed';
  is_edited: boolean;
  deleted: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
  sender?: User;
  // eslint-disable-next-line
  attachments?: any[];
  // eslint-disable-next-line
  reactions?: any[];
}
export interface ChatMessage {
  id: string;
  name: string;
  type: string;
  last_message_id: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  last_message: LastMessage;
  participants: Participant[];
  messages: Messages[];
}

// const chats: ChatMessage[] = [
//   {
//     id: "f4430c3e-0867-4f09-8870-fd950791898d",
//     type: "direct",
//     name: "",
//     last_message_id: "7aa535ab-7d34-4cb8-8c92-3beddf37be39",
//     created_at: "2025-10-22T00:27:21.551564+07:00",
//     updated_at: "2025-10-22T00:27:21.589271+07:00",
//     deleted_at: null,
//     last_message: {
//       id: "7aa535ab-7d34-4cb8-8c92-3beddf37be39",
//       conversation_id: "f4430c3e-0867-4f09-8870-fd950791898d",
//       content: "Đây là cuộc trò chuyện riêng tư hãy chào nhau đi!",
//       type: "system",
//       status: "sent",
//       is_edited: false,
//       deleted: false,
//       created_at: "2025-10-22T00:27:21.571602+07:00",
//       updated_at: "2025-10-22T00:27:21.571602+07:00",
//       deleted_at: null
//     },
//     participants: [
//       {
//         id: "41d71f2e-4c50-454a-91c9-5d8cabaf72e2",
//         user_id: "cf9900bd-efd8-482c-9740-8f560c919f09",
//         conversation_id: "f4430c3e-0867-4f09-8870-fd950791898d",
//         role: "member",
//         joined_at: "2025-10-22T00:27:21.601319+07:00",
//         deleted_at: null,
//         user: {
//           id: "cf9900bd-efd8-482c-9740-8f560c919f09",
//           name: "Haru Nice",
//           email: "awrazer0019@gmail.com",
//           avatar: "avatar/cf9900bd-efd8-482c-9740-8f560c919f09.jpg",
//           provider: "google",
//           created_at: "2025-10-22T00:27:01.891691+07:00",
//           status: "offline",
//           last_seen: "0001-01-01T07:00:00+07:00",
//           updated_at: "2025-10-22T00:27:01.891691+07:00",
//           devices: null,
//           participants: null
//         }
//       },
//       {
//         id: "1dd5ac43-0274-4ef5-ac77-dd8cbc8ed853",
//         user_id: "0fed337b-59d7-45a1-99db-bbf02866f524",
//         conversation_id: "f4430c3e-0867-4f09-8870-fd950791898d",
//         role: "member",
//         joined_at: "2025-10-22T00:27:21.601319+07:00",
//         deleted_at: null,
//         user: {
//           id: "0fed337b-59d7-45a1-99db-bbf02866f524",
//           name: "ho thinh",
//           email: "thinhho0019@gmail.com",
//           avatar: "avatar/0fed337b-59d7-45a1-99db-bbf02866f524.jpg",
//           provider: "google",
//           created_at: "2025-10-22T00:26:45.073365+07:00",
//           status: "offline",
//           last_seen: "0001-01-01T07:00:00+07:00",
//           updated_at: "2025-10-22T00:26:45.073365+07:00",
//           devices: null,
//           participants: null
//         }
//       }
//     ],
//   }
// ];


export default function HomeChat() {
  // Start with no selection so we can show a friendly empty state
  const [selected, setSelected] = useState<ChatMessage | null>(null);
  const [showProfile, setShowProfile] = useState(false);
  const { lastMessage } = useSocket();
  const [showAddFriend, setShowAddFriend] = useState(false);
  const isMobile = useIsMobile();
  const urlImage = localStorage.getItem("avatar")

  const { chats, setConversations, reFetchConversation } = useConversation();

  // Listen to ALL incoming messages to update conversation previews (last_message)
  // This should run regardless of which conversation is currently selected
  useEffect(() => {
    console.log("Last message received in HomeChat:", lastMessage);
    if (!lastMessage) return;

    try {
      const eventData = JSON.parse(lastMessage.data);

      switch (eventData.type) {
        case 'receive_message': {
          const conversation_id = eventData.conversation;
          const newMessage: Messages = eventData.message;
          // Always update last_message preview for the conversation, even if it's not selected
          handleUpdateLastMessage(conversation_id, newMessage);
          break;
        }
        case 'update_friend': {
          reFetchConversation();
          break;
        }

        default: {
          break;
        }
      }
    } catch (error) {
      console.log(error);
      console.log("Received non-JSON or malformed data:", lastMessage.data);
    }
  }, [lastMessage]);
  const handleUpdateLastMessage = (conversationId: string, newMessage: Messages) => {
    setConversations(prev =>
      prev.map(conv => {

        if (conv.id === conversationId) {
          console.log("Updating last message for conversation:", conversationId, newMessage);
          return {
            ...conv,
            last_message: {
              id: newMessage.id,
              conversation_id: newMessage.conversation_id,
              content: newMessage.content,
              type: newMessage.type,
              status: newMessage.status,
              is_edited: newMessage.is_edited,
              deleted: newMessage.deleted,
              created_at: newMessage.created_at,
              updated_at: newMessage.updated_at,
              deleted_at: newMessage.deleted_at ?? null,
            },
          };
        }
        return conv;
      })
    );
  };
  const navigate = useNavigate();
  const { user, loading, isAuthenticated } = useAuth(true);
  console.log("User data:", user);
  const { loadingImage, src } = useImage(urlImage || "");
  console.log("Image src:", src);
  const handlerAddFriend = async (friend_id: string) => {
    await fetchAddFriend(friend_id);
  }
  const handlerClickChat = (chat: ChatMessage) => {
    navigate(`/t/${chat.id}`);
    setSelected(chat);

  }
  const mockUsers = [
    { id: "u1", name: "Alice Nguyen", email: "alice@gmail.com" },
    { id: "u2", name: "Bob Tran", email: "bob@gmail.com" },
    { id: "u3", name: "Thinh Ho", email: "thinh@example.com", added: true },
  ];
  // Auto-select conversation based on URL param
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const conversationId = window.location.pathname.split("/t/")[1] || urlParams.get("conversation_id");
    if (conversationId && chats.length > 0) {
      const chat = chats.find(c => c.id === conversationId);
      if (chat) {
        setSelected(chat);
      }
    }
  }, [chats]);
  if (!isAuthenticated) {
    return;
  } else if (loading) {
    return (<LoadingFullScreen />)
  }
  if (isMobile) {
    return (

      <>
        <PopupFriendsManager
          show={showAddFriend}
          onClose={() => setShowAddFriend(false)}
          onAddFriend={handlerAddFriend}
          friendsList={mockUsers}
        />
        <PopupProfile user={user ?? undefined}
          is_profile_owner={true} onAvatarChange={() => { }} show={showProfile} onClose={() => { setShowProfile(false); }}></PopupProfile>
        <div className="h-screen w-screen flex flex-col bg-gray-100 dark:bg-gray-900 chat-scroll">

          <div className="p-4 border-b border-gray-200 dark:border-gray-800">
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-xl font-bold text-gray-800 dark:text-white">{user?.name || "Chats"}</h1>
              <div className="flex items-center gap-2">
                <button onClick={() => setShowAddFriend(true)} className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiUsers className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
                <button className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiEdit className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
                <NotificationBell />
              </div>
            </div>
            <div className="flex space-x-4 overflow-x-auto pb-2 pt-1 -mx-4 px-4">
              {chats.map(chat => (
                <div key={chat.id} onClick={() => handlerClickChat(chat)} className="flex-shrink-0 flex flex-col items-center gap-2 cursor-pointer">
                  <div className={`relative p-0.5 rounded-full ${selected?.id === chat.id ? 'ring-2 ring-blue-500' : ''}`}>

                    <Avatar src={imgAvatar} size="lg" online={chat.participants.find(p => p.user_id === user?.id)?.user.status === "true"} />
                  </div>
                  <span className={`text-xs truncate w-16 text-center ${selected?.id === chat.id ? 'font-semibold text-blue-500' : 'text-gray-600 dark:text-gray-400'}`}>{chat.name}</span>
                </div>
              ))}
            </div>
          </div>
          <div className="flex-1 min-h-0">
            {selected ? (
              <ChatView onUpdateLastMessage={handleUpdateLastMessage}
                id={selected.id.toString()} chats={selected?.messages} name={selected.name} img={
                  defaultProxyImageUrl(selected.participants.find(p => p.user_id !== user?.id)?.user.avatar || "")
                } is_mobile={isMobile} />
            ) : (
              <div className="flex flex-col items-center justify-center h-full text-center px-6">
                <div className="p-6 bg-white/5 rounded-full mb-4">
                  <FiSmile className="w-16 h-16 text-pink-400" />
                </div>
                <p className="mt-2 text-gray-300 text-lg font-medium">Chọn người bạn muốn nhắn tin nhé 💬</p>
                <p className="mt-2 text-gray-500 text-sm">Bấm nút bên dưới để tìm và thêm bạn mới.</p>
                <button
                  onClick={() => setShowAddFriend(true)}
                  className="mt-5 inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-full shadow"
                >
                  Tạo cuộc trò chuyện mới
                </button>
              </div>
            )}
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <PopupFriendsManager
        show={showAddFriend}
        onClose={() => setShowAddFriend(false)}
        onAddFriend={handlerAddFriend}
        friendsList={mockUsers}
      />
      <PopupProfile user={user ?? undefined} is_profile_owner={true} onAvatarChange={() => { }} show={showProfile} onClose={() => { setShowProfile(false); }}></PopupProfile>
      <div className="h-screen w-screen bg-gray-100 dark:bg-gray-900 chat-scroll">

        <PanelGroup direction="horizontal" className="flex h-full w-full">
          <Panel defaultSize={25} minSize={20} maxSize={35} className="flex flex-col">
            <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-800">
              <div className="flex items-center gap-3">

                <div onClick={() => setShowProfile(true)} className="relative rounded-full">
                  {loadingImage && <LoadingComponent />}
                  <Avatar src={(src) ?? undefined} size="md" online={false} />
                </div>

                <h1 className="text-xl font-bold text-gray-800 dark:text-white">{user?.name || "User"}</h1>
              </div>
              <div className="flex items-center gap-2">
                <button className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiMoreHorizontal className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
                <button onClick={() => setShowAddFriend(true)} className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiUsers className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
                <NotificationBell />
              </div>
            </div>

            <div className="p-4">
              <div className="relative">
                <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search or start a new chat"
                  className="w-full h-10 pl-10 pr-4 rounded-full bg-gray-200 dark:bg-gray-800 text-gray-800 dark:text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            <div className="flex-1 min-h-0 overflow-y-auto px-2">
              {chats.length > 0 && chats.map(chat => (
                <div
                  key={chat.id}
                  onClick={() => handlerClickChat(chat)}
                  className={`flex items-center gap-3 p-3 cursor-pointer rounded-xl transition-colors duration-200 ${selected?.id === chat.id
                    ? "bg-blue-500 text-white"
                    : "hover:bg-gray-200 dark:hover:bg-gray-800"
                    }`}
                >
                  <Avatar src={defaultProxyImageUrl(chat.participants.find(p => p.user_id !== user?.id)?.user.avatar || "")
                  } size="lg" online={chat.participants.find(p => p.user_id === user?.id)?.user.status === "true"} />
                  <div className="flex-1 min-w-0">
                    <div className="flex justify-between items-center">
                      <span className={`font-semibold truncate ${selected?.id === chat.id ? 'text-white' : 'text-gray-800 dark:text-white'}`}>
                        {chat.participants.find(p => p.user_id !== user?.id)?.user.name || "Unknown User"}
                      </span>
                      <span className={`text-xs ${selected?.id === chat.id ? 'text-blue-200' : 'text-gray-500 dark:text-gray-400'}`}>
                        {convertUtcToDatePart(chat.last_message.updated_at, TypeDate.Hours, 7)?.toString().padStart(2, "0")}:{convertUtcToDatePart(chat.last_message.updated_at, TypeDate.Minutes, 7)?.toString().padStart(2, "0")}
                      </span>
                    </div>
                    <p className={`text-sm truncate ${selected?.id === chat.id ? 'text-blue-100' : 'text-gray-500 dark:text-gray-400'}`}>
                      {chat.last_message.content}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </Panel>
          <PanelResizeHandle className="w-1 bg-gray-200 dark:bg-gray-800 hover:bg-blue-500 transition-colors" />
          <Panel minSize={30}>
            {selected ? (
              <ChatView id={selected.id.toString()} userInfor={selected.participants.find(p => p.user_id !== user?.id)?.user} name={selected.participants.find(p => p.user_id !== user?.id)?.user.name || "Unknown User"}
                img={
                  defaultProxyImageUrl(selected.participants.find(p => p.user_id !== user?.id)?.user.avatar || "")
                } onUpdateLastMessage={handleUpdateLastMessage}
              />
            ) : (
              <div className="flex-col items-center justify-center h-full text-center bg-gray-100 dark:bg-gray-900 hidden md:flex">
                <div className="p-6 bg-white dark:bg-gray-800 rounded-full">
                  <FiSend className="w-16 h-16 text-blue-500" />
                </div>
                <h2 className="mt-6 text-2xl font-semibold text-gray-800 dark:text-white">
                  Your Messages
                </h2>
                <p className="mt-2 text-gray-500 dark:text-gray-400">
                  Select a chat to start a secure and private conversation.
                </p>
              </div>
            )}
          </Panel>
        </PanelGroup>
      </div>
    </>
  );
}


