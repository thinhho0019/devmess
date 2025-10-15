package repository

// MockUserRepository là struct mô phỏng UserRepository (dùng cho unit test)
import (
	"project/models"

	"github.com/google/uuid"
)

type MockUserRepository struct {
	MockCreateUser          func(user *models.User) (*models.User, error)
	MockGetUserByID         func(id uuid.UUID) (*models.User, error)
	MockGetUserByEmail      func(email string) (*models.User, error)
	MockLoginPassword       func(email string, password string, provider string) (*models.User, error)
	MockUpdateUser          func(user *models.User) error
	MockDeleteUser          func(id uint) error
	MockGetAllUsers         func() ([]models.User, error)
	MockFindUserWithFriends func(email string, user_id uuid.UUID) (*models.User, string, error)
	MockGetUserByAccesToken func(accessToken string) (*models.User, error)
}

// Implement interface UserRepository ↓↓↓

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	if m.MockCreateUser != nil {
		return m.MockCreateUser(user)
	}
	return nil, nil
}

func (m *MockUserRepository) LoginPassword(email string, password string, provider string) (*models.User, error) {
	if m.MockLoginPassword != nil {
		return m.MockLoginPassword(email, password, provider)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	if m.MockGetUserByID != nil {
		return m.MockGetUserByID(id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	if m.MockGetUserByEmail != nil {
		return m.MockGetUserByEmail(email)
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	if m.MockUpdateUser != nil {
		return m.MockUpdateUser(user)
	}
	return nil
}

func (m *MockUserRepository) DeleteUser(id uint) error {
	if m.MockDeleteUser != nil {
		return m.MockDeleteUser(id)
	}
	return nil
}

func (m *MockUserRepository) GetAllUsers() ([]models.User, error) {
	if m.MockGetAllUsers != nil {
		return m.MockGetAllUsers()
	}
	return nil, nil
}

func (m *MockUserRepository) FindUserWithStatusFriend(email string, user_id uuid.UUID) (*models.User, string, error) {
	if m.MockFindUserWithFriends != nil {
		return m.MockFindUserWithFriends(email, user_id)
	}
	return nil, "", nil
}
func (m *MockUserRepository) GetUserByAccesToken(accessToken string) (*models.User, error) {
	if m.MockGetUserByAccesToken != nil {
		return m.MockGetUserByAccesToken(accessToken)
	}
	return nil, nil
}
