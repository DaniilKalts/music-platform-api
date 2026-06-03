package favorite_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DaniilKalts/music-platform-api/internal/domain/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	service "github.com/DaniilKalts/music-platform-api/internal/service/favorite"
)

type mockFavRepo struct{ mock.Mock }

func (m *mockFavRepo) AddFavorite(ctx context.Context, f *favorite.Favorite) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}
func (m *mockFavRepo) RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error {
	args := m.Called(ctx, userID, trackID)
	return args.Error(0)
}
func (m *mockFavRepo) ListFavoritesByUserID(ctx context.Context, userID uuid.UUID) ([]*track.Track, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*track.Track), args.Error(1)
}
func (m *mockFavRepo) CountFavoritesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func TestAddFavorite(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	trID := uuid.New()
	limit := 2

	t.Run("Success_Premium", func(t *testing.T) {
		mFav := new(mockFavRepo)
		mUser := new(mockUserRepo)
		s := service.NewService(mFav, mUser, limit)

		mUser.On("GetByID", ctx, uID).Return(&user.User{ID: uID, Subscription: user.SubscriptionPremium}, nil)
		mFav.On("AddFavorite", ctx, mock.Anything).Return(nil)

		err := s.AddFavorite(ctx, uID, trID)
		assert.NoError(t, err)
	})

	t.Run("Success_Free_UnderLimit", func(t *testing.T) {
		mFav := new(mockFavRepo)
		mUser := new(mockUserRepo)
		s := service.NewService(mFav, mUser, limit)

		mUser.On("GetByID", ctx, uID).Return(&user.User{ID: uID, Subscription: user.SubscriptionFree}, nil)
		mFav.On("CountFavoritesByUserID", ctx, uID).Return(int64(1), nil)
		mFav.On("AddFavorite", ctx, mock.Anything).Return(nil)

		err := s.AddFavorite(ctx, uID, trID)
		assert.NoError(t, err)
	})

	t.Run("Error_Free_LimitExceeded", func(t *testing.T) {
		mFav := new(mockFavRepo)
		mUser := new(mockUserRepo)
		s := service.NewService(mFav, mUser, limit)

		mUser.On("GetByID", ctx, uID).Return(&user.User{ID: uID, Subscription: user.SubscriptionFree}, nil)
		mFav.On("CountFavoritesByUserID", ctx, uID).Return(int64(2), nil)

		err := s.AddFavorite(ctx, uID, trID)
		assert.ErrorIs(t, err, favorite.ErrFavoriteLimitReached)
	})
}
