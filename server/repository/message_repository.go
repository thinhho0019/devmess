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

type MessageRepository interface {
	CreateMessage(message *models.Message) (*models.Message, error)
	GetMessageByID(messageID uuid.UUID) (*models.Message, error)
	GetMessagesByConversationID(conversationID uuid.UUID, limit, offset int) ([]*models.Message, error)
	GetLastMessageByConversationID(conversationID uuid.UUID) (*models.Message, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
	DeleteMessage(messageID uuid.UUID) error
	DeleteMessagesByConversationID(conversationID uuid.UUID) error
	SoftDeleteMessage(messageID uuid.UUID, deletedBy uuid.UUID) error

	// Search and filter
	SearchMessages(conversationID uuid.UUID, query string, limit, offset int) ([]*models.Message, error)
	GetMessagesByUser(userID uuid.UUID, limit, offset int) ([]*models.Message, error)
	GetUnreadMessages(conversationID, userID uuid.UUID) ([]*models.Message, error)

	// Message stats
	CountMessagesByConversationID(conversationID uuid.UUID) (int64, error)
	CountUnreadMessages(conversationID, userID uuid.UUID) (int64, error)

	// Message reactions and status
	UpdateMessageStatus(messageID uuid.UUID, status string) error
	GetMessagesAfterTime(conversationID uuid.UUID, after time.Time) ([]*models.Message, error)
	GetMessagesBeforeTime(conversationID uuid.UUID, userID uuid.UUID, before *time.Time, limit int) ([]*models.Message, error)
	GetMessagesBetweenDates(conversationID uuid.UUID, startDate, endDate time.Time) ([]*models.Message, error)
}

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepository() MessageRepository {
	return &messageRepo{
		db: database.DB,
	}
}

// CreateMessage - Tạo message mới
func (r *messageRepo) CreateMessage(message *models.Message) (*models.Message, error) {
	if message == nil {
		return nil, errors.New("message is nil")
	}

	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	if err := r.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Load relationships
	if err := r.db.Preload("Sender").First(message, message.ID).Error; err != nil {
		return message, nil // Return message even if preload fails
	}

	return message, nil
}

// GetMessageByID - Lấy message theo ID
func (r *messageRepo) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
	var message models.Message

	err := r.db.Preload("Sender").
		Preload("Conversation").
		Where("id = ?", messageID).
		First(&message).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

// GetMessagesByConversationID - Lấy messages của conversation với pagination
func (r *messageRepo) GetMessagesByConversationID(conversationID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message

	err := r.db.Preload("Sender").
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Reverse để có order tăng dần theo thời gian
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetLastMessageByConversationID - Lấy message cuối cùng của conversation
func (r *messageRepo) GetLastMessageByConversationID(conversationID uuid.UUID) (*models.Message, error) {
	var message models.Message

	err := r.db.Preload("Sender").
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Order("created_at DESC").
		First(&message).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No messages found
		}
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}

	return &message, nil
}

// UpdateMessage - Cập nhật message
func (r *messageRepo) UpdateMessage(message *models.Message) (*models.Message, error) {
	if message == nil {
		return nil, errors.New("message is nil")
	}

	message.UpdatedAt = time.Now()

	if err := r.db.Save(message).Error; err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	// Load relationships
	if err := r.db.Preload("Sender").First(message, message.ID).Error; err != nil {
		return message, nil
	}

	return message, nil
}

