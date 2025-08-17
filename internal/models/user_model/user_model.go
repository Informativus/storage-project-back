package user_model

import "github.com/google/uuid"

const (
	TableName = "users"
	TokenLen  = 255
)

type UserModel struct {
	ID      uuid.UUID `json:"id"`
	Token   string    `json:"token"`
	Blocked bool      `json:"blocked"`
}
