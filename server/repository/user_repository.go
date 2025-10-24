package repository

import (
	"errors"
	"log"
	"project/database"
	"project/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	LoginPassword(email string, password string, provider string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetAllUsers() ([]models.User, error)
	FindUserWithStatusFriend(email string, user_id uuid.UUID) (*models.User, string, error)
	GetUserByAccesToken(accessToken string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}
type UserWithStatusFriend struct {
	ID           uuid.UUID `gorm:"column:id"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	Avatar       string    `gorm:"column:avatar"`
	StatusFriend string    `gorm:"column:status_friend"`
}

func NewUserRepository() UserRepository {
	return &userRepo{
		db: database.DB, // database.DB là connection GORM chính
	}
}

// 🧱 Create user
func (r *userRepo) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) LoginPassword(email string, password string, provider string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Provider != provider {
		return nil, errors.New("please login with " + user.Provider)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// 🔍 Find user by ID
func (r *userRepo) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 🔍 Find user by email
func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 📝 Update user info
func (r *userRepo) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// ❌ Delete user
func (r *userRepo) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// 📃 Get all users
func (r *userRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) FindUserWithStatusFriend(email string, userID uuid.UUID) (*models.User, string, error) {
	log.Printf("[FindUserWithStatusFriend][Repo] Input email: %s, userID: %s", email, userID)

	var result UserWithStatusFriend
	err := r.db.Debug().Table("users").
		Select(`users.*,friendships.status as status_friend`).
		Joins("LEFT JOIN friendships ON (users.id = friendships.requested_by OR users.id = friendships.friend_id)").
		Where("users.email = ?", email).
		Scan(&result).Error

	if err != nil {
		log.Printf("[FindUserWithStatusFriend][Repo] SQL error: %v", err)
		return nil, "", err
	}

	log.Printf("[FindUserWithStatusFriend][Repo] Scan result: %+v", result)

	if result.ID != uuid.Nil {
		log.Printf("[FindUserWithStatusFriend][Repo] User found: ID=%s, Status=%s", result.ID, result.StatusFriend)
		user := models.User{
			ID:     result.ID,
			Name:   result.Name,
			Email:  result.Email,
			Avatar: result.Avatar,
		}
		return &user, result.StatusFriend, nil
	}

	log.Printf("[FindUserWithStatusFriend][Repo] User not found for email: %s", email)
	return nil, "", nil
}

func (r *userRepo) GetUserByAccesToken(accessToken string) (*models.User, error) {
	var result models.User
	err := r.db.Table("users").
		Select("users.*").
		Joins("JOIN devices ON devices.user_id = users.id").
		Joins("JOIN tokens ON tokens.device_id = devices.id").
		Where("tokens.access_token = ?", accessToken).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil

}
