package httpws

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
)

func TestHandleDeleteScoresReturnsNoContent(t *testing.T) {
	store := &mockScoreStore{}

	req := httptest.NewRequest(http.MethodDelete, "/api/scores", nil)
	w := httptest.NewRecorder()
	HandleDeleteScores(clong.NewService(store))(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestHandleDeleteScoresStoreError(t *testing.T) {
	store := &mockScoreStore{rmErr: errors.New("db unavailable")}

	req := httptest.NewRequest(http.MethodDelete, "/api/scores", nil)
	w := httptest.NewRecorder()
	HandleDeleteScores(clong.NewService(store))(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
