package repository

import (
	"fmt"
	"project/database"
	"project/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository interface {
	CreateToken(token *models.Token) error
	UpdateToken(token *models.Token) error
	GetTokenByRefresh(refresh string) (*models.Token, error)
	GetTokenByAccess(access string) (*models.Token, error)
	GetTokensByUserID(device_id string) (*models.Token, error)
	GetUserForAccessToken(accessToken string) (*models.User, string, error)
	DeleteTokenByID(ID uuid.UUID) error
	DeleteTokensByUserID(userID uint) error
	GetTokensByDeviceID(userID uuid.UUID, deviceID string) (*models.Token, error)
}

type tokenRepo struct {
	db *gorm.DB
}

var NewTokenRepository = func() TokenRepository {
	return &tokenRepo{
		db: database.DB,
	}
}

func (r *tokenRepo) UpdateToken(token *models.Token) error {
	token.UpdatedAt = time.Now()
	return r.db.Save(token).Error
}

// ✅ Tạo token mới (thêm bản ghi)
func (r *tokenRepo) CreateToken(token *models.Token) error {
	return r.db.Create(token).Error
}

// ✅ Lấy token theo refresh token
func (r *tokenRepo) GetTokenByRefresh(refresh string) (*models.Token, error) {
	var token models.Token
	if err := r.db.Where("refresh_token = ?", refresh).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
func (r *tokenRepo) GetTokenByAccess(access string) (*models.Token, error) {
	var token models.Token
	if err := r.db.Where("access_token = ?", access).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepo) GetTokensByDeviceID(userID uuid.UUID, deviceID string) (*models.Token, error) {
	var token models.Token
	fmt.Printf("user_id is :%v", userID)
	err := r.db.
		Table("tokens").Where("tokens.device_id = ? AND devices.user_id = ?", deviceID, userID).
		Joins("JOIN devices ON tokens.device_id = devices.id").
		First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // không tìm thấy → trả về nil
		}
		return nil, err // lỗi khác thì trả về
	}
	return &token, nil
}

// ✅ Lấy toàn bộ token theo user (nếu user có nhiều device)
func (r *tokenRepo) GetTokensByUserID(device_id string) (*models.Token, error) {
	var tokens models.Token
	if err := r.db.Table("tokens").Where("tokens.device_id = ?", device_id).First(&tokens).Error; err != nil {
		return nil, err
	}
	return &tokens, nil
}

// ✅ Xóa token cụ thể (logout 1 thiết bị)
func (r *tokenRepo) DeleteTokenByID(ID uuid.UUID) error {
	return r.db.Delete(&models.Token{}, ID).Error
}
func (r *tokenRepo) DeleteTokenByDeviceID(deviceID uuid.UUID) error {
	return r.db.Unscoped().Where("device_id = ?", deviceID).Delete(&models.Token{}).Error
}

// ✅ Xóa tất cả token của user (logout all)
func (r *tokenRepo) DeleteTokensByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Token{}).Error
}

func (r *tokenRepo) GetUserForAccessToken(accessToken string) (*models.User, string, error) {
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
	err := r.db.
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
	fmt.Println(time.Now())
	fmt.Println(time.Unix(int64(result.ExpiresAt), 0))
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
