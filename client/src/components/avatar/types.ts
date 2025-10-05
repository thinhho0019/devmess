export interface AvatarProps {
  src?: string;          // link ảnh avatar
  alt?: string;          // alt text
  size?: "sm" | "md" | "lg"; // kích thước
  rounded?: boolean;     // bo tròn hay không
  online?: boolean;      // trạng thái online
}

export interface AvatarGroupProps {
  avatars: AvatarProps[];
  max?: number; // số avatar hiển thị tối đa
}