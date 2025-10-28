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

	DB = db

	// 3️⃣ Reset database nếu có flag
	if os.Getenv("RESET_DB") == "true" {
		log.Println("🔄 Resetting database...")
		if err := ResetDatabase(db); err != nil {
			log.Fatalf("❌ Lỗi reset database: %v", err)
		}
		log.Println("✅ Database đã được reset thành công!")
		return
	}

	// 4️⃣ Auto migrate theo thứ tự dependencies
	if err := AutoMigrateInOrder(db); err != nil {
		log.Fatalf("❌ Lỗi auto migrate: %v", err)
	}

	log.Println("✅ Kết nối PostgreSQL thành công!")
}

// AutoMigrateInOrder migrates tables in correct dependency order
func AutoMigrateInOrder(db *gorm.DB) error {
	// 1️⃣ Independent tables first
	log.Println("🔄 Migrating users table...")
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate users: %w", err)
	}
	log.Println("✅ Users table migrated")

	log.Println("🔄 Migrating tokens table...")
	if err := db.AutoMigrate(&models.Token{}); err != nil {
		return fmt.Errorf("failed to migrate tokens: %w", err)
	}
	log.Println("✅ Tokens table migrated")

	log.Println("🔄 Migrating devices table...")
	if err := db.AutoMigrate(&models.Device{}); err != nil {
		return fmt.Errorf("failed to migrate devices: %w", err)
	}
	log.Println("✅ Devices table migrated")

	// 2️⃣ Conversation (depends on User but no circular dependency)
	log.Println("🔄 Migrating conversations table...")
	if err := db.AutoMigrate(&models.Conversation{}); err != nil {
		return fmt.Errorf("failed to migrate conversations: %w", err)
	}
	log.Println("✅ Conversations table migrated")

	// 3️⃣ Messages (depends on Conversation and User)
	log.Println("🔄 Migrating messages table...")
	if err := db.AutoMigrate(&models.Message{}); err != nil {
		return fmt.Errorf("failed to migrate messages: %w", err)
	}
	log.Println("✅ Messages table migrated")

	// 4️⃣ Participants (depends on Conversation and User)
	log.Println("🔄 Migrating participants table...")
	if err := db.AutoMigrate(&models.Participant{}); err != nil {
		return fmt.Errorf("failed to migrate participants: %w", err)
	}
	log.Println("✅ Participants table migrated")

	// 5️⃣ Friendships (depends on User)
	log.Println("🔄 Migrating friendships table...")
	if err := db.AutoMigrate(&models.Friendship{}); err != nil {
		return fmt.Errorf("failed to migrate friendships: %w", err)
	}
	log.Println("✅ Friendships table migrated")

	log.Println("✅ All tables migrated successfully!")
	return nil
}

func ResetDatabase(db *gorm.DB) error {
	log.Println("🗑️ Dropping existing schema...")
	if err := db.Exec("DROP SCHEMA public CASCADE;").Error; err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

	log.Println("📦 Creating new schema...")
	if err := db.Exec("CREATE SCHEMA public;").Error; err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("🔄 Running migrations...")
	return AutoMigrateInOrder(db)
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
	}

	log.Printf("✅ Database '%s' sẵn sàng!", dbname)
}

func isDuplicateDatabaseError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
