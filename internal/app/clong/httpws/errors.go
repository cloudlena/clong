package httpws

import (
	"log"
	"net/http"
)

// UnauthorizedError occurs when a user does something they are not authorized for.
type UnauthorizedError struct {
	msg string
}

// NewUnauthorizedError creates a new unauthorized-error.
func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{msg}
}

// Error returns the error string.
func (e *UnauthorizedError) Error() string {
	return e.msg
}

// handleHTTPError handles HTTP errors.
func handleHTTPError(w http.ResponseWriter, err error) {
	var code int
	switch err.(type) {
	case *UnauthorizedError:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}

	http.Error(w, http.StatusText(code), code)

	// Log if server error
	if code >= 500 {
		log.Println(err)
	}
}
