package repository

import (
	"errors"
	"project/database"
	"project/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendshipRepository interface {
	CreateFriendship(f *models.Friendship) (*models.Friendship, error)
	GetByID(id uuid.UUID) (*models.Friendship, error)
	GetFriendship(userID, friendID uuid.UUID) (*models.Friendship, error)
	ListFriends(userID uuid.UUID, status string) ([]models.Friendship, error)
	ListPendingRequests(userID uuid.UUID) ([]models.Friendship, error)
	UpdateFriendship(f *models.Friendship) (*models.Friendship, error)
	DeleteFriendship(id uuid.UUID) error
	AreFriends(userID, friendID uuid.UUID) (bool, error)
	GetFriendshipBetweenUsers(userID uuid.UUID, friendID uuid.UUID) (*models.Friendship, error)
}

type friendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository() FriendshipRepository {
	return &friendshipRepository{
		db: database.DB,
	}
}

func (r *friendshipRepository) GetFriendshipBetweenUsers(userID uuid.UUID, friendID uuid.UUID) (*models.Friendship, error) {
	var f models.Friendship
	err := r.db.
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID).
		Preload("User").
		Preload("Friend").
		First(&f).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *friendshipRepository) CreateFriendship(f *models.Friendship) (*models.Friendship, error) {
	if f == nil {
		return nil, errors.New("friendship is nil")
	}
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	if err := r.db.Create(f).Error; err != nil {
		return nil, err
	}
	// preload relations for convenience
	_ = r.db.Preload("User").Preload("Friend").First(f, "id = ?", f.ID).Error
	return f, nil
}

func (r *friendshipRepository) GetByID(id uuid.UUID) (*models.Friendship, error) {
	var f models.Friendship
	if err := r.db.Preload("User").Preload("Friend").First(&f, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *friendshipRepository) GetFriendship(userID, friendID uuid.UUID) (*models.Friendship, error) {
	var f models.Friendship
	err := r.db.
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID).
		Preload("User").
		Preload("Friend").
		First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *friendshipRepository) ListFriends(userID uuid.UUID, status string) ([]models.Friendship, error) {
	var list []models.Friendship
	query := r.db.Preload("Friend").Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *friendshipRepository) ListPendingRequests(userID uuid.UUID) ([]models.Friendship, error) {
	var list []models.Friendship
	// pending requests received by user (friend_id = userID and status = 'pending')
	if err := r.db.Preload("User").Where("friend_id = ? AND status = ?", userID, "pending").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *friendshipRepository) UpdateFriendship(f *models.Friendship) (*models.Friendship, error) {
	if f == nil {
		return nil, errors.New("friendship is nil")
	}
	if err := r.db.Save(f).Error; err != nil {
		return nil, err
	}
	_ = r.db.Preload("User").Preload("Friend").First(f, "id = ?", f.ID).Error
	return f, nil
}

func (r *friendshipRepository) DeleteFriendship(id uuid.UUID) error {
	if err := r.db.Delete(&models.Friendship{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *friendshipRepository) AreFriends(userID, friendID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Friendship{}).
		Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?", userID, friendID, friendID, userID, "accepted").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
