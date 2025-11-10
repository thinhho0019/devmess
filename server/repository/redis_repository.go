package repository

import (
	"encoding/json"
	"fmt"
	"log"
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
	CreateKeyConversationTwoUserID(userID1, userID2 string) string
	SetKeyConversationTwoUserID(userID1, userID2 string, conversation *models.Conversation, ttl time.Duration) error
	GetConversationIDByTwoUserID(userID1, userID2 string) (*models.Conversation, error)
}

type redisRepo struct{}

func NewRedisRepository() RedisRepository {
	return &redisRepo{}
}

func (r *redisRepo) CreateKeyConversationTwoUserID(userID1, userID2 string) string {
	var key string
	if userID1 < userID2 {
		key = fmt.Sprintf("conversation:%s:%s", userID1, userID2)
	} else {
		key = fmt.Sprintf("conversation:%s:%s", userID2, userID1)
	}

	log.Printf("🔑 [Redis] CreateKey: %s + %s → %s", userID1, userID2, key)
	return key
}

func (r *redisRepo) SetKeyConversationTwoUserID(userID1, userID2 string, conversation *models.Conversation, ttl time.Duration) error {
	if ttl == 0 {
		ttl = time.Hour
	}

	key := r.CreateKeyConversationTwoUserID(userID1, userID2)

	// Serialize conversation
	data, err := json.Marshal(conversation)
	if err != nil {
		log.Printf("❌ [Redis.SetConversation] Marshal error: %v", err)
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	log.Printf("💾 [Redis.SetConversation] Key: %s, TTL: %v, ConvoID: %s", key, ttl, conversation.ID)

	if err := database.RDB.Set(database.Ctx, key, data, ttl).Err(); err != nil {
		log.Printf("❌ [Redis.SetConversation] Set error: %v", err)
		return fmt.Errorf("failed to set conversation in redis: %w", err)
	}

	log.Printf("✅ [Redis.SetConversation] Success - Key: %s", key)
	return nil
}

func (r *redisRepo) GetConversationIDByTwoUserID(userID1, userID2 string) (*models.Conversation, error) {
	key := r.CreateKeyConversationTwoUserID(userID1, userID2)

	log.Printf("🔍 [Redis.GetConversation] Fetching key: %s", key)

	val, err := database.RDB.Get(database.Ctx, key).Result()
	if err != nil {
		log.Printf("❌ [Redis.GetConversation] Get error: %v", err)
		return nil, fmt.Errorf("failed to get conversation from redis: %w", err)
	}

	var conversation models.Conversation
	if err := json.Unmarshal([]byte(val), &conversation); err != nil {
		log.Printf("❌ [Redis.GetConversation] Unmarshal error: %v", err)
		return nil, fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	log.Printf("✅ [Redis.GetConversation] Found - ConvoID: %s, Type: %s", conversation.ID, conversation.Type)
	return &conversation, nil
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
