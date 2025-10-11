package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.Engine) {
	protected := r.Group("/api", middleware.VerifyAccessToken)
	{
		protected.GET("/images/:filename", handler.ServeImage)
	}
}
