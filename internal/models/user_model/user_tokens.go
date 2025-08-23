package user_model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TokenTableName = "user_tokens"
	TokenLen       = 512
)

type UserTokensModel struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
