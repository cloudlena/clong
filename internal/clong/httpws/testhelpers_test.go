package httpws_test

import (
	"context"

	"github.com/cloudlena/clong/internal/clong"
)

type mockScoreStore struct {
	scores  []*clong.Score
	listErr error
	addErr  error
	rmErr   error
}

func (m *mockScoreStore) ListAll(_ context.Context) ([]*clong.Score, error) {
	return m.scores, m.listErr
}

func (m *mockScoreStore) Add(_ context.Context, s *clong.Score) error {
	if m.addErr != nil {
		return m.addErr
	}
	m.scores = append(m.scores, s)
	return nil
}

func (m *mockScoreStore) RemoveAll(_ context.Context) error {
	return m.rmErr
}
