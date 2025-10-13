package handler

import (
	"net/http"
	"project/models"
	"project/repository"

	"github.com/gin-gonic/gin"
)

func FindUserByEmail(c *gin.Context) {
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
	user, err := repository.NewUserRepository().GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
