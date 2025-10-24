package service

import (
	"project/models"
	"project/repository"
	"project/utils"
)

type ParticipantService struct {
	repoParticipant repository.ParticipantRepository
	redisRepo       repository.RedisRepository
}

func NewParticipantService(repoParticipant repository.ParticipantRepository, redisRepo repository.RedisRepository) *ParticipantService {
	return &ParticipantService{
		repoParticipant: repoParticipant,
		redisRepo:       redisRepo,
	}
}

func (s *ParticipantService) GetParticipantsByConversationID(conversationID string) (*[]models.Participant, error) {
	// check redis cache first
	if parPtr, err := s.redisRepo.GetParticipantByToken(conversationID); err == nil {
		if parPtr != nil {
			return parPtr, nil
		}
	}
	conversationIDUUID, err := utils.StringToUUID(conversationID)
	if err != nil {
		return nil, err
	}
	participants, err := s.repoParticipant.GetParticipantsByConversationID(conversationIDUUID)
	if err != nil {
		return nil, err
	}
	if participants != nil {
		// set to redis cache
		s.redisRepo.SetTokenParticipant(conversationID, participants, 0)
	}
	return participants, nil
}
