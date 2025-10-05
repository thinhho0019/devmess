package main

import (
	"os"
	"project/router"
	"project/websocket"
	// import Gin
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()
	// Khởi tạo Gin router
	r := router.SetupRouter(hub)

	// Chạy server trên port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
