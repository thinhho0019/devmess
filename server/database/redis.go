package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println("✅ Database 'redis' sẵn sàng!")
}
func CloseRedis() {
	RDB.Close()
	fmt.Println("✅ Database 'redis' đã đóng!")
}
