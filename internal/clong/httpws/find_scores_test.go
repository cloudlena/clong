package httpws_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
)

func TestHandleFindScoresReturnsJSON(t *testing.T) {
	store := &mockScoreStore{
		scores: []*clong.Score{
			{ID: "1", Player: clong.User{ID: "u1", Name: "Alice"}, FinalScore: 10, Color: "red"},
			{ID: "2", Player: clong.User{ID: "u2", Name: "Bob"}, FinalScore: 20, Color: "blue"},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/scores", nil)
	w := httptest.NewRecorder()
	httpws.HandleFindScores(store)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("unexpected Content-Type: %s", ct)
	}

	var got []*clong.Score
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(got))
	}
	if got[0].FinalScore != 10 || got[1].FinalScore != 20 {
		t.Errorf("unexpected scores: %+v", got)
	}
}

func TestHandleFindScoresEmptyList(t *testing.T) {
	store := &mockScoreStore{scores: []*clong.Score{}}

	req := httptest.NewRequest(http.MethodGet, "/scores", nil)
	w := httptest.NewRecorder()
	httpws.HandleFindScores(store)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got []*clong.Score
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty list, got %d items", len(got))
	}
}

func TestHandleFindScoresStoreError(t *testing.T) {
	store := &mockScoreStore{listErr: errors.New("db unavailable")}

	req := httptest.NewRequest(http.MethodGet, "/scores", nil)
	w := httptest.NewRecorder()
	httpws.HandleFindScores(store)(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
