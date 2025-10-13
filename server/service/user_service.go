package service

import (
	"errors"
	"fmt"
	"project/models"
	"project/repository"
	"project/utils"
)

type UserService struct {
	repo repository.UserRepository
}

func (s *UserService) RegisterUser(name string, email string, password string) (*models.User, error) {

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
		Avatar:   "img.jpg",
		Provider: "local",
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {

		return nil, err
	}

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
