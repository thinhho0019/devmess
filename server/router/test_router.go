package router

import (
	"project/handler"

	"github.com/gin-gonic/gin"
)

func TestRouter(r *gin.Engine) {
	test := r.Group("/api")
	{
		test.GET("/", handler.TestHandle)
	}
}
