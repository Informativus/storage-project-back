package fld_dto

import "github.com/google/uuid"

type DelFld struct {
	FldID uuid.UUID `json:"fldID" validate:"required"`
}

type CreateFldBody struct {
	Name string `json:"name" validate:"required,fld_max,fld_valid"`
}

type CreateFldDto struct {
	Name     string    `json:"name" validate:"required,fld_max,fld_valid"`
	ParentID uuid.UUID `json:"parentID" validate:"required"`
}

type CreateFldRes struct {
	FldID uuid.UUID `json:"fldID"`
}
