package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Device struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"` // 1 user nhiều device
	Type      string    `json:"type"`                                    // Android, Web, iOS
	Name      string    `json:"name"`                                    // Chrome, Safari, App name
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Token     Token     `gorm:"foreignKey:DeviceID;constraint:OnDelete:CASCADE" json:"token"`
}

type Token struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DeviceID     uuid.UUID      `gorm:"type:uuid;uniqueIndex" json:"device_id"` // 1 token chỉ cho 1 device
	AccessToken  string         `gorm:"uniqueIndex;not null" json:"access_token"`
	RefreshToken string         `gorm:"uniqueIndex;not null" json:"refresh_token"`
	ExpiresAt    int64          `gorm:"not null" json:"expires_at"`
	TokenType    string         `gorm:"type:varchar(50);not null;default:'Bearer'" json:"token_type"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Device) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}
