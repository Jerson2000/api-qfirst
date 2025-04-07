package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jerson2000/api-qfirst/enum"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"not null;unique" json:"email"`
	Password  string         `gorm:"not null" json:"password"`
	Age       *int           `gorm:"default:null" json:"age,omitempty"`
	Gender    *string        `gorm:"default:null" json:"gender,omitempty"`
	Address   *string        `gorm:"default:null" json:"address,omitempty"`
	Phone     *string        `gorm:"default:null" json:"phone,omitempty"`
	Role      *enum.Role     `gorm:"default:user;" json:"role,omitempty"`
	Bookings  []Booking      `gorm:"foreignKey:UserId" json:"bookings,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
