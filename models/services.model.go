package models

import (
	"time"

	"gorm.io/gorm"
)

type Services struct {
	Id          uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Duration    *string        `gorm:"default:null" json:"duration,omitempty"`
	Bookings    []Booking      `gorm:"foreignKey:ServiceId" json:"bookings,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