// DeleteMessage - Xóa message vĩnh viễn
func (r *messageRepo) DeleteMessage(messageID uuid.UUID) error {
	result := r.db.Unscoped().Delete(&models.Message{}, "id = ?", messageID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete message: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// SoftDeleteMessage - Xóa message soft delete
func (r *messageRepo) SoftDeleteMessage(messageID uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now()

	result := r.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to soft delete message: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// DeleteMessagesByConversationID - Xóa tất cả messages của conversation
func (r *messageRepo) DeleteMessagesByConversationID(conversationID uuid.UUID) error {
	if err := r.db.Unscoped().Delete(&models.Message{}, "conversation_id = ?", conversationID).Error; err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	return nil
}

// SearchMessages - Tìm kiếm messages trong conversation
func (r *messageRepo) SearchMessages(conversationID uuid.UUID, query string, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message

	err := r.db.Preload("Sender").
		Where("conversation_id = ? AND content ILIKE ? AND deleted_at IS NULL",
			conversationID, "%"+query+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	return messages, nil
}

// GetMessagesByUser - Lấy messages của user
func (r *messageRepo) GetMessagesByUser(userID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message

	err := r.db.Preload("Sender").
		Preload("Conversation").
		Where("sender_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user messages: %w", err)
	}

	return messages, nil
}

// GetUnreadMessages - Lấy messages chưa đọc
func (r *messageRepo) GetUnreadMessages(conversationID, userID uuid.UUID) ([]*models.Message, error) {
	var messages []*models.Message

	// Subquery để lấy last_read_message_id từ participants
	subQuery := r.db.Model(&models.Participant{}).
		Select("last_read_message_id").
		Where("conversation_id = ? AND user_id = ?", conversationID, userID)

	err := r.db.Preload("Sender").
		Where("conversation_id = ? AND sender_id != ? AND deleted_at IS NULL", conversationID, userID).
		Where("id NOT IN (?) OR created_at > (SELECT created_at FROM messages WHERE id = (?))",
			subQuery, subQuery).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get unread messages: %w", err)
	}

	return messages, nil
}

// CountMessagesByConversationID - Đếm số messages trong conversation
func (r *messageRepo) CountMessagesByConversationID(conversationID uuid.UUID) (int64, error) {
	var count int64

	err := r.db.Model(&models.Message{}).
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

// CountUnreadMessages - Đếm số messages chưa đọc
func (r *messageRepo) CountUnreadMessages(conversationID, userID uuid.UUID) (int64, error) {
	var count int64

	// Get participant's last read message
	var participant models.Participant
	err := r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no participant record, count all messages
			return r.CountMessagesByConversationID(conversationID)
		}
		return 0, fmt.Errorf("failed to get participant: %w", err)
	}

	// Try to read last_read_message_id column directly in case the Participant struct doesn't include it
	var lastReadMessageID *uuid.UUID
	_ = r.db.Model(&models.Participant{}).
		Select("last_read_message_id").
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Scan(&lastReadMessageID).Error

	query := r.db.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND deleted_at IS NULL", conversationID, userID)

	if lastReadMessageID != nil {
		// Count messages after last read message
		var lastReadTime time.Time
		r.db.Model(&models.Message{}).
			Select("created_at").
			Where("id = ?", *lastReadMessageID).
			Scan(&lastReadTime)

		query = query.Where("created_at > ?", lastReadTime)
	}

	err = query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	return count, nil
}

// UpdateMessageStatus - Cập nhật status của message
func (r *messageRepo) UpdateMessageStatus(messageID uuid.UUID, status string) error {
	result := r.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update message status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// GetMessagesAfterTime - Lấy messages sau thời điểm cụ thể
func (r *messageRepo) GetMessagesAfterTime(conversationID uuid.UUID, after time.Time) ([]*models.Message, error) {
	var messages []*models.Message

	err := r.db.Preload("Sender").
		Where("conversation_id = ? AND created_at > ? AND deleted_at IS NULL", conversationID, after).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get messages after time: %w", err)
	}

	return messages, nil
}

func (r *messageRepo) GetMessagesBeforeTime(conversationID uuid.UUID, userID uuid.UUID, before *time.Time, limit int) ([]*models.Message, error) {
	var messages []*models.Message

	if limit <= 0 {
		limit = 50
	}

	query := r.db.Table("messages").
		Select("messages.*").
		Joins("JOIN conversations ON conversations.id = messages.conversation_id").
		Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("participants.user_id = ? and participants.conversation_id = ?", userID, conversationID).
		Where("messages.deleted_at IS NULL")

	if before != nil {
		query = query.Where("messages.created_at < ?", *before)
	}

	err := query.
		Preload("Sender").
		Order("messages.created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get messages by user: %w", err)
	}

	return messages, nil
}

// GetMessagesBetweenDates - Lấy messages trong khoảng thời gian
func (r *messageRepo) GetMessagesBetweenDates(conversationID uuid.UUID, startDate, endDate time.Time) ([]*models.Message, error) {
	var messages []*models.Message

	err := r.db.Preload("Sender").
		Where("conversation_id = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL",
			conversationID, startDate, endDate).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get messages between dates: %w", err)
	}

	return messages, nil
}
