package models

import (
	"gorm.io/gorm"
)

type User struct {
	Name     string  `gorm:"not null"`
	Email    string  `gorm:"not null;unique"`
	Password string  `gorm:"not null"`
	Age      *int    `gorm:"default:null"`
	Gender   *string `gorm:"default:null"`
	Address  *string `gorm:"default:null"`
	gorm.Model
}
