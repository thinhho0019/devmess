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


type AuthGoogleHandler struct {
	authService *service.AuthService
	GoogleOAuthConfig *oauth2.Config
}

func NewAuthGoogleHandler(authService *service.AuthService, GoogleOAuthConfig *oauth2.Config) *AuthGoogleHandler {
	return &AuthGoogleHandler{authService: authService, GoogleOAuthConfig: GoogleOAuthConfig}
}

func (h *AuthGoogleHandler) InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	fmt.Println("Initializing Google OAuth with Client ID:", clientID)
	fmt.Println("Client Secret:", clientSecret)
	h.GoogleOAuthConfig = &oauth2.Config{
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

func (h *AuthGoogleHandler) GoogleLoginHandler(c *gin.Context) {
	url := h.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	fmt.Println("Redirecting to:", url)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func (h *AuthGoogleHandler) GoogleCallBackHandler(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'code' parameter"})
		return
	}

 

	_, tokenModel, _, err := h.authService.HandleGoogleCallback(
		code,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		h.GoogleOAuthConfig,
	)
	if err != nil {
		// log error server-side để debug, trả lỗi cho client bằng redirect hoặc json
		fmt.Println("Google callback error:", err)
		redirectURLError := fmt.Sprintf("/error?msg=%s", url.QueryEscape(err.Error()))
		c.Redirect(http.StatusTemporaryRedirect, redirectURLError)
		return
	}

	redirectURL := fmt.Sprintf(
		"%s/auth/success?token=%s",
		os.Getenv("DEFAULT_URL"),
		tokenModel.AccessToken,
	)

	fmt.Println("Redirecting to:", redirectURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
