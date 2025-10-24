package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Avatar       string         `gorm:"type:varchar(255)" json:"avatar,omitempty"`
	Password     string         `gorm:"type:varchar(255);not null" json:"-"`
	Provider     string         `gorm:"type:varchar(50);not null;default:'local'" json:"provider"`
	CreatedAt    time.Time      `json:"created_at"`
	Status       string         `gorm:"type:varchar(10);default:'offline'" json:"status"`
	LastSeen     time.Time      `json:"last_seen"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Devices      []Device       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Participants []Participant  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	 
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
