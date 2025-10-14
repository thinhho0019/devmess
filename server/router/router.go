package router

import (
	"project/handler"
	"project/middleware"
	"project/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(hub *websocket.Hub,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	imageHandler *handler.ImageHandler) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			// "http://localhost",
			// "http://localhost:80",
			// "http://localhost:5173",
			// "http://127.0.0.1",
			// "http://127.0.0.1:80",
			"https://devmess.cloud",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RateLimitMiddleware())
	// Thêm route cho WebSocket
	r.GET("/ws", func(c *gin.Context) {
		handler.ServeWs(hub, c.Writer, c.Request)
	})

	// Gọi các module router
	AuthRouter(r, authHandler, authMiddleware)
	UserRouter(r, authHandler, authMiddleware)
	ImageRouter(r, authMiddleware, imageHandler)
	FriendshipRouter(r, authMiddleware)
	return r
}
