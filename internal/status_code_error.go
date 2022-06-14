package internal

type StatusCodeError struct {
	Code int64
	error
}

// NewStatusCodeError creates a wrapper for a standard error with a status code.
func NewStatusCodeError(code int64, err error) *StatusCodeError {
	return &StatusCodeError{
		Code:  code,
		error: err,
	}
}

// StatusCode retuns the StatusCodeError status code
func (e *StatusCodeError) StatusCode() int64 {
	return e.Code
}
