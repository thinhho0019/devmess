package repository

import (

	"encoding/json"
	"time"
	"project/database"
	"project/models"
)

type RedisRepository interface {
	SetToken(token string, user *models.User, ttl time.Duration) error
	GetUserByToken(token string) (*models.User, error)
	DeleteToken(token string) error
}

type redisRepo struct{}

func NewRedisRepository() RedisRepository {
	return &redisRepo{}
}

func (r *redisRepo) SetToken(token string, user *models.User, ttl time.Duration) error {
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

func (r *redisRepo) DeleteToken(token string) error {
	return database.RDB.Del(database.Ctx, "token:"+token).Err()
}
