package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	_ = Validate.RegisterValidation("token_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= user_model.TokenLen
	})

	_ = Validate.RegisterValidation("fld_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= folder_model.FolderNameLen
	})

	_ = Validate.RegisterValidation("fld_valid", func(fl validator.FieldLevel) bool {
		return !strings.ContainsAny(fl.Field().String(), `/\:*?"<>|`)
	})
}
