package errors

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrNotFound     = &Error{Code: "NOT_FOUND", Message: "Resource not found", Status: 404}
	ErrUnauthorized = &Error{Code: "UNAUTHORIZED", Message: "Unauthorized", Status: 401}
	ErrForbidden    = &Error{Code: "FORBIDDEN", Message: "Forbidden", Status: 403}
	ErrBadRequest   = &Error{Code: "BAD_REQUEST", Message: "Bad request", Status: 400}
	ErrInternal     = &Error{Code: "INTERNAL_ERROR", Message: "Internal server error", Status: 500}
)
