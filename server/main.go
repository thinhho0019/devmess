package main

import (
	"fmt"
	"log"
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
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Không tìm thấy file .env hoặc không load được:", err)
	} else {
		fmt.Println("✅ Đã load file .env thành công!")
	}

	// Initialize database connections
	database.ConnectDB()
	database.InitRedis()

	// Defer cleanup operations
	defer func() {
		database.CloseRedis()
		fmt.Println("✅ Đã đóng kết nối Redis")
	}()

	var GoogleOAuthConfig *oauth2.Config

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewDeviceRepository()
	tokenRepo := repository.NewTokenRepository()
	redisRepo := repository.NewRedisRepository()
	friendRepo := repository.NewFriendshipRepository()
	conversationRepo := repository.NewConversationRepository()
	participantRepo := repository.NewParticipantRepository()
	messageRepo := repository.NewMessageRepository()
	// Initialize services
	googleService := service.NewGoogleService(GoogleOAuthConfig)
	googleService.InitGoogleOAuth(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)
	authService := service.NewAuthService(userRepo, deviceRepo, tokenRepo, redisRepo, googleService.OAuthConfig)
	userService := service.NewUserService(userRepo)
	conversationService := service.NewConversationService(conversationRepo, participantRepo, messageRepo)
	messageService := service.NewMessageService(messageRepo)
	friendService := service.NewInitFriendService(friendRepo, userRepo)
	participantService := service.NewParticipantService(participantRepo, redisRepo)
	// Initialize WebSocket hub
	hub := websocket.NewHub()
	wsHandler := websocket.NewWsHandler(hub, authService)
	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	authHandler := handler.NewAuthHandler(userService, authService, googleService)
	friendHandler := handler.NewFriendHandler(friendService, hub, conversationService)
	conversationHandler := handler.NewConversationHandler(conversationService)
	imageHandler := handler.NewImageHandler(authService)
	authGoogleHandler := handler.NewAuthGoogleHandler(authService, googleService.OAuthConfig)
	messageHandler := handler.NewMessageHandler(*messageService, hub, *participantService)
	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	go hub.Run()

	// Setup router với tất cả handlers
	r := router.SetupRouter(
		hub,
		authHandler,
		userHandler,
		authMiddleware,
		friendHandler,
		wsHandler,

		authService,
		imageHandler,
		conversationHandler,
		messageHandler,
	)

	// Google OAuth routes
	r.GET("/api/auth/google", authGoogleHandler.GoogleLoginHandler)
	r.GET("/api/auth/google/callback", authGoogleHandler.GoogleCallBackHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server đang chạy trên port %s\n", port)
	fmt.Printf("📡 WebSocket hub đã được khởi động\n")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Không thể khởi động server:", err)
	}
}
