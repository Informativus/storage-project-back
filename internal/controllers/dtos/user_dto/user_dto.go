package user_dto

type CreateUserDto struct {
	UrsName       string `json:"usr_name" validate:"required,fld_max,fld_valid"`
	ConnUserToFld *bool  `json:"conn_user_to_fld"`
}
