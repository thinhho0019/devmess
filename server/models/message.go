package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"conversation_id"`
	SenderID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"sender_id"`
	Content        string         `json:"content"`
	Type           string         `gorm:"type:varchar(20);default:'text';check:type IN ('text','image','file','video','system')" json:"type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Deleted        bool           `gorm:"default:false" json:"deleted"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Sender       User         `gorm:"foreignKey:SenderID" json:"sender"`
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
}
