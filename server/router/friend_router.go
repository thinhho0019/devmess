package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func FriendshipRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		friendship := v1.Group("/friendships", middleware.VerifyAccessToken)
		{
			friendship.POST("/send-invite", handler.SendInviteFriend)
		}
	}
}
