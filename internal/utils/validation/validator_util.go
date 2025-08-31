package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ivan/storage-project-back/internal/models/file_model"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	registerUserValidations(Validate)
	registerFolderValidations(Validate)
	registerFileValidations(Validate)
}

func registerUserValidations(v *validator.Validate) {
	_ = v.RegisterValidation("token_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= user_model.TokenLen
	})
}

func registerFolderValidations(v *validator.Validate) {
	_ = v.RegisterValidation("fld_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= folder_model.FolderNameLen
	})

	_ = v.RegisterValidation("fld_valid", func(fl validator.FieldLevel) bool {
		return !strings.ContainsAny(fl.Field().String(), `/\:*?"<>|`)
	})
}

func registerFileValidations(v *validator.Validate) {
	_ = v.RegisterValidation("file_name_valid", validateFileName)
	_ = v.RegisterValidation("file_name_max", validateFileNameMax)
	_ = v.RegisterValidation("file_name_min", validateFileNameMin)
}

func validateFileName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if strings.TrimSpace(name) == "" {
		return false
	}

	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return false
	}

	if name == "." || name == ".." {
		return false
	}

	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return false
	}

	return true
}

func validateFileNameMax(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) <= file_model.MaxFileNameLen
}

func validateFileNameMin(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) >= file_model.MinFileNameLen
}
