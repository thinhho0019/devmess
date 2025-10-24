package service

import (
	"errors"
	"fmt"
	"log"
	"project/models"
	"project/repository"

	"time"

	"project/utils"
)

type FriendService struct {
	repoFriend repository.FriendshipRepository
	repoUser   repository.UserRepository
}

func NewInitFriendService(repoFriend repository.FriendshipRepository, repoUser repository.UserRepository) *FriendService {
	return &FriendService{
		repoFriend: repoFriend,
		repoUser:   repoUser,
	}
}

func (f *FriendService) SendInviteFriend(requesterID string, friendID string) error {
	if requesterID == "" || friendID == "" {
		return errors.New("all fields are required")
	}
	// Chuyển đổi requesterID và friendID từ string sang uuid.UUID
	requesterUUID, err := utils.StringToUUID(requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester ID: %w", err)
	}
	friendUUID, err := utils.StringToUUID(friendID)
	if err != nil {
		return fmt.Errorf("invalid friend ID: %w", err)
	}
	friendShip := models.Friendship{
		UserID:      requesterUUID,
		FriendID:    friendUUID,
		RequestedBy: requesterUUID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// Kiểm tra nếu đã có mối quan hệ bạn bè hoặc lời mời kết bạn tồn tại
	existing, err := f.repoFriend.GetFriendshipBetweenUsers(requesterUUID, friendUUID)
	if err != nil {
		return fmt.Errorf("failed to check existing friendship: %w", err)

	}
	// Nếu đã có mối quan hệ bạn bè hoặc lời mời kết bạn tồn tại
	if existing != nil {
		if existing.Status == "no_friend" || existing.Status == "" {

			existing.Status = "pending"
			existing.FriendID = friendUUID
			existing.UserID = requesterUUID
			existing.RequestedBy = requesterUUID
			existing.UpdatedAt = time.Now()
			log.Println("[log exsting friend]", existing)
			// Lưu thay đổi vào database
			if _, err := f.repoFriend.SaveFriendship(existing); err != nil {
				return fmt.Errorf("failed to update friendship: %w", err)
			}
			return nil
		}
		return errors.New("friendship or invitation already exists")
	}
	// Lưu vào database
	if _, err := f.repoFriend.CreateFriendship(&friendShip); err != nil {
		return fmt.Errorf("failed to create friendship: %w", err)
	}
	return nil
}

func (f *FriendService) CancelInviteFriend(requesterID string, friendID string) error {
	if requesterID == "" || friendID == "" {
		return errors.New("all fields are required")
	}
	// Chuyển đổi requesterID và friendID từ string sang uuid.UUID
	requesterUUID, err := utils.StringToUUID(requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester ID: %w", err)
	}
	friendUUID, err := utils.StringToUUID(friendID)
	if err != nil {
		return fmt.Errorf("invalid friend ID: %w", err)
	}
	// Kiểm tra nếu lời mời kết bạn tồn tại
	existing, err := f.repoFriend.GetFriendshipBetweenUsers(requesterUUID, friendUUID)
	if err != nil {
		return fmt.Errorf("failed to check existing friendship: %w", err)
	}
	// Nếu không có lời mời kết bạn tồn tại
	if existing == nil || existing.Status != "pending" || existing.RequestedBy != requesterUUID {
		return errors.New("no pending invitation found to cancel")
	}
	// change status to "no friend"
	existing.Status = "no_friend"
	if _, err := f.repoFriend.UpdateFriendship(existing); err != nil {
		return fmt.Errorf("failed to update friendship status: %w", err)
	}
	return nil
}
func (f *FriendService) GetListsFriendInvite(userID string) ([]models.User, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	// Chuyển đổi userID từ string sang uuid.UUID
	userUUID, err := utils.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	// Lấy danh sách lời mời kết bạn
	invites, err := f.repoFriend.ListPendingRequests(userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending requests: %w", err)
	}
	// Trả về danh sách lời mời kết bạn
	// filter list only return user friends
	var users []models.User
	for i := range invites {
		if invites[i].User.ID == userUUID {
			users = append(users, invites[i].Friend)
			continue
		}
		users = append(users, invites[i].User)
	}
	return users, nil
}

func (f *FriendService) AcceptInviteFriend(requesterID string, friendID string) error {
	if requesterID == "" || friendID == "" {
		return errors.New("all fields are required")
	}
	// Chuyển đổi requesterID và friendID từ string sang uuid.UUID
	requesterUUID, err := utils.StringToUUID(requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester ID: %w", err)
	}
	friendUUID, err := utils.StringToUUID(friendID)
	if err != nil {
		return fmt.Errorf("invalid friend ID: %w", err)
	}
	// Kiểm tra nếu lời mời kết bạn tồn tại
	existing, err := f.repoFriend.GetFriendshipBetweenUsers(requesterUUID, friendUUID)
	if err != nil {
		return fmt.Errorf("failed to check existing friendship: %w", err)
	}
	// Nếu không có lời mời kết bạn tồn tại
	if existing == nil || existing.Status != "pending" || existing.RequestedBy != friendUUID {
		return errors.New("no pending invitation found to accept")
	}
	// change status to "friends"
	existing.Status = "friend"
	if _, err := f.repoFriend.SaveFriendship(existing); err != nil {
		return fmt.Errorf("failed to update friendship status: %w", err)
	}

	return nil
}
func (f *FriendService) RejectInviteFriend(requesterID string, friendID string) error {
	if requesterID == "" || friendID == "" {
		return errors.New("all fields are required")
	}
	// Chuyển đổi requesterID và friendID từ string sang uuid.UUID
	requesterUUID, err := utils.StringToUUID(requesterID)
	if err != nil {
		return fmt.Errorf("invalid requester ID: %w", err)
	}
	friendUUID, err := utils.StringToUUID(friendID)
	if err != nil {
		return fmt.Errorf("invalid friend ID: %w", err)
	}
	// Kiểm tra nếu lời mời kết bạn tồn tại
	existing, err := f.repoFriend.GetFriendshipBetweenUsers(requesterUUID, friendUUID)
	if err != nil {
		return fmt.Errorf("failed to check existing friendship: %w", err)
	}
	// Nếu không có lời mời kết bạn tồn tại
	if existing == nil || existing.Status != "pending" || existing.RequestedBy != friendUUID {
		return errors.New("no pending invitation found to reject")
	}
	// change status to "no friend"
	existing.Status = "no_friend"
	if _, err := f.repoFriend.UpdateFriendship(existing); err != nil {
		return fmt.Errorf("failed to update friendship status: %w", err)
	}
	return nil
}
func (f *FriendService) GetListFriends(userID string) ([]models.User, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	// Chuyển đổi userID từ string sang uuid.UUID
	userUUID, err := utils.StringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	// Lấy danh sách bạn bè
	friends, err := f.repoFriend.ListFriends(userUUID, "friend")
	if err != nil {
		return nil, fmt.Errorf("failed to list friends: %w", err)
	}
	fmt.Println("List friend", friends)
	// Trả về danh sách bạn bè
	var users []models.User
	for i := range friends {
		if friends[i].User.ID == userUUID {
			users = append(users, friends[i].Friend)
		} else {
			users = append(users, friends[i].User)
		}
	}
	return users, nil
}

func (f *FriendService) RemoveFriend(userID string, friendID string) error {
	if userID == "" || friendID == "" {
		return errors.New("all fields are required")
	}
	// Chuyển đổi userID và friendID từ string sang uuid.UUID
	userUUID, err := utils.StringToUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	friendUUID, err := utils.StringToUUID(friendID)
	if err != nil {
		return fmt.Errorf("invalid friend ID: %w", err)
	}
	// Kiểm tra nếu mối quan hệ bạn bè tồn tại
	existing, err := f.repoFriend.GetFriendshipBetweenUsers(userUUID, friendUUID)
	if err != nil {
		return fmt.Errorf("failed to check existing friendship: %w", err)
	}
	// Nếu không có mối quan hệ bạn bè tồn tại
	if existing == nil || existing.Status != "friend" {
		return errors.New("no friendship found to remove")
	}
	// change status to "no friend"
	existing.Status = "no_friend"
	if _, err := f.repoFriend.UpdateFriendship(existing); err != nil {
		return fmt.Errorf("failed to update friendship status: %w", err)
	}
	return nil
}
