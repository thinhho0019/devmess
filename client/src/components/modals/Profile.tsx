import { motion, AnimatePresence } from "framer-motion";
import { useRef } from "react";
import { Camera, LogOut } from "lucide-react";
import type { UserResponse } from "../../types/UserResponse";
import { Avatar } from "../avatar";
import { useImage } from "../../hooks/api/useImage";
import { logout } from "../../utils/Auth";

export const PopupProfile = ({ show, onClose, user, is_profile_owner, onAvatarChange }: {
    show: boolean;
    onClose: () => void;
    user?: UserResponse;
    is_profile_owner?: boolean;
    
    onAvatarChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}) => {
  const fileRef = useRef<HTMLInputElement | null>(null);
  const {src} = useImage(user?.avatar || "");
  return (
    <AnimatePresence>
      {show && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 bg-black/60 backdrop-blur-sm flex justify-center items-center z-50"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.8, y: 30 }}
            animate={{ scale: 1, y: 0 }}
            exit={{ scale: 0.8, y: 30 }}
            transition={{ duration: 0.25 }}
            onClick={(e) => e.stopPropagation()}
            className="relative bg-[#111827] text-white rounded-3xl p-8 w-[90%] max-w-md shadow-2xl border border-white/10"
          >
            {/* Avatar */}
            <div className="flex flex-col items-center gap-4">
              <div className="relative group">
                <Avatar src={src || undefined} size="md" online={false} />
                <button
                  onClick={() => fileRef.current?.click()}
                  className="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover:opacity-100 rounded-full transition"
                >
                  <Camera className="w-6 h-6 text-white" />
                </button>
                <input
                  type="file"
                  accept="image/*"
                  className="hidden"
                  ref={fileRef}
                  onChange={onAvatarChange}
                />
              </div>

              {/* Name & Gmail */}
              <div className="text-center space-y-1">
                <h2 className="text-xl font-semibold">{user?.name || "Unnamed User"}</h2>
                <p className="text-gray-400">{user?.email || "no-email@example.com"}</p>
              </div>

              {/* Buttons */}
              <div className="mt-6 flex gap-3">
                <button
                  onClick={onClose}
                  className="px-5 py-2 rounded-xl bg-gray-800 hover:bg-gray-700 text-gray-200 font-medium transition"
                >
                  Close
                </button>
               {is_profile_owner && (<button
                  onClick={() => {
                    // Clear auth and reload app; close modal first for immediate UI feedback
                    onClose();
                    logout(true);
                  }}
                  className="px-4 py-2 rounded-xl bg-red-600 hover:bg-red-500 text-white font-medium transition flex items-center gap-2"
                >
                  <LogOut className="w-4 h-4" />
                  Đăng xuất
                </button>)}
                {/* <button
                  onClick={() => alert("Saved!")}
                  className="px-5 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-medium transition"
                >
                  Save
                </button> */}
              </div>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
