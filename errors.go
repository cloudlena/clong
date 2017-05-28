package clong

import (
	"log"
	"net/http"
)

// Error codes commonly used throughout the application.
const (
	errUpgradingConnection = "error upgrading connection"
	errClosingConnection   = "error closing WebSocket connection"
	errReadingMessage      = "error reading message"
)

// UnauthorizedError occurs when a user tries to do something they are not authorized to do.
type UnauthorizedError struct {
	msg string
}

// Error returns the error string.
func (e UnauthorizedError) Error() string {
	return e.msg
}

// Errors commonly used throughout the application.
var (
	ErrUserIDMissing = &UnauthorizedError{"user ID missing"}
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
