package errors

import "errors"

var (
	CSRFTokenNotFoundError = errors.New("CSRF token not found")
)
