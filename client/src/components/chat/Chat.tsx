import React, { useEffect, useRef, useState } from "react";
import type { Chat, ChatProps, EmojiData, MessageReaction } from "./types";
import { convertTimeMessage, getTimeIsoCurrent } from "../../utils/date"
import ChatInput from "../inputchat/InputChat";
import { HeartIcon as HeartSolid } from "@heroicons/react/24/solid";
import EmojiBox from "../emoji/emojiBox";
import { EmojiPopup } from "../emoji/emojiPopUp";
const emojis = [
    { emoji: "😀", type: "smile", description: "Cười vui vẻ" },
    { emoji: "😂", type: "laugh", description: "Cười chảy nước mắt" },
    { emoji: "😍", type: "love", description: "Thả tim, yêu thích" },
    { emoji: "😎", type: "cool", description: "Ngầu lòi" },
    { emoji: "👍", type: "like", description: "Đồng ý, thích" },
    { emoji: "🔥", type: "fire", description: "Tuyệt vời, nóng bỏng" },
    { emoji: "🎉", type: "party", description: "Ăn mừng" },
    { emoji: "❤️", type: "heart", description: "Thả tim, yêu thương" },
];
const ChatView: React.FC<ChatProps> = ({
    id,
    name = "thinhho",
    img = "",
    chats = [
        {
            id: "msg_001",
            message: "Hello team, please review the PR #123 Hello team, please review the PR #123Hello team, please review the PR #123Hello team, please review the PR #123Hello team, please review the PR #123",
            created_at: "2025-09-21T10:00:00Z",
            user_id: "user_001",
            type: "text",
            status: "read",
            mentions: ["user_002", "user_003"],
            hashtags: ["projectX", "urgent"],
            links: ["https://github.com/example/repo/pull/123"],
            reactions: [
                { user_ids: ["user_002"], emoji: "👍", count: 20, type: "like" },
                { user_ids: ["user_003"], emoji: "🔥", count: 20, type: "fire" }
            ],
            thread_id: "thread_001",
        },
        {
            id: "msg_002",
            message: "Here is the design file.",
            created_at: "2025-09-21T10:05:00Z",
            user_id: "user_002",
            type: "image",
            status: "delivered",
            attachments: [
                {
                    id: "att_001",
                    type: "image",
                    url: "https://cdn.example.com/design.png",
                    name: "UI Design",
                },
            ],
            reply: {
                message_id: "msg_001",
                user_id: "user_002",
                snippet: "Hello team, please review...",
            },
        },
        {
            id: "msg_003",
            message: "O",
            created_at: "2025-09-21T10:10:00Z",
            updated_at: "2025-09-21T10:12:00Z",
            user_id: "user_002",
            type: "text",
            status: "sent",
            is_edited: true,
        },
        {
            id: "msg_004",
            message: "This message was deleted",
            created_at: "2025-09-21T10:15:00Z",
            user_id: "user_003",
            is_deleted: true,
            deleted_at: "2025-09-21T10:16:00Z",
            reactions: [
                { user_ids: ["user_002"], emoji: "👍", count: 20, type: "like" },
                { user_ids: ["user_003"], emoji: "🔥", count: 20, type: "fire" },
                { user_ids: ["user_002"], emoji: "❤️", count: 1, type: "heart" },
            ],
        },
        {
            id: "msg_005",
            message: "Spam content link http://spam.example.com",
            created_at: "2025-09-21T10:20:00Z",
            user_id: "user_999",
            type: "text",
            is_flagged: true,
            flagged_reason: "Spam link detected",
        }
    ]
}) => {
    const [showEmoji, setShowEmoji] = useState(false);
    const [posEmoji, setPosEmoji] = useState({ top: 0, left: 0 });
    const [idItemClick, setIdItemClick] = useState("");
    const boxRef = useRef<{ [key: string]: HTMLDivElement | null }>({});;
    const [mess, setMess] = useState<Chat[]>(chats as Chat[]);
    const [clickedHeart, setClickedHeart] = useState<boolean>(false);
    const scrollRef = useRef<HTMLDivElement>(null);
    const handleClickedHeart = (id?: string) => {

        setMess(prevChats => {

            const newChats = [...prevChats];


            const idx = newChats.findIndex(c => c.id === id);

            if (idx === -1) return prevChats;


            const chat = { ...newChats[idx] };


            if (!chat.reactions) {
                chat.reactions = [
                    { user_ids: ["user_002"], emoji: "❤️", count: 1, type: "heart" },
                ];
            } else {
                const heart = chat.reactions.find(c => c.type === "heart");
                console.log(heart?.count)
                if (heart) {
                    console.log(heart?.count)
                    heart.count += 1;
                    console.log(heart?.count)
                    if (!heart.user_ids.includes("user_002")) {
                        heart.user_ids = [...heart.user_ids, "user_002"];
                    }
                } else {
                    chat.reactions = [
                        ...chat.reactions,
                        { user_ids: ["user_002"], emoji: "❤️", count: 1, type: "heart" },
                    ];
                }
            }

            // gán lại vào mảng clone
            newChats[idx] = chat;

            // trả về mảng mới -> React render lại
            return newChats;
        });
        setClickedHeart(!clickedHeart);
    }
    const onSend = (trimed: string, file: File | null) => {
        console.log(file);
        const new_user_id = "msg_" + crypto.randomUUID().replace(/-/g, '');
        const newMessage: Chat = {
            id: new_user_id,
            message: trimed,
            created_at: getTimeIsoCurrent(),
            user_id: "user_999",
            type: "text",
            status: "sent",
            mentions: [],
            hashtags: [],
            links: [],
            reactions: [],
            is_flagged: true,
            flagged_reason: "Spam link detected",
        };

        setMess((prev) => [...prev, newMessage])

        console.log(trimed);
    };
    const handleClickEmoji = (emoji: MessageReaction, id?: string) => {
        setMess(prevChats => {

            const newChats = [...prevChats];


            const idx = newChats.findIndex(c => c.id === id);

            if (idx === -1) return prevChats;


            const chat = { ...newChats[idx] };


            if (!chat.reactions) {
                chat.reactions = [
                    { user_ids: ["user_002"], emoji: emoji.emoji, count: 1, type: emoji.type },
                ];
            } else {
                const heart = chat.reactions.find(c => c.type === emoji.type);

                if (heart) {

                    heart.count += 1;

                    if (!heart.user_ids.includes("user_002")) {
                        heart.user_ids = [...heart.user_ids, "user_002"];
                    }
                } else {
                    chat.reactions = [
                        ...chat.reactions,
                        { user_ids: ["user_002"], emoji: emoji.emoji, count: 1, type: emoji.type },
                    ];
                }
            }

            // gán lại vào mảng clone
            newChats[idx] = chat;

            // trả về mảng mới -> React render lại
            return newChats;
        });
    }
    const handleSelectEmoji = async (emoji: EmojiData, id: string) => {
        setMess(prevChats => {

            const newChats = [...prevChats];


            const idx = newChats.findIndex(c => c.id === id);

            if (idx === -1) return prevChats;


            const chat = { ...newChats[idx] };


            if (!chat.reactions) {
                chat.reactions = [
                    { user_ids: ["user_002"], emoji: emoji.emoji, count: 1, type: emoji.type },
                ];
            } else {
                const heart = chat.reactions.find(c => c.type === emoji.type);

                if (heart) {

                    heart.count += 1;

                    if (!heart.user_ids.includes("user_002")) {
                        heart.user_ids = [...heart.user_ids, "user_002"];
                    }
                } else {
                    chat.reactions = [
                        ...chat.reactions,
                        { user_ids: ["user_002"], emoji: emoji.emoji, count: 1, type: emoji.type },
                    ];
                }
            }

            // gán lại vào mảng clone
            newChats[idx] = chat;

            // trả về mảng mới -> React render lại
            return newChats;
        });
        await new Promise(r => setTimeout(r, 500));
        setShowEmoji(false);
    };
    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            const el = boxRef.current[idItemClick];
            if (boxRef.current && !el?.contains(e.target as Node)) {
                setShowEmoji(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [mess]);
    const handleOpenEmoji = (pre: boolean, id: string, is_send: boolean) => {
        setIdItemClick(id);
        setShowEmoji(!pre);

        // 🔹 Lấy ref theo id
        const el = boxRef.current[id];
        if (el) {
            const rect = el.getBoundingClientRect();
            console.log(rect);

            setPosEmoji({
                top: rect.top - 40,
                left: is_send ? rect.left + 100 : rect.left - 200,
            });
        }
    };
    return (
        <div
            key={id}
            className="p-2 h-full grid"
            style={{ gridTemplateRows: "10% 80% 10%" }}
        >
            <div className="">
                <h3 className="font-semibold">{name}</h3>
                {img && <img src={img} alt={name} className="w-10 h-10 rounded-full" />}
            </div>
            <div ref={scrollRef} className="w-full  overflow-y-auto chat-scroll">
                <div className="h-full flex-col max-w-4xl mx-auto px-6 py-4 shadow space-y-1" >

                    {mess.map((c, i) => (

                        <div key={i} className={`flex ${c.user_id === "user_002" ? "justify-start" : "justify-end"}  `}>

                            <div ref={(el) => { boxRef.current[c.id] = el; }} onClick={() => handleOpenEmoji(showEmoji, c.id, c.user_id === "user_002")} className={`relative w-fit max-w-[400px] text-sm bg-blue-500 pl-2 pt-2 pr-2 pb-0.2 rounded-xl`}>
                                <div key={i} className="text-[15px]">{c.message}</div>
                                <EmojiBox emojis={c.reactions} onSelect={(emoji) => handleClickEmoji(emoji, c.id)} />
                                <div className="flex justify-end">
                                    <span className="text-[10px] text-gray-200">{convertTimeMessage(c.created_at, 7)}</span>
                                </div>
                                <span
                                    className={`absolute ${c.user_id === "user_002" ? "-bottom-2 -right-3" : "-bottom-2 -left-3"
                                        } h-6 w-6 flex items-center justify-center`}
                                >
                                    <HeartSolid onClick={(e) => {
                                        e.stopPropagation();
                                        handleClickedHeart(c.id);
                                    }} className={`z-10 ${clickedHeart ? "h-10 w-10" : "h-6 w-6"}  ease-in-out text-red-500 opacity-0 transition-opacity duration-500 hover:opacity-100`} fill="currentColor" />
                                </span>
                                {showEmoji && c.id === idItemClick && (
                                    <EmojiPopup show={showEmoji} position={posEmoji} emojis={emojis} onSelect={(emoji) => handleSelectEmoji(emoji, c.id)} />
                                )
                                }


                            </div>

                        </div>
                    ))}
                </div>
            </div>

            <div className="bg-blue-200 h-full">
                <ChatInput onSend={onSend as (text: string, file?: File | null) => void} />
            </div>
        </div>
    );
}

export default ChatView;