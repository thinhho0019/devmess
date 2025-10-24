package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"project/models"
	"project/repository"
	"project/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var ErrInvalidOrExpired = errors.New("invalid or expired access token")
var TokenSourceFunc = func(cfg *oauth2.Config, ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	return cfg.TokenSource(ctx, t)
}

type AuthService struct {
	userRepo   repository.UserRepository
	deviceRepo repository.DeviceRepository
	tokenRepo  repository.TokenRepository
	redisRepo  repository.RedisRepository
	authGoogle *oauth2.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	deviceRepo repository.DeviceRepository,
	tokenRepo repository.TokenRepository,
	redisRepo repository.RedisRepository,
	authGoogle *oauth2.Config,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		tokenRepo:  tokenRepo,
		redisRepo:  redisRepo,
		authGoogle: authGoogle,
	}
}

// LoginWithGoogle xử lý logic login Google
func (a *AuthService) LoginWithGoogle(
	userInfo *models.GoogleUserInfo,
	ip, userAgent string,
) (*models.Token, *models.Device, error) {
	fmt.Println("🔹 Start LoginWithGoogle")

	// 1️⃣ Kiểm tra user đã tồn tại
	user, err := a.userRepo.GetUserByEmail(userInfo.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		// save avater to server
		uuid_user := uuid.New()
		urlImage, err := a.SaveImageGoogle(userInfo.Picture, fmt.Sprintf("%s.jpg", uuid_user.String()))
		if err != nil {
			fmt.Println("❌ Failed to save image:", err)
			urlImage = "avatar/img.jpg"
		}
		user = &models.User{
			ID:       uuid_user,
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Avatar:   urlImage,
			Provider: "google",
		}
		user, err = a.userRepo.CreateUser(user)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Gọi hàm tạo session chung
	return a.CreateSession(user, ip, userAgent, userInfo.AccessToken, userInfo.RefreshToken, int64(userInfo.ExpiresIn), "google")
}
func (h *AuthService) SaveImageGoogle(imageURL string, fileName string) (string, error) {
	// Implement your logic to save the image URL to Google Drive or any other service
	//download image from imageURL and save to google drive
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Save the image to Google Drive or any other service
	// ...
	if err := os.MkdirAll("uploads/avatar", os.ModePerm); err != nil {
		return "", err
	}
	file, err := os.Create("uploads/avatar/" + fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	return "avatar/" + fileName, nil
}

// CreateSession tạo hoặc cập nhật device, tạo token và lưu vào Redis
func (a *AuthService) CreateSession(

	user *models.User,
	ip, userAgent, accessToken, refreshToken string,
	expiresIn int64,
	provider string,
) (*models.Token, *models.Device, error) {
	// 2️⃣ Nhận dạng thiết bị
	deviceType, deviceName := detectDevice(userAgent)
	user.Provider = provider
	existingDevice, err := a.deviceRepo.GetDeviceByInfo(user.ID, ip, userAgent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check existing device: %w", err)
	}

	var device *models.Device
	if existingDevice == nil {
		device = &models.Device{
			ID:        uuid.New(),
			UserID:    user.ID,
			Type:      deviceType,
			Name:      deviceName,
			IP:        ip,
			UserAgent: userAgent,
		}
		if _, err := a.deviceRepo.CreateDevice(device); err != nil {
			return nil, nil, err
		}
	} else {
		device = existingDevice
		if _, err := a.deviceRepo.UpdateDevice(device); err != nil {
			return nil, nil, err
		}
	}

	// 3️⃣ Token
	existingToken, err := a.tokenRepo.GetTokensByDeviceID(user.ID, device.ID.String())
	if err != nil {
		return nil, nil, err
	}
	if existingToken != nil {
		if err := a.tokenRepo.DeleteTokenByID(existingToken.ID); err != nil {
			return nil, nil, err
		}
	}

	// Nếu không có accessToken (trường hợp đăng nhập bằng mật khẩu), tạo một cái mới.
	if accessToken == "" {
		// Đặt thời gian hết hạn cho access token, ví dụ: 24 giờ
		accessTokenDuration := 24 * time.Hour
		newAccessToken, err := utils.GenerateToken(user.ID, accessTokenDuration)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
		}
		accessToken = newAccessToken // Gán giá trị cho biến accessToken bên ngoài
		expiresIn = time.Now().Add(accessTokenDuration).Unix()
	}

	// Nếu không có refreshToken, tạo một cái mới.
	if refreshToken == "" {
		newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
		}
		refreshToken = newRefreshToken // Gán giá trị cho biến refreshToken bên ngoài
	}

	token := &models.Token{
		ID:           uuid.New(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresIn,
		DeviceID:     device.ID,
		TokenType:    provider,
	}
	fmt.Printf("Token struct: %+v\n", token)
	if err := a.tokenRepo.CreateToken(token); err != nil {
		return nil, nil, err
	}
	time_duration := time.Until(time.Unix(token.ExpiresAt, 0))
	if err := a.redisRepo.SetToken(token.AccessToken, user, time_duration); err != nil {
		fmt.Println("❌ Failed to save access token to Redis:", err)
		return token, device, errors.New("error save access token")
	}
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

func (a *AuthService) VerifyAccessToken(accessToken string) (*models.User, string, error) {
	// check redis exits token
	if user, err := a.redisRepo.GetUserByToken(accessToken); err == nil {

		return user, "", nil
	}
	user, refreshToken, err := a.tokenRepo.GetUserForAccessToken(accessToken)

	// 1️⃣ Nếu có lỗi nhưng repo vẫn cấp refresh token mới
	if err != nil {
		return nil, "", err
	}

	// 2️⃣ Nếu user hợp lệ → token OK
	if user != nil {
		return user, "", nil
	}

	// 3️⃣ Nếu không có user nhưng có refresh token (trường hợp hiếm)
	if refreshToken != "" {
		return nil, refreshToken, ErrInvalidOrExpired
	}

	// 4️⃣ Trường hợp không có gì hợp lệ
	return nil, "", ErrInvalidOrExpired
}

func RefreshAccessToken(refreshToken string, googleAuthConfig *oauth2.Config) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	newToken, err := TokenSourceFunc(googleAuthConfig, context.Background(), token).Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
func (a *AuthService) RefreshToken(refreshToken string, accessToken string, googleAuthConfig *oauth2.Config) (*models.Token, error) {
	token, err := a.tokenRepo.GetTokenByRefresh(refreshToken)
	var user *models.User
	if err != nil {
		a.redisRepo.DeleteToken(accessToken)
		return nil, err
	}

	user, err = a.redisRepo.GetUserByToken(accessToken)
	if err != nil {
		// check in sql
		user, err = a.userRepo.GetUserByAccesToken(accessToken)
		if err != nil {
			return nil, errors.New("can not get user, please login again")
		}
	}
	if token.TokenType == "google" {
		// Xử lý làm mới token Google nếu cần thiết
		// Giả sử bạn có hàm `RefreshGoogleToken` để làm việc này
		newToken, err := RefreshAccessToken(token.RefreshToken,
			googleAuthConfig)
		if err != nil {
			a.redisRepo.DeleteToken(accessToken)
			return nil, err
		}
		if newToken != nil {
			token.AccessToken = newToken.AccessToken
			token.ExpiresAt = newToken.Expiry.Unix()
			token.RefreshToken = newToken.RefreshToken
			if err := a.tokenRepo.UpdateToken(token); err != nil {

				a.redisRepo.DeleteToken(accessToken)
				return nil, err
			}
			a.redisRepo.SetToken(token.AccessToken, user, 0)
			return token, nil
		}
	} else {
		// kiểm tra refresh token còn thời gian không
		claims, err := utils.ValidateRefreshToken(token.RefreshToken)
		if err != nil {
			a.redisRepo.DeleteToken(accessToken)
			return nil, err
		}
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			a.redisRepo.DeleteToken(accessToken)
			return nil, errors.New("refresh token has expired")
		}
		// Xử lý làm mới token thông thường
		// Đặt thời gian hết hạn cho access token, ví dụ: 24 giờ
		accessTokenDuration := 1 * time.Hour
		newAccessToken, err := utils.GenerateToken(token.ID, accessTokenDuration)
		if err != nil {
			return nil, err
		}
		token.AccessToken = newAccessToken

		if err := a.tokenRepo.UpdateToken(token); err != nil {
			return nil, err
		}
		a.redisRepo.SetToken(token.AccessToken, user, 0)
		return token, nil
	}
	return nil, nil
}
func (a *AuthService) HandleGoogleCallback(code, ip, userAgent string, GoogleOAuthConfig *oauth2.Config) (*models.GoogleUserInfo, *models.Token, *models.Device, error) {
	if GoogleOAuthConfig == nil {
		return nil, nil, nil, errors.New("GoogleOAuthConfig not initialized")
	}

	token, err := GoogleOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("exchange token failed: %v", err)
	}

	client := GoogleOAuthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user info failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("google returned status %d", resp.StatusCode)
	}

	var userInfo models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, nil, nil, fmt.Errorf("parse user info failed: %v", err)
	}

	userInfo.AccessToken = token.AccessToken
	userInfo.RefreshToken = token.RefreshToken
	// userInfo.ExpiresIn = token.Expiry.Unix()
	userInfo.ExpiresIn = time.Now().Add(24 * time.Hour).Unix()

	tokenModel, device, err := a.LoginWithGoogle(&userInfo, ip, userAgent)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("login failed: %v", err)
	}

	return &userInfo, tokenModel, device, nil
}
