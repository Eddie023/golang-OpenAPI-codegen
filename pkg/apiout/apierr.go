package apiout

import "errors"

// APIError is used to pass an error during the request through the
// application with web specific context.
type APIError struct {
	Err    error
	Status int
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &APIError{err, status}
}

func (e *APIError) Error() string {
	return e.Err.Error()
}

func (e *APIError) GetHttpStatus() int {
	return e.Status
}

func IsApiError(err error) bool {
	var be *APIError

	return errors.As(err, &be)
}
