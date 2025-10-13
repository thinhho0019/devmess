import { motion, AnimatePresence } from "framer-motion";
import { useRef } from "react";
import { Camera } from "lucide-react";
import type { UserResponse } from "../../types/UserResponse";

export const PopupProfile = ({ show, onClose, user, onAvatarChange }: {
    show: boolean;
    onClose: () => void;
    user?: UserResponse;
    onAvatarChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}) => {
  const fileRef = useRef<HTMLInputElement | null>(null);

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
                <img
                  src={user?.avatar || "https://ui-avatars.com/api/?name=User&background=111827&color=fff"}
                  alt="avatar"
                  className="w-28 h-28 rounded-full object-cover border-4 border-white/10 shadow-lg"
                />
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
                <button
                  onClick={() => alert("Saved!")}
                  className="px-5 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-medium transition"
                >
                  Save
                </button>
              </div>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
