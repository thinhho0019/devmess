package handler

import (
	"net/http"
	"project/models"
	"project/repository"
	"project/service"

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

func AuthHandle(c *gin.Context) {

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

func Register(c *gin.Context) {
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

	// --- Tạo người dùng mới ---
	// Giả sử bạn có một phương thức `RegisterUser` trong `userService`
	// để xử lý việc băm mật khẩu và tạo người dùng.
	user, err := userService.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	// --- Tự động đăng nhập và tạo session sau khi đăng ký ---
	deviceRepo := repository.NewDeviceRepository()
	tokenRepo := repository.NewTokenRepository()
	redisRepo := repository.NewRedisRepository()
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	token, _, err := service.CreateSession(userRepo, deviceRepo, tokenRepo, redisRepo, user, ip, userAgent, "", "", 0, "local")
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

func LoginPassword(c *gin.Context) {
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

	// --- Tạo Session (Device & Token) ---
	deviceRepo := repository.NewDeviceRepository()
	tokenRepo := repository.NewTokenRepository()
	redisRepo := repository.NewRedisRepository()

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Gọi service để tạo session.
	// Truyền chuỗi rỗng và 0 vì đây là đăng nhập bằng mật khẩu, không có token/expiresIn từ bên ngoài.
	token, _, err := service.CreateSession(userRepo, deviceRepo, tokenRepo, redisRepo, user, ip, userAgent, "", "", 0, "local")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session: " + err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"access_token": token.AccessToken})
}

func CheckEmailExist(c *gin.Context) {
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
