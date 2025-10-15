package handler

import (
	"net/http"
	"project/models"

	"project/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) FindUserWithStatusFriends(c *gin.Context) {
	email := c.Query("email")
	user_id := c.Query("user_id")
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userValue.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}
	if user.Email == email {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot search for yourself"})
		return
	}
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or UserID query parameter is required"})
		return
	}
	// service check email and status friends
	result, err := h.userService.FindUserWithStatusFriend(email, user_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error find user :" + err.Error()})
		return
	}
	if result != nil {
		c.JSON(http.StatusAccepted, result)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func (h *UserHandler) FindUserByEmail(c *gin.Context) {
	email := c.Query("email")
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userValue.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}
	if user.Email == email {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot search for yourself"})
		return
	}
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"name":       user.Name,
		"avatar":     user.Avatar,
	})
}
