package main

import (
	"os"
	"project/router"
	// import Gin
)

func main() {

	// Khởi tạo Gin router
	r := router.SetupRouter()

	// Chạy server trên port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
