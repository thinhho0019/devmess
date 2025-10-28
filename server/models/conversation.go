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

	LastMessage  *Message       `json:"last_message,omitempty" gorm:"-"` // Changed: remove foreignKey
	Participants []*Participant `json:"participants,omitempty" gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
	Messages     []*Message     `json:"messages,omitempty" gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// AfterFind hook để load LastMessage manually
func (c *Conversation) AfterFind(tx *gorm.DB) error {
	if c.LastMessageID != nil {
		var msg Message
		if err := tx.First(&msg, "id = ?", c.LastMessageID).Error; err == nil {
			c.LastMessage = &msg
		}
	}
	return nil
}
