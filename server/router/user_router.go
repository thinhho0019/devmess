package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/check-email", handler.CheckEmailExist)
			auth.POST("/login", handler.LoginPassword)
			auth.POST("/register", handler.Register)
			auth.POST("/forgot-password", handler.ForgotPassword)
			auth.POST("/reset-password", handler.ResetPassword)
		}

		// for user
		users := v1.Group("/users", middleware.VerifyAccessToken)
		{
			users.GET("/search", handler.FindUserByEmail)
		}
	}
}
