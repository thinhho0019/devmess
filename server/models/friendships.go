package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Friendship represents a friend relation / request between two users.
// Status can be: "pending", "accepted", "blocked", "rejected".
type Friendship struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_friend" json:"user_id"`
	FriendID    uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_friend" json:"friend_id"`
	Status      string         `gorm:"type:varchar(20);default:'pending';check:status IN ('no_friend','pending','friend','blocked')" json:"status"`
	RequestedBy uuid.UUID      `gorm:"type:uuid;not null" json:"requested_by"` // who sent the request
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   User `gorm:"foreignKey:UserID" json:"user"`
	Friend User `gorm:"foreignKey:FriendID" json:"friend"`
}
