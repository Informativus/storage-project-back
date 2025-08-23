package user_model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TableName = "users"
)

type UserModel struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Blocked   bool      `json:"blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
