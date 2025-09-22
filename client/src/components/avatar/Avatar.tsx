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
}) => {
  return (
    <img
      src={src || imgAvatar}
      alt={alt}
      className={`${sizeClasses[size]} ${
        rounded ? "rounded-full" : "rounded-md"
      } object-cover border border-gray-700`}
    />
  );
};

export default Avatar;
