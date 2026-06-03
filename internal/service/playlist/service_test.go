package playlist_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	service "github.com/DaniilKalts/music-platform-api/internal/service/playlist"
)

type mockPlaylistRepo struct{ mock.Mock }

func (m *mockPlaylistRepo) CreatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (m *mockPlaylistRepo) ListPlaylistsByUserID(ctx context.Context, userID uuid.UUID) ([]*playlist.Playlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*playlist.Playlist), args.Error(1)
}

func (m *mockPlaylistRepo) GetPlaylistByIDForUser(ctx context.Context, id, userID uuid.UUID) (*playlist.Playlist, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (m *mockPlaylistRepo) UpdatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (m *mockPlaylistRepo) DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockPlaylistRepo) CountPlaylistsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPlaylistRepo) AddTrackToPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error {
	args := m.Called(ctx, playlistID, trackID, userID)
	return args.Error(0)
}

func (m *mockPlaylistRepo) RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error {
	args := m.Called(ctx, playlistID, trackID, userID)
	return args.Error(0)
}

func (m *mockPlaylistRepo) ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) ([]*track.Track, error) {
	args := m.Called(ctx, playlistID, userID)
	return args.Get(0).([]*track.Track), args.Error(1)
}

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func TestCreatePlaylist(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	limit := 2

	t.Run("Success_Premium", func(t *testing.T) {
		mRepo := new(mockPlaylistRepo)
		mUser := new(mockUserRepo)
		s := service.NewService(mRepo, mUser, limit)

		mUser.On("GetByID", ctx, uID).Return(&user.User{ID: uID, Subscription: user.SubscriptionPremium}, nil)
		mRepo.On("CreatePlaylist", ctx, mock.Anything).Return(&playlist.Playlist{ID: uuid.New()}, nil)

		p, err := s.CreatePlaylist(ctx, service.CreateInput{UserID: uID, Name: "My Playlist"})
		assert.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("Success_Free_UnderLimit", func(t *testing.T) {
		mRepo := new(mockPlaylistRepo)
		mUser := new(mockUserRepo)
		s := service.NewService(mRepo, mUser, limit)

		mUser.On("GetByID", ctx, uID).Return(&user.User{ID: uID, Subscription: user.SubscriptionFree}, nil)
		mRepo.On("CountPlaylistsByUserID", ctx, uID).Return(int64(1), nil)
		mRepo.On("CreatePlaylist", ctx, mock.Anything).Return(&playlist.Playlist{ID: uuid.New()}, nil)

		p, err := s.CreatePlaylist(ctx, service.CreateInput{UserID: uID, Name: "My Playlist"})
		assert.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("Error_Free_LimitExceeded", func(t *testing.T) {
		mRepo := new(mockPlaylistRepo)
		mUser := new(mockUserRepo)
		s := service.NewService(mRepo, mUser, limit)

		mUser.On("GetByID", ctx, uID).Return(&user.User{ID: uID, Subscription: user.SubscriptionFree}, nil)
		mRepo.On("CountPlaylistsByUserID", ctx, uID).Return(int64(2), nil)

		p, err := s.CreatePlaylist(ctx, service.CreateInput{UserID: uID, Name: "My Playlist"})
		assert.ErrorIs(t, err, playlist.ErrPlaylistLimitReached)
		assert.Nil(t, p)
	})
}

func TestUpdatePlaylist(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	pID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mRepo := new(mockPlaylistRepo)
		s := service.NewService(mRepo, nil, 2)

		p := &playlist.Playlist{ID: pID, UserID: uID, Name: "Old Name"}
		mRepo.On("GetPlaylistByIDForUser", ctx, pID, uID).Return(p, nil)
		mRepo.On("UpdatePlaylist", ctx, mock.Anything).Return(p, nil)

		updated, err := s.UpdatePlaylist(ctx, service.UpdateInput{ID: pID, UserID: uID, Name: "New Name"})
		assert.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		mRepo := new(mockPlaylistRepo)
		s := service.NewService(mRepo, nil, 2)

		mRepo.On("GetPlaylistByIDForUser", ctx, pID, uID).Return(nil, playlist.ErrPlaylistNotFound)

		updated, err := s.UpdatePlaylist(ctx, service.UpdateInput{ID: pID, UserID: uID, Name: "New Name"})
		assert.ErrorIs(t, err, playlist.ErrPlaylistNotFound)
		assert.Nil(t, updated)
	})
}
