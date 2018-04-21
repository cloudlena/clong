package clong

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// FindScoresHandler returns all scores as JSON.
func FindScoresHandler(db DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve scores from DB
		scrs, err := findScores(db)
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
	})
}

// findScores retrieves all scores from the DB.
func findScores(db DB) ([]Score, error) {
	var scrs []Score
	rows, err := db.Query("SELECT id, playerID, playerName, finalScore, color FROM scores")
	if err != nil {
		return []Score{}, errors.Wrap(err, "error getting scores from DB")
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB rows"))
		}
	}()
	for rows.Next() {
		var s Score
		err = rows.Scan(&s.ID, &s.Player.ID, &s.Player.Name, &s.FinalScore, &s.Color)
		if err != nil {
			return []Score{}, errors.Wrap(err, "error scanning DB rows")
		}
		scrs = append(scrs, s)
	}
	err = rows.Err()
	if err != nil {
		return []Score{}, errors.Wrap(err, "error in DB rows")
	}

	return scrs, nil
}
