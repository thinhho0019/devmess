package service

import (
	"context"
	"errors"
	// "project/models"
	// "project/repository"

	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestDetectDevice(t *testing.T) {
	tests := []struct {
		name         string
		userAgent    string
		expectedType string
		expectedName string
	}{
		{
			name:         "Android Chrome",
			userAgent:    "Mozilla/5.0 (Linux; Android 12; Pixel 6) Chrome/110.0.5481.65 Mobile Safari/537.36",
			expectedType: "Android",
			expectedName: "Chrome",
		},
		{
			name:         "iPhone Safari",
			userAgent:    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile Safari/604.1",
			expectedType: "iOS",
			expectedName: "Safari",
		},
		{
			name:         "Windows Firefox",
			userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
			expectedType: "Web",
			expectedName: "Firefox",
		},
		{
			name:         "Mac Safari",
			userAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15",
			expectedType: "Web",
			expectedName: "Safari",
		},
		{
			name:         "Unknown Browser",
			userAgent:    "CustomDevice/1.0",
			expectedType: "Web",
			expectedName: "Browser",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deviceType, deviceName := detectDevice(tt.userAgent)
			assert.Equal(t, tt.expectedType, deviceType)
			assert.Equal(t, tt.expectedName, deviceName)
		})
	}
}

// func TestVerifyAccessToken(t *testing.T) {
// 	// Giữ lại repo gốc và phục hồi sau khi test xong
// 	original := repository.NewTokenRepository
// 	defer func() { repository.NewTokenRepository = original }()

// 	tests := []struct {
// 		name                 string
// 		mockRepo             repository.TokenRepository
// 		expectedUser         *models.User
// 		expectedRefreshToken string
// 		expectedErr          error
// 	}{
// 		{
// 			name: "✅ User hợp lệ, không có refresh token",
// 			mockRepo: &repository.MockTokenRepository{
// 				MockGetUserForAccessToken: func(accessToken string) (*models.User, string, error) {
// 					return &models.User{Name: "thinhho"}, "", nil
// 				},
// 			},
// 			expectedUser:         &models.User{Name: "thinhho"},
// 			expectedRefreshToken: "",
// 			expectedErr:          nil,
// 		},
// 		{
// 			name: "✅ User hợp lệ, có refresh token mới",
// 			mockRepo: &repository.MockTokenRepository{
// 				MockGetUserForAccessToken: func(accessToken string) (*models.User, string, error) {
// 					return nil, "new-refresh-token", errors.New("invalid or expired access token")
// 				},
// 			},
// 			expectedUser:         nil,
// 			expectedRefreshToken: "new-refresh-token",
// 			expectedErr:          errors.New("invalid or expired access token"),
// 		},
// 		{
// 			name: "❌ Token không tồn tại",
// 			mockRepo: &repository.MockTokenRepository{
// 				MockGetUserForAccessToken: func(accessToken string) (*models.User, string, error) {
// 					return nil, "", errors.New("token not found")
// 				},
// 			},
// 			expectedUser:         nil,
// 			expectedRefreshToken: "",
// 			expectedErr:          errors.New("token not found"),
// 		},
// 		{
// 			name: "❌ Lỗi truy cập database",
// 			mockRepo: &repository.MockTokenRepository{
// 				MockGetUserForAccessToken: func(accessToken string) (*models.User, string, error) {
// 					return nil, "", errors.New("database connection failed")
// 				},
// 			},
// 			expectedUser:         nil,
// 			expectedRefreshToken: "",
// 			expectedErr:          errors.New("database connection failed"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			repository.NewTokenRepository = func() repository.TokenRepository {
// 				return tt.mockRepo
// 			}

// 			user, refresh, err := VerifyAccessToken("dummy-access-token")

// 			// So sánh kết quả
// 			assert.Equal(t, tt.expectedUser, user)
// 			assert.Equal(t, tt.expectedRefreshToken, refresh)

// 			// So sánh lỗi (so sánh message thay vì instance)
// 			if tt.expectedErr != nil {
// 				assert.EqualError(t, err, tt.expectedErr.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			t.Logf("✅ [%s] PASS: user=%+v refresh=%s err=%v", tt.name, user, refresh, err)
// 		})
// 	}
// }

func TestRefreshAccessToken(t *testing.T) {
	// Lưu lại function gốc để restore sau test
	origFunc := TokenSourceFunc
	defer func() { TokenSourceFunc = origFunc }()

	tests := []struct {
		name          string
		mockFunc      func(config *oauth2.Config, ctx context.Context, t *oauth2.Token) oauth2.TokenSource
		expectedToken *oauth2.Token
		expectedErr   error
	}{
		{
			name: "✅ Refresh token thành công",
			mockFunc: func(config *oauth2.Config, ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
				return oauth2.ReuseTokenSource(nil, mockTokenSource{
					token: &oauth2.Token{AccessToken: "new-access-token"},
					err:   nil,
				})
			},
			expectedToken: &oauth2.Token{AccessToken: "new-access-token"},
			expectedErr:   nil,
		},
		{
			name: "❌ Refresh token thất bại",
			mockFunc: func(config *oauth2.Config, ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
				return oauth2.ReuseTokenSource(nil, mockTokenSource{
					token: nil,
					err:   errors.New("refresh failed"),
				})
			},
			expectedToken: nil,
			expectedErr:   errors.New("refresh failed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			TokenSourceFunc = tc.mockFunc

			token, err := RefreshAccessToken("client_id", "client_secret", "refresh_token", &oauth2.Config{})

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedToken.AccessToken, token.AccessToken)
			}
		})
	}
}

// mockTokenSource implement oauth2.TokenSource interface
type mockTokenSource struct {
	token *oauth2.Token
	err   error
}

func (m mockTokenSource) Token() (*oauth2.Token, error) {
	return m.token, m.err
}
