import { useState, useEffect } from "react";
import { FaSearch } from "react-icons/fa";
import { FiEdit, FiMoreHorizontal, FiSend } from "react-icons/fi";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import imgAvatar from "../assets/img.jpg";
import { Avatar } from "../components/avatar";
import { ChatView } from "../components/chat";

interface ChatMessage {
  id: number;
  name: string;
  lastMessage: string;
  time: string;
  online: boolean;
}

const chats: ChatMessage[] = [
  { id: 1, name: "Alice", lastMessage: "Hey, how are you?", time: "10:30 AM", online: true },
  { id: 2, name: "Bob", lastMessage: "Let's catch up later.", time: "9:45 AM", online: false },
  { id: 3, name: "Project Team", lastMessage: "Meeting at 2 PM.", time: "Yesterday", online: true },
  { id: 4, name: "Family Group", lastMessage: "Dinner tonight?", time: "Yesterday", online: true },
  { id: 5, name: "John Doe", lastMessage: "Can you send me the file?", time: "Sep 18", online: false },
  { id: 6, name: "Sarah", lastMessage: "See you tomorrow!", time: "Sep 17", online: true },
  { id: 7, name: "Chris", lastMessage: "Check this out!", time: "Sep 16", online: false },
  { id: 8, name: "Design Team", lastMessage: "New mockups available.", time: "Sep 15", online: true },
];

const useIsMobile = () => {
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768);

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return isMobile;
};

export default function HomeChat() {
  const [selected, setSelected] = useState<ChatMessage | null>(chats[0]);
  const isMobile = useIsMobile();

  if (isMobile) {
    return (
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
                <div className={`p-0.5 rounded-full ${selected?.id === chat.id ? 'ring-2 ring-blue-500' : ''}`}>
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
    );
  }

  return (
    <div className="h-screen w-screen bg-gray-100 dark:bg-gray-900">
      <PanelGroup direction="horizontal" className="flex h-full w-full">
        <Panel defaultSize={25} minSize={20} maxSize={35} className="flex flex-col">
          <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-800">
            <div className="flex items-center gap-3">
              <Avatar src={imgAvatar} size="md" online={true} />
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
                className={`flex items-center gap-3 p-3 cursor-pointer rounded-xl transition-colors duration-200 ${
                  selected?.id === chat.id
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
            <ChatView id={selected.id.toString()} name={selected.name} />
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
  );
}
