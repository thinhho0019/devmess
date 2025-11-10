package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func ConversationRouter(r *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	conversationHandler *handler.ConversationHandler) {

	v1 := r.Group("/api/v1")
	{
		conversations := v1.Group("/conversations", authMiddleware.VerifyAccessToken)
		{
			{
				conversations.GET("/", conversationHandler.GetUserConversationsByUserID)
				conversations.GET("/messages/", conversationHandler.GetMessageByConversationID)
				conversations.POST("/find-conversation", conversationHandler.FindConversationByUser)
			}
		}
	}
}
