package main

import (
	"log"
	"net/http"
)

// Common error messages within the app
const (
	ErrUpgradingConnection = "error upgrading connection"
	ErrReadingMessage      = "error reading message"
)

// handleHTTPError handles HTTP errors
func handleHTTPError(w http.ResponseWriter, statusCode int, err error, msg string) {
	extMsg := http.StatusText(statusCode)
	if msg != "" {
		extMsg = extMsg + ": " + msg
	}
	http.Error(w, extMsg, statusCode)

	logMsg := extMsg
	if err != nil && err.Error() != msg {
		logMsg = logMsg + ": " + err.Error()
	}
	log.Println(logMsg)
}
