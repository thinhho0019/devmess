package service

import (
	"errors"
	"fmt"
	"os"
	"project/models"
	"project/repository"
	"project/utils"
)

type UserService struct {
	repo repository.UserRepository
}

func (s *UserService) RegisterUser(name string, email string, password string) (*models.User, error) {
	// Validate input
	if name == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Avatar:   fmt.Sprintf("%s/api/?name=%s", os.Getenv("DEFAULT_URL_SERVER"), "img.jpg"),
		Provider: "local",
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		// Nếu có lỗi khi tạo người dùng, trả về lỗi đó.
		return nil, err
	}

	// Nếu không có lỗi, trả về người dùng đã được tạo (có thể đã được cập nhật ID từ DB).
	return createdUser, nil
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CheckEmail(email string) (bool, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
func (s *UserService) LoginPassword(email string, password string) (*models.User, error) {
	fmt.Println(email, password)
	user, err := s.repo.LoginPassword(email, password, "local")
	if user == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
