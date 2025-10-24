package handler

import (
	"encoding/json"
	"project/models"
	"project/service"
	"project/utils"
	"project/websocket"

	"github.com/gin-gonic/gin"
)

type InviteFriendRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FriendID string `json:"friend_id" binding:"required"`
}

type FriendHandler struct {
	friendService       *service.FriendService
	conversationService *service.ConversationService
	hub                 *websocket.Hub
}

func NewFriendHandler(friendService *service.FriendService, hub *websocket.Hub,conversationService *service.ConversationService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
		conversationService: conversationService,
		hub:           hub,
	}
}

func (r *FriendHandler) SendInviteFriend(c *gin.Context) {
	user, _ := c.Get("user")
	println("user", user)
	if user == nil {
		c.JSON(400, gin.H{"error": "user not found in context"})
		return
	}
	userModel := user.(*models.User)
	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidUUID(req.UserID) || !utils.IsValidUUID(req.FriendID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	err := r.friendService.SendInviteFriend(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var message_send = map[string]interface{}{
		"type": "friend_invite",
		"user": userModel,
	}
	payload, err := json.Marshal(message_send)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to encode message"})
		return
	}
	r.hub.NotifyInviteFriend(req.FriendID, json.RawMessage(payload))
	c.JSON(200, gin.H{"message": "Invite sent successfully"})
}

func (r *FriendHandler) CancelInviteFriend(c *gin.Context) {
	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidUUID(req.UserID) || !utils.IsValidUUID(req.FriendID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	err := r.friendService.CancelInviteFriend(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Invite canceled successfully"})
}
func (r *FriendHandler) GetListInviteFriend(c *gin.Context) {
	user, exits := c.Get("user")
	if !exits {
		c.JSON(400, gin.H{"error": "user_id query parameter is required"})
		return
	}
	model := user.(*models.User)
	userID := model.ID.String()
	if !utils.IsValidUUID(userID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	list, err := r.friendService.GetListsFriendInvite(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (r *FriendHandler) AcceptInviteFriend(c *gin.Context) {
	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidUUID(req.UserID) || !utils.IsValidUUID(req.FriendID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	err := r.friendService.AcceptInviteFriend(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var message_send = map[string]interface{}{
		"type": "update_friend",
	}
	var direct_json = map[string]interface{}{
		"type": "direct_conversation",
	}
	payload, err := json.Marshal(message_send)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to encode message"})
		return
	}
	payload_conversation, err := json.Marshal(direct_json)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to encode message"})
		return
	}
	r.hub.SendToUsers([]string{req.UserID, req.FriendID}, json.RawMessage(payload))
	// create conversation
	_, err = r.conversationService.CreateDirectConversation(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	r.hub.SendToUsers([]string{req.UserID, req.FriendID}, json.RawMessage(payload_conversation))
	c.JSON(200, gin.H{"message": "Invite accepted successfully"})
}

func (r *FriendHandler) RejectInviteFriend(c *gin.Context) {
	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidUUID(req.UserID) || !utils.IsValidUUID(req.FriendID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	err := r.friendService.RejectInviteFriend(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Invite rejected successfully"})
}
func (r *FriendHandler) RemoveFriend(c *gin.Context) {
	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsValidUUID(req.UserID) || !utils.IsValidUUID(req.FriendID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	err := r.friendService.RemoveFriend(req.UserID, req.FriendID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Friend removed successfully"})
}

func (r *FriendHandler) GetListFriends(c *gin.Context) {
	user, exits := c.Get("user")
	if !exits {
		c.JSON(400, gin.H{"error": "user_id query parameter is required"})
		return
	}
	model := user.(*models.User)
	userID := model.ID.String()
	if !utils.IsValidUUID(userID) {
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}
	list, err := r.friendService.GetListFriends(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}
