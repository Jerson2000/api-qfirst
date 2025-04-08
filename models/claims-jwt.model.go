package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jerson2000/api-qfirst/enum"
)

type JwtClaims struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role enum.Role `json:"role"`
	jwt.RegisteredClaims
}
