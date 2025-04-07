package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	jwt.RegisteredClaims
}
