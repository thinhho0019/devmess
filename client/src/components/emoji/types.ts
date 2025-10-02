import type { MessageReaction } from "../chat";

export interface EmojiBoxProps {
  emojis?: MessageReaction[];                        
  onSelect?: (emoji: MessageReaction) => void;       
  title?: string;                        
}