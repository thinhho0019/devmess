package repository

import (
	"project/models"

	"github.com/google/uuid"
)

// ✅ MockTokenRepository — mô phỏng TokenRepository cho unit test
type MockTokenRepository struct {
	MockCreateToken           func(token *models.Token) error
	MockGetTokenByRefresh     func(refresh string) (*models.Token, error)
	MockGetTokenByAccess      func(access string) (*models.Token, error)
	MockGetTokensByUserID     func(deviceID string) (*models.Token, error)
	MockGetUserForAccessToken func(accessToken string) (*models.User, string, error)
	MockDeleteTokenByID       func(ID uuid.UUID) error
	MockDeleteTokensByUserID  func(userID uint) error
	MockGetTokensByDeviceID   func(userId uuid.UUID, deviceID string) (*models.Token, error)
}

// Implement interface TokenRepository ↓↓↓

func (m *MockTokenRepository) CreateToken(token *models.Token) error {
	if m.MockCreateToken != nil {
		return m.MockCreateToken(token)
	}
	return nil
}

func (m *MockTokenRepository) GetTokenByRefresh(refresh string) (*models.Token, error) {
	if m.MockGetTokenByRefresh != nil {
		return m.MockGetTokenByRefresh(refresh)
	}
	return nil, nil
}

func (m *MockTokenRepository) GetTokenByAccess(access string) (*models.Token, error) {
	if m.MockGetTokenByAccess != nil {
		return m.MockGetTokenByAccess(access)
	}
	return nil, nil
}

func (m *MockTokenRepository) GetTokensByUserID(deviceID string) (*models.Token, error) {
	if m.MockGetTokensByUserID != nil {
		return m.MockGetTokensByUserID(deviceID)
	}
	return nil, nil
}

func (m *MockTokenRepository) GetUserForAccessToken(accessToken string) (*models.User, string, error) {
	if m.MockGetUserForAccessToken != nil {
		return m.MockGetUserForAccessToken(accessToken)
	}
	return nil, "", nil
}

func (m *MockTokenRepository) DeleteTokenByID(ID uuid.UUID) error {
	if m.MockDeleteTokenByID != nil {
		return m.MockDeleteTokenByID(ID)
	}
	return nil
}

func (m *MockTokenRepository) DeleteTokensByUserID(userID uint) error {
	if m.MockDeleteTokensByUserID != nil {
		return m.MockDeleteTokensByUserID(userID)
	}
	return nil
}

func (m *MockTokenRepository) GetTokensByDeviceID(userId uuid.UUID, deviceID string) (*models.Token, error) {
	if m.MockGetTokensByDeviceID != nil {
		return m.MockGetTokensByDeviceID(userId, deviceID)
	}
	return nil, nil
}
