package playlistrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type Repository interface {
	CreatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error)
	ListPlaylistsByUserID(ctx context.Context, userID uuid.UUID) ([]*playlist.Playlist, error)
	GetPlaylistByIDForUser(ctx context.Context, id, userID uuid.UUID) (*playlist.Playlist, error)
	UpdatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error)
	DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error
	CountPlaylistsByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

	AddTrackToPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error
	RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error
	ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) ([]*track.Track, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) CreatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error) {
	row, err := r.q.CreatePlaylist(ctx, sqlc.CreatePlaylistParams{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: toPgText(p.Description),
	})
	if err != nil {
		return nil, err
	}
	return toDomainPlaylistFromCreate(row), nil
}

func (r *repository) ListPlaylistsByUserID(ctx context.Context, userID uuid.UUID) ([]*playlist.Playlist, error) {
	rows, err := r.q.ListPlaylistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	playlists := make([]*playlist.Playlist, len(rows))
	for i, row := range rows {
		playlists[i] = toDomainPlaylistFromList(row)
	}
	return playlists, nil
}

func (r *repository) GetPlaylistByIDForUser(ctx context.Context, id, userID uuid.UUID) (*playlist.Playlist, error) {
	row, err := r.q.GetPlaylistByIDForUser(ctx, sqlc.GetPlaylistByIDForUserParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if isNoRows(err) {
			return nil, playlist.ErrPlaylistNotFound
		}
		return nil, err
	}
	return toDomainPlaylistFromGet(row), nil
}

func (r *repository) UpdatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error) {
	row, err := r.q.UpdatePlaylist(ctx, sqlc.UpdatePlaylistParams{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: toPgText(p.Description),
	})
	if err != nil {
		if isNoRows(err) {
			return nil, playlist.ErrPlaylistNotFound
		}
		return nil, err
	}
	return toDomainPlaylistFromUpdate(row), nil
}

func (r *repository) DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.q.DeletePlaylist(ctx, sqlc.DeletePlaylistParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if isNoRows(err) {
			return playlist.ErrPlaylistNotFound
		}
		return err
	}
	return nil
}

func (r *repository) CountPlaylistsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.q.CountPlaylistsByUserID(ctx, userID)
}

func (r *repository) AddTrackToPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error {
	_, err := r.q.AddTrackToPlaylist(ctx, sqlc.AddTrackToPlaylistParams{
		PlaylistID: playlistID,
		TrackID:    trackID,
		UserID:     userID,
	})
	if err != nil {
		if isNoRows(err) {
			return playlist.ErrPlaylistNotFound
		}
		return err
	}
	return nil
}

func (r *repository) RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error {
	_, err := r.q.RemoveTrackFromPlaylist(ctx, sqlc.RemoveTrackFromPlaylistParams{
		PlaylistID: playlistID,
		TrackID:    trackID,
		UserID:     userID,
	})
	if err != nil {
		if isNoRows(err) {
			return playlist.ErrPlaylistNotFound
		}
		return err
	}
	return nil
}

func (r *repository) ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) ([]*track.Track, error) {
	rows, err := r.q.ListPlaylistTracks(ctx, sqlc.ListPlaylistTracksParams{
		PlaylistID: playlistID,
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	tracks := make([]*track.Track, len(rows))
	for i, row := range rows {
		tracks[i] = toDomainTrackFromPlaylistList(row)
	}
	return tracks, nil
}

func isNoRows(err error) bool {
	return err.Error() == "no rows in result set"
}

func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}
