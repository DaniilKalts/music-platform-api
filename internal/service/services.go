package service

import (
	"context"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository"
	"github.com/DaniilKalts/music-platform-api/internal/service/auth"
)

type AuthService interface {
	Register(ctx context.Context, input auth.RegisterInput) (*user.User, error)
	Login(ctx context.Context, input auth.LoginInput) (*auth.TokenPair, error)
}

type Services struct {
	Auth AuthService
}

func NewServices(repositories *repository.Repositories, tokenManager auth.TokenManager) *Services {
	return &Services{
		Auth: auth.NewService(repositories.User, tokenManager),
	}
}
