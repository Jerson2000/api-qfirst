package models

import (
	"time"

	"gorm.io/gorm"
)

// Services represents a service offered by the platform
// @Description Represents a service with name, description, price, and duration.
// @Param id path int true "Service ID"
// @Success 200 {object} Services
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
type Services struct {
	// ID of the service, automatically incremented
	// @example 1
	Id uint `gorm:"primaryKey;autoIncrement;not null" json:"id"`

	// Name of the service
	// @example "Web Development"
	Name string `gorm:"not null" json:"name"`

	// Description of the service
	// @example "Comprehensive website design and development service."
	Description string `gorm:"not null" json:"description"`

	// Price of the service
	// @example 199.99
	Price float64 `gorm:"not null" json:"price"`

	// Duration of the service, optional
	// @example "3 hours"
	Duration *string `gorm:"default:null" json:"duration,omitempty"`

	// List of bookings related to this service
	// @example [{ "id": 1, "user_id": 1, "status": "confirmed" }]
	Bookings []Booking `gorm:"foreignKey:ServiceId" json:"bookings,omitempty"`

	// Timestamp when the service was created
	// @example "2025-04-10T14:22:00Z"
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the service was last updated
	// @example "2025-04-10T14:22:00Z"
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
