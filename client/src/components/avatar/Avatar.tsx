import React from "react";
import type { AvatarProps } from "./types";
import imgAvatar from "../../assets/img.jpg"

const sizeClasses: Record<"sm" | "md" | "lg", string> = {
  sm: "w-8 h-8 text-sm",
  md: "w-12 h-12 text-base",
  lg: "w-16 h-16 text-lg",
};

const Avatar: React.FC<AvatarProps> = ({
  src,
  alt = "avatar",
  size = "md",
  rounded = true,
  online = false,
}) => {
  return (
    <div className="relative">
      <img
        src={src || imgAvatar}
        alt={alt}
        className={`${sizeClasses[size]} ${
          rounded ? "rounded-full" : "rounded-md"
        } object-cover border-2 border-gray-200 dark:border-gray-700`}
      />
      {online && (
        <span
          className={`absolute bottom-1 right-1 block h-3 w-3 rounded-full bg-green-500 border-2 border-white dark:border-gray-900`}
        />
      )}
    </div>
  );
};

export default Avatar;
