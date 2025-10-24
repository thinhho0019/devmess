package router

import (
	"project/handler"
	"project/middleware"
	"project/service"
	"project/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(hub *websocket.Hub,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	friendHandler *handler.FriendHandler,
	wsHandler *websocket.WsHandler,
	authService *service.AuthService,
	imageHandler *handler.ImageHandler,
	conversationHandler *handler.ConversationHandler,
	messageHandler *handler.MessageHanlder,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost",
			"http://localhost:80",
			"http://localhost:5173",
			"http://localhost:3000",
			"http://127.0.0.1",
			"http://127.0.0.1:80",
			"https://devmess.cloud",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RateLimitMiddleware())

	// Thêm route cho WebSocket
	r.GET("/ws", wsHandler.ServeWs())

	// Gọi các module router
	AuthRouter(r, authHandler, authMiddleware)
	UserRouter(r, authHandler, authMiddleware, userHandler)
	ImageRouter(r, authMiddleware, imageHandler)
	FriendshipRouter(r, authMiddleware, friendHandler)
	ConversationRouter(r, authMiddleware, conversationHandler)
	MessageRouter(r, authMiddleware, messageHandler)
	return r
}
