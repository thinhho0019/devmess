import { useState, useRef } from 'react';

export const useEmojiPopup = () => {
  const [showEmojiPopup, setShowEmojiPopup] = useState(false);
  const [emojiPopupPosition, setEmojiPopupPosition] = useState({ top: 0, left: 0 });
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const messageRefs = useRef<{ [key: string]: HTMLDivElement | null }>({});

  const openEmojiPopup = (messageId: string, event: React.MouseEvent<HTMLElement, MouseEvent>) => {
    const rect = event.currentTarget.getBoundingClientRect();
    // Adjust position based on sender/receiver if needed
    setEmojiPopupPosition({ top: rect.top - 60, left: rect.left - 100 });
    setSelectedMessageId(messageId);
    setShowEmojiPopup(true);
  };

  const closeEmojiPopup = () => {
    setShowEmojiPopup(false);
    setSelectedMessageId(null);
  };

  return {
    showEmojiPopup,
    emojiPopupPosition,
    selectedMessageId,
    openEmojiPopup,
    closeEmojiPopup,
    messageRefs,
  };
};
