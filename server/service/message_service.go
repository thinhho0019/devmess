package service

import (
	"project/models"
	"project/repository"
	"project/utils"
	"time"

	"github.com/google/uuid"
)

type MessageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(messageRepo repository.MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}
func (s *MessageService) SendMessageToConversation(userID string, conversationID string, content string) (*models.Message, error) {
	// convert userID and conversationID from string to uuid.UUID
	uid, err := utils.StringToUUID(userID)
	if err != nil {
		return nil, err
	}
	cid, err := utils.StringToUUID(conversationID)
	if err != nil {
		return nil, err
	}
	message := &models.Message{
		ID:             uuid.New(),
		ConversationID: cid,
		SenderID:       &uid,
		Content:        content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	sentMessage, err := s.messageRepo.CreateMessage(message)
	if err != nil {
		return nil, err
	}
	return sentMessage, nil
}

func (s *MessageService) GetAllMessageToConversation(conversationID string,userID string, limit int, before *time.Time) ([]*models.Message, error) {
	cid, err := utils.StringToUUID(conversationID)
	if err != nil {
		return nil, err
	}
	uid, err := utils.StringToUUID(userID)
	if err != nil {
		return nil, err
	}
	messages, err := s.messageRepo.GetMessagesBeforeTime(cid,uid, before,limit)
	if err != nil {
		return nil, err
	}
	return messages, nil

}
