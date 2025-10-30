package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type          string     `json:"type" gorm:"type:varchar(20);not null;default:'direct'"`
	Name          string     `json:"name" gorm:"type:varchar(255)"`
	Description   string     `json:"description,omitempty" gorm:"type:text"`
	Avatar        string     `json:"avatar,omitempty" gorm:"type:varchar(500)"`
	LastMessageID *uuid.UUID `json:"last_message_id,omitempty" gorm:"type:uuid"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Participants []*Participant `json:"participants,omitempty" gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
	// Messages     []*Message     `json:"messages,omitempty" gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`

	LastMessage *Message `json:"last_message,omitempty" gorm:"foreignKey:LastMessageID;references:ID"`
}
