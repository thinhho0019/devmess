package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine) {
	protected := r.Group("/api/v1", middleware.VerifyAccessToken)
	{
		protected.GET("/auth-me", handler.AuthHandle)
	}

}
