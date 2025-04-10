package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jerson2000/api-qfirst/enum"
	"gorm.io/gorm"
)

// User represents the user model in the database
// @Description User represents a user with details such as name, email, password, etc.
// @Param id path string true "User ID"
// @Success 200 {object} User
type User struct {
	// ID of the user, automatically generated
	// @example 123e4567-e89b-12d3-a456-426614174000
	Id uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`

	// Name of the user
	// @example "John Doe"
	Name string `gorm:"not null" json:"name"`

	// Email of the user, unique
	// @example "johndoe@example.com"
	Email string `gorm:"not null;unique" json:"email"`

	// Password of the user
	// @example "password123"
	Password string `gorm:"not null" json:"password"`

	// Age of the user, optional
	// @example 25
	Age *int `gorm:"default:null" json:"age,omitempty"`

	// Gender of the user, optional
	// @example "male"
	Gender *string `gorm:"default:null" json:"gender,omitempty"`

	// Address of the user, optional
	// @example "123 Main St, Springfield, IL"
	Address *string `gorm:"default:null" json:"address,omitempty"`

	// Phone number of the user, optional
	// @example "+1234567890"
	Phone *string `gorm:"default:null" json:"phone,omitempty"`

	// Role of the user, default is 'user'
	// @example "admin"
	Role *enum.Role `gorm:"default:user;" json:"role,omitempty"`

	// Whether the user's email has been verified
	// @example false
	IsVerified bool `gorm:"default:false" json:"is_verified"`

	// List of bookings made by the user
	// @example [{ "id": 1, "date": "2025-04-10", "status": "confirmed" }]
	Bookings []Booking `gorm:"foreignKey:UserId" json:"bookings,omitempty" swaggerignore:"true"`

	// List of devices associated with the user
	// @example [{ "id": 1, "device_name": "iPhone", "device_type": "smartphone" }]
	Devices []Devices `gorm:"foreignKey:UserId" json:"devices,omitempty" swaggerignore:"true"`

	// List of refresh tokens associated with the user
	// @example [{ "id": 1, "token": "abcd1234" }]
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserId" json:"refresh_tokens,omitempty" swaggerignore:"true"`

	// Timestamp when the user was created
	// @example "2025-04-10T14:22:00Z"
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the user was last updated
	// @example "2025-04-10T14:22:00Z"
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
