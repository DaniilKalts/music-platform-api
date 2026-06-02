package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository"
	"github.com/DaniilKalts/music-platform-api/internal/service/auth"
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

type Services struct {
	Auth AuthService
	User UserService
}

func NewServices(
	repositories *repository.Repositories,
	tokenManager auth.TokenManager,
	blacklist auth.Blacklist,
	refresh auth.RefreshTokens,
) *Services {
	return &Services{
		Auth: auth.NewService(repositories.User, tokenManager, blacklist, refresh),
		User: serviceuser.NewService(repositories.User),
	}
}
