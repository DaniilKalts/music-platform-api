package track

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"

	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type Repository interface {
	GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error)
	ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error)
	SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error)
	ListGenres(ctx context.Context) ([]*track.Genre, error)
	TrackExists(ctx context.Context, id uuid.UUID) (bool, error)
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

type PopularCache interface {
	Get(ctx context.Context) ([]*track.Track, error)
	Set(ctx context.Context, tracks []*track.Track) error
}

type Service struct {
	repo         Repository
	historyRepo  HistoryRepository
	trackCache   TrackCache
	genreCache   GenreCache
	searchCache  SearchCache
	popularCache PopularCache

	sg singleflight.Group
}

func NewService(
	repo Repository,
	historyRepo HistoryRepository,
	trackCache TrackCache,
	genreCache GenreCache,
	searchCache SearchCache,
	popularCache PopularCache,
) *Service {
	return &Service{
		repo:         repo,
		historyRepo:  historyRepo,
		trackCache:   trackCache,
		genreCache:   genreCache,
		searchCache:  searchCache,
		popularCache: popularCache,
	}
}

var errTrackNotFoundCached = errors.New("track not found (cached)")

func (s *Service) GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	// 1. Try Cache
	t, err := s.trackCache.Get(ctx, id)
	if err == nil {
		return t, nil
	}
	if errors.Is(err, errTrackNotFoundCached) {
		return nil, track.ErrTrackNotFound
	}

	// 2. Use singleflight to prevent thundering herd
	res, err, _ := s.sg.Do(fmt.Sprintf("get_track:%s", id), func() (interface{}, error) {
		// Double check cache inside singleflight
		if t, err := s.trackCache.Get(ctx, id); err == nil {
			return t, nil
		}

		t, err := s.repo.GetTrackByID(ctx, id)
		if err != nil {
			if errors.Is(err, track.ErrTrackNotFound) {
				_ = s.trackCache.SetNotFound(ctx, id)
			}
			return nil, err
		}

		_ = s.trackCache.Set(ctx, t)
		return t, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*track.Track), nil
}

func (s *Service) ListTracks(ctx context.Context, page, limit int32) ([]*track.Track, error) {
	offset := (page - 1) * limit
	return s.repo.ListTracks(ctx, limit, offset)
}

func (s *Service) SearchTracks(ctx context.Context, query string, page, limit int32) ([]*track.Track, error) {
	cacheKey := fmt.Sprintf("%s:%d:%d", query, page, limit)

	// 1. Try Cache
	if tracks, err := s.searchCache.Get(ctx, cacheKey); err == nil {
		return tracks, nil
	}

	// 2. Singleflight for search
	res, err, _ := s.sg.Do(fmt.Sprintf("search:%s", cacheKey), func() (interface{}, error) {
		offset := (page - 1) * limit
		tracks, err := s.repo.SearchTracks(ctx, query, limit, offset)
		if err != nil {
			return nil, err
		}

		_ = s.searchCache.Set(ctx, cacheKey, tracks)
		return tracks, nil
	})

	if err != nil {
		return nil, err
	}

	return res.([]*track.Track), nil
}

func (s *Service) ListGenres(ctx context.Context) ([]track.Genre, error) {
	// 1. Try Cache
	if genres, err := s.genreCache.Get(ctx); err == nil {
		return genres, nil
	}

	// 2. Singleflight
	res, err, _ := s.sg.Do("list_genres", func() (interface{}, error) {
		genrePtrs, err := s.repo.ListGenres(ctx)
		if err != nil {
			return nil, err
		}

		genres := make([]track.Genre, len(genrePtrs))
		for i, g := range genrePtrs {
			genres[i] = *g
		}

		_ = s.genreCache.Set(ctx, genres)
		return genres, nil
	})

	if err != nil {
		return nil, err
	}

	return res.([]track.Genre), nil
}

func (s *Service) PlayTrack(ctx context.Context, userID, trackID uuid.UUID) (*track.Track, error) {
	t, err := s.GetTrackByID(ctx, trackID)
	if err != nil {
		return nil, err
	}

	record, err := history.NewHistoryRecord(userID, trackID)
	if err != nil {
		return nil, err
	}

	if err := s.historyRepo.CreateListeningHistory(ctx, record); err != nil {
		return nil, err
	}

	return t, nil
}
