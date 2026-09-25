package httpws

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleServerError(t *testing.T) {
	w := httptest.NewRecorder()
	handleServerError(w, errors.New("something broke"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
