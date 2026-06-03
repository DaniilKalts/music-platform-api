package favoriterepo

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type Repository interface {
	AddFavorite(ctx context.Context, f *favorite.Favorite) error
	RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	ListFavoritesByUserID(ctx context.Context, userID uuid.UUID) ([]*track.Track, error)
	CountFavoritesByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	FavoriteExists(ctx context.Context, userID, trackID uuid.UUID) (bool, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) Repository {
	return &repository{q: sqlc.New(db)}
}

func (r *repository) AddFavorite(ctx context.Context, f *favorite.Favorite) error {
	_, err := r.q.AddFavorite(ctx, sqlc.AddFavoriteParams{
		UserID:  f.UserID,
		TrackID: f.TrackID,
	})
	if err != nil {
		if isNoRows(err) {
			return track.ErrTrackNotFound
		}
		return err
	}
	return nil
}

func (r *repository) RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error {
	_, err := r.q.RemoveFavorite(ctx, sqlc.RemoveFavoriteParams{
		UserID:  userID,
		TrackID: trackID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) ListFavoritesByUserID(ctx context.Context, userID uuid.UUID) ([]*track.Track, error) {
	rows, err := r.q.ListFavoritesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tracks := make([]*track.Track, len(rows))
	for i, row := range rows {
		tracks[i] = toDomainTrackFromFavoriteList(row)
	}
	return tracks, nil
}

func (r *repository) CountFavoritesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.q.CountFavoritesByUserID(ctx, userID)
}

func (r *repository) FavoriteExists(ctx context.Context, userID, trackID uuid.UUID) (bool, error) {
	return r.q.FavoriteExists(ctx, sqlc.FavoriteExistsParams{
		UserID:  userID,
		TrackID: trackID,
	})
}

func isNoRows(err error) bool {
	return err.Error() == "no rows in result set"
}
