package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/models"
	"project/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	fmt.Println("Initializing Google OAuth with Client ID:", clientID)
	fmt.Println("Client Secret:", clientSecret)
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLoginHandler(c *gin.Context) {
	url := GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	fmt.Println("Redirecting to:", url)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func GoogleCallBackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'code' parameter"})
		return
	}
	fmt.Println("🔹 Received code:", code)
	// Đổi code lấy access token
	token, err := GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}
	fmt.Println("✅ Token exchanged successfully")

	// Lấy thông tin người dùng từ Google
	client := GoogleOAuthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Google returned status: %d", resp.StatusCode)})
		return
	}

	// Parse JSON response từ Google
	userInfo := &models.GoogleUserInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info: " + err.Error()})
		return
	}
	userInfo.AccessToken = token.AccessToken
	userInfo.RefreshToken = token.RefreshToken
	userInfo.ExpiresIn = int(token.Expiry.Unix())
	userInfo.TokenType = token.TokenType
	userInfo.Scope = token.Extra("scope").(string)
	fmt.Println("🔹 User info:", userInfo.Email, userInfo.Name)

	// Lưu token và user vào DB
	tokenModel, device, err := service.LoginWithGoogle(userInfo, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed: " + err.Error()})
		return
	}

	// TODO: Generate session or JWT token for app use
	c.SetCookie(
		"access_token",
		tokenModel.AccessToken,
		int(tokenModel.ExpiresAt),
		"/",
		"devmess.cloud",
		true,
		true,
	)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"user":    userInfo,
		"device":  device,
		"token":   tokenModel.AccessToken,
	})
}
