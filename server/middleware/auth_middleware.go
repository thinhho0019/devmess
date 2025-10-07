package middleware

import (
	"errors"

	"net/http"
	"project/handler"
	"project/service"

	"github.com/gin-gonic/gin"
)

func VerifyAccessToken(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid access token"})
		c.Abort()
		return
	}

	user, rf, err := service.VerifyAccessToken(accessToken)

	if err != nil && !errors.Is(err, service.ErrInvalidOrExpired) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify access token: " + err.Error()})
		c.Abort()
		return
	}
	 
	if user == nil {
		// check refresh token
		new_token, err := service.RefreshAccessToken(handler.GoogleOAuthConfig.ClientID, handler.GoogleOAuthConfig.ClientSecret, rf, handler.GoogleOAuthConfig)
		if new_token != nil && err == nil {
			c.SetCookie("access_token", new_token.AccessToken, 3600, "/", "", false, true)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token expired or invalid"})
		c.Abort()
		return
	}
 
	 
	c.Set("user", user)
	c.Next()
}
