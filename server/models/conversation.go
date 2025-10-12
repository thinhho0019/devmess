package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string    `gorm:"size:100" json:"name"`
	IsGroup    bool      `gorm:"default:false" json:"is_group"`
	CreatedBy  uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Members  []ConversationMember `gorm:"foreignKey:ConversationID" json:"members"`
	Messages []Message            `gorm:"foreignKey:ConversationID" json:"messages"`
}
