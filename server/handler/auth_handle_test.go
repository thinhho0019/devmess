package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"project/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Test AuthHandle when user is present in context
func TestAuthHandle_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	now := time.Now()
	u := &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Avatar:    "avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}
	c.Set("user", u)

	AuthHandle(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, body["email"])
	assert.Equal(t, u.Name, body["name"])
}

// Test LoginPassword returns 400 for invalid JSON and for empty password
func TestLoginPassword_BadRequests(t *testing.T) {
	// invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("not-a-json"))
	c.Request = req

	LoginPassword(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// empty password
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	body := map[string]string{"email": "a@b.com", "password": ""}
	bs, _ := json.Marshal(body)
	req2 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bs))
	req2.Header.Set("Content-Type", "application/json")
	c2.Request = req2

	LoginPassword(c2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// Test ForgotPassword returns 400 for invalid body and 200 for well-formed (non-disclosure)
func TestForgotPassword_BadRequestAndOk(t *testing.T) {
	// invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/auth/forgot", bytes.NewBufferString("bad"))
	c.Request = req

	ForgotPassword(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// well-formed but email may not exist -> handler always returns 200 (non-disclosure)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	body := map[string]string{"email": "notfound@example.com"}
	bs, _ := json.Marshal(body)
	req2 := httptest.NewRequest(http.MethodPost, "/auth/forgot", bytes.NewBuffer(bs))
	req2.Header.Set("Content-Type", "application/json")
	c2.Request = req2

	ForgotPassword(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

// Test ResetPassword invalid token returns 400
func TestResetPassword_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{"token": "invalid-or-bad-token", "password": "newpass123"}
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/reset", bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ResetPassword(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
