package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func TestRouter(r *gin.Engine) {
	protected := r.Group("/api", middleware.VerifyAccessToken)
	{
		protected.GET("/", handler.TestHandle)
	}

}
