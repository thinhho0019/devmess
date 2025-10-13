import { motion, AnimatePresence } from "framer-motion";
import { useEffect, useState } from "react";
import { Search, UserPlus, Check, Loader2, X } from "lucide-react";
import type { UserResponse } from "../../types/UserResponse";
import { findUserByEmail } from "../../services/user/userService";
import { getAuthToken } from "../../utils/Auth";

interface PopupFindFriendsProps {
  show: boolean;
  onClose: () => void;
  onAddFriend: (email: string) => void;
  friendsList?: { name: string; email: string; avatar?: string; added?: boolean, token?: string }[];
}

export const PopupFindFriends = ({
  show,
  onClose,
  onAddFriend,
  friendsList = [],
}: PopupFindFriendsProps) => {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<typeof friendsList>([]);
  // per-user UI status: idle | sending | sended | canceled | accepted
  const [statuses, setStatuses] = useState<Record<string, string>>(() => {
    const map: Record<string, string> = {};
    friendsList.forEach((f) => {
      if (f.email && f.added) map[f.email] = "accepted";
    });
    return map;
  });

  useEffect(() => {
    // sync statuses when friendsList changes (e.g., initial data)
    setStatuses((prev) => {
      const next = { ...prev };
      friendsList.forEach((f) => {
        if (f.email && f.added) next[f.email] = "accepted";
      });
      return next;
    });
  }, [friendsList]);

  const setStatus = (email: string, status: string) => {
    setStatuses((s) => ({ ...s, [email]: status }));
  };

  const handleSearch = async () => {
    if (!query.trim()) {
      setResults([]);
      return;
    }
    const userFound: UserResponse = await findUserByEmail(query.trim());
    if (userFound && userFound.email) {
      setResults([{
        name: userFound.name || "Unknown",
        email: userFound.email,
        avatar: userFound.avatar,
        token: getAuthToken() || undefined,
        added: friendsList.some((f) => f.email === userFound.email)
      }]);
    }
  };

  const handleAddClick = async (email: string) => {
    try {
      setStatus(email, "sending");
      // support promise or void: await regardless
      await Promise.resolve(onAddFriend(email));
      setStatus(email, "sended");
    } catch (err) {
      console.error("Add friend failed", err);
      setStatus(email, "idle");
    }
  };

  const handleCancelClick = async (email: string) => {
    // local UI cancel — if you have an API to cancel request, call it here
    setStatus(email, "canceled");
  };

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
            {/* Header */}
            <h2 className="text-2xl font-semibold mb-4 text-center">
              🔍 Find Friends
            </h2>

            {/* Search box */}
            <div className="flex items-center gap-2 bg-gray-800 rounded-xl p-2 mb-5">
              <Search className="text-gray-400 w-5 h-5 ml-2" />
              <input
                type="text"
                placeholder="Enter friend's Gmail..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="bg-transparent text-white w-full outline-none px-2"
              />
              <button
                onClick={handleSearch}
                className="bg-indigo-600 hover:bg-indigo-500 px-3 py-1.5 rounded-lg text-sm font-medium"
              >
                Search
              </button>
            </div>

            {/* Search results */}
            <div className="max-h-64 overflow-y-auto space-y-3">
              {results.length === 0 ? (
                <p className="text-gray-400 text-center">
                  No users found. Try another Gmail.
                </p>
              ) : (
                results.map((user) => (
                  <div
                    key={user.email}
                    className="flex items-center justify-between bg-gray-800/50 p-3 rounded-xl hover:bg-gray-800 transition"
                  >
                    <div className="flex items-center gap-3">
                      {(() => {
                        const hasAvatar = !!user.avatar;
                        const avatarSrc = hasAvatar
                          ? `${import.meta.env.VITE_API_URL}/protected/?filename=${encodeURIComponent(
                            String(user.avatar)
                          )}&token=${encodeURIComponent(String(user.token ?? ""))}`
                          : `https://ui-avatars.com/api/?name=${encodeURIComponent(
                            String(user.name || "")
                          )}&background=1e293b&color=fff`;
                        return (
                          <img
                            src={avatarSrc}
                            alt="avatar"
                            className="w-10 h-10 rounded-full object-cover"
                          />
                        );
                      })()}
                      <div>
                        <h4 className="font-semibold text-[15px]">{user.name}</h4>
                        <p className="text-sm text-gray-400">{user.email}</p>
                      </div>
                    </div>
                    {user.added || statuses[user.email] === "accepted" ? (
                      <div className="flex items-center gap-1 text-green-400 text-sm font-medium">
                        <Check className="w-4 h-4" /> Added
                      </div>
                    ) : statuses[user.email] === "sending" ? (
                      <div className="flex items-center gap-2">
                        <button
                          disabled
                          className="flex items-center gap-2 bg-indigo-600/80 px-3 py-1.5 rounded-lg text-sm font-medium opacity-80"
                        >
                          <Loader2 className="w-4 h-4 animate-spin" />
                          Sending
                        </button>
                        <button
                          onClick={() => handleCancelClick(user.email)}
                          className="flex items-center gap-1 bg-gray-700 hover:bg-gray-600 px-3 py-1 rounded-lg text-sm font-medium"
                        >
                          <X className="w-4 h-4" />
                          Cancel
                        </button>
                      </div>
                    ) : statuses[user.email] === "sended" ? (
                      <div className="flex items-center gap-1 text-indigo-400 text-sm font-medium">
                        <Check className="w-4 h-4" /> Sent
                      </div>
                    ) : statuses[user.email] === "canceled" ? (
                      <button
                        onClick={() => handleAddClick(user.email)}
                        className="flex items-center gap-1 bg-indigo-600 hover:bg-indigo-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                      >
                        <UserPlus className="w-4 h-4" />
                        Retry
                      </button>
                    ) : (
                      <button
                        onClick={() => handleAddClick(user.email)}
                        className="flex items-center gap-1 bg-indigo-600 hover:bg-indigo-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                      >
                        <UserPlus className="w-4 h-4" />
                        Add
                      </button>
                    )}
                  </div>
                ))
              )}
            </div>

            {/* Buttons */}
            <div className="mt-6 flex justify-center">
              <button
                onClick={onClose}
                className="px-5 py-2 rounded-xl bg-gray-800 hover:bg-gray-700 text-gray-200 font-medium transition"
              >
                Close
              </button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
