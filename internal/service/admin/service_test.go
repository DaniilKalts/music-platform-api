package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	service "github.com/DaniilKalts/music-platform-api/internal/service/admin"
)

type mockTrackRepo struct{ mock.Mock }

func (m *mockTrackRepo) CreateTrackWithDependencies(ctx context.Context, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error) {
	args := m.Called(ctx, title, artistName, albumName, genreID, durationSeconds, fileURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackRepo) UpdateTrackWithDependencies(ctx context.Context, id uuid.UUID, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error) {
	args := m.Called(ctx, id, title, artistName, albumName, genreID, durationSeconds, fileURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackRepo) SoftDeleteTrack(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) UpdateSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error) {
	args := m.Called(ctx, id, sub)
	return args.Get(0).(*user.User), args.Error(1)
}

type mockTrackCache struct{ mock.Mock }

func (m *mockTrackCache) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTrack(t *testing.T) {
	ctx := context.Background()
	mTrack := new(mockTrackRepo)
	mUser := new(mockUserRepo)
	mCache := new(mockTrackCache)
	s := service.NewService(mTrack, mUser, mCache)

	input := service.CreateTrackInput{
		Title:           "Song",
		ArtistName:      "Artist",
		AlbumName:       "Album",
		GenreID:         uuid.New(),
		DurationSeconds: 180,
		FileURL:         "http://example.com/file.mp3",
	}

	t.Run("Success", func(t *testing.T) {
		mTrack.On("CreateTrackWithDependencies", ctx, input.Title, input.ArtistName, input.AlbumName, input.GenreID, input.DurationSeconds, input.FileURL).
			Return(&track.Track{ID: uuid.New()}, nil)

		res, err := s.CreateTrack(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestDeleteTrack(t *testing.T) {
	ctx := context.Background()
	mTrack := new(mockTrackRepo)
	mUser := new(mockUserRepo)
	mCache := new(mockTrackCache)
	s := service.NewService(mTrack, mUser, mCache)
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mTrack.On("SoftDeleteTrack", ctx, id).Return(nil)
		mCache.On("Delete", ctx, id).Return(nil)

		err := s.DeleteTrack(ctx, id)
		assert.NoError(t, err)
	})
}
