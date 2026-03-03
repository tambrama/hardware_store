package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Session struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
	AppID     uuid.UUID
	HashToken string
	ExpiresAt time.Time
}

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	AppID  uuid.UUID `json:"add_id"`
	jwt.RegisteredClaims
}
