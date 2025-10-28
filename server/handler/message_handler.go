package handler

import (
	"encoding/json"
	"net/http"
	"project/models"
	"project/service"
	"project/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHanlder struct {
	messageService     service.MessageService
	participantService service.ParticipantService
	hub                *websocket.Hub
}

func NewMessageHandler(messageService service.MessageService, hub *websocket.Hub, participantService service.ParticipantService) *MessageHanlder {
	return &MessageHanlder{
		messageService:     messageService,
		hub:                hub,
		participantService: participantService,
	}
}

func (h *MessageHanlder) SendMessageToConversation(c *gin.Context) {
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
	message, err := h.messageService.SendMessageToConversation(user.ID.String(), req.ConversationID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	participants, err := h.participantService.GetParticipantsByConversationID(req.ConversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if participants == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No participants found"})
		return
	}

	var userid_to string
	for _, participant := range *participants {
		if participant.UserID == user.ID {
			continue
		}
		userid_to = participant.UserID.String()
	}
	payload := map[string]interface{}{
		"type":         "receive_message",
		"message":      message,
		"conversation": req.ConversationID,
		"sender_id":    user.ID.String(),
	}
	payloadd, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode message"})
		return
	}
	println("payload", string(payloadd))
	h.hub.SendToUser(userid_to, payloadd)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h *MessageHanlder) GetAllMessageToConversation(c *gin.Context) {

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
	limit := c.Query("limit")
	//convert limit to int
	// default limit is 50
	limitInt := 50
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id is required"})
		return
	}
	messages, err := h.messageService.GetAllMessageToConversation(conversationID, user.ID.String(), limitInt, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}
