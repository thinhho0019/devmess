package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func MessageRouter(r *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	messageHandler *handler.MessageHanlder) {

	v1 := r.Group("/api/v1")
	{
		messages := v1.Group("/messages", authMiddleware.VerifyAccessToken)
		{
			{
				messages.GET("/", messageHandler.GetAllMessageToConversation)
				messages.POST("/send", messageHandler.SendMessageToConversation)
			}
		}
	}
}
