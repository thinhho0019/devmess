package middleware

import (
	"errors"
	"strings"

	"net/http"

	"project/service"

	"github.com/gin-gonic/gin"
)

func VerifyAccessToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}
	// "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
		return
	}
	accessToken := parts[1]

	user, _, err := service.VerifyAccessToken(accessToken)

	if err != nil && !errors.Is(err, service.ErrInvalidOrExpired) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify access token: " + err.Error()})
		c.Abort()
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token expired or invalid"})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
