package historyrepo

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type Repository interface {
	CreateListeningHistory(ctx context.Context, h *history.HistoryRecord) error
	ListListeningHistoryByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) Repository {
	return &repository{q: sqlc.New(db)}
}

func (r *repository) CreateListeningHistory(ctx context.Context, h *history.HistoryRecord) error {
	_, err := r.q.CreateListeningHistory(ctx, sqlc.CreateListeningHistoryParams{
		ID:      h.ID,
		UserID:  h.UserID,
		TrackID: h.TrackID,
	})
	if err != nil {
		if isNoRows(err) {
			return track.ErrTrackNotFound
		}
		return err
	}
	return nil
}

func (r *repository) ListListeningHistoryByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error) {
	rows, err := r.q.ListListeningHistoryByUserID(ctx, sqlc.ListListeningHistoryByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*history.HistoryRecord, len(rows))
	for i, row := range rows {
		records[i] = toDomainHistoryRecordFromList(row)
	}
	return records, nil
}

func isNoRows(err error) bool {
	return err.Error() == "no rows in result set"
}
