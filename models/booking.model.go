package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jerson2000/api-qfirst/enum"
	"gorm.io/gorm"
)

type Booking struct {
	Id            uint                `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	StartDate     time.Time           `gorm:"not null" json:"start_date"`
	EndDate       time.Time           `gorm:"not null" json:"end_date"`
	TotalPrice    float64             `gorm:"not null" json:"total_price"`
	Status        *enum.BookingStatus `gorm:"default:pending" json:"status"`
	PaymentStatus *enum.PaymentStatus `gorm:"default:unpaid" json:"payment_status"`
	ServiceId     uint                `gorm:"not null,index" json:"service_id"`
	Service       Services            `gorm:"foreignKey:ServiceId" json:"service"`
	UserId        uuid.UUID           `gorm:"not null,index" json:"user_id"`
	User          User                `gorm:"foreignKey:UserId" json:"user"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deleted_at"`
}
