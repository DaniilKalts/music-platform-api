package trackrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type Repository interface {
	CreateTrack(ctx context.Context, t *track.Track) (*track.Track, error)
	GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error)
	ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error)
	SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error)
	UpdateTrack(ctx context.Context, t *track.Track) (*track.Track, error)
	SoftDeleteTrack(ctx context.Context, id uuid.UUID) error
	TrackExists(ctx context.Context, id uuid.UUID) (bool, error)

	CreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (*track.Artist, error)
	FindOrCreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error)

	CreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error)
	GetAlbumByID(ctx context.Context, id uuid.UUID) (*track.Album, error)
	FindOrCreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error)

	CreateGenre(ctx context.Context, g *track.Genre) (*track.Genre, error)
	ListGenres(ctx context.Context) ([]*track.Genre, error)
	GetGenreByID(ctx context.Context, id uuid.UUID) (*track.Genre, error)

	CreateTrackWithDependencies(ctx context.Context, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error)
	UpdateTrackWithDependencies(ctx context.Context, id uuid.UUID, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error)
}

type repository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
		q:  sqlc.New(db),
	}
}

func (r *repository) CreateTrack(ctx context.Context, t *track.Track) (*track.Track, error) {
	row, err := r.q.CreateTrack(ctx, sqlc.CreateTrackParams{
		ID:              t.ID,
		Title:           t.Title,
		ArtistID:        t.ArtistID,
		AlbumID:         t.AlbumID,
		GenreID:         t.GenreID,
		DurationSeconds: int32(t.DurationSeconds),
		FileUrl:         t.FileURL,
	})
	if err != nil {
		return nil, err
	}
	return toDomainTrackFromCreate(row), nil
}

func (r *repository) GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	row, err := r.q.GetTrackByID(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrTrackNotFound
		}
		return nil, err
	}
	return toDomainTrackFromGet(row), nil
}

func (r *repository) ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error) {
	rows, err := r.q.ListTracks(ctx, sqlc.ListTracksParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	tracks := make([]*track.Track, len(rows))
	for i, row := range rows {
		tracks[i] = toDomainTrackFromList(row)
	}
	return tracks, nil
}

func (r *repository) SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error) {
	rows, err := r.q.SearchTracks(ctx, sqlc.SearchTracksParams{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	tracks := make([]*track.Track, len(rows))
	for i, row := range rows {
		tracks[i] = toDomainTrackFromSearch(row)
	}
	return tracks, nil
}

func (r *repository) UpdateTrack(ctx context.Context, t *track.Track) (*track.Track, error) {
	row, err := r.q.UpdateTrack(ctx, sqlc.UpdateTrackParams{
		ID:              t.ID,
		Title:           t.Title,
		ArtistID:        t.ArtistID,
		AlbumID:         t.AlbumID,
		GenreID:         t.GenreID,
		DurationSeconds: int32(t.DurationSeconds),
		FileUrl:         t.FileURL,
	})
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrTrackNotFound
		}
		return nil, err
	}
	return toDomainTrackFromUpdate(row), nil
}

func (r *repository) SoftDeleteTrack(ctx context.Context, id uuid.UUID) error {
	_, err := r.q.SoftDeleteTrack(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return track.ErrTrackNotFound
		}
		return err
	}
	return nil
}

func (r *repository) TrackExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.q.TrackExists(ctx, id)
}

func (r *repository) CreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error) {
	row, err := r.q.CreateArtist(ctx, sqlc.CreateArtistParams{
		ID:   a.ID,
		Name: a.Name,
	})
	if err != nil {
		return nil, err
	}
	return toDomainArtistFromCreate(row), nil
}

func (r *repository) GetArtistByID(ctx context.Context, id uuid.UUID) (*track.Artist, error) {
	row, err := r.q.GetArtistByID(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrArtistNotFound
		}
		return nil, err
	}
	return toDomainArtistFromGet(row), nil
}

func (r *repository) FindOrCreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error) {
	row, err := r.q.FindOrCreateArtist(ctx, sqlc.FindOrCreateArtistParams{
		ID:   a.ID,
		Name: a.Name,
	})
	if err != nil {
		return nil, err
	}
	return toDomainArtistFromFindOrCreate(row), nil
}

