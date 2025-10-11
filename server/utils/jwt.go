package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Lấy secret key từ biến môi trường. Rất quan trọng để giữ bí mật này an toàn!
// Trong môi trường development, nó sẽ dùng một giá trị mặc định.
var jwtSecretKey = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		// Giá trị mặc định này KHÔNG an toàn cho production.
		// Hãy chắc chắn rằng bạn đã đặt biến môi trường JWT_SECRET_KEY.
		return "my-super-secret-key-for-dev"
	}
	return secret
}

// CustomClaims chứa dữ liệu của token (payload) và thời gian hết hạn.
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken tạo ra một JWT mới cho một người dùng cụ thể.
func GenerateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	// Thiết lập claims
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Thời gian hết hạn của token
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			// Thời gian token được tạo
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Issuer (người phát hành)
			Issuer: "web_chat_app",
		},
	}

	// Tạo token với claims và phương thức ký HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với secret key và trả về chuỗi token
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	// Refresh token thường có thời gian sống dài, ví dụ 30 ngày.
	duration := 14 * 24 * time.Hour

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Thời gian hết hạn của token
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			// Thời gian token được tạo
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Issuer khác để phân biệt với access token
			Issuer: "web_chat_app_refresh",
			// Thêm một ID duy nhất cho token (JTI - JWT ID)
			ID: uuid.New().String(),
		},
	}

	// Tạo token với claims và phương thức ký HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với secret key và trả về chuỗi token
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func ValidateRefreshToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("refresh token has expired")
		}
		return nil, errors.New("invalid refresh token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Kiểm tra issuer để đảm bảo đây là refresh token
		if claims.Issuer != "web_chat_app_refresh" {
			return nil, errors.New("invalid token issuer: expected refresh token")
		}
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

// ValidateToken xác thực một chuỗi token và trả về claims nếu hợp lệ.
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Đảm bảo rằng thuật toán ký là HS256 như chúng ta đã dùng để tạo token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		// Xử lý các lỗi phổ biến như token hết hạn
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("invalid token")
	}

	// Nếu token hợp lệ, trích xuất claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
