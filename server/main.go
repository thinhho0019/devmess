package main

import (
	"fmt"
	"os"

	"project/database"
	"project/handler"
	"project/middleware"
	"project/repository"
	"project/service"

	"project/router"
	"project/websocket"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	// import Gin
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Không tìm thấy file .env hoặc không load được:", err)
	} else {
		fmt.Println("✅ Đã load file .env thành công!")
	}
	database.ConnectDB()
	database.InitRedis()
	var GoogleOAuthConfig *oauth2.Config
	// new repo
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewDeviceRepository()
	tokenRepo := repository.NewTokenRepository()
	redisRepo := repository.NewRedisRepository()
	// new service
	authService := service.NewAuthService(userRepo, deviceRepo, tokenRepo, redisRepo)
	userService := service.NewUserService(userRepo)
	//
	authHandler := handler.NewAuthHandler(userService, authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	imageHandler := handler.NewImageHandler(authService)
	authGoogleHandler := handler.NewAuthGoogleHandler(authService, GoogleOAuthConfig)

	authGoogleHandler.InitGoogleOAuth(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)

	hub := websocket.NewHub()
	go hub.Run()
	// Khởi tạo Gin router
	r := router.SetupRouter(hub, authHandler, authMiddleware, imageHandler)
	r.GET("api/auth/google", authGoogleHandler.GoogleLoginHandler)
	r.GET("api/auth/google/callback", authGoogleHandler.GoogleCallBackHandler)
	// Chạy server trên port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
	defer database.CloseRedis()
}
