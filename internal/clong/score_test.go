package clong_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
)

func TestScoreJSONHidesPlayerID(t *testing.T) {
	b, err := json.Marshal(clong.Score{PlayerID: "secret-id", PlayerName: "Alice", FinalScore: 42, Color: "#ff0000"})
	if err != nil {
		t.Fatalf("error encoding score: %v", err)
	}
	if strings.Contains(string(b), "secret-id") {
		t.Errorf("expected player ID to be hidden, got %s", b)
	}
}
