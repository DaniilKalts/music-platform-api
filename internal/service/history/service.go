package history

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
)

type HistoryRepository interface {
	ListListeningHistoryByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error)
}

type Service struct {
	repo HistoryRepository
}

func NewService(repo HistoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error) {
	return s.repo.ListListeningHistoryByUserID(ctx, userID, limit, offset)
}
