import { motion, AnimatePresence } from "framer-motion";
import { useState } from "react";
import { Search, UserPlus, Check } from "lucide-react";

interface PopupFindFriendsProps {
  show: boolean;
  onClose: () => void;
  onAddFriend: (email: string) => void;
  friendsList?: { name: string; email: string; avatar?: string; added?: boolean }[];
}

export const PopupFindFriends = ({
  show,
  onClose,
  onAddFriend,
  friendsList = [],
}: PopupFindFriendsProps) => {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<typeof friendsList>([]);

  const handleSearch = () => {
    if (!query.trim()) {
      setResults([]);
      return;
    }

    // 🔍 Fake filter — thay bằng API thực tế trong backend
    const filtered = friendsList.filter((f) =>
      f.email.toLowerCase().includes(query.toLowerCase())
    );
    setResults(filtered);
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
                results.map((user, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between bg-gray-800/50 p-3 rounded-xl hover:bg-gray-800 transition"
                  >
                    <div className="flex items-center gap-3">
                      <img
                        src={
                          user.avatar ||
                          `https://ui-avatars.com/api/?name=${user.name}&background=1e293b&color=fff`
                        }
                        alt="avatar"
                        className="w-10 h-10 rounded-full object-cover"
                      />
                      <div>
                        <h4 className="font-semibold text-[15px]">{user.name}</h4>
                        <p className="text-sm text-gray-400">{user.email}</p>
                      </div>
                    </div>

                    {user.added ? (
                      <div className="flex items-center gap-1 text-green-400 text-sm font-medium">
                        <Check className="w-4 h-4" /> Added
                      </div>
                    ) : (
                      <button
                        onClick={() => onAddFriend(user.email)}
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
