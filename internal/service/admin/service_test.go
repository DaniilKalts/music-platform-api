package admin_test

import (
	"context"
	"io"
	"strings"
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

type mockFileStorage struct{ mock.Mock }

func (m *mockFileStorage) Upload(ctx context.Context, filename string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, filename, reader, size, contentType)
	return args.String(0), args.Error(1)
}
func (m *mockFileStorage) Delete(ctx context.Context, fileURL string) error {
	args := m.Called(ctx, fileURL)
	return args.Error(0)
}

func TestCreateTrack(t *testing.T) {
	ctx := context.Background()
	mTrack := new(mockTrackRepo)
	mUser := new(mockUserRepo)
	mCache := new(mockTrackCache)
	mStorage := new(mockFileStorage)
	s := service.NewService(mTrack, mUser, mCache, mStorage)

	input := service.CreateTrackInput{
		Title:           "Song",
		ArtistName:      "Artist",
		AlbumName:       "Album",
		GenreID:         uuid.New(),
		DurationSeconds: 180,
		File:            strings.NewReader("audio-bytes"),
		FileSize:        11,
		ContentType:     "audio/mpeg",
		Filename:        "song.mp3",
	}

	t.Run("Success", func(t *testing.T) {
		fileURL := "http://example.com/song.mp3"
		// Ключ объекта генерируется сервисом: <uuid>.mp3, а не исходное имя файла.
		mStorage.On("Upload", ctx, mock.MatchedBy(func(key string) bool {
			if !strings.HasSuffix(key, ".mp3") {
				return false
			}
			_, err := uuid.Parse(strings.TrimSuffix(key, ".mp3"))
			return err == nil
		}), input.File, input.FileSize, input.ContentType).
			Return(fileURL, nil)
		mTrack.On("CreateTrackWithDependencies", ctx, input.Title, input.ArtistName, input.AlbumName, input.GenreID, input.DurationSeconds, fileURL).
			Return(&track.Track{ID: uuid.New()}, nil)

		res, err := s.CreateTrack(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mStorage.AssertExpectations(t)
		mTrack.AssertExpectations(t)
	})
}

func TestDeleteTrack(t *testing.T) {
	ctx := context.Background()
	mTrack := new(mockTrackRepo)
	mUser := new(mockUserRepo)
	mCache := new(mockTrackCache)
	mStorage := new(mockFileStorage)
	s := service.NewService(mTrack, mUser, mCache, mStorage)
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mTrack.On("SoftDeleteTrack", ctx, id).Return(nil)
		mCache.On("Delete", ctx, id).Return(nil)

		err := s.DeleteTrack(ctx, id)
		assert.NoError(t, err)
	})
}
