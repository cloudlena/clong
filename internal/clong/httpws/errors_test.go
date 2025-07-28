package httpws

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("not allowed")
	if err.Error() != "not allowed" {
		t.Errorf("expected 'not allowed', got %q", err.Error())
	}
}

func TestHandleHTTPErrorUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	handleHTTPError(w, NewUnauthorizedError("forbidden"))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleHTTPErrorInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	handleHTTPError(w, errors.New("something broke"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestHandleHTTPErrorWrappedUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	wrapped := errors.Join(NewUnauthorizedError("no access"), errors.New("context"))
	handleHTTPError(w, wrapped)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrapped UnauthorizedError, got %d", w.Code)
	}
}
