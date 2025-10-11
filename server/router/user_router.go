package router

import (
	"project/handler"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	r.POST("api/check-email", handler.CheckEmailExist)
	r.POST("api/login", handler.LoginPassword)
	r.POST("api/register", handler.Register)
}
