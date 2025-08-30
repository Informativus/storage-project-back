package errsvc

import (
	"fmt"
	"runtime"
	"strings"
)

type AppErrorTemplate struct {
	Key     string
	Message string
	Code    int
}

func (t AppErrorTemplate) New() *AppError {
	return &AppError{
		Key:     t.Key,
		Message: t.Message,
		Code:    t.Code,
		trace:   captureStack(3),
	}
}

var (
	UsrErr = struct {
		NotFound          AppErrorTemplate
		BadReq            AppErrorTemplate
		AlreadyExists     AppErrorTemplate
		GenerateToken     AppErrorTemplate
		Forbidden         AppErrorTemplate
		InconsistentState AppErrorTemplate
		Internal          AppErrorTemplate
	}{
		NotFound:          AppErrorTemplate{"user_not_found", "user not found", 404},
		BadReq:            AppErrorTemplate{"bad_request", "bad request", 400},
		AlreadyExists:     AppErrorTemplate{"user_already_exists", "user already exists", 400},
		GenerateToken:     AppErrorTemplate{"token_generate_failed", "internal server error", 500},
		Forbidden:         AppErrorTemplate{"forbidden", "forbidden", 403},
		InconsistentState: AppErrorTemplate{"inconsistent_state", "inconsistent state", 500},
		Internal:          AppErrorTemplate{"internal", "internal server error", 500},
	}

	FldErr = struct {
		NotFound       AppErrorTemplate
		BadReq         AppErrorTemplate
		CantDelMainFld AppErrorTemplate
		CreateFailed   AppErrorTemplate
		DelFailed      AppErrorTemplate
		AlreadyExists  AppErrorTemplate
		Internal       AppErrorTemplate
	}{
		NotFound:       AppErrorTemplate{"folder_not_found", "folder not found", 404},
		BadReq:         AppErrorTemplate{"bad_request", "bad request", 400},
		CreateFailed:   AppErrorTemplate{"folder_create_failed", "folder create failed", 500},
		CantDelMainFld: AppErrorTemplate{"cannot_delete_main_folder", "cannot delete main folder", 400},
		DelFailed:      AppErrorTemplate{"folder_delete_failed", "folder delete failed", 500},
		AlreadyExists:  AppErrorTemplate{"folder_already_exists", "folder already exists", 400},
		Internal:       AppErrorTemplate{"internal", "internal server error", 500},
	}

	SecurityErr = struct {
		NotFound      AppErrorTemplate
		AlreadyExists AppErrorTemplate
		Internal      AppErrorTemplate
	}{
		NotFound:      AppErrorTemplate{"folder_not_found", "folder not found", 404},
		AlreadyExists: AppErrorTemplate{"folder_already_exists", "folder already exists", 400},
		Internal:      AppErrorTemplate{"internal", "internal server error", 500},
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
	pc := make([]uintptr, 50)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])

	var b strings.Builder
	for {
		frame, more := frames.Next()
		b.WriteString(frame.Function)
		b.WriteString("\n\t")
		b.WriteString(frame.File)
		b.WriteString(":")
		b.WriteString(fmt.Sprint(frame.Line))
		b.WriteString("\n")
		if !more {
			break
		}
	}
	return b.String()
}
