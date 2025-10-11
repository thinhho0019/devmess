import { useState } from "react";
import { FaSearch } from "react-icons/fa";
import { FiEdit, FiMoreHorizontal, FiSend } from "react-icons/fi";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import imgAvatar from "../assets/img.jpg";
import { Avatar } from "../components/avatar";
import { ChatView, type Chat } from "../components/chat";
import { useIsMobile } from "../hooks";
import { useAuth } from "../hooks/auth/is_login";
import LoadingFullScreen from "../components/loading/LoadingFullScreen";
import { useImage } from "../hooks/api/useImage";
import LoadingComponent from "../components/loading/LoadingComponent";
import { SocketProvider } from "../contexts/SocketContext";

interface ChatMessage {
  id: number;
  name: string;
  lastMessage: string;
  time: string;
  online: boolean;
}


interface ChatMessage {
  id: number;
  name: string;
  lastMessage: string;
  time: string;
  online: boolean;
  messages: Chat[];
}

const chats: ChatMessage[] = [
  {
    id: 1,
    name: "Thịnh hồ",
    lastMessage: "Hey, how are you?",
    time: "10:30 AM",
    online: true,
    messages: [
      {
        id: "msg_1",
        message: "Hey, how are you?",
        created_at: "2025-10-05T10:30:00Z",
        user_id: "user_alice",
        type: "text",
        status: "read",
        is_edited: false,
        is_deleted: false

      }],
  }
];


export default function HomeChat() {
  const [selected, setSelected] = useState<ChatMessage | null>(chats[0]);
  const isMobile = useIsMobile();
  const urlImage = localStorage.getItem("avatar")
  const { loading, isAuthenticated } = useAuth(true);
  const { loadingImage, src } = useImage(urlImage || "");
  if (!isAuthenticated) {
    return;
  } else if (loading) {
    return (<LoadingFullScreen />)
  }
  if (isMobile) {
    return (
      <SocketProvider>
        <div className="h-screen w-screen flex flex-col bg-gray-100 dark:bg-gray-900">

          <div className="p-4 border-b border-gray-200 dark:border-gray-800">
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-xl font-bold text-gray-800 dark:text-white">Chats</h1>
              <div className="flex items-center gap-2">
                <button className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FaSearch className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
                <button className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiEdit className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
              </div>
            </div>
            <div className="flex space-x-4 overflow-x-auto pb-2 pt-1 -mx-4 px-4">
              {chats.map(chat => (
                <div key={chat.id} onClick={() => setSelected(chat)} className="flex-shrink-0 flex flex-col items-center gap-2 cursor-pointer">
                  <div className={`relative p-0.5 rounded-full ${selected?.id === chat.id ? 'ring-2 ring-blue-500' : ''}`}>

                    <Avatar src={imgAvatar} size="lg" online={chat.online} />
                  </div>
                  <span className={`text-xs truncate w-16 text-center ${selected?.id === chat.id ? 'font-semibold text-blue-500' : 'text-gray-600 dark:text-gray-400'}`}>{chat.name}</span>
                </div>
              ))}
            </div>
          </div>
          <div className="flex-1 min-h-0">
            {selected ? (
              <ChatView id={selected.id.toString()} name={selected.name} />
            ) : (
              <div className="flex flex-col items-center justify-center h-full text-center">
                <FiSend className="w-16 h-16 text-gray-400" />
                <p className="mt-4 text-gray-500">Select a chat to start messaging</p>
              </div>
            )}
          </div>
        </div>
      </SocketProvider>
    );
  }

  return (
    <SocketProvider>
      <div className="h-screen w-screen bg-gray-100 dark:bg-gray-900">

        <PanelGroup direction="horizontal" className="flex h-full w-full">
          <Panel defaultSize={25} minSize={20} maxSize={35} className="flex flex-col">
            <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-800">
              <div className="flex items-center gap-3">

                <div className="relative rounded-full">
                  {loadingImage && <LoadingComponent />}
                  <Avatar src={src || imgAvatar} size="md" online={false} />
                </div>

                <h1 className="text-xl font-bold text-gray-800 dark:text-white">Chats</h1>
              </div>
              <div className="flex items-center gap-2">
                <button className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiMoreHorizontal className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
                <button className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
                  <FiEdit className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                </button>
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
              {chats.map(chat => (
                <div
                  key={chat.id}
                  onClick={() => setSelected(chat)}
                  className={`flex items-center gap-3 p-3 cursor-pointer rounded-xl transition-colors duration-200 ${selected?.id === chat.id
                    ? "bg-blue-500 text-white"
                    : "hover:bg-gray-200 dark:hover:bg-gray-800"
                    }`}
                >
                  <Avatar src={imgAvatar} size="lg" online={chat.online} />
                  <div className="flex-1 min-w-0">
                    <div className="flex justify-between items-center">
                      <span className={`font-semibold truncate ${selected?.id === chat.id ? 'text-white' : 'text-gray-800 dark:text-white'}`}>
                        {chat.name}
                      </span>
                      <span className={`text-xs ${selected?.id === chat.id ? 'text-blue-200' : 'text-gray-500 dark:text-gray-400'}`}>
                        {chat.time}
                      </span>
                    </div>
                    <p className={`text-sm truncate ${selected?.id === chat.id ? 'text-blue-100' : 'text-gray-500 dark:text-gray-400'}`}>
                      {chat.lastMessage}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </Panel>
          <PanelResizeHandle className="w-1 bg-gray-200 dark:bg-gray-800 hover:bg-blue-500 transition-colors" />
          <Panel minSize={30}>
            {selected ? (
              <ChatView id={selected.id.toString()} name={selected.name} chats={selected.messages} />
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
    </SocketProvider>
  );
}
