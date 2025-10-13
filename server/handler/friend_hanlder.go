package handler

import (
	"project/models"
	"project/repository"
	"project/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type InviteFriendRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FriendID string `json:"friend_id" binding:"required"`
}

func SendInviteFriend(c *gin.Context) {
	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidUUID(req.UserID) || !utils.IsValidUUID(req.FriendID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	userID, err := utils.StringToUUID(req.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user_id"})
		return
	}
	friendID, err := utils.StringToUUID(req.FriendID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid friend_id"})
		return
	}
	friendShip := models.Friendship{
		UserID:      userID,
		FriendID:    friendID,
		RequestedBy: userID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// Kiểm tra nếu đã có mối quan hệ bạn bè hoặc lời mời kết bạn tồn tại
	existing, err := repository.NewFriendshipRepository().GetFriendshipBetweenUsers(userID, friendID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check existing friendship: " + err.Error()})
		return
	}
	// Nếu đã có mối quan hệ bạn bè hoặc lời mời kết bạn tồn tại
	if existing != nil {
		c.JSON(400, gin.H{"error": "Friendship or invitation already exists"})
		return
	}
	// Lưu vào database
	if _, err := repository.NewFriendshipRepository().CreateFriendship(&friendShip); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create friendship: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Invite friend sent", "data": friendShip})
}
