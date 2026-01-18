package apperror

import "net/http"

type HTTPError struct {
	Code       int
	Message    string
	DebugError error
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewInternal(msg string, err error) error {
	return HTTPError{Code: http.StatusInternalServerError, Message: msg, DebugError: err}
}

func NewBadReq(msg string) error {
	return HTTPError{Code: http.StatusBadRequest, Message: msg}
}

func NewUnauthorized(msg string) error {
	return HTTPError{Code: http.StatusUnauthorized, Message: msg}
}

func NewNotFound(msg string) error {
	return HTTPError{Code: http.StatusNotFound, Message: msg}
}
