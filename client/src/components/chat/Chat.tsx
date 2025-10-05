
import type { ChatProps, EmojiData } from "./types";
import { convertTimeMessage } from "../../utils/date";
import ChatInput from "../inputchat/InputChat";
import { EmojiPopup } from "../emoji/emojiPopUp";
import { useChatManager } from "../../hooks/chat/useChatManager";
import { useEmojiManager } from "../../hooks/chat/useEmojiManager";
import EmojiBox from "../emoji/emojiBox";

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
            <div className="">
                <h3 className="font-semibold text-white">{name}</h3>
                {img && <img src={img} alt={name} className="w-10 h-10 rounded-full" />}
            </div>
            <div ref={chatContainerRef} className="w-full overflow-y-auto chat-scroll">
                <div className="h-full flex-col max-w-4xl mx-auto px-6 py-4 shadow space-y-1">
                    {messages.map((c) => {
                        const isCurrentUser = c.user_id === current_user_id;
                        return (
                            <div key={c.id} className={`flex ${isCurrentUser ? "justify-end" : "justify-start"}`}>
                                <div
                                    ref={(el) => { messageRefs.current[c.id] = el; }}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        openEmojiPopup(c.id, isCurrentUser);
                                    }}
                                    className="relative w-fit max-w-[400px] text-sm bg-blue-500 pl-2 pt-2 pr-2 pb-0.2 rounded-xl cursor-pointer"
                                >
                                    <div className="text-[15px]">{c.message}</div>
                                    <EmojiBox emojis={c.reactions} onSelect={(emoji) => addReaction(c.id, emoji)} />
                                    <div className="flex justify-end items-center">
                                        <span className="text-[10px] text-gray-200 mr-1">{convertTimeMessage(c.created_at, 7)}</span>
                                        {isCurrentUser && (
                                            <span className="text-[10px]">
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

            <div className="bg-blue-200 h-full">
                <ChatInput onSend={sendMessage} />
            </div>
        </div>
    );
}

export default ChatView;