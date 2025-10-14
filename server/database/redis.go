package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	// 🔍 Kiểm tra kết nối
	if err := RDB.Ping(Ctx).Err(); err != nil {
		fmt.Println("❌ Lỗi kết nối Redis:", err)
	} else {
		fmt.Println("✅ Database 'redis' sẵn sàng!")
	}
}

func CloseRedis() {
	if RDB != nil {
		RDB.Close()
		fmt.Println("✅ Database 'redis' đã đóng!")
	}
}
