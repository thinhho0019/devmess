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
	UpdateUserOnline(user_id string) error
	IsUserOnline(user_id string) (bool, time.Duration, error)
}

type redisRepo struct{}

func NewRedisRepository() RedisRepository {
	return &redisRepo{}
}

func (r *redisRepo) UpdateUserOnline(user_id string) error {
	now := time.Now().Unix()
	// save timetamp to hash
	err := database.RDB.HSet(database.Ctx, "user_last_seen", user_id, now).Err()
	if err != nil {
		return fmt.Errorf("failed to update user online status: %w", err)
	}
	err = database.RDB.SetEx(database.Ctx, "online:"+user_id, 1, 300*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set user online key: %w", err)
	}
	return nil
}
func (r *redisRepo) IsUserOnline(user_id string) (bool, time.Duration, error) {
	ttl, err := database.RDB.TTL(database.Ctx, "online:"+user_id).Result()
	if err != nil {
		return false, 0, fmt.Errorf("failed to check user online status: %w", err)
	}
	lastScreenInt, err := database.RDB.HGet(database.Ctx, "user_last_seen", user_id).Int64()
	if err != nil && err.Error() != "redis: nil" {
		return false, 0, fmt.Errorf("failed to get user last seen: %w", err)
	}
	lastSceen := time.Unix(lastScreenInt, 0)
	offlineDuration := time.Since(lastSceen)
	return ttl > 1, offlineDuration, nil
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
