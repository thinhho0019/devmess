package service

import (
	"context"
	"errors"
	"fmt"
	"project/database"
	"project/models"
	"project/repository"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"gorm.io/gorm"
)

var ErrInvalidOrExpired = errors.New("invalid or expired access token")

// LoginWithGoogle xử lý logic login Google
func LoginWithGoogle(userInfo *models.GoogleUserInfo, ip, userAgent string) (*models.Token, *models.Device, error) {
	fmt.Println("🔹 Start LoginWithGoogle")
	fmt.Printf("UserInfo: Email=%s, Name=%s\n", userInfo.Email, userInfo.Name)
	fmt.Printf("Client: IP=%s, User-Agent=%s\n", ip, userAgent)

	// 1️⃣ Kiểm tra user đã tồn tại
	user, err := repository.GetUserByEmail(userInfo.Email)
	if err != nil {
		fmt.Println("❌ Error getting user:", err)
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		fmt.Println("🆕 User not found, creating new user...")
		newUser := &models.User{
			Name:   userInfo.Name,
			Email:  userInfo.Email,
			Avatar: userInfo.Picture,
		}
		user, err = repository.CreateUser(newUser)
		if err != nil {
			fmt.Println("❌ Error creating user:", err)
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
		fmt.Println("✅ User created:", user.Email)
	} else {
		fmt.Println("✅ User found:", user.Email)
	}

	// 2️⃣ Nhận dạng thiết bị
	deviceType, deviceName := detectDevice(userAgent)
	fmt.Printf("Detected Device: Type=%s, Name=%s\n", deviceType, deviceName)

	// Kiểm tra device đã tồn tại chưa
	existingDevice, err := repository.GetDeviceByInfo(user.ID, ip, userAgent)
	if err != nil {
		fmt.Println("❌ Error checking existing device:", err)
		return nil, nil, fmt.Errorf("failed to check existing device: %w", err)
	}
	if existingDevice == nil {
		fmt.Println("🆕 Device not found, will create new device")
	} else {
		fmt.Println("♻ Existing device found:", existingDevice.Type, existingDevice.Name)
	}

	// 3️⃣ Transaction tạo token + gắn device
	var token *models.Token
	var device *models.Device
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Xử lý device
		if existingDevice == nil {
			device = &models.Device{
				ID:        uuid.New(),
				Type:      deviceType,
				UserID:    user.ID,
				Name:      deviceName,
				IP:        ip,
				UserAgent: userAgent,
			}
			if err := tx.Create(device).Error; err != nil {
				fmt.Println("❌ Error creating device:", err)
				return fmt.Errorf("failed to create device: %w", err)
			}
			fmt.Println("✅ Device created:", device.Type, device.Name)
		} else {
			device = existingDevice

			if err := tx.Save(device).Error; err != nil {
				fmt.Println("❌ Error updating device token:", err)
				return fmt.Errorf("failed to update device token: %w", err)
			}
			fmt.Println("♻ Device token updated:", device.Type, device.Name)
		}
		// Tạo token mới
		token = &models.Token{
			ID:           uuid.New(),
			AccessToken:  userInfo.AccessToken,
			RefreshToken: userInfo.RefreshToken,
			ExpiresAt:    int64(userInfo.ExpiresIn),
			DeviceID:     device.ID,
			TokenType:    "Bearer",
		}
		if err := tx.Create(token).Error; err != nil {
			fmt.Println("❌ Error creating token:", err)
			return fmt.Errorf("failed to create token: %w", err)
		}
		fmt.Println("✅ Token created:", token.AccessToken)

		return nil
	})

	if err != nil {
		fmt.Println("❌ Login transaction failed:", err)
		return nil, nil, fmt.Errorf("login transaction failed: %w", err)
	}

	fmt.Println("✅ LoginWithGoogle finished successfully")
	fmt.Printf("Result: Token=%s, Device=%s\n", token.AccessToken, device.Name)
	return token, device, nil
}

// detectDevice đơn giản parse UserAgent
func detectDevice(userAgent string) (string, string) {
	ua := userAgent
	deviceType := "Web"
	deviceName := "Browser"

	lowerUA := strings.ToLower(ua)
	if strings.Contains(lowerUA, "android") {
		deviceType = "Android"
	} else if strings.Contains(lowerUA, "iphone") {
		deviceType = "iOS"
	}

	if strings.Contains(lowerUA, "chrome") {
		deviceName = "Chrome"
	} else if strings.Contains(lowerUA, "safari") {
		deviceName = "Safari"
	} else if strings.Contains(lowerUA, "firefox") {
		deviceName = "Firefox"
	}

	return deviceType, deviceName
}

func VerifyAccessToken(accessToken string) (*models.User, string, error) {
	user, refresh_token, err := repository.GetUserForAccessToken(accessToken)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, refresh_token, ErrInvalidOrExpired
	}
	return user, "", nil
}
func RefreshAccessToken(clientID, clientSecret, refreshToken string, googleAuthConfig *oauth2.Config) (*oauth2.Token, error) {
	print("rf", refreshToken)
	token := &oauth2.Token{RefreshToken: refreshToken}
	print("rf", token.RefreshToken)
	newToken, err := googleAuthConfig.TokenSource(context.Background(), token).Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
