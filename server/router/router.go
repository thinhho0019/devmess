package router

import (
	"project/handler"
	"project/middleware"
	"project/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(hub *websocket.Hub) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
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
	AuthRouter(r)
	UserRouter(r)
	ImageRouter(r)
	FriendshipRouter(r)
	return r
}
