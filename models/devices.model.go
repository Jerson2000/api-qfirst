package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// for multi device authentication still not implemented
type Devices struct {
	Id             uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	IsMobile       bool           `gorm:"default:false" json:"is_mobile,omitempty"`
	IsBot          bool           `gorm:"default:false" json:"is_bot,omitempty"`
	Mozilla        *string        `gorm:"default:null" json:"mozilla,omitempty"`
	Model          *string        `gorm:"default:null" json:"model,omitempty"`
	Platform       *string        `gorm:"default:null" json:"platform,omitempty"`
	OS             *string        `gorm:"default:null" json:"os,omitempty"`
	EngineName     *string        `gorm:"default:null" json:"engine_name,omitempty"`
	EngineVersion  *string        `gorm:"default:null" json:"engine_version,omitempty"`
	BrowserName    *string        `gorm:"default:null" json:"browser_name,omitempty"`
	BrowserVersion *string        `gorm:"default:null" json:"browser_version,omitempty"`
	IP             *string        `gorm:"default:null" json:"ip,omitempty"`
	UserId         uuid.UUID      `gorm:"not null,index" json:"user_id"`
	User           User           `gorm:"foreignKey:UserId" json:"user"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
