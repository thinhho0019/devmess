package repository

import (
	"project/database"
	"project/models"

	"gorm.io/gorm"
)

// 🧱 Create user
func CreateUser(user *models.User) (*models.User, error) {
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 🔍 Find user by ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 🔍 Find user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 📝 Update user info
func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

// ❌ Delete user
func DeleteUser(id uint) error {
	return database.DB.Delete(&models.User{}, id).Error
}

// 📃 Get all users
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
