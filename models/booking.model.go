package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jerson2000/api-qfirst/enum"
	"gorm.io/gorm"
)

// Booking represents a booking made by a user for a service
// @Description Represents a booking with associated user, service, start date, end date, total price, status, and payment status.
// @Param id path int true "Booking ID"
// @Success 200 {object} Booking
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
type Booking struct {
	// ID of the booking, automatically incremented
	// @example 1
	Id uint `gorm:"primaryKey;autoIncrement;not null" json:"id"`

	// The start date of the booking
	// @example "2025-04-15T09:00:00Z"
	StartDate time.Time `gorm:"not null" json:"start_date"`

	// The end date of the booking
	// @example "2025-04-15T17:00:00Z"
	EndDate time.Time `gorm:"not null" json:"end_date"`

	// Total price of the booking
	// @example 250.75
	TotalPrice float64 `gorm:"not null" json:"total_price"`

	// Status of the booking (e.g., pending, confirmed, etc.)
	// @example "pending"
	Status *enum.BookingStatus `gorm:"default:pending" json:"status"`

	// Payment status of the booking (e.g., unpaid, paid, etc.)
	// @example "unpaid"
	PaymentStatus *enum.PaymentStatus `gorm:"default:unpaid" json:"payment_status"`

	// ID of the associated service for the booking
	// @example 1
	ServiceId uint `gorm:"not null,index" json:"service_id"`

	// Service associated with the booking
	// @example { "id": 1, "name": "Web Development", "price": 199.99 }
	Service Services `gorm:"foreignKey:ServiceId" json:"service" swaggerignore:"true"`

	// ID of the user making the booking
	// @example "123e4567-e89b-12d3-a456-426614174000"
	UserId uuid.UUID `gorm:"not null,index" json:"user_id"`

	// User who made the booking
	// @example { "id": "123e4567-e89b-12d3-a456-426614174000", "name": "John Doe", "email": "johndoe@example.com" }
	User User `gorm:"foreignKey:UserId" json:"user" swaggerignore:"true"`

	// Timestamp when the booking was created
	// @example "2025-04-10T14:22:00Z"
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the booking was last updated
	// @example "2025-04-10T14:22:00Z"
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
