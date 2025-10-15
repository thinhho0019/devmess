package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"project/models"
	"project/repository"

	"project/utils"
)

type UserService struct {
	repoUser repository.UserRepository
}

func NewInitUserService(repoUser *repository.UserRepository) *UserService {
	return &UserService{
		repoUser: *repoUser,
	}
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

	createdUser, err := s.repoUser.CreateUser(user)
	if err != nil {

		return nil, err
	}

	return createdUser, nil
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repoUser: r}
}

func (s *UserService) CheckEmail(email string) (bool, error) {
	user, err := s.repoUser.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
func (s *UserService) LoginPassword(email string, password string) (*models.User, error) {
	fmt.Println(email, password)
	user, err := s.repoUser.LoginPassword(email, password, "local")
	if user == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindUserWithStatusFriend(email string, user_id string) (map[string]interface{}, error) {
	// convert user_id to uuid
	uuid_user_id, err := utils.StringToUUID(user_id)
	if err != nil {
		return nil, err
	}
	user, status, err := s.repoUser.FindUserWithStatusFriend(email, uuid_user_id)
	if user != nil {
		//convert to map
		var userMap map[string]interface{}
		data, _ := json.Marshal(user)
		json.Unmarshal(data, &userMap)
		userMap["status"] = status
		return userMap, nil
	}
	return nil, err
}
