package service

import (
	"fmt"
	 

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	OAuthConfig *oauth2.Config
}

func NewGoogleService(oauthConfig *oauth2.Config) *GoogleService {
	return &GoogleService{
		OAuthConfig: oauthConfig,
	}
}
func (h *GoogleService) InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	fmt.Println("Initializing Google OAuth with Client ID:", clientID)
	fmt.Println("Client Secret:", clientSecret)
	h.OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
