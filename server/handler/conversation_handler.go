package handler

import (
	"net/http"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FindUserConversationRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler(conversationService *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}
func (h *ConversationHandler) GetUserConversationsByUserID(c *gin.Context) {
	// 1️⃣ Lấy user từ context
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userValue.(*models.User)
	if !ok || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	// 2️⃣ Lấy limit
	limitStr := c.DefaultQuery("limit", "20")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// 3️⃣ Lấy before (timestamp)
	beforeStr := c.Query("before")
	var beforePtr *time.Time
	if beforeStr != "" {
		beforeTime, err := utils.ConvertMilisecondToTime(beforeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid before timestamp"})
			return
		}
		beforePtr = beforeTime
	}

	// 4️⃣ Gọi service
	conversations, err := h.conversationService.GetUserConversations(user.ID.String(), limitInt, beforePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5️⃣ Trả kết quả
	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
	})
}

func (h *ConversationHandler) FindConversationByUser(c *gin.Context) {
	var req FindUserConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userValue.(*models.User)
	if !ok || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	conversation, err := h.conversationService.FindConversationBytwoUserIDs(user.ID.String(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"conversation_id": conversation.ID,
	})
}

func (h *ConversationHandler) GetMessageByConversationID(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userValue.(*models.User)
	if !ok || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}
	conversationID := c.Query("conversation_id")
	limitStr := c.Query("limit")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	beforeStr := c.Query("before")
	var beforePtr *time.Time
	if beforeStr != "" {
		beforeTime, err := utils.ConvertMilisecondToTime(beforeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid before timestamp"})
			return
		}
		beforePtr = beforeTime
	}
	messages, err := h.conversationService.GetMessageByConversationID(conversationID, limitInt, beforePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func (h *ConversationHandler) SendMessage(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userValue.(*models.User)
	if !ok || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}
	var req struct {
		ConversationID string `json:"conversation_id" binding:"required"`
		Content        string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content cannot be empty"})
		return
	}

	message, err := h.conversationService.SendMessage(req.ConversationID, user.ID.String(), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
