import { useState } from "react";
import { FaSearch } from "react-icons/fa";
import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import imgAvatar from "../assets/img.jpg"
import { Avatar } from "../components/avatar";
import  { ChatView  } from "../components/chat";
interface ChatMessage {
  id: number;
  name: string;
  lastMessage: string;
  time: string;
}

const chats: ChatMessage[] = [
  { id: 1, name: "Temp", lastMessage: "SEED created the group «Temp»", time: "Aug 14" },
  { id: 2, name: "Telegram", lastMessage: "Login code: ••••", time: "Jul 27" },
  { id: 3, name: "GROUP_FILE_BACKUP", lastMessage: "backup_live_novel.rar", time: "Mar 11" },
  { id: 4, name: "Wallet", lastMessage: "Say hello to Wallet in Vietnam", time: "Mar 6" },
  { id: 5, name: "More chat", lastMessage: "Message content here", time: "Sep 15" },
  // ...thêm nhiều chat test scroll
];

export default function Home() {
  const [selected, setSelected] = useState<ChatMessage | null>(null);

  return (
    <div className="h-screen w-screen"> <PanelGroup direction="horizontal" className="flex h-full w-full overflow-hidden bg-gray-900 text-white">
      <Panel defaultSize={25} minSize={18} maxSize={30}>
        <div className="flex flex-col w-full h-full border-r border-gray-700">
          {/* Search */}
          <div className="flex items-center p-3">
            <FaSearch className="w-5 h-5 text-gray-400 mr-2 flex-shrink-0" />
            <input
              type="text"
              placeholder="Search"
              className="flex-1 h-8 rounded bg-gray-800 text-white placeholder-gray-400 px-3 focus:outline-none"
            />
          </div>

          {/* Chat list */}
          <div className="flex-1 min-h-0 p-2 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-700 scrollbar-track-gray-900">
            {chats.map(chat => (
              <div
                key={chat.id}
                onClick={() => setSelected(chat)}
                className={`flex relative z-10 cursor-pointer rounded-xl p-2  hover:bg-[#5b52da]
    ${selected?.id === chat.id ? "bg-[#5b52da]" : "hover:bg-gray-700 transition-colors duration-230"}`}
              >
                <Avatar src={imgAvatar} size="lg" />
                <div className="w-full h-full pl-2">

                  <div className="flex justify-between items-center min-w-0">
                    <span className="lg:text-base sm:text-sm font-semibold truncate overflow-hidden whitespace-nowrap">
                      {chat.name}
                    </span>
                    <span className="text-xs text-gray-400">{chat.time}</span>
                  </div>
                  <p className="lg:text-base sm:text-sm text-gray-400 truncate overflow-hidden whitespace-nowrap">
                    {chat.lastMessage}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </Panel>
      <PanelResizeHandle />
      <Panel minSize={30}>
        <ChatView id="ak47" name={selected?.name} />
        {/* <div className="flex-1 min-w-0 flex items-center justify-center overflow-hidden">
          {selected ? (
            <div className="text-center px-4">
              <h2 className="text-2xl font-bold">{selected.name}</h2>
              <p className="text-gray-400 mt-2">Chat content here...</p>
            </div>
          ) : (
            <p className="text-gray-600">Select a chat to start messaging</p>
          )}
        </div> */}
      </Panel>
    </PanelGroup></div>

  );
}
