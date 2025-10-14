package handler

import (
	"fmt"
	"net/http"
	"os"
	"project/models"
	"project/repository"
	"project/service"
	"project/utils"

	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	Email string `json:"email"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
type AuthHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewAuthHandler(userService *service.UserService, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{userService: userService, authService: authService}
}

func (a *AuthHandler) AuthHandle(c *gin.Context) {

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
	// ✅ Trả JSON gọn gàng
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"name":      user.Name,
		"avatar":    user.Avatar,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

func (a *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- Khởi tạo repository và service ---
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)

	// --- Kiểm tra email đã tồn tại chưa ---
	exists, err := userService.CheckEmail(req.Email)
	println("ERROR:", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check email existence"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	user, err := userService.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	// --- Tự động đăng nhập và tạo session sau khi đăng ký ---

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	token, _, err := a.authService.CreateSession(user, ip, userAgent, "", "", 0, "local")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session after registration: " + err.Error()})
		return
	}

	// --- Trả về token ---
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

func (a *AuthHandler) LoginPassword(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password Empty"})
		return
	}

	// --- Xác thực người dùng ---
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	user, err := userService.LoginPassword(req.Email, req.Password)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Gọi service để tạo session.
	// Truyền chuỗi rỗng và 0 vì đây là đăng nhập bằng mật khẩu, không có token/expiresIn từ bên ngoài.
	token, _, err := a.authService.CreateSession(user, ip, userAgent, "", "", 0, "local")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session: " + err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"access_token": token.AccessToken})
}

func (a *AuthHandler) CheckEmailExist(c *gin.Context) {
	var req EmailRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	repo := repository.NewUserRepository()
	userService := service.NewUserService(repo)

	exists, err := userService.CheckEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusOK, gin.H{"message": "user exists"})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "email not exist"})
		return
	}
}

// ForgotPassword: tạo reset token và gửi link (không tiết lộ email tồn tại)
func (a *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByEmail(req.Email)
	if user != nil && user.Provider != "local" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This email is registered via " + user.Provider})
		return
	}
	if err != nil || user == nil {
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link was sent"})
		return
	}

	// Tạo reset token
	token, err := utils.GenerateResetToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	// Build reset link -> frontend sẽ có route nhận token
	frontendURL := os.Getenv("DEFAULT_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	// Gửi email (placeholder)
	_ = utils.SendResetEmail(user.Email, resetLink)

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link was sent"})
}

// ResetPassword: verify token và cập nhật mật khẩu
func (a *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.ValidateResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByID(claims.UserID)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Hash password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Cập nhật mật khẩu. Giả sử repo cung cấp UpdateUser hoặc UpdatePassword
	user.Password = hashed
	if err := userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset"})
}
