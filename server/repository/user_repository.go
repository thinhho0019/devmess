package repository

import (
	"errors"
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
	User   models.User
	Status string
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

func (r *userRepo) FindUserWithStatusFriend(email string, user_id uuid.UUID) (*models.User, string, error) {
	var result UserWithStatusFriend
	err := r.db.Table("users").
		Joins("left join friendships on (users.id = friendships.requested_by OR users.id = friendships.friend_id)").
		Where("users.email = ? and users.id = ?", email, user_id).Scan(&result).Error
	if err != nil {
		return nil, "", err
	}

	if result.User.ID != uuid.Nil {
		return &result.User, result.Status, nil
	}

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
