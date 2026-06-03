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

func (m *mockTrackRepo) CreateTrack(ctx context.Context, t *track.Track) (*track.Track, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackRepo) UpdateTrack(ctx context.Context, t *track.Track) (*track.Track, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackRepo) SoftDeleteTrack(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockTrackRepo) GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*track.Track), args.Error(1)
}
func (m *mockTrackRepo) FindOrCreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*track.Artist), args.Error(1)
}
func (m *mockTrackRepo) FindOrCreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*track.Album), args.Error(1)
}
func (m *mockTrackRepo) GetGenreByID(ctx context.Context, id uuid.UUID) (*track.Genre, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*track.Genre), args.Error(1)
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
		artist := &track.Artist{ID: uuid.New(), Name: input.ArtistName}
		album := &track.Album{ID: uuid.New(), Name: input.AlbumName}
		genre := &track.Genre{ID: input.GenreID, Name: "Rock"}

		mTrack.On("FindOrCreateArtist", ctx, mock.Anything).Return(artist, nil)
		mTrack.On("FindOrCreateAlbum", ctx, mock.Anything).Return(album, nil)
		mTrack.On("GetGenreByID", ctx, input.GenreID).Return(genre, nil)
		mTrack.On("CreateTrack", ctx, mock.Anything).Return(&track.Track{ID: uuid.New()}, nil)

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
