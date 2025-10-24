import { useEffect, useState } from "react";
import { fetchListInviteFriends } from "../../api/friends";
import type { UserResponse } from "../../types/UserResponse";


export function useListInvite() {
  // Giả sử bạn có một API để lấy danh sách lời mời kết bạn
  const [invites, setInvites] = useState<UserResponse[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

    const fetchInvites = async () => {
        try {
            setLoading(true);
            setError(null);
             
            const response:UserResponse[] = await fetchListInviteFriends();
            setInvites(response);
        // eslint-disable-next-line
        } catch (err: any) {
            setError(err.message || "Lỗi khi tải lời mời kết bạn");
        }
        finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchInvites();
    }, []);

  return { invites, loading, error, refetch: fetchInvites };
}