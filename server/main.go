package main

import (
	"fmt"
	"os"
	"project/database"
	"project/handler"
	"project/router"
	"project/websocket"

	"github.com/joho/godotenv"
	// import Gin
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Không tìm thấy file .env hoặc không load được:", err)
	} else {
		fmt.Println("✅ Đã load file .env thành công!")
	}

	handler.InitGoogleOAuth(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)
	database.ConnectDB()
	hub := websocket.NewHub()
	go hub.Run()
	// Khởi tạo Gin router
	r := router.SetupRouter(hub)
	r.GET("api/auth/google", handler.GoogleLoginHandler)
	r.GET("api/auth/google/callback", handler.GoogleCallBackHandler)
	// Chạy server trên port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
