package models

import (
	"time"

	"github.com/google/uuid"
)

type ConversationMember struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index" json:"conversation_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Role           string    `gorm:"type:varchar(10);default:'member';check:role IN ('admin','member')" json:"role"`
	JoinedAt       time.Time `json:"joined_at"`

	User         User         `gorm:"foreignKey:UserID" json:"user"`
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
}
