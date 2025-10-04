import { createPortal } from "react-dom";
import type { EmojiData } from "../chat";
import { motion, AnimatePresence } from "framer-motion";
interface EmojiPopupProps {
    show: boolean;
    position: { top: number; left: number };
    emojis: EmojiData[];
    onSelect: (emoji: EmojiData) => void;
}

export const EmojiPopup: React.FC<EmojiPopupProps> = ({ show, position, emojis, onSelect }) => {
    if (!show) return null;

    return createPortal(
        <AnimatePresence>
            {show && (
                <motion.div
                    initial={{ opacity: 0, scale: 0.8, y: 10 }}
                    animate={{ opacity: 1, scale: 1, y: 0 }}
                    exit={{ opacity: 0, scale: 0.8, y: 10 }}
                    transition={{ duration: 0.25, ease: "easeOut" }}
                    style={{
                        position: "fixed",
                        top: position.top,
                        left: position.left,
                        zIndex: 9999,
                    }}
                    className="grid grid-cols-8 gap-2 bg-black pl-2 pr-2 shadow-lg rounded-xl"
                    onClick={(e) => e.stopPropagation()}
                >
                    {emojis.map((emoji: any, idx: number) => (
                        <motion.button
                            key={idx}
                            whileHover={{ scale: 1.7, rotate: 15 }}
                            whileTap={{ scale: 1.4 }}
                            transition={{ duration: 0.08, ease: "easeOut" }}
                            onMouseDown={(e) => {
                                e.stopPropagation();
                                onSelect(emoji);
                            }}
                            className="transition-transform h-15 w-8 "
                        >
                            {emoji.emoji}
                        </motion.button>
                    ))}
                </motion.div>
            )}
        </AnimatePresence>,
        document.body
    );
};
