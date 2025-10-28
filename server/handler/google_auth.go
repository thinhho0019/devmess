package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"project/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthGoogleHandler struct {
	authService       *service.AuthService
	GoogleOAuthConfig *oauth2.Config
}

func NewAuthGoogleHandler(authService *service.AuthService, GoogleOAuthConfig *oauth2.Config) *AuthGoogleHandler {
	handler := &AuthGoogleHandler{
		authService:       authService,
		GoogleOAuthConfig: GoogleOAuthConfig,
	}

	// Log OAuth configuration on startup
	log.Println("✅ [GoogleAuth] Handler initialized")
	log.Printf("   Client ID: %s", GoogleOAuthConfig.ClientID)
	log.Printf("   Redirect URL: %s", GoogleOAuthConfig.RedirectURL)

	return handler
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
	// Store return URL if provided
	returnURL := c.Query("return_url")
	if returnURL != "" {
		// Secure cookie for HTTPS
		isSecure := os.Getenv("ENV") == "production"
		c.SetCookie("return_url", returnURL, 300, "/", "", isSecure, true)
		log.Printf("🔖 [GoogleAuth] Stored return_url: %s", returnURL)
	}

	// Generate OAuth URL
	authURL := h.GoogleOAuthConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	log.Printf("🔗 [GoogleAuth] Redirecting to Google OAuth: %s", authURL)
	http.Redirect(c.Writer, c.Request, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthGoogleHandler) GoogleCallBackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	log.Printf("📥 [GoogleAuth] Callback received - Code: %s, State: %s",
		truncateString(code, 20), state)

	// Validate authorization code
	if code == "" {
		log.Println("❌ [GoogleAuth] Missing authorization code")
		h.redirectToFrontendError(c, "Thiếu mã xác thực từ Google")
		return
	}

	// Process OAuth callback
	log.Println("🔄 [GoogleAuth] Processing OAuth callback...")
	user, tokenModel, device, err := h.authService.HandleGoogleCallback(
		code,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		h.GoogleOAuthConfig,
	)

	if err != nil {
		log.Printf("❌ [GoogleAuth] Callback error: %v", err)
		h.redirectToFrontendError(c, "Đăng nhập thất bại: "+err.Error())
		return
	}

	log.Printf("✅ [GoogleAuth] Login successful - User: %s, Device: %s",
		user.Email, device.ID)

	// Get return URL from cookie
	returnURL, err := c.Cookie("return_url")
	if err != nil || returnURL == "" {
		returnURL = "/dashboard" // default page after login
	}

	// Clear return URL cookie
	isSecure := os.Getenv("ENV") == "production"
	c.SetCookie("return_url", "", -1, "/", "", isSecure, true)

	// Get frontend URL from environment
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Println("⚠️ [GoogleAuth] FRONTEND_URL not set, using DEFAULT_URL")
		frontendURL = os.Getenv("DEFAULT_URL")
	}

	// Construct redirect URL with token and return path
	redirectURL := fmt.Sprintf(
		"%s/auth/success?token=%s&refresh_token=%s&return_url=%s",
		frontendURL,
		tokenModel.AccessToken,
		tokenModel.RefreshToken,
		url.QueryEscape(returnURL),
	)

	log.Printf("🔗 [GoogleAuth] Redirecting to: %s", maskToken(redirectURL))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// redirectToFrontendError redirects to frontend error page with message
func (h *AuthGoogleHandler) redirectToFrontendError(c *gin.Context, errMsg string) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = os.Getenv("DEFAULT_URL")
	}

	redirectURL := fmt.Sprintf(
		"%s/auth/error?msg=%s",
		frontendURL,
		url.QueryEscape(errMsg),
	)

	log.Printf("🔗 [GoogleAuth] Redirecting to error page: %s", redirectURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// Helper functions for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func maskToken(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	q := u.Query()
	if token := q.Get("token"); token != "" {
		q.Set("token", truncateString(token, 10)+"...")
	}
	if refreshToken := q.Get("refresh_token"); refreshToken != "" {
		q.Set("refresh_token", truncateString(refreshToken, 10)+"...")
	}

	u.RawQuery = q.Encode()
	return u.String()
}
