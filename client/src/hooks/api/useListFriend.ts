import {   useEffect, useState } from "react";
import { fetchListFriends } from "../../api/friends";
import type { UserResponse } from "../../types/UserResponse";


export function useListFriend() {
    // Hook logic for managing friend list
    const [listFriends, setListFriends] = useState<UserResponse[]>([]);
    const [activeTab, setActiveTab] = useState('search');
    const [isLoading, setIsLoading] = useState(false);


    const fetchFriends = async () => {
        setIsLoading(true); 
        try {
            const data:UserResponse[] = await fetchListFriends();
            //filter change status friends
            const filteredData:UserResponse[] = data.map((user: UserResponse) => {
                user.status = "friend";
                return user;
            });
            console.log("Fetched friends:", filteredData);
            setListFriends(filteredData);
        }
        catch (error) {
            console.error("Failed to fetch friends:", error);
        }

        finally {
            setIsLoading(false);
        }
    };
    useEffect(() => {
        fetchFriends();
    }, []);
    
    return {
        listFriends,
        activeTab,
        setActiveTab,
        isLoading,
        refetchFriends: fetchFriends,
    };

}
