package repository

import (
	"fmt"
	"project/database"
	"project/models"
	"time"

	"github.com/google/uuid"
)

// ✅ Tạo token mới (thêm bản ghi)
func CreateToken(token *models.Token) error {
	return database.DB.Create(token).Error
}

// ✅ Lấy token theo refresh token
func GetTokenByRefresh(refresh string) (*models.Token, error) {
	var token models.Token
	if err := database.DB.Where("refresh_token = ?", refresh).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
func GetTokenByAccess(access string) (*models.Token, error) {
	var token models.Token
	if err := database.DB.Where("access_token = ?", access).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// ✅ Lấy toàn bộ token theo user (nếu user có nhiều device)
func GetTokensByUserID(userID uint) ([]models.Token, error) {
	var tokens []models.Token
	if err := database.DB.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// ✅ Xóa token cụ thể (logout 1 thiết bị)
func DeleteTokenByID(ID uint) error {
	return database.DB.Delete(&models.Token{}, ID).Error
}

// ✅ Xóa tất cả token của user (logout all)
func DeleteTokensByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.Token{}).Error
}

func GetUserForAccessToken(accessToken string) (*models.User, string, error) {
	print("at", accessToken)
	var result struct {
		ID           uuid.UUID
		Email        string
		Name         string
		Avatar       string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		ExpiresAt    int64
		RefreshToken string
	}
	err := database.DB.
		Table("users").
		Select("users.*, tokens.expires_at, tokens.refresh_token").
		Joins("JOIN devices ON devices.user_id = users.id").
		Joins("JOIN tokens ON tokens.device_id = devices.id").
		Where("tokens.access_token = ?", accessToken).
		First(&result).Error
	fmt.Println(result)
	if err != nil {
		return nil, "", err
	}

	if time.Now().After(time.Unix(int64(result.ExpiresAt), 0)) {

		return nil, result.RefreshToken, nil
	}
	user := &models.User{
		ID:        result.ID,
		Email:     result.Email,
		Name:      result.Name,
		Avatar:    result.Avatar,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	return user, "", nil
}
