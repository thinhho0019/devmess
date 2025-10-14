package handler

import (
	// "errors"
	// "net/http"
	// "net/http/httptest"

	// "project/models"
	// // "project/service"
	// "strings"
	// "testing"

	// "github.com/gin-gonic/gin"
	// "github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var GoogleOAuthConfigTest = &oauth2.Config{
	ClientID:     "fake-client-id",
	ClientSecret: "fake-secret",
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{"email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
}

// func mockHandleGoogleCallbackInvalidCode(
// 	code, ip, ua string,
// 	cfg *oauth2.Config,
// ) (*models.GoogleUserInfo, *models.Token, *models.Device, error) {
// 	return nil, nil, nil, errors.New("invalid code")
// }

// func TestGoogleLoginHandler(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()

// 	// Gán config thật vào biến toàn cục
// 	GoogleOAuthConfig = &oauth2.Config{
// 		ClientID:     "fake-client-id",
// 		ClientSecret: "fake-secret",
// 		RedirectURL:  "http://localhost:8080/callback",
// 		Scopes:       []string{"email"},
// 		Endpoint: oauth2.Endpoint{
// 			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
// 			TokenURL: "https://oauth2.googleapis.com/token",
// 		},
// 	}

// 	r.GET("/login/google", )

// 	req, _ := http.NewRequest(http.MethodGet, "/login/google", nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
// 	location := w.Header().Get("Location")
// 	assert.NotEmpty(t, location)
// 	assert.True(t, strings.HasPrefix(location, "https://accounts.google.com/o/oauth2/auth"))
// }

// func  (a *AuthHandler) TestGoogleCallBackHandler_MissingCode(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	req, _ := http.NewRequest(http.MethodGet, "/callback", nil)
// 	c.Request = req

// 	GoogleCallBackHandler(c)
// 	assert.Equal(t, http.StatusBadRequest, w.Code)

// }
// func (a *AuthHandler) TestGoogleCallBackHandler_InvalidCode(t *testing.T) {
// 	// ⚙️ Setup gin test context
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	// ⚙️ Add query param "code" (giả lập Google redirect về)
// 	req, _ := http.NewRequest(http.MethodGet, "/callback?code=bad_code", nil)
// 	c.Request = req

// 	// ⚙️ Thay tạm hàm service.HandleGoogleCallback bằng mock
// 	original := a.authService.HandleGoogleCallback
// 	a.authService.HandleGoogleCallback = mockHandleGoogleCallbackInvalidCode
// 	defer func() { a.authService.HandleGoogleCallback = original }()

// 	// 🧪 Gọi handler
// 	a.GoogleCallBackHandler(c)

// 	// ✅ Kiểm tra phản hồi
// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "invalid code")
// }
 
// func TestGoogleCallBackHandler_Error(t *testing.T) {
// 	// ⚙️ Setup gin test context
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	// ⚙️ Add query param "code" (giả lập Google redirect về)
// 	req, _ := http.NewRequest(http.MethodGet, "/callback?code=bad_code", nil)
// 	c.Request = req

// 	// ⚙️ Thay tạm hàm service.HandleGoogleCallback bằng mock
// 	original := service.HandleGoogleCallback
// 	service.HandleGoogleCallback = mockHandleGoogleCallbackInvalidCode
// 	defer func() { service.HandleGoogleCallback = original }()

// 	// 🧪 Gọi handler
// 	GoogleCallBackHandler(c)

// 	// ✅ Kiểm tra phản hồi
// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "invalid code")
// }
