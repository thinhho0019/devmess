package router

import (
	"project/handler"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func FriendshipRouter(r *gin.Engine, authMiddleware *middleware.AuthMiddleware,
	friendHandler *handler.FriendHandler) {
	v1 := r.Group("/api/v1")
	{
		friendship := v1.Group("/friendships", authMiddleware.VerifyAccessToken)
		{
			friendship.POST("/send-invite", friendHandler.SendInviteFriend)
			friendship.POST("/cancel-invite", friendHandler.CancelInviteFriend)
			friendship.GET("/list-invite-friends", friendHandler.GetListInviteFriend)
			friendship.POST("/accept-invite", friendHandler.AcceptInviteFriend)
			friendship.POST("/reject-invite", friendHandler.RejectInviteFriend)
			friendship.GET("/list-friends", friendHandler.GetListFriends)
			friendship.DELETE("/remove-friends", friendHandler.RemoveFriend)
		}
	}
}
