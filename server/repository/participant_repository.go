package repository

import (
	"errors"
	"fmt"
	"project/database"
	"project/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	CreateParticipant(participant *models.Participant) error
	GetParticipant(conversationID, userID uuid.UUID) (*models.Participant, error)
	GetParticipantsByConversationID(conversationID uuid.UUID) (*[]models.Participant, error)
	GetParticipantsByUserID(userID uuid.UUID) ([]*models.Participant, error)
	UpdateParticipant(participant *models.Participant) error
	DeleteParticipant(conversationID, userID uuid.UUID) error
	DeleteParticipantsByConversationID(conversationID uuid.UUID) error
	IsParticipant(conversationID, userID uuid.UUID) (bool, error)
}

type participantRepo struct {
	db *gorm.DB
}

func NewParticipantRepository() ParticipantRepository {
	return &participantRepo{
		db: database.DB,
	}
}

// CreateParticipant - Tạo participant mới
func (r *participantRepo) CreateParticipant(participant *models.Participant) error {
	if participant == nil {
		return errors.New("participant is nil")
	}

	if participant.ID == uuid.Nil {
		participant.ID = uuid.New()
	}

	participant.JoinedAt = time.Now()

	if err := r.db.Create(participant).Error; err != nil {
		return fmt.Errorf("failed to create participant: %w", err)
	}

	return nil
}

// GetParticipant - Lấy participant theo conversation ID và user ID
func (r *participantRepo) GetParticipant(conversationID, userID uuid.UUID) (*models.Participant, error) {
	var participant models.Participant

	err := r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("participant not found")
		}
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	return &participant, nil
}

// GetParticipantsByConversationID - Lấy tất cả participants của conversation
func (r *participantRepo) GetParticipantsByConversationID(conversationID uuid.UUID) (*[]models.Participant, error) {
	var participants *[]models.Participant

	err := r.db.Preload("User").
		Where("conversation_id = ? ", conversationID).
		Order("joined_at ASC").
		Find(&participants).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return participants, nil
}

// GetParticipantsByUserID - Lấy tất cả participants của user
func (r *participantRepo) GetParticipantsByUserID(userID uuid.UUID) ([]*models.Participant, error) {
	var participants []*models.Participant

	err := r.db.Preload("Conversation").
		Where("user_id = ? AND left_at IS NULL", userID).
		Order("joined_at DESC").
		Find(&participants).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user participants: %w", err)
	}

	return participants, nil
}

// UpdateParticipant - Cập nhật participant
func (r *participantRepo) UpdateParticipant(participant *models.Participant) error {
	if participant == nil {
		return errors.New("participant is nil")
	}

	if err := r.db.Save(participant).Error; err != nil {
		return fmt.Errorf("failed to update participant: %w", err)
	}

	return nil
}

// DeleteParticipant - Xóa participant khỏi conversation
func (r *participantRepo) DeleteParticipant(conversationID, userID uuid.UUID) error {
	result := r.db.Delete(&models.Participant{},
		"conversation_id = ? AND user_id = ?", conversationID, userID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete participant: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("participant not found")
	}

	return nil
}

// DeleteParticipantsByConversationID - Xóa tất cả participants của conversation
func (r *participantRepo) DeleteParticipantsByConversationID(conversationID uuid.UUID) error {
	if err := r.db.Delete(&models.Participant{}, "conversation_id = ?", conversationID).Error; err != nil {
		return fmt.Errorf("failed to delete participants: %w", err)
	}

	return nil
}

// IsParticipant - Kiểm tra user có phải participant không
func (r *participantRepo) IsParticipant(conversationID, userID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Model(&models.Participant{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NULL", conversationID, userID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check participant: %w", err)
	}

	return count > 0, nil
}
