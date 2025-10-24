package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (s *UserService) FindUserWithStatusFriend(email string, userID string) (map[string]interface{}, error) {
	log.Printf("[FindUserWithStatusFriend] Input email: %s, userID: %s", email, userID)

	// Convert user_id sang uuid
	uuidUserID, err := utils.StringToUUID(userID)
	if err != nil {
		log.Printf("[FindUserWithStatusFriend] Invalid userID: %v", err)
		return nil, err
	}
	log.Printf("[FindUserWithStatusFriend] Converted userID to UUID: %s", uuidUserID)

	// Gọi repo để lấy user + status
	user, status, err := s.repoUser.FindUserWithStatusFriend(email, uuidUserID)
	if err != nil {
		log.Printf("[FindUserWithStatusFriend] Repo error: %v", err)
		return nil, err
	}

	if user == nil {
		log.Printf("[FindUserWithStatusFriend] User not found for email: %s", email)
		return nil, nil
	}

	// Chuyển user sang map và thêm status
	userMap := make(map[string]interface{})
	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("[FindUserWithStatusFriend] Error marshaling user: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(data, &userMap); err != nil {
		log.Printf("[FindUserWithStatusFriend] Error unmarshaling to map: %v", err)
		return nil, err
	}

	userMap["status"] = status
	log.Printf("[FindUserWithStatusFriend] Result userMap: %+v", userMap)

	return userMap, nil
}
