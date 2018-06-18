package clong

import (
	"log"
	"net/http"
)

// unauthorizedError occurs when a user does something they are not authorized for.
type unauthorizedError struct {
	msg string
}

// Error returns the error string.
func (e *unauthorizedError) Error() string {
	return e.msg
}

// Errors commonly used throughout the application.
var (
	errUserIDMissing   = &unauthorizedError{msg: "user ID missing"}
	errUserNameMissing = &unauthorizedError{msg: "user name missing"}
)

// handleHTTPError handles HTTP errors.
func handleHTTPError(w http.ResponseWriter, err error) {
	var code int
	switch err.(type) {
	case *unauthorizedError:
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
