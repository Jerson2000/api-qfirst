package models

import (
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Id        uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Code      int       `gorm:"not null" json:"code"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	Expiry    time.Time `gorm:"not null" json:"expiry"`
	UserId    uuid.UUID `gorm:"not null,index" json:"user_id"`
}
