package clong

import (
	"log"
	"net/http"
)

// Error codes commonly used throughout the application.
const (
	errUpgradingConnection = "error upgrading connection"
	errReadingJSON         = "error reading JSON"
)

// UnauthorizedError occurs when a user does something they are not authorized for.
type UnauthorizedError struct {
	msg string
}

// Error returns the error string.
func (e UnauthorizedError) Error() string {
	return e.msg
}

// Errors commonly used throughout the application.
var (
	ErrUserIDMissing   = &UnauthorizedError{"user ID missing"}
	ErrUserNameMissing = &UnauthorizedError{"user name missing"}
)

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
