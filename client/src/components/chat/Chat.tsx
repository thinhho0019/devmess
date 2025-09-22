import React, { useEffect, useRef, useState } from "react";
import type { Chat, ChatProps } from "./types";
import { convertTimeMessage, getTimeIsoCurrent } from "../../utils/date"
import ChatInput from "../inputchat/InputChat";

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
                { user_id: "user_002", emoji: "👍" },
                { user_id: "user_003", emoji: "🔥" },
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
            message: "Oops, wrong file. Let me re-upload.",
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
    const [mess, setMess] = useState<Chat[]>(chats as Chat[]);
    const scrollRef = useRef<HTMLDivElement>(null)
    const onSend = (trimed: string, file: File | null) => {
        const newMessage: Chat = {
            id: "msg_005",
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
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [mess])
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
                <div className="h-full flex-col max-w-4xl mx-auto px-6 py-4 shadow">

                    {mess.map((c, i) => (
                        <div key={i} className={`flex ${c.user_id === "user_002" ? "justify-start" : "justify-end"}`}>
                            <div className={`w-fit max-w-[400px] text-sm bg-blue-500 mb-2 px-3 py-2 rounded-xl`}>
                                <div key={i} className="text-[15px]">{c.message}</div>
                                <div className="flex justify-end">
                                    <span className="text-[10px] text-gray-200">{convertTimeMessage(c.created_at, 7)}</span>
                                </div>
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