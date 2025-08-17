package errsvc

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidFolder     = errors.New("invalid_folder")
	ErrInvalidFolderName = errors.New("invalid_folder_name")
	ErrFolderExist       = errors.New("folder_exist")
	ErrTokenSignFailed   = errors.New("token_sign_failed")
	ErrGenFolderFailed   = errors.New("gen_folder_failed")
	ErrFldDeleteFailed   = errors.New("fld_delete_failed")
	ErrFldNotFound       = errors.New("fld_not_found")
	ErrBadRequest        = errors.New("bad_request")
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
			Message: "Invalid folder",
			Code:    http.StatusBadRequest,
		}
	case ErrInvalidFolderName:
		return HttpError{
			Message: "Folder name is not valid",
			Code:    http.StatusBadRequest,
		}
	case ErrFolderExist:
		return HttpError{
			Message: "Folder with this name alrady exist",
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
	case ErrBadRequest:
		return HttpError{
			Message: "Bad request",
			Code:    http.StatusBadRequest,
		}
	case ErrFldDeleteFailed:
		return HttpError{
			Message: "Failed to delete folder",
			Code:    http.StatusInternalServerError,
		}
	case ErrFldNotFound:
		return HttpError{
			Message: "Folder not found",
			Code:    http.StatusNotFound,
		}
	default:
		return HttpError{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}
}
