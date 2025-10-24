package repository

import (
	"encoding/json"
	"fmt"
	"project/database"
	"project/models"
	"time"
)

type RedisRepository interface {
	SetToken(token string, user *models.User, ttl time.Duration) error
	GetUserByToken(token string) (*models.User, error)
	DeleteToken(token string) error
	SetTokenParticipant(token string, participants *[]models.Participant, ttl time.Duration) error
	GetParticipantByToken(token string) (*[]models.Participant, error)
}

type redisRepo struct{}

func NewRedisRepository() RedisRepository {
	return &redisRepo{}
}

func (r *redisRepo) SetToken(token string, user *models.User, ttl time.Duration) error {
	if ttl == 0 {
		ttl = time.Hour
	}
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return database.RDB.Set(database.Ctx, "token:"+token, data, ttl).Err()
}

func (r *redisRepo) GetUserByToken(token string) (*models.User, error) {
	val, err := database.RDB.Get(database.Ctx, "token:"+token).Result()
	if err != nil {
		return nil, err
	}
	var user models.User
	json.Unmarshal([]byte(val), &user)
	return &user, nil
}

// set-get participant in redis (token is conversationID)
func (r *redisRepo) SetTokenParticipant(token string, participants *[]models.Participant, ttl time.Duration) error {
	if ttl == 0 {
		ttl = time.Hour
	}
	data, err := json.Marshal(participants)
	if err != nil {
		return err
	}
	return database.RDB.Set(database.Ctx, "token-participants:"+token, data, ttl).Err()
}
func (r *redisRepo) GetParticipantByToken(token string) (*[]models.Participant, error) {
	val, err := database.RDB.Get(database.Ctx, "token-participants:"+token).Result()
	if err != nil {
		return nil, err
	}
	var participants []models.Participant
	if err := json.Unmarshal([]byte(val), &participants); err != nil {
		return nil, fmt.Errorf("failed to unmarshal participants: %w", err)
	}
	return &participants, nil
}

func (r *redisRepo) DeleteToken(token string) error {
	return database.RDB.Del(database.Ctx, "token:"+token).Err()
}
