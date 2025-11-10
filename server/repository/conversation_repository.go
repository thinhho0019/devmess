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

type ConversationRepository interface {
	CreateConversation(conversation *models.Conversation) (*models.Conversation, error)
	CreateDirectConversation(userID1 uuid.UUID, userID2 uuid.UUID) (*models.Conversation, error)
	GetConversationByID(conversationID uuid.UUID) (*models.Conversation, error)

	GetConversationsByUserID(
		userID uuid.UUID,
		limit int,
		before *time.Time, // con trỏ thời gian: chỉ lấy các cuộc trò chuyện cũ hơn thời điểm này
	) ([]*models.Conversation, error)
	CountConversationsByUserID(userID uuid.UUID) (int64, error)
	UpdateLastMessageIDInConversation(conversationID uuid.UUID, messageID uuid.UUID) error
	UpdateConversation(conversation *models.Conversation) (*models.Conversation, error)
	DeleteConversation(conversationID uuid.UUID) error
	GetDirectConversation(userID1, userID2 uuid.UUID) (*models.Conversation, error)
	GetConversationWithParticipants(conversationID uuid.UUID) (*models.Conversation, error)
	GetMessageByConversationID(conversationID uuid.UUID, limit int, before *time.Time) ([]*models.Message, error)
	FindConversationByTwoUserID(userID1 uuid.UUID, userID2 uuid.UUID) (*models.Conversation, error)
}

type conversationRepo struct {
	db *gorm.DB
}

func NewConversationRepository() ConversationRepository {
	return &conversationRepo{
		db: database.DB,
	}
}

// CreateConversation - Tạo conversation mới
func (r *conversationRepo) CreateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	if conversation == nil {
		return nil, errors.New("conversation is nil")
	}

	if conversation.ID == uuid.Nil {
		conversation.ID = uuid.New()
	}

	conversation.CreatedAt = time.Now()
	conversation.UpdatedAt = time.Now()

	if err := r.db.Create(conversation).Error; err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation, nil
}

func (r *conversationRepo) UpdateLastMessageIDInConversation(conversationID uuid.UUID, messageID uuid.UUID) error {
	result := r.db.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_id", messageID)
	if result.Error != nil {
		return fmt.Errorf("failed to update last message ID: %w", result.Error)
	}
	return nil
}

func (r *conversationRepo) GetMessageByConversationID(conversationID uuid.UUID, limit int, before *time.Time) ([]*models.Message, error) {
	var messages []*models.Message
	query := r.db.Where("conversation_id = ?", conversationID)
	// Nếu có cursor "before" (load thêm các message cũ hơn)
	if before != nil {
		query = query.Where("created_at < ?", before)
	}
	err := query.Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by conversation ID: %w", err)
	}
	return messages, nil
}

func (r *conversationRepo) CreateDirectConversation(userID1 uuid.UUID, userID2 uuid.UUID) (*models.Conversation, error) {
	var existing models.Conversation

	// 🔍 Kiểm tra đã có conversation direct giữa 2 người này chưa
	err := r.db.
		Joins("JOIN participants p1 ON p1.conversation_id = conversations.id").
		Joins("JOIN participants p2 ON p2.conversation_id = conversations.id").
		Where("conversations.type = ?", "direct").
		Where("p1.user_id = ? AND p2.user_id = ?", userID1, userID2).
		First(&existing).Error

	if err == nil {
		// ✅ Đã tồn tại → trả về luôn
		return &existing, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}
	// create message new direct conversation
	message_id := uuid.New()
	conversation_id := uuid.New()

	// ⚙️ Tạo mới conversation
	conversation := &models.Conversation{
		ID:        conversation_id,
		Type:      "direct",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ⚡ Transaction để đảm bảo toàn vẹn dữ liệu
	tx := r.db.Begin()
	if err := tx.Create(conversation).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	message := &models.Message{
		ID:             message_id,
		ConversationID: conversation_id,
		SenderID:       nil,
		Type:           "system",
		Content:        "Đây là cuộc trò chuyện riêng tư hãy chào nhau đi!",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create initial message: %w", err)
	}
	conversation.LastMessageID = &message_id
	if err := tx.Save(conversation).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update conversation with last message: %w", err)
	}
	// 👥 Thêm 2 participant (user1 và user2)
	participants := []models.Participant{
		{
			ConversationID: conversation.ID,
			UserID:         userID1,
			JoinedAt:       time.Now(),
		},
		{
			ConversationID: conversation.ID,
			UserID:         userID2,
			JoinedAt:       time.Now(),
		},
	}

	if err := tx.Create(&participants).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add participants: %w", err)
	}

	tx.Commit()
	return conversation, nil
}

