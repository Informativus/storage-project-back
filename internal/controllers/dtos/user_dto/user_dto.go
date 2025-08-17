package user_dto

type CreateUserDto struct {
	FldName       string `json:"fldName" validate:"required,fld_max,fld_valid"`
	ConnUserToFld bool   `json:"connUserToFld"`
}
