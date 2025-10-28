package main

import (
	"log"
	"net/http"
	"os"

	"project/database"
	"project/handler"
	"project/middleware"
	"project/repository"
	"project/router"
	"project/service"
	"project/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	} else {
		log.Println("✅ Đã load file .env thành công!")
	}

	// Validate required environment variables
	validateEnv()

	// Log environment configuration
	logEnvironment()

	// Initialize database connections
	database.ConnectDB()
	database.InitRedis()

	defer func() {
		database.CloseRedis()
		log.Println("✅ Đã đóng kết nối Redis")
	}()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewDeviceRepository()
	tokenRepo := repository.NewTokenRepository()
	redisRepo := repository.NewRedisRepository()
	friendRepo := repository.NewFriendshipRepository()
	conversationRepo := repository.NewConversationRepository()
	participantRepo := repository.NewParticipantRepository()
	messageRepo := repository.NewMessageRepository()

	// Initialize Google OAuth Config
	googleOAuthConfig := initGoogleOAuth()

	// Initialize services
	authService := service.NewAuthService(userRepo, deviceRepo, tokenRepo, redisRepo, googleOAuthConfig)
	userService := service.NewUserService(userRepo)
	conversationService := service.NewConversationService(conversationRepo, participantRepo, messageRepo)
	messageService := service.NewMessageService(messageRepo, conversationRepo)
	friendService := service.NewInitFriendService(friendRepo, userRepo)
	participantService := service.NewParticipantService(participantRepo, redisRepo)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	wsHandler := websocket.NewWsHandler(hub, authService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService, authService, nil)
	friendHandler := handler.NewFriendHandler(friendService, hub, conversationService)
	conversationHandler := handler.NewConversationHandler(conversationService)
	imageHandler := handler.NewImageHandler(authService)
	authGoogleHandler := handler.NewAuthGoogleHandler(authService, googleOAuthConfig)
	messageHandler := handler.NewMessageHandler(*messageService, hub, *participantService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Start WebSocket hub
	go hub.Run()
	log.Println("📡 WebSocket hub đã được khởi động")

	// Setup router
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
	r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
	// Google OAuth routes
	r.GET("/api/v1/auth/google", authGoogleHandler.GoogleLoginHandler)
	r.GET("/api/auth/google/callback", authGoogleHandler.GoogleCallBackHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server đang chạy trên port %s", port)
	log.Printf("🌍 Environment: %s", getEnvironment())

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Không thể khởi động server:", err)
	}
}

// initGoogleOAuth initializes Google OAuth configuration from environment
func initGoogleOAuth() *oauth2.Config {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	// Auto-construct redirect URL if not set
	if redirectURL == "" {
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = os.Getenv("DEFAULT_URL_SERVER")
		}
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		redirectURL = baseURL + "/api/auth/google/callback"
		log.Printf("⚠️ [GoogleOAuth] GOOGLE_REDIRECT_URL not set, using: %s", redirectURL)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	log.Println("✅ [GoogleOAuth] Configuration initialized")
	log.Printf("   Client ID: %s", maskClientID(clientID))
	log.Printf("   Redirect URL: %s", redirectURL)

	return config
}

// validateEnv checks for required environment variables
func validateEnv() {
	required := []string{
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"JWT_SECRET_KEY",
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}

	missing := []string{}
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Printf("⚠️ Missing required environment variables: %v", missing)
	}
}

// logEnvironment logs current environment configuration
func logEnvironment() {
	log.Println("📋 Environment Configuration:")
	log.Printf("   ENV: %s", getEnvironment())
	log.Printf("   BASE_URL: %s", os.Getenv("BASE_URL"))
	log.Printf("   DEFAULT_URL_SERVER: %s", os.Getenv("DEFAULT_URL_SERVER"))
	log.Printf("   FRONTEND_URL: %s", getFrontendURL())
	log.Printf("   DB_HOST: %s", os.Getenv("DB_HOST"))
	log.Printf("   REDIS_HOST: %s", os.Getenv("REDIS_HOST"))
}

// getEnvironment returns current environment (production/development)
func getEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	return env
}

// getFrontendURL returns frontend URL with fallback
func getFrontendURL() string {
	url := os.Getenv("FRONTEND_URL")
	if url == "" {
		url = os.Getenv("DEFAULT_URL")
	}
	if url == "" {
		url = "http://localhost:3000"
	}
	return url
}

// maskClientID masks sensitive client ID for logging
func maskClientID(id string) string {
	if len(id) < 20 {
		return "***"
	}
	return id[:10] + "..." + id[len(id)-10:]
}
