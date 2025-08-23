package user_dto

import "github.com/google/uuid"

type CreateUserDto struct {
	UrsName       string `json:"usr_name" validate:"required,fld_max,fld_valid"`
	ConnUserToFld *bool  `json:"conn_user_to_fld"`
}

type DeleteUserDto struct {
	UrsID uuid.UUID `json:"id"`
}
