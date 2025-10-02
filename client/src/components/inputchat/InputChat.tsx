import React, { useState, useRef, useEffect } from "react";

interface ChatInputProps {
    value?: string;
    placeholder?: string;
    disabled?: boolean;
    onSend?: (text: string, file?: File | null) => void;
    maxRows?: number;
}

export default function ChatInput({
    value = "",
    placeholder = "Nhập tin nhắn...",
    disabled = false,
    onSend,
}: ChatInputProps) {
    const [text, setText] = useState(value);
    const [file, setFile] = useState<File | null>(null);
    const [isSending, setIsSending] = useState(false);
    const textareaRef = useRef<HTMLTextAreaElement | null>(null);

    useEffect(() => {
        setText(value);
    }, [value]);

    const handleSend = async () => {
        if (disabled || isSending) return;
        const trimmed = text.trim();
        if (!trimmed && !file) return; // nothing to send

        setIsSending(true);
        try {
            await Promise.resolve(onSend ? onSend(trimmed, file) : Promise.resolve());
        } finally {
            setIsSending(false);
            setText("");
            setFile(null);
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    return (
        <div className="w-full h-full border-t bg-white dark:bg-slate-900 dark:border-slate-800 p-3 content-center justify-self-center-safe   ">
            <div className="flex gap-2 max-w-4xl mx-auto">
                <div className="flex bg-[#212121]/90 w-full  content-center rounded-s-xl rounded-lg border border-slate-700">
                    <div className="flex items-center gap-2">
                        <button
                            type="button"
                            className="p-2 rounded-md hover:bg-slate-100 dark:hover:bg-slate-800"
                            title="Emoji (placeholder)"
                            onClick={() => {
                                // placeholder: open emoji picker
                                setText((t) => t + "🙂");
                            }}
                        >
                            <svg className="w-6 h-6 text-slate-600 dark:text-slate-300" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <path d="M12 17c2.21 0 4-1.79 4-4H8c0 2.21 1.79 4 4 4z" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                                <path d="M9 9h.01M15 9h.01" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                                <path d="M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0z" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                            </svg>
                        </button>
                    </div>

                    <div className="flex-1">
                        <textarea
                            ref={textareaRef}
                            value={text}
                            onChange={(e) => setText(e.target.value)}
                            onKeyDown={handleKeyDown}
                            placeholder={placeholder}
                            disabled={disabled}
                            className="w-full resize-none h-[20px] min-h-[40px] max-h-[100px] overflow-hidden  p-3 text-slate-900 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-0 focus:border-transparent"
                            aria-label="Chat message input"
                        />
                    </div>
                </div>


                <div className="flex items-center gap-2">
                    <button
                        onClick={handleSend}
                        disabled={disabled || isSending}
                        className={`p-3 rounded-full flex items-center justify-center ${disabled || isSending
                            ? "bg-slate-300 text-slate-600 cursor-not-allowed"
                            : "bg-blue-600 text-white hover:bg-blue-700"
                            }`}
                        aria-label="Send message"
                        title="Send"
                    >
                        {isSending ? (
                            <svg className="w-5 h-5 animate-spin" viewBox="0 0 24 24">
                                <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" strokeLinecap="round" />
                            </svg>
                        ) : (
                            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <path d="M22 2L11 13" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                                <path d="M22 2L15 22L11 13L2 9L22 2Z" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                            </svg>
                        )}
                    </button>
                </div>
            </div>

        </div>
    );
}
