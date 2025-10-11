package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"project/models"

	"strings"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 1️⃣ Tạo DB nếu chưa có
	CreateDBIfNotExists()

	// 2️⃣ Sau đó mới mở GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Không thể kết nối database: %v", err)
	}

	// 3️⃣ Auto migrate
	err = db.AutoMigrate(&models.User{}, &models.Token{}, &models.Device{})
	if err != nil {
		log.Fatalf("❌ Lỗi auto migrate: %v", err)
	}

	log.Println("✅ Auto migrate thành công!")
	log.Println("✅ Kết nối PostgreSQL thành công!")

	DB = db
}

func CreateDBIfNotExists() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Kết nối vào DB mặc định (postgres)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Lỗi kết nối postgres: %v", err)
	}
	defer db.Close()

	// Tạo database nếu chưa có
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil && !isDuplicateDatabaseError(err) {
		log.Fatalf("❌ Lỗi tạo database: %v", err)
	} else {
		log.Printf("✅ Database '%s' sẵn sàng!", dbname)
	}
}

func isDuplicateDatabaseError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
 