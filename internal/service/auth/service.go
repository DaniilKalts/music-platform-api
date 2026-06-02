package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserRepository interface {
	Create(ctx context.Context, u user.User, password user.Password) (*user.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*user.User, user.Password, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type TokenManager interface {
	GeneratePair(userID uuid.UUID, role string) (*jwt.Pair, error)
	ParseToken(tokenStr string, tokenType jwt.TokenType) (*jwt.Claims, error)
}

type Blacklist interface {
	Revoke(ctx context.Context, token string, expiresAt time.Time) error
}

type Service struct {
	users        UserRepository
	tokenManager TokenManager
	blacklist    Blacklist
}

func NewService(users UserRepository, tokenManager TokenManager, blacklist Blacklist) *Service {
	return &Service{
		users:        users,
		tokenManager: tokenManager,
		blacklist:    blacklist,
	}
}

func (s *Service) Logout(ctx context.Context, token string) error {
	claims, err := s.tokenManager.ParseToken(token, jwt.TokenTypeAccess)
	if err != nil {
		return err
	}

	return s.blacklist.Revoke(ctx, token, claims.ExpiresAt.Time)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.tokenManager.ParseToken(refreshToken, jwt.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	u, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	pair, err := s.tokenManager.GeneratePair(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           pair.AccessToken,
		AccessTokenExpiresAt:  pair.AccessTokenExpiresAt,
		RefreshToken:          pair.RefreshToken,
		RefreshTokenExpiresAt: pair.RefreshTokenExpiresAt,
	}, nil
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*user.User, error) {

	password, err := user.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	u, err := user.NewUser(input.Email, input.Username)
	if err != nil {
		return nil, err
	}

	return s.users.Create(ctx, *u, password)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*TokenPair, error) {
	email := user.NormalizeEmail(input.Email)
	if email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	u, password, err := s.users.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !password.Matches(input.Password) {
		return nil, ErrInvalidCredentials
	}

	pair, err := s.tokenManager.GeneratePair(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           pair.AccessToken,
		AccessTokenExpiresAt:  pair.AccessTokenExpiresAt,
		RefreshToken:          pair.RefreshToken,
		RefreshTokenExpiresAt: pair.RefreshTokenExpiresAt,
	}, nil
}