// GetConversationByID - Lấy conversation theo ID
func (r *conversationRepo) GetConversationByID(conversationID uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation

	err := r.db.Where("id = ?", conversationID).First(&conversation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}

// GetConversationsByUserID - Lấy danh sách conversations của user với pagination (cursor-based)
func (r *conversationRepo) GetConversationsByUserID(
	userID uuid.UUID,
	limit int,
	before *time.Time, // con trỏ thời gian: chỉ lấy các cuộc trò chuyện cũ hơn thời điểm này
) ([]*models.Conversation, error) {
	fmt.Println("userId", userID)
	var conversations []*models.Conversation
	query := r.db.Table("conversations").
		Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("participants.user_id = ?", userID)

	// Nếu có cursor "before" (load thêm các conversation cũ hơn)
	if before != nil {
		query = query.Where("conversations.updated_at < ?", before)
	}
	err := query.
		Order("conversations.updated_at DESC").
		Limit(limit).
		Preload("Participants").Preload("Participants.User").
		Preload("LastMessage").
		Find(&conversations).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user conversations: %w", err)
	}
	return conversations, nil
}

// CountConversationsByUserID - Đếm số lượng conversations của user
func (r *conversationRepo) CountConversationsByUserID(userID uuid.UUID) (int64, error) {
	var count int64

	err := r.db.Table("conversations").
		Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("participants.user_id = ? AND participants.left_at IS NULL", userID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count user conversations: %w", err)
	}

	return count, nil
}

// UpdateConversation - Cập nhật conversation
func (r *conversationRepo) UpdateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	if conversation == nil {
		return nil, errors.New("conversation is nil")
	}

	conversation.UpdatedAt = time.Now()

	if err := r.db.Save(conversation).Error; err != nil {
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	return conversation, nil
}

// DeleteConversation - Xóa conversation
func (r *conversationRepo) DeleteConversation(conversationID uuid.UUID) error {
	result := r.db.Delete(&models.Conversation{}, "id = ?", conversationID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete conversation: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// GetDirectConversation - Lấy direct conversation giữa 2 users
func (r *conversationRepo) GetDirectConversation(userID1, userID2 uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation

	// Tìm conversation có type = 'direct' và có cả 2 users làm participants
	err := r.db.Table("conversations").
		Joins("JOIN participants p1 ON p1.conversation_id = conversations.id").
		Joins("JOIN participants p2 ON p2.conversation_id = conversations.id").
		Where(`conversations.type = 'direct' 
               AND p1.user_id = ? AND p1.left_at IS NULL
               AND p2.user_id = ? AND p2.left_at IS NULL
               AND p1.user_id != p2.user_id`, userID1, userID2).
		First(&conversation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Không tìm thấy direct conversation
		}
		return nil, fmt.Errorf("failed to get direct conversation: %w", err)
	}

	return &conversation, nil
}

// GetConversationWithParticipants - Lấy conversation kèm participants
func (r *conversationRepo) GetConversationWithParticipants(conversationID uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation

	err := r.db.Preload("Participants").
		Preload("Participants.User").
		Where("id = ?", conversationID).
		First(&conversation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation with participants: %w", err)
	}

	return &conversation, nil
}

func (r *conversationRepo) FindConversationByTwoUserID(userID1 uuid.UUID, userID2 uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.Table("conversations").
		Joins("JOIN participants p1 ON p1.conversation_id = conversations.id").
		Joins("JOIN participants p2 ON p2.conversation_id = conversations.id").
		Where("conversations.type = ?", "direct").
		Where("p1.user_id = ? AND p2.user_id = ?", userID1, userID2).
		First(&conversation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find conversation by two user IDs: %w", err)
	}
	return &conversation, nil
}
