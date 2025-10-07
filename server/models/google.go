package models

// GoogleUserInfo lưu thông tin người dùng Google trả về
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
}
