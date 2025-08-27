package fld_dto

type DelFld struct {
	Name string `json:"name" validate:"required,fld_max,fld_valid"`
}
