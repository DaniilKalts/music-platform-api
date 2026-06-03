package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository"
	"github.com/DaniilKalts/music-platform-api/internal/service/auth"
	servicefavorite "github.com/DaniilKalts/music-platform-api/internal/service/favorite"
	servicehistory "github.com/DaniilKalts/music-platform-api/internal/service/history"
	servicetrack "github.com/DaniilKalts/music-platform-api/internal/service/track"
	serviceuser "github.com/DaniilKalts/music-platform-api/internal/service/user"
)

type AuthService interface {
	Register(ctx context.Context, input auth.RegisterInput) (*user.User, error)
	Login(ctx context.Context, input auth.LoginInput) (*auth.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*auth.TokenPair, error)
}

type UserService interface {
	GetMe(ctx context.Context, id uuid.UUID) (*user.User, error)
	UpdateMe(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*user.User, error)
}

type TrackService interface {
	GetTrack(ctx context.Context, id uuid.UUID) (*track.Track, error)
	ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error)
	SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error)
	ListGenres(ctx context.Context) ([]*track.Genre, error)
	PlayTrack(ctx context.Context, userID, trackID uuid.UUID) (*track.Track, error)
}

type FavoriteService interface {
	AddFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	ListFavorites(ctx context.Context, userID uuid.UUID) ([]*track.Track, error)
}

type HistoryService interface {
	ListHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error)
}

type Services struct {
	Auth     AuthService
	User     UserService
	Track    TrackService
	Favorite FavoriteService
	History  HistoryService
}

func NewServices(
	repositories *repository.Repositories,
	tokenManager auth.TokenManager,
	blacklist auth.Blacklist,
	refresh auth.RefreshTokens,
	tCache servicetrack.TrackCache,
	gCache servicetrack.GenreCache,
	sCache servicetrack.SearchCache,
	freeFavLimit int,
) *Services {
	return &Services{
		Auth: auth.NewService(repositories.User, tokenManager, blacklist, refresh),
		User: serviceuser.NewService(repositories.User),
		Track: servicetrack.NewService(
			repositories.Track,
			repositories.History,
			tCache,
			gCache,
			sCache,
		),
		Favorite: servicefavorite.NewService(repositories.Favorite, repositories.User, freeFavLimit),
		History:  servicehistory.NewService(repositories.History),
	}
}
