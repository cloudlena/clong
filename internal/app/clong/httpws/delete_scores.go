package httpws

import (
	"fmt"
	"net/http"

	"github.com/cloudlena/clong/internal/app/clong"
)

// HandleDeleteScores deletes all scores and resets the scoreboard.
func HandleDeleteScores(scores clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := scores.RemoveAll(ctx)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error removing all scores from store: %w", err))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
