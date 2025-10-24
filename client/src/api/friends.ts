 
import type { UserResponse } from "../types/UserResponse";
import api from "./api";
import type { AxiosError } from "axios";

// Types for friend API requests and responses
interface AddFriendRequest {
    user_id: string;
    friend_id: string;
}

interface AddFriendResponse {
    message: string;
    friendship_id?: string;
    status?: string;
}

interface ErrorResponse {
    error: string;
    message?: string;
}


export const fetchAddFriend = async (friend_id: string): Promise<AddFriendResponse> => {
    try {
        const dataReq: AddFriendRequest = {
            user_id: localStorage.getItem("user_id") || "",
            friend_id: friend_id,
        };
        console.log("➡️ Sending add friend request:", dataReq);
        const res = await api.post<AddFriendResponse>("/friendships/send-invite", dataReq);
        
        return res.data;
    } catch (error) {
        const axiosError = error as AxiosError<ErrorResponse>;
        console.error("❌ Add friend error:", axiosError.response?.data || axiosError.message);

        // Ném ra lỗi gọn gàng để component xử lý UI
        const errorMessage = axiosError.response?.data?.error ||
            axiosError.response?.data?.message ||
            "Failed to send friend request";
        throw new Error(errorMessage);
    }
};

export const fetchListInviteFriends = async (): Promise<UserResponse[]> => {
    try {
        const res = await api.get(`/friendships/list-invite-friends`);
        if (!res.data) {
            console.warn("⚠️ No data in friend invites response");
            return [];
        }
        const listUser: UserResponse[] = res.data.map((user: UserResponse) => {
            user.status = "pending";
            return user;
        });
        return listUser;
    } catch (error) {
        const axiosError = error as AxiosError<ErrorResponse>;
        console.error("❌ Fetch friend invites error:", axiosError.response?.data || axiosError.message);
        throw new Error("Failed to fetch friend invites");
    }
};
export const fetchListFriends = async (): Promise<UserResponse[]> => {
    try {
        const res = await api.get(`/friendships/list-friends`);
        if (!res.data) {
            console.warn("⚠️ No data in friend list response");
            return [];
        }
        const listUser: UserResponse[] = res.data.map((user: UserResponse) => {
            user.status = "friend";
            return user;
        });
        return listUser;
    } catch (error) {
        const axiosError = error as AxiosError<ErrorResponse>;
        console.error("❌ Fetch friends error:", axiosError.response?.data || axiosError.message);
        throw new Error("Failed to fetch friends");
    }
}

export const rejectFriendRequest = async (friend_id: string): Promise<AddFriendResponse> => {
    try {
        const dataReq: AddFriendRequest = {
            user_id: localStorage.getItem("user_id") || "",
            friend_id: friend_id,
        }
        console.log("➡️ Rejecting friend request:", dataReq);
        const res = await api.post<AddFriendResponse>("/friendships/reject-invite", dataReq);
        if(!res.data){
            throw new Error("No data in reject friend request response");
        }
        return res.data;
    }
    catch (error) {
        const axiosError = error as AxiosError<ErrorResponse>;
        console.error("❌ Reject friend request error:", axiosError.response?.data || axiosError.message);
        throw new Error("Failed to reject friend request");
    }
}

export const cancelFriendRequest = async (friend_id: string): Promise<AddFriendResponse> => {
    try {
        const dataReq: AddFriendRequest = {
            user_id: localStorage.getItem("user_id") || "",
            friend_id: friend_id,
        };
        console.log("➡️ Canceling friend request:", dataReq);
        const res = await api.post<AddFriendResponse>("/friendships/cancel-invite", dataReq);
        return res.data;
    } catch (error) {
        const axiosError = error as AxiosError<ErrorResponse>;
        console.error("❌ Cancel friend request error:", axiosError.response?.data || axiosError.message);

        const errorMessage = axiosError.response?.data?.error ||
            axiosError.response?.data?.message ||
            "Failed to cancel friend request";
        throw new Error(errorMessage);
    }
};

export const acceptFriendRequest = async (friend_id: string): Promise<AddFriendResponse> => {
    try {
        const dataReq: AddFriendRequest = {
            user_id: localStorage.getItem("user_id") || "",
            friend_id: friend_id,
        };
        console.log("➡️ Accepting friend request:", dataReq);
        const res = await api.post<AddFriendResponse>("/friendships/accept-invite", dataReq);
        return res.data;
    } catch (error) {
        const axiosError = error as AxiosError<ErrorResponse>;
        console.error("❌ Accept friend request error:", axiosError.response?.data || axiosError.message);

        const errorMessage = axiosError.response?.data?.error ||
            axiosError.response?.data?.message ||
            "Failed to accept friend request";
        throw new Error(errorMessage);
    }
};
