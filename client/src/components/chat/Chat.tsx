
import type { ChatProps, EmojiData } from "./types";
import { convertTimeMessage } from "../../utils/date";
import ChatInput from "../inputchat/InputChat";
import { EmojiPopup } from "../emoji/emojiPopUp";
import { useChatManager } from "../../hooks/chat/useChatManager";
import { useEmojiManager } from "../../hooks/chat/useEmojiManager";
import EmojiBox from "../emoji/emojiBox";
import { Phone, Video, MoreHorizontal } from "lucide-react";

const emojis: EmojiData[] = [
    { emoji: "😀", type: "smile", description: "Cười vui vẻ" },
    { emoji: "😂", type: "laugh", description: "Cười chảy nước mắt" },
    { emoji: "😍", type: "love", description: "Thả tim, yêu thích" },
    { emoji: "😎", type: "cool", description: "Ngầu lòi" },
    { emoji: "👍", type: "like", description: "Đồng ý, thích" },
    { emoji: "🔥", type: "fire", description: "Tuyệt vời, nóng bỏng" },
    { emoji: "🎉", type: "party", description: "Ăn mừng" },
    { emoji: "❤️", type: "heart", description: "Thả tim, yêu thương" },
];
const ramdomUserId = Math.random().toString(36).substring(2, 15);
const ChatView: React.FC<ChatProps> = ({
    id,
    name = "thinhho",
    img = "",
    chats = [],
    is_mobile = false,
    current_user_id = ramdomUserId // Example current user
}) => {
    const { messages, chatContainerRef, sendMessage, addReaction } = useChatManager(chats, current_user_id);
    const { showEmojiPopup, popupPosition, selectedMessageId, messageRefs, openEmojiPopup, closeEmojiPopup } = useEmojiManager();

    const handleEmojiSelect = (emoji: EmojiData) => {
        if (selectedMessageId) {
            addReaction(selectedMessageId, emoji);
        }
        closeEmojiPopup();
    };

    return (
        <div
            key={id}
            className="p-2 h-full grid"
            style={{ gridTemplateRows: "10% 80% 10%" }}
        >
            <div className="flex items-center justify-between px-4 py-2 bg-gradient-to-r from-gray-900/60 to-gray-800/60 backdrop-blur-lg border-b border-white/10 rounded-xl mb-2 shadow-sm">
                <div className="flex items-center gap-3">
                    <div className="relative">
                        <img
                            src={img || "https://ui-avatars.com/api/?name=" + name}
                            alt={name}
                            className="w-10 h-10 rounded-full border border-gray-700 shadow-md object-cover"
                        />
                        {/* Online indicator */}
                        <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-gray-900 rounded-full"></span>
                    </div>

                    <div>
                        <h3 className="font-semibold text-white text-sm sm:text-base flex items-center gap-1">
                            {name}
                            <span className="text-xs text-gray-400 font-normal">• online</span>
                        </h3>
                        {!is_mobile && (
                            <p className="text-xs text-gray-400">Chat securely with {name}</p>)}

                    </div>
                </div>

                {/* Right-side action buttons (optional) */}
                <div className="flex items-center gap-3">
                    <button
                        className="p-2 rounded-lg hover:bg-gray-700/40 transition flex items-center justify-center"
                        title="Voice Call"
                        aria-label="Voice Call"
                    >
                        <Phone className="w-5 h-5 text-gray-200" />
                    </button>
                    <button
                        className="p-2 rounded-lg hover:bg-gray-700/40 transition flex items-center justify-center"
                        title="Video Call"
                        aria-label="Video Call"
                    >
                        <Video className="w-5 h-5 text-gray-200" />
                    </button>
                    <button
                        className="p-2 rounded-lg hover:bg-gray-700/40 transition flex items-center justify-center"
                        title="More Options"
                        aria-label="More Options"
                    >
                        <MoreHorizontal className="w-5 h-5 text-gray-200" />
                    </button>
                </div>
            </div>
            <div ref={chatContainerRef} className="w-full overflow-y-auto chat-scroll">
                <div className="h-full flex-col max-w-4xl mx-auto px-6 py-4  ">
                    {messages.map((c, idx) => {
                        const isCurrentUser = c.user_id === current_user_id;
                        const bubbleBase = "relative w-fit max-w-[520px] text-sm rounded-2xl px-4 py-3 shadow-lg cursor-pointer transform transition-all";
                        // Muted/darker palettes for a calmer 'fantasy' look
                        const bubbleClasses = isCurrentUser
                            ? `${bubbleBase} bg-gradient-to-br from-slate-800 to-slate-700 text-gray-100 self-end ring-1 ring-black/20`
                            : `${bubbleBase} bg-gradient-to-br from-indigo-900 to-indigo-800 text-gray-100 self-start ring-1 ring-black/20`;

                        return (
                            <div key={c.id} className={`flex ${isCurrentUser ? "justify-end" : "justify-start"} py-1`}>
                                <div className="relative">
                                    <div
                                        ref={(el) => { messageRefs.current[c.id] = el; }}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            openEmojiPopup(c.id, isCurrentUser);
                                        }}
                                        className={bubbleClasses}
                                    >
                                        {/* tail removed as requested */}
                                        <div className="text-[15px] leading-relaxed tracking-wide font-medium">{c.message}</div>
                                        <div className="mt-2">
                                            <EmojiBox emojis={c.reactions} onSelect={(emoji) => addReaction(c.id, emoji)} />
                                        </div>
                                        <div className="flex justify-end items-center mt-2 space-x-2">
                                            <span className="text-[10px] text-white/70">{convertTimeMessage(c.created_at, 7)}</span>
                                            {isCurrentUser && idx === messages.length - 1 && (
                                                <span className="text-[10px] text-white/90">
                                                    {c.status === 'sending' && '🕒'}
                                                    {c.status === 'sent' && '✓'}
                                                    {c.status === 'delivered' && '✓✓'}
                                                    {c.status === 'read' && <span style={{ color: '#34B7F1' }}>✓✓</span>}
                                                    {c.status === 'failed' && '❌'}
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {showEmojiPopup && selectedMessageId && (
                <EmojiPopup
                    show={showEmojiPopup}
                    position={popupPosition}
                    emojis={emojis}
                    onSelect={handleEmojiSelect}
                />
            )}
            <ChatInput onSend={sendMessage} />

        </div>
    );
}

export default ChatView;