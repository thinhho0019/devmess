package service

import (
	"log"
	"project/models"
	"project/repository"
	"project/utils"
	"time"
)

type ConversationService struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	participantRepo  repository.ParticipantRepository
}

func NewConversationService(conversationRepo repository.ConversationRepository, participantRepo repository.ParticipantRepository, messageRepo repository.MessageRepository) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
		messageRepo:      messageRepo,
	}
}

func (s *ConversationService) GetUserConversations(userID string, limit int, before *time.Time) ([]*models.Conversation, error) {
	// convert userID from string to uuid.UUID
	uid, err := utils.StringToUUID(userID)
	if err != nil {
		return nil, err
	}
	r, err := s.conversationRepo.GetConversationsByUserID(uid, limit, before)
	log.Println("[Get Conversation for user] ", r)
	log.Println("[Get Conversation for user] ", err)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *ConversationService) CreateDirectConversation(userID1 string, userID2 string) (*models.Conversation, error) {
	uid1, err := utils.StringToUUID(userID1)
	if err != nil {
		return nil, err
	}
	uid2, err := utils.StringToUUID(userID2)
	if err != nil {
		return nil, err
	}
	conversation, err := s.conversationRepo.CreateDirectConversation(uid1, uid2)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (s *ConversationService) GetMessageByConversationID(conversationID string, limit int, before *time.Time) ([]*models.Message, error) {
	cid, err := utils.StringToUUID(conversationID)
	if err != nil {
		return nil, err
	}
	messages, err := s.conversationRepo.GetMessageByConversationID(cid, limit, before)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *ConversationService) SendMessage(conversationID string, senderID string, content string) (*models.Message, error) {
	cid, err := utils.StringToUUID(conversationID)
	if err != nil {
		return nil, err
	}
	sid, err := utils.StringToUUID(senderID)
	if err != nil {
		return nil, err
	}
	mess := &models.Message{
		ConversationID: cid,
		SenderID:       &sid,
		Content:        content,
		Type:           "text",
		Status:         "sent",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	message, err := s.messageRepo.CreateMessage(mess)
	if err != nil {
		return nil, err
	}
	return message, nil
}
