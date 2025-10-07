package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestHandle(c *gin.Context) {
	email := c.Query("username")
	pass := c.Query("password")
	c.JSON(http.StatusOK, gin.H{"message": email + " có tuổi là " + pass})
}
