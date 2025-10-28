import { motion, AnimatePresence } from "framer-motion";
import { useEffect, useState } from "react";
import { Search, UserPlus, Check, Loader2, X, Clock, Users, UserMinus, MessageCircle } from "lucide-react";
import type { UserResponse } from "../../types/UserResponse";
import { findUserByEmail } from "../../services/user/userService";
import { getAuthToken } from "../../utils/Auth";
import { fetchAddFriend, cancelFriendRequest, acceptFriendRequest, rejectFriendRequest } from "../../api/friends";
import { useListInvite } from "../../hooks/api/useListInvite";
import { useListFriend } from "../../hooks/api/useListFriend";
import { useSocket } from "../../hooks/socket/useSocket";

// Friend status types
const FriendStatus = {
    NoFriends: 'no_friends',
    Frind: 'friend',
    Pending: 'pending',
    Accepted: 'accepted',
    Canceled: 'canceled',
    Sending: 'sending',
} as const;

export type FriendStatus = typeof FriendStatus[keyof typeof FriendStatus];

// Tab types
type TabType = 'search' | 'incoming' | 'friends';

interface PopupFriendsManagerProps {
    show: boolean;
    onClose: () => void;
    onAddFriend?: (id: string) => void; // Optional fallback
    friendsList?: {
        id: string;
        name: string;
        email: string;
        avatar?: string;
        status?: FriendStatus;
        token?: string;
        receiverType?: 'incoming' | 'outgoing'; // Whether user sent or received the friend request
    }[];
}

