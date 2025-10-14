package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	protected := r.Group("/api/v1", authMiddleware.VerifyAccessToken)
	{
		protected.GET("/auth-me", authHandler.AuthHandle)
	}

}
