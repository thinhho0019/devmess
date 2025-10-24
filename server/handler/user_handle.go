package handler

import (
	"fmt"
	"log"
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
	// Lấy query params
	email := c.Query("email")
	if email == "" {
		log.Println("[FindUserWithStatusFriends] Missing email query parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email   query parameter is required"})
		return
	}

	// Lấy user từ context
	userValue, exists := c.Get("user")
	if !exists {
		log.Println("[FindUserWithStatusFriends] User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userValue.(*models.User)
	if !ok {
		log.Printf("[FindUserWithStatusFriends] Invalid user type: %T", userValue)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	log.Printf("[FindUserWithStatusFriends] Authenticated user: %s (%s)", user.Name, user.Email)

	// Check không search chính mình
	if user.Email == email {
		log.Printf("[FindUserWithStatusFriends] User tried to search themselves: %s", email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot search for yourself"})
		return
	}

	// Check email rỗng

	// Call service
	result, err := h.userService.FindUserWithStatusFriend(email, user.ID.String())
	if err != nil {
		log.Printf("[FindUserWithStatusFriends] Error finding user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Error finding user: " + err.Error()})
		return
	}

	// Response
	if result != nil {
		log.Printf("[FindUserWithStatusFriends] User found: %+v", result)
		c.JSON(http.StatusOK, result)
		return
	}

	log.Printf("[FindUserWithStatusFriends] User not found for email: %s", email)
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func (h *UserHandler) FindUserByEmail(c *gin.Context) {
	email := c.Query("email")
	fmt.Println("email user find:" + email)
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
