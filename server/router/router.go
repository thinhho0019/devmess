package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Gọi các module router
	TestRouter(r)

	return r
}
