package track

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type TrackRepository interface {
	GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error)
	ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error)
	SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error)
	TrackExists(ctx context.Context, id uuid.UUID) (bool, error)
	ListGenres(ctx context.Context) ([]*track.Genre, error)
}

type HistoryRepository interface {
	CreateListeningHistory(ctx context.Context, h *history.HistoryRecord) error
}

type TrackCache interface {
	Get(ctx context.Context, id uuid.UUID) (*track.Track, error)
	Set(ctx context.Context, t *track.Track) error
	SetNotFound(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GenreCache interface {
	Get(ctx context.Context) ([]track.Genre, error)
	Set(ctx context.Context, genres []track.Genre) error
}

type SearchCache interface {
	Get(ctx context.Context, query string) ([]*track.Track, error)
	Set(ctx context.Context, query string, tracks []*track.Track) error
}

type Service struct {
	tracks  TrackRepository
	history HistoryRepository
	tCache  TrackCache
	gCache  GenreCache
	sCache  SearchCache
}

func NewService(
	tracks TrackRepository,
	history HistoryRepository,
	tCache TrackCache,
	gCache GenreCache,
	sCache SearchCache,
) *Service {
	return &Service{
		tracks:  tracks,
		history: history,
		tCache:  tCache,
		gCache:  gCache,
		sCache:  sCache,
	}
}

func (s *Service) GetTrack(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	t, err := s.tCache.Get(ctx, id)
	if err == nil {
		return t, nil
	}

	if errors.Is(err, track.ErrTrackNotFound) {
		return nil, track.ErrTrackNotFound
	}

	t, err = s.tracks.GetTrackByID(ctx, id)
	if err != nil {
		if errors.Is(err, track.ErrTrackNotFound) {
			_ = s.tCache.SetNotFound(ctx, id)
			return nil, err
		}
		return nil, fmt.Errorf("get track from repo: %w", err)
	}

	_ = s.tCache.Set(ctx, t)

	return t, nil
}

func (s *Service) ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error) {
	return s.tracks.ListTracks(ctx, limit, offset)
}

func (s *Service) SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error) {
	if offset == 0 && limit <= 20 {
		if cached, err := s.sCache.Get(ctx, query); err == nil {
			return cached, nil
		}
	}

	tracks, err := s.tracks.SearchTracks(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search tracks: %w", err)
	}

	if offset == 0 && limit <= 20 {
		_ = s.sCache.Set(ctx, query, tracks)
	}

	return tracks, nil
}

func (s *Service) ListGenres(ctx context.Context) ([]*track.Genre, error) {
	if cached, err := s.gCache.Get(ctx); err == nil {
		res := make([]*track.Genre, len(cached))
		for i := range cached {
			res[i] = &cached[i]
		}
		return res, nil
	}

	genres, err := s.tracks.ListGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("list genres: %w", err)
	}

	cacheData := make([]track.Genre, len(genres))
	for i, g := range genres {
		cacheData[i] = *g
	}
	_ = s.gCache.Set(ctx, cacheData)

	return genres, nil
}

func (s *Service) PlayTrack(ctx context.Context, userID, trackID uuid.UUID) (*track.Track, error) {
	t, err := s.GetTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}

	record, err := history.NewHistoryRecord(userID, trackID)
	if err != nil {
		return nil, err
	}

	if err := s.history.CreateListeningHistory(ctx, record); err != nil {
		return nil, fmt.Errorf("create history record: %w", err)
	}

	return t, nil
}
