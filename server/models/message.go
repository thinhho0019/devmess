package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ConversationID uuid.UUID      `json:"conversation_id" gorm:"type:uuid;not null;index"`
	SenderID       *uuid.UUID     `json:"sender_id,omitempty" gorm:"type:uuid;index"`
	Content        string         `json:"content" gorm:"type:text"`
	Type           string         `json:"type" gorm:"type:varchar(20);default:'text';check:type IN ('text','image','file','video','system')"`
	Status         string         `json:"status" gorm:"type:varchar(20);default:'sent';check:status IN ('sent','delivered','read')"`
	IsEdited       bool           `json:"is_edited" gorm:"default:false"`
	ReplyToID      *uuid.UUID     `json:"reply_to_id,omitempty" gorm:"type:uuid;index"`
	Deleted        bool           `json:"deleted" gorm:"default:false"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Sender       *User         `json:"sender,omitempty" gorm:"foreignKey:SenderID;references:ID;constraint:OnDelete:SET NULL;"`
	Conversation *Conversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE;"`
}
