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

	_ = Validate.RegisterValidation("token_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= user_model.TokenLen
	})

	_ = Validate.RegisterValidation("fld_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= folder_model.FolderNameLen
	})

	_ = Validate.RegisterValidation("fld_valid", func(fl validator.FieldLevel) bool {
		return !strings.ContainsAny(fl.Field().String(), `/\:*?"<>|`)
	})

	_ = Validate.RegisterValidation("file_name_valid", func(fl validator.FieldLevel) bool {
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
	})

	_ = Validate.RegisterValidation("file_name_max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= file_model.MaxFileNameLen
	})

	_ = Validate.RegisterValidation("file_name_min", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= file_model.MinFileNameLen
	})
}
