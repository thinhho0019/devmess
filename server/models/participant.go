package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Participant struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_conversation_user"`
	ConversationID uuid.UUID      `json:"conversation_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_conversation_user"`
	Role           string         `json:"role" gorm:"type:varchar(20);default:'member'"` // member | admin | owner
	LastReadAt     *time.Time     `json:"last_read_at,omitempty"`
	JoinedAt       time.Time      `json:"joined_at" gorm:"autoCreateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	User           User           `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
