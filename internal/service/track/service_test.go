package track_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	service "github.com/DaniilKalts/music-platform-api/internal/service/track"
)

type mockTrackRepo struct{ mock.Mock }

func (m *mockTrackRepo) GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackRepo) ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*track.Track), args.Error(1)
}
func (m *mockTrackRepo) SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]*track.Track), args.Error(1)
}
func (m *mockTrackRepo) TrackExists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
func (m *mockTrackRepo) ListGenres(ctx context.Context) ([]*track.Genre, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*track.Genre), args.Error(1)
}

type mockHistoryRepo struct{ mock.Mock }

func (m *mockHistoryRepo) CreateListeningHistory(ctx context.Context, h interface{}) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

type mockTrackCache struct{ mock.Mock }

func (m *mockTrackCache) Get(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackCache) Set(ctx context.Context, t *track.Track) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *mockTrackCache) SetNotFound(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockGenreCache struct{ mock.Mock }

func (m *mockGenreCache) Get(ctx context.Context) ([]track.Genre, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]track.Genre), args.Error(1)
}
func (m *mockGenreCache) Set(ctx context.Context, genres []track.Genre) error {
	args := m.Called(ctx, genres)
	return args.Error(0)
}

type mockSearchCache struct{ mock.Mock }

func (m *mockSearchCache) Get(ctx context.Context, query string) ([]*track.Track, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*track.Track), args.Error(1)
}
func (m *mockSearchCache) Set(ctx context.Context, query string, tracks []*track.Track) error {
	args := m.Called(ctx, query, tracks)
	return args.Error(0)
}

func TestGetTrack(t *testing.T) {
	ctx := context.Background()
	trID := uuid.New()
	mockT := &track.Track{ID: trID, Title: "Test"}

	t.Run("CacheHit", func(t *testing.T) {
		mCache := new(mockTrackCache)
		mRepo := new(mockTrackRepo)
		s := service.NewService(mRepo, nil, mCache, nil, nil)

		mCache.On("Get", ctx, trID).Return(mockT, nil)

		res, err := s.GetTrack(ctx, trID)
		assert.NoError(t, err)
		assert.Equal(t, mockT, res)
		mCache.AssertExpectations(t)
	})

	t.Run("CacheMiss_RepoSuccess", func(t *testing.T) {
		mCache := new(mockTrackCache)
		mRepo := new(mockTrackRepo)
		s := service.NewService(mRepo, nil, mCache, nil, nil)

		mCache.On("Get", ctx, trID).Return(nil, errors.New("miss"))
		mRepo.On("GetTrackByID", ctx, trID).Return(mockT, nil)
		mCache.On("Set", ctx, mockT).Return(nil)

		res, err := s.GetTrack(ctx, trID)
		assert.NoError(t, err)
		assert.Equal(t, mockT, res)
		mRepo.AssertExpectations(t)
		mCache.AssertExpectations(t)
	})

	t.Run("NegativeCaching", func(t *testing.T) {
		mCache := new(mockTrackCache)
		mRepo := new(mockTrackRepo)
		s := service.NewService(mRepo, nil, mCache, nil, nil)

		mCache.On("Get", ctx, trID).Return(nil, errors.New("miss"))
		mRepo.On("GetTrackByID", ctx, trID).Return(nil, track.ErrTrackNotFound)
		mCache.On("SetNotFound", ctx, trID).Return(nil)

		res, err := s.GetTrack(ctx, trID)
		assert.ErrorIs(t, err, track.ErrTrackNotFound)
		assert.Nil(t, res)
		mRepo.AssertExpectations(t)
		mCache.AssertExpectations(t)
	})
}

func TestPlayTrack(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	trID := uuid.New()
	mockT := &track.Track{ID: trID, Title: "Test"}

	mCache := new(mockTrackCache)
	mRepo := new(mockTrackRepo)
	mHist := new(mockHistoryRepo)
	s := service.NewService(mRepo, mHist, mCache, nil, nil)

	t.Run("Success", func(t *testing.T) {
		mCache.On("Get", ctx, trID).Return(mockT, nil)
		mHist.On("CreateListeningHistory", ctx, mock.Anything).Return(nil)

		res, err := s.PlayTrack(ctx, uID, trID)
		assert.NoError(t, err)
		assert.Equal(t, mockT, res)
		mHist.AssertExpectations(t)
	})
}
