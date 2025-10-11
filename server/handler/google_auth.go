package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	// init repo redis save cache accesstoken

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'code' parameter"})
		return
	}

	_, tokenModel, _, err := service.HandleGoogleCallback(
		code,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		GoogleOAuthConfig,
	)
	if err != nil {
		redirectURLError := fmt.Sprintf("/error?msg=%s", url.QueryEscape(err.Error()))
		c.Redirect(http.StatusTemporaryRedirect, redirectURLError)
		return
	}

	redirectURL := fmt.Sprintf(
		"%s/auth/success?token=%s",
		os.Getenv("DEFAULT_URL"),
		tokenModel.AccessToken,
	)
	// save accesstoken to redis

	fmt.Println("Redirecting to:", redirectURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

 
