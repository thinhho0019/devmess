package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine, authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	userHanlder *handler.UserHandler) {

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/check-email", authHandler.CheckEmailExist)
			auth.POST("/login", authHandler.LoginPassword)
			auth.POST("/register", authHandler.Register)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.GET("/refresh-token", authHandler.AuthRefreshToken)
		}

		// for user
		users := v1.Group("/users", authMiddleware.VerifyAccessToken)
		{
			users.GET("/search", userHanlder.FindUserWithStatusFriends)
		}
		
	}
}