func (r *repository) CreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error) {
	row, err := r.q.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		ID:   a.ID,
		Name: a.Name,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAlbumFromCreate(row), nil
}

func (r *repository) GetAlbumByID(ctx context.Context, id uuid.UUID) (*track.Album, error) {
	row, err := r.q.GetAlbumByID(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrAlbumNotFound
		}
		return nil, err
	}
	return toDomainAlbumFromGet(row), nil
}

func (r *repository) FindOrCreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error) {
	row, err := r.q.FindOrCreateAlbum(ctx, sqlc.FindOrCreateAlbumParams{
		ID:   a.ID,
		Name: a.Name,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAlbumFromFindOrCreate(row), nil
}

func (r *repository) CreateGenre(ctx context.Context, g *track.Genre) (*track.Genre, error) {
	row, err := r.q.CreateGenre(ctx, sqlc.CreateGenreParams{
		ID:   g.ID,
		Name: g.Name,
	})
	if err != nil {
		return nil, err
	}
	return toDomainGenreFromCreate(row), nil
}

func (r *repository) ListGenres(ctx context.Context) ([]*track.Genre, error) {
	rows, err := r.q.ListGenres(ctx)
	if err != nil {
		return nil, err
	}

	genres := make([]*track.Genre, len(rows))
	for i, row := range rows {
		genres[i] = toDomainGenreFromList(row)
	}
	return genres, nil
}

func (r *repository) GetGenreByID(ctx context.Context, id uuid.UUID) (*track.Genre, error) {
	row, err := r.q.GetGenreByID(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrGenreNotFound
		}
		return nil, err
	}
	return toDomainGenreFromGet(row), nil
}

func (r *repository) CreateTrackWithDependencies(ctx context.Context, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	artRow, err := qtx.FindOrCreateArtist(ctx, sqlc.FindOrCreateArtistParams{
		ID:   uuid.New(),
		Name: artistName,
	})
	if err != nil {
		return nil, err
	}

	albRow, err := qtx.FindOrCreateAlbum(ctx, sqlc.FindOrCreateAlbumParams{
		ID:   uuid.New(),
		Name: albumName,
	})
	if err != nil {
		return nil, err
	}

	_, err = qtx.GetGenreByID(ctx, genreID)
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrGenreNotFound
		}
		return nil, err
	}

	trackRow, err := qtx.CreateTrack(ctx, sqlc.CreateTrackParams{
		ID:              uuid.New(),
		Title:           title,
		ArtistID:        artRow.ID,
		AlbumID:         albRow.ID,
		GenreID:         genreID,
		DurationSeconds: int32(durationSeconds),
		FileUrl:         fileURL,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return toDomainTrackFromCreate(trackRow), nil
}

func (r *repository) UpdateTrackWithDependencies(ctx context.Context, id uuid.UUID, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	artRow, err := qtx.FindOrCreateArtist(ctx, sqlc.FindOrCreateArtistParams{
		ID:   uuid.New(),
		Name: artistName,
	})
	if err != nil {
		return nil, err
	}

	albRow, err := qtx.FindOrCreateAlbum(ctx, sqlc.FindOrCreateAlbumParams{
		ID:   uuid.New(),
		Name: albumName,
	})
	if err != nil {
		return nil, err
	}

	_, err = qtx.GetGenreByID(ctx, genreID)
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrGenreNotFound
		}
		return nil, err
	}

	trackRow, err := qtx.UpdateTrack(ctx, sqlc.UpdateTrackParams{
		ID:              id,
		Title:           title,
		ArtistID:        artRow.ID,
		AlbumID:         albRow.ID,
		GenreID:         genreID,
		DurationSeconds: int32(durationSeconds),
		FileUrl:         fileURL,
	})
	if err != nil {
		if isNoRows(err) {
			return nil, track.ErrTrackNotFound
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return toDomainTrackFromUpdate(trackRow), nil
}

func isNoRows(err error) bool {
	return err.Error() == "no rows in result set"
}
