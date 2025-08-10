package errsvc

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidFolder   = errors.New("invalid_folder")
	ErrTokenSignFailed = errors.New("token_sign_failed")
	ErrGenFolderFailed = errors.New("gen_folder_failed")
)

type HttpError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ErrorService struct{}

func NewErrorService() *ErrorService {
	return &ErrorService{}
}

func (e *ErrorService) MapError(err error) HttpError {
	switch err {
	case ErrInvalidFolder:
		return HttpError{
			Message: "Folder name is not valid",
			Code:    http.StatusBadRequest,
		}
	case ErrTokenSignFailed:
		return HttpError{
			Message: "Failed to sign JWT token",
			Code:    http.StatusInternalServerError,
		}
	case ErrGenFolderFailed:
		return HttpError{
			Message: "Failed to generate folder",
			Code:    http.StatusInternalServerError,
		}
	default:
		return HttpError{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}
}
