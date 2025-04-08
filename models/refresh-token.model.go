package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id           uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Token        string    `gorm:"not null" json:"token"`
	RefreshToken string    `gorm:"not null" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	UserId       uuid.UUID `gorm:"not null;index" json:"user_id"` // Foreign key to User model
	User         User      `gorm:"foreignKey:UserId" json:"user"` // RefreshToken belongs to User
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
