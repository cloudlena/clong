package httpws

import (
	"fmt"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleDeleteScores deletes all scores and resets the scoreboard.
func HandleDeleteScores(scores clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := scores.RemoveAll(r.Context()); err != nil {
			handleServerError(w, fmt.Errorf("error removing all scores from store: %w", err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
