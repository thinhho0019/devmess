package service

import (
	"errors"
	"project/models"
	"project/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockRepo      func(email string) (*models.User, error)
		expectedEmail bool
		expectError   bool
	}{
		{
			name:  "Email exists",
			email: "a@example.com",
			mockRepo: func(email string) (*models.User, error) {
				return &models.User{Email: email}, nil
			},
			expectedEmail: true,
			expectError:   false,
		},
		{
			name:  "Email not exists",
			email: "b@example.com",
			mockRepo: func(email string) (*models.User, error) {
				return nil, nil
			},
			expectedEmail: false,
			expectError:   false,
		},
		{
			name:  "Repo error",
			email: "c@example.com",
			mockRepo: func(email string) (*models.User, error) {
				return nil, errors.New("db error")
			},
			expectedEmail: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &repository.MockUserRepository{
				MockGetUserByEmail: tt.mockRepo,
			}

			service := &UserService{
				repo: mockRepo,
			}

			actual, err := service.CheckEmail(tt.email)

			assert.Equal(t, tt.expectedEmail, actual)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
