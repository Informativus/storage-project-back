package fld_dto

type DeleteFld struct {
	Name string `json:"name" validate:"required,fld_max,fld_valid"`
}
