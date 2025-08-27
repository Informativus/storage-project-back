package errsvc

import (
	"fmt"
	"runtime"
)

var (
	UsrErr = struct {
		NotFound          *AppError
		BadReq            *AppError
		AlreadyExists     *AppError
		GenerateToken     *AppError
		InconsistentState *AppError
		Internal          *AppError
	}{
		NotFound:          NewAppError("user_not_found", "user not found", 404),
		BadReq:            NewAppError("bad_request", "bad request", 400),
		AlreadyExists:     NewAppError("user_already_exists", "user already exists", 400),
		InconsistentState: NewAppError("inconsistent_state", "inconsistent state", 500),
		GenerateToken:     NewAppError("token_generate_failed", "internal server error", 500),
		Internal:          NewAppError("internal", "internal server error", 500),
	}

	FldErr = struct {
		NotFound       *AppError
		BadReq         *AppError
		CantDelMainFld *AppError
		CreateFailed   *AppError
		DelFailed      *AppError
		AlreadyExists  *AppError
	}{
		NotFound:       NewAppError("folder_not_found", "folder not found", 404),
		BadReq:         NewAppError("bad_request", "bad request", 400),
		CreateFailed:   NewAppError("folder_create_failed", "folder create failed", 500),
		CantDelMainFld: NewAppError("cannot_delete_main_folder", "cannot delete main folder", 400),
		DelFailed:      NewAppError("folder_delete_failed", "folder delete failed", 500),
		AlreadyExists:  NewAppError("folder_already_exists", "folder already exists", 400),
	}
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key"`
	trace   string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Trace() string {
	return e.trace
}

func NewAppError(key, msg string, code int) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Key:     key,
		trace:   captureStack(3),
	}
}

func captureStack(skip int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)

	frames := runtime.CallersFrames(pc[:n])
	var stack string

	for {
		frame, more := frames.Next()
		stack += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return stack
}
