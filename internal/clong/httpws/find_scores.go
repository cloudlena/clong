package httpws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleFindScores returns all scores as JSON.
func HandleFindScores(scores clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scrs, err := scores.ListAll(r.Context())
		if err != nil {
			handleServerError(w, fmt.Errorf("error finding scores: %w", err))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(scrs); err != nil {
			log.Printf("error encoding JSON: %v\n", err)
		}
	}
}
