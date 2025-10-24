export interface UserResponse {
  id?: string;
  email?: string;
  name?: string;
  avatar?: string;
  createdAt?: string;
  updatedAt?: string;
  status?: string; // e.g., 'no_friends', 'pending', 'accepted', etc.
}