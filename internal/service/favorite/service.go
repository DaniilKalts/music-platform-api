package favorite

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, f *favorite.Favorite) error
	RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	ListFavoritesByUserID(ctx context.Context, userID uuid.UUID) ([]*track.Track, error)
	CountFavoritesByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type Service struct {
	favorites FavoriteRepository
	users     UserRepository
	freeLimit int
}

func NewService(favorites FavoriteRepository, users UserRepository, freeLimit int) *Service {
	return &Service{
		favorites: favorites,
		users:     users,
		freeLimit: freeLimit,
	}
}

func (s *Service) AddFavorite(ctx context.Context, userID, trackID uuid.UUID) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if u.Subscription == user.SubscriptionFree {
		count, err := s.favorites.CountFavoritesByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("count favorites: %w", err)
		}

		if int(count) >= s.freeLimit {
			return favorite.ErrFavoriteLimitReached
		}
	}

	f, err := favorite.NewFavorite(userID, trackID)
	if err != nil {
		return fmt.Errorf("create favorite domain object: %w", err)
	}

	if err := s.favorites.AddFavorite(ctx, f); err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}

	return nil
}

func (s *Service) RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error {
	return s.favorites.RemoveFavorite(ctx, userID, trackID)
}

func (s *Service) ListFavorites(ctx context.Context, userID uuid.UUID) ([]*track.Track, error) {
	return s.favorites.ListFavoritesByUserID(ctx, userID)
}
