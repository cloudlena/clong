package clong

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// HandleFindScores returns all scores as JSON.
func HandleFindScores(db ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve scores from DB
		scrs, err := db.Scores()
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error finding scores"))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(scrs)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error encoding JSON"))
			return
		}
	}
}
