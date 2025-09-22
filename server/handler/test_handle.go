package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestHandle(c *gin.Context) {
	var req struct {
		Email    string `json: "username`
		Password string `json: "password`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": req.Email + " có tuổi là " + fmt.Sprint(req.Password)})
}
