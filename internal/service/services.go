package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/cache"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository"
	serviceadmin "github.com/DaniilKalts/music-platform-api/internal/service/admin"
	"github.com/DaniilKalts/music-platform-api/internal/service/auth"
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
	GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error)
	ListTracks(ctx context.Context, page, limit int32) ([]*track.Track, error)
	SearchTracks(ctx context.Context, query string, page, limit int32) ([]*track.Track, error)
	ListGenres(ctx context.Context) ([]track.Genre, error)
	PlayTrack(ctx context.Context, userID, trackID uuid.UUID) (*track.Track, error)
}

type AdminService interface {
	CreateTrack(ctx context.Context, input serviceadmin.CreateTrackInput) (*track.Track, error)
	UpdateTrack(ctx context.Context, input serviceadmin.UpdateTrackInput) (*track.Track, error)
	DeleteTrack(ctx context.Context, id uuid.UUID) error
	UpdateUserSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error)
}

type Services struct {
	Auth  AuthService
	User  UserService
	Track TrackService
	Admin AdminService
}

func NewServices(
	repositories *repository.Repositories,
	tokenManager auth.TokenManager,
	caches *cache.Caches,
) *Services {
	return &Services{
		Auth: auth.NewService(repositories.User, tokenManager, caches.Blacklist, caches.Refresh),
		User: serviceuser.NewService(repositories.User),
		Track: servicetrack.NewService(
			repositories.Track,
			repositories.History,
			caches.Track,
			caches.Genre,
			caches.Search,
			caches.Popular,
		),
		Admin: serviceadmin.NewService(
			repositories.Track,
			repositories.User,
			caches.Track,
		),
	}
}
