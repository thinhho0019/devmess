import { useState, useRef, useCallback, useEffect } from 'react';

export const useEmojiManager = () => {
    const [showEmojiPopup, setShowEmojiPopup] = useState(false);
    const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
    const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
    const messageRefs = useRef<{ [key: string]: HTMLDivElement | null }>({});

    const openEmojiPopup = useCallback((messageId: string, isCurrentUser: boolean) => {
        const element = messageRefs.current[messageId];
        if (element) {
            const rect = element.getBoundingClientRect();
            setPopupPosition({
                top: rect.top - 50, // Position above the message
                left: isCurrentUser ? rect.left - 200 : rect.right - 80, // Adjust based on who sent it
            });
            setSelectedMessageId(messageId);
            setShowEmojiPopup(true);
        }
    }, []);

    const closeEmojiPopup = useCallback(() => {
        setShowEmojiPopup(false);
        setSelectedMessageId(null);
    }, []);

    useEffect(() => {
        const handleClickOutside = () => {
            if (showEmojiPopup) {
                closeEmojiPopup();
            }
        };

        document.addEventListener('click', handleClickOutside);
        return () => {
            document.removeEventListener('click', handleClickOutside);
        };
    }, [showEmojiPopup, closeEmojiPopup]);


    return {
        showEmojiPopup,
        popupPosition,
        selectedMessageId,
        messageRefs,
        openEmojiPopup,
        closeEmojiPopup,
    };
};
