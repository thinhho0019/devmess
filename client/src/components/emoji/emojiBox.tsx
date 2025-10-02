import type { EmojiBoxProps } from "./types";




export default function EmojiBox({ emojis, onSelect, title }: EmojiBoxProps) {
    console.log(emojis?.length)
    return (
        <div className=" ">
            {title && (
                <h3 className="text-gray-700 font-semibold mb-2">{title}</h3>
            )}
            
            <div className={`flex gap-3 text-2xs`}>
                {emojis?.map((emoji) => (
                    <button
                        key={emoji.id}
                        onClick={() => onSelect && onSelect(emoji)}
                        className="hover:scale-110 text-xs min-w-[50px] max-w-[100px] flex items-center justify-center  transition-transform duration-150 bg-white rounded-2xl pl-1 pr-1 pt-1 pb-1  "
                    >
                        {emoji.emoji}
                        {emoji.count && (
                            <p className=" text-black font-bold  ">{emoji.count} </p>
                        )}
                    </button>
                ))}
            </div>
        </div>
    );
}
