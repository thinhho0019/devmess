import React, { useState, useRef, useEffect } from "react";
import { FiPaperclip, FiSmile, FiX, FiSend, FiLoader } from "react-icons/fi";
import EmojiPicker from "../emoji/EmojiPicker";

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
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
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

    const handleEmojiSelect = (emoji: string) => {
        setText(prev => prev + emoji);
        // Focus back to textarea after selecting emoji
        setTimeout(() => textareaRef.current?.focus(), 100);
    };

    const toggleEmojiPicker = () => {
        setShowEmojiPicker(prev => !prev);
    };

    return (
        <div className="w-full h-full border-t bg-white dark:bg-slate-900 dark:border-slate-800 p-3 content-center justify-self-center-safe   ">
            <div className="flex gap-2 max-w-4xl mx-auto">
                <div className="flex bg-[#212121]/90 w-full  content-center rounded-s-xl rounded-lg border border-slate-700">
                    <div className="flex items-center gap-2 ml-2">

                        <label className="p-2 rounded-md hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer flex items-center" title="Attach file">
                            <input
                                type="file"
                                className="hidden"
                                disabled={disabled}
                                onChange={(e) => {
                                    const f = e.target.files?.[0] ?? null;
                                    setFile(f);
                                }}
                                aria-label="Attach file"
                            />
                            <FiPaperclip className="w-5 h-5 text-slate-600 dark:text-slate-300" />


                        </label>

                        {file && (
                            <div className="flex items-center gap-2 px-2">
                                <span className="text-xs text-slate-600 dark:text-slate-300 max-w-[140px] truncate" title={file.name}>
                                    {file.name}
                                </span>
                                <button
                                    type="button"
                                    onClick={() => setFile(null)}
                                    className="p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
                                    aria-label="Remove attached file"
                                    title="Remove file"
                                >
                                    <FiX className="w-4 h-4 text-slate-600 dark:text-slate-300" />
                                </button>
                            </div>
                        )}

                        <div className="relative">
                            <button
                                type="button"
                                className={`p-2 rounded-md transition-colors ${
                                    showEmojiPicker 
                                        ? "bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400" 
                                        : "hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
                                }`}
                                title="Emoji picker"
                                onClick={toggleEmojiPicker}
                                disabled={disabled}
                            >
                                <FiSmile className="w-5 h-5" />
                            </button>

                            <EmojiPicker
                                isOpen={showEmojiPicker}
                                onClose={() => setShowEmojiPicker(false)}
                                onEmojiSelect={handleEmojiSelect}
                            />
                        </div>
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
                            <FiLoader className="w-5 h-5 animate-spin" />
                        ) : (
                            <FiSend className="w-5 h-5" />
                        )}
                    </button>
                </div>
            </div>

        </div>
    );
}