export const PopupFriendsManager = ({
    show,
    onClose,
    onAddFriend,
    friendsList = [],
}: PopupFriendsManagerProps) => {
    const [activeTab, setActiveTab] = useState<TabType>('search');
    const { invites, loading, error, refetch } = useListInvite();
    const {listFriends,refetchFriends} = useListFriend();
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<typeof friendsList>([]);
    // eslint-disable-next-line
    const { lastMessage  } = useSocket();
    // receive websocket message to update friend list in real-time
    useEffect(() => {
        if (lastMessage !== null) {
            try {
                const eventData = JSON.parse(lastMessage.data);
                switch (eventData.type) {
                    case 'friend_invite': 
                        refetch();
                        return;
                    case 'update_friend':
                        refetchFriends();
                        return;
                }
                
            } catch (error) {
                console.log(error);
                console.log("Received non-JSON or malformed data:", lastMessage.data);
            }
        }
    }, [lastMessage]);
    // Track friend status per user ID
    const [statuses, setStatuses] = useState<Record<string, FriendStatus>>(() => {
        const map: Record<string, FriendStatus> = {};
        friendsList.forEach((f) => {
            if (f.id && f.status) {
                map[f.id] = f.status;
            }
        });
        return map;
    });

    useEffect(() => {
        // Sync statuses when friendsList, invites, or listFriends changes
        setStatuses((prev) => {
            const next = { ...prev };
            
            // Sync from friendsList prop
            friendsList.forEach((f) => {
                if (f.id && f.status) {
                    next[f.id] = f.status;
                }
            });
            
            // Add invites status (incoming requests)
            (invites || []).forEach((invite) => {
                if (invite.id) {
                    next[invite.id] = 'pending';
                }
            });
            
            // Add friends status (accepted friends)
            (listFriends || []).forEach((friend) => {
                if (friend.id) {
                    next[friend.id] = 'friend';
                }
            });
            
            return next;
        });
    }, [friendsList, invites, listFriends]);

    const setStatus = (userId: string, status: FriendStatus) => {
        setStatuses((s) => ({ ...s, [userId]: status }));
    };

    // Clear search when modal closes
    const handleClose = () => {
        setQuery("");
        setResults([]);
        setActiveTab('search');
        onClose();
    };

    const handleSearch = async () => {
        if (!query.trim()) {
            setResults([]);
            return;
        }

        try {
            const userFound: UserResponse = await findUserByEmail(query.trim());
            console.log("Search result:", userFound);
            if (userFound && userFound.email) {
                // Determine initial status based on friendsList
                const existingFriend = friendsList.find(f => f.email === userFound.email);
                console.log("Existing friend:", existingFriend);
                const initialStatus: FriendStatus = (userFound.status as FriendStatus) || 'no_friends';

                setResults([{
                    id: userFound.id || "",
                    name: userFound.name || "Unknown",
                    email: userFound.email,
                    avatar: userFound.avatar,
                    token: getAuthToken() || undefined,
                    status: initialStatus
                }]);
            }
        } catch (error) {
            console.error("Search failed:", error);
            setResults([]);
        }
    };
    const handleRejectRequest = async (userId: string) => {
        try {
            setStatus(userId, "sending");
            await rejectFriendRequest(userId);
            setStatus(userId, "no_friends");
            // Refresh invites list
            await refetch();
        } catch (err) {
            console.error("Reject friend request failed:", err);
            setStatus(userId, "pending");
        }
    };
    // Send friend request
    const handleSendRequest = async (userId: string) => {
        try {
            setStatus(userId, "sending");
            if (onAddFriend) {
                await Promise.resolve(onAddFriend(userId));
            } else {
                // Use direct API call
                await fetchAddFriend(userId);
            }

            setStatus(userId, "pending");
        } catch (err) {
            console.error("Send friend request failed:", err);
            setStatus(userId, "no_friends");
        }
    };

    // Cancel friend request
    const handleCancelRequest = async (userId: string) => {
        try {
            setStatus(userId, "sending");
            await cancelFriendRequest(userId);
            setStatus(userId, "no_friends");
            // Refresh invites list if we're rejecting an incoming request
            if (activeTab === 'incoming') {
                await refetch();
            }
        } catch (err) {
            console.error("Cancel friend request failed:", err);
            setStatus(userId, "pending");
        }
    };

    // Accept friend request
    const handleAcceptRequest = async (userId: string) => {
        try {
            setStatus(userId, "sending");
            await acceptFriendRequest(userId);
            setStatus(userId, "accepted");
            // Refresh invites list
            await refetch();
        } catch (err) {
            console.error("Accept friend request failed:", err);
            setStatus(userId, "pending");
        }
    };

    // Remove friend
    const handleRemoveFriend = async (userId: string) => {
        try {
            setStatus(userId, "sending");
            // TODO: Call API to remove friend
            // await removeFriend(userId);
            setStatus(userId, "no_friends");
        } catch (err) {
            console.error("Remove friend failed:", err);
            setStatus(userId, "accepted");
        }
    };

    // Start chat with friend
    const handleStartChat = (userId: string) => {
        // TODO: Navigate to chat with this user
        console.log("Starting chat with user:", userId);
        // Close modal and potentially navigate to chat
        handleClose();
    };

    // Get status for display
    const getUserStatus = (user: typeof results[0]): FriendStatus => {
        return statuses[user.id] || user.status || 'no_friends';
    };

    // Transform invites data to match component format (for incoming requests)
    const transformInviteToUser = (invite: UserResponse) => ({
        id: invite.id || "",
        name: invite.name || "Unknown",
        email: invite.email || "",
        avatar: invite.avatar,
        status: 'pending' as FriendStatus,
        token: getAuthToken() || undefined,
        receiverType: 'incoming' as const
    });

    // Transform friends data to match component format (for accepted friends)
    const transformFriendToUser = (friend: UserResponse) => ({
        id: friend.id || "",
        name: friend.name || "Unknown",
        email: friend.email || "",
        avatar: friend.avatar,
        status: 'friend' as FriendStatus, // Use 'friend' status for accepted friends
        token: getAuthToken() || undefined,
        receiverType: undefined
    });

    // Filter friends by status
    const acceptedFriends = (listFriends || []).map(transformFriendToUser);
    const pendingRequests = friendsList.filter(f => f.status === 'pending' && f.receiverType === 'outgoing');
    // Use data from useListInvite hook for incoming requests
    const incomingRequests = (invites || []).map(transformInviteToUser);

    // Render action button based on status
    const renderActionButton = (user: typeof results[0], isInFriendsList = false, isIncomingRequest = false) => {
        const status = getUserStatus(user);
        console.log(`Rendering button for ${user.email} with status: ${status}`);

        switch (status) {
            case 'accepted':
                return isInFriendsList ? (
                    <button
                        onClick={() => handleRemoveFriend(user.id)}
                        className="flex items-center gap-1 bg-red-600 hover:bg-red-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                    >
                        <UserMinus className="w-4 h-4" />
                        Hủy kết bạn
                    </button>
                ) : (
                    <div className="flex items-center gap-1 text-green-400 text-sm font-medium">
                        <Check className="w-4 h-4" /> Bạn bè
                    </div>
                );
            case 'friend':
                return(
                    <button
                        onClick={() => handleStartChat(user.id)}
                        className="flex items-center gap-1 bg-blue-600 hover:bg-blue-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                    >
                        <MessageCircle className="w-4 h-4" />
                        Trò chuyện
                    </button>
                )
            case 'pending':
                return isIncomingRequest ? (
                    <div className="flex items-center gap-1">
                        <button
                            onClick={() => handleAcceptRequest(user.id)}
                            className="flex items-center gap-1 bg-green-600 hover:bg-green-500 px-2 py-1 rounded text-xs font-medium"
                        >
                            <Check className="w-3 h-3" />
                            Chấp nhận
                        </button>
                        <button
                            onClick={() => handleRejectRequest(user.id)}
                            className="flex items-center gap-1 bg-red-600 hover:bg-red-500 px-2 py-1 rounded text-xs font-medium"
                        >
                            <X className="w-3 h-3" />
                            Từ chối
                        </button>
                    </div>
                ) : (
                    <div className="flex items-center gap-2">
                        <div className="flex items-center gap-1 text-yellow-400 text-sm font-medium">
                            <Clock className="w-4 h-4" />
                        </div>
                        <button
                            onClick={() => handleCancelRequest(user.id)}
                            className="flex items-center gap-1 bg-red-600 hover:bg-red-500 px-2 py-1 rounded text-xs font-medium"
                        >
                            <X className="w-3 h-3" />
                            Hủy
                        </button>
                    </div>
                );

            case 'canceled':
                return (
                    <div className="flex items-center gap-2">
                        <span className="text-gray-400 text-sm">Đã hủy</span>
                        <button
                            onClick={() => handleSendRequest(user.id)}
                            className="flex items-center gap-1 bg-indigo-600 hover:bg-indigo-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                        >
                            <UserPlus className="w-4 h-4" />
                            Gửi lại
                        </button>
                    </div>
                );

            case 'sending':
                return (
                    <div className="flex items-center gap-2 bg-indigo-600/80 px-3 py-1.5 rounded-lg text-sm font-medium opacity-80">
                        <Loader2 className="w-4 h-4 animate-spin" />
                        Đang xử lý...
                    </div>
                );

            case 'no_friends':
            default:
                return (
                    <button
                        onClick={() => handleSendRequest(user.id)}
                        className="flex items-center gap-1 bg-indigo-600 hover:bg-indigo-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                    >
                        <UserPlus className="w-4 h-4" />
                        Kết bạn
                    </button>
                );
        }
    };

    // Render user card
    const renderUserCard = (user: typeof results[0], isInFriendsList = false, isIncomingRequest = false) => (
        <div
            key={user.id}
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
            {renderActionButton(user, isInFriendsList, isIncomingRequest)}
        </div>
    );

    return (
        <AnimatePresence>
            {show && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    className="fixed inset-0 bg-black/60 backdrop-blur-sm flex justify-center items-center z-50"
                    onClick={handleClose}
                >
                    <motion.div
                        initial={{ scale: 0.8, y: 30 }}
                        animate={{ scale: 1, y: 0 }}
                        exit={{ scale: 0.8, y: 30 }}
                        transition={{ duration: 0.25 }}
                        onClick={(e) => e.stopPropagation()}
                        className="relative bg-[#111827] text-white rounded-3xl p-8 w-[90%] max-w-lg shadow-2xl border border-white/10"
                    >
                        {/* Header */}
                        <h2 className="text-2xl font-semibold mb-6 text-center">
                            👥 Quản lý bạn bè
                        </h2>

                        {/* Tabs */}
                        <div className="flex bg-gray-800 rounded-xl p-1 mb-6">
                            <button
                                onClick={() => setActiveTab('search')}
                                className={`flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-medium transition ${activeTab === 'search'
                                    ? 'bg-indigo-600 text-white'
                                    : 'text-gray-400 hover:text-white'
                                    }`}
                            >
                                <Search className="w-4 h-4" />
                                Tìm bạn
                            </button>
                            <button
                                onClick={() => setActiveTab('incoming')}
                                className={`flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-medium transition ${activeTab === 'incoming'
                                    ? 'bg-indigo-600 text-white'
                                    : 'text-gray-400 hover:text-white'
                                    }`}
                            >
                                <Clock className="w-4 h-4" />
                                Nhận lời mời ({incomingRequests.length})
                            </button>
                            <button
                                onClick={() => setActiveTab('friends')}
                                className={`flex-1 flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-medium transition ${activeTab === 'friends'
                                    ? 'bg-indigo-600 text-white'
                                    : 'text-gray-400 hover:text-white'
                                    }`}
                            >
                                <Users className="w-4 h-4" />
                                Bạn bè ({acceptedFriends.length})
                            </button>
                        </div>

                        {/* Tab Content */}
                        <div className="min-h-[300px]">
                            {activeTab === 'search' ? (
                                <motion.div
                                    initial={{ opacity: 0, x: -20 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    {/* Search box */}
                                    <div className="flex items-center gap-2 bg-gray-800 rounded-xl p-2 mb-5">
                                        <Search className="text-gray-400 w-5 h-5 ml-2" />
                                        <input
                                            type="text"
                                            placeholder="Nhập Gmail của bạn bè..."
                                            value={query}
                                            onChange={(e) => setQuery(e.target.value)}
                                            onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                                            className="bg-transparent text-white w-full outline-none px-2"
                                        />
                                        <button
                                            onClick={handleSearch}
                                            className="bg-indigo-600 hover:bg-indigo-500 px-3 py-1.5 rounded-lg text-sm font-medium"
                                        >
                                            Tìm
                                        </button>
                                    </div>

                                    {/* Search results */}
                                    <div className="max-h-64 overflow-y-auto space-y-3">
                                        {results.length === 0 ? (
                                            <p className="text-gray-400 text-center py-8">
                                                {query.trim() ? 'Không tìm thấy người dùng.' : 'Nhập Gmail để tìm kiếm bạn bè.'}
                                            </p>
                                        ) : (
                                            results.map((user) => renderUserCard(user))
                                        )}
                                    </div>

                                    {/* Pending requests section (outgoing only) */}
                                    {pendingRequests.length > 0 && (
                                        <div className="mt-6">
                                            <h3 className="text-lg font-semibold mb-3 text-yellow-400">
                                                Lời mời đã gửi ({pendingRequests.length})
                                            </h3>
                                            <div className="space-y-2 max-h-32 overflow-y-auto">
                                                {pendingRequests.map((user) => renderUserCard(user))}
                                            </div>
                                        </div>
                                    )}
                                </motion.div>
                            ) : activeTab === 'incoming' ? (
                                <motion.div
                                    initial={{ opacity: 0, x: 0 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    {/* Incoming requests */}
                                    <div className="max-h-80 overflow-y-auto space-y-3">
                                        {loading ? (
                                            <div className="text-center py-12">
                                                <Loader2 className="w-8 h-8 text-indigo-500 mx-auto mb-4 animate-spin" />
                                                <p className="text-gray-400">Đang tải lời mời...</p>
                                            </div>
                                        ) : error ? (
                                            <div className="text-center py-12">
                                                <X className="w-16 h-16 text-red-500 mx-auto mb-4" />
                                                <p className="text-red-400 text-lg">Lỗi tải dữ liệu</p>
                                                <p className="text-gray-500 text-sm mt-2">{error}</p>
                                                <button
                                                    onClick={refetch}
                                                    className="mt-4 bg-indigo-600 hover:bg-indigo-500 px-4 py-2 rounded-lg text-sm font-medium"
                                                >
                                                    Thử lại
                                                </button>
                                            </div>
                                        ) : incomingRequests.length === 0 ? (
                                            <div className="text-center py-12">
                                                <Clock className="w-16 h-16 text-gray-500 mx-auto mb-4" />
                                                <p className="text-gray-400 text-lg">Không có lời mời nào</p>
                                                <p className="text-gray-500 text-sm mt-2">
                                                    Khi có người gửi lời mời kết bạn, bạn sẽ thấy ở đây.
                                                </p>
                                            </div>
                                        ) : (
                                            incomingRequests.map((user) => renderUserCard(user, false, true))
                                        )}
                                    </div>
                                </motion.div>
                            ) : (
                                <motion.div
                                    initial={{ opacity: 0, x: 20 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    transition={{ duration: 0.2 }}
                                >
                                    {/* Friends list */}
                                    <div className="max-h-80 overflow-y-auto space-y-3">
                                        {acceptedFriends.length === 0 ? (
                                            <div className="text-center py-12">
                                                <Users className="w-16 h-16 text-gray-500 mx-auto mb-4" />
                                                <p className="text-gray-400 text-lg">Chưa có bạn bè nào</p>
                                                <p className="text-gray-500 text-sm mt-2">
                                                    Hãy chuyển sang tab "Tìm bạn" để tìm kiếm bạn bè mới!
                                                </p>
                                                <button
                                                    onClick={() => setActiveTab('search')}
                                                    className="mt-4 bg-indigo-600 hover:bg-indigo-500 px-4 py-2 rounded-lg text-sm font-medium"
                                                >
                                                    Tìm bạn bè
                                                </button>
                                            </div>
                                        ) : (
                                            acceptedFriends.map((user) => renderUserCard(user, true))
                                        )}
                                    </div>
                                </motion.div>
                            )}
                        </div>

                        {/* Close button */}
                        <div className="mt-6 flex justify-center">
                            <button
                                onClick={handleClose}
                                className="px-6 py-2 rounded-xl bg-gray-800 hover:bg-gray-700 text-gray-200 font-medium transition"
                            >
                                Đóng
                            </button>
                        </div>
                    </motion.div>
                </motion.div>
            )}
        </AnimatePresence>
    );
};