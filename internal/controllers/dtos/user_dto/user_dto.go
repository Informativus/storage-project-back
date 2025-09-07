package user_dto

type CreateUserDto struct {
	UsrName string `json:"usrName" validate:"required,fld_max,fld_valid"`
}

type CreateUserResponse struct {
	Token string `json:"token"`
}

type GenTokenReq struct {
	UsrName string `json:"usrName" validate:"required,fld_max,fld_valid"`
}

type GenTokenRes struct {
	Token string `json:"token"`
}

type BlockUserReq struct {
	UsrName string `json:"usrName" validate:"required,fld_max,fld_valid"`
	Blocked bool   `json:"blocked"`
}
