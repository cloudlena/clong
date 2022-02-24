package httpws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudlena/clong/internal/app/clong"
)

// HandleFindScores returns all scores as JSON.
func HandleFindScores(scores clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		scrs, err := scores.ListAll(ctx)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error finding scores: %w", err))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(scrs)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error encoding JSON: %w", err))
			return
		}
	}
}
