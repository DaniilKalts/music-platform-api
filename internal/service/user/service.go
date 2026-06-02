package user

import (
	"context"
	"strings"

	"github.com/google/uuid"

	domainuser "github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, email, username *string) (*domainuser.User, error)
}

type Service struct {
	users Repository
}

type UpdateInput struct {
	Email    *string
	Username *string
}

func NewService(users Repository) *Service {
	return &Service{users: users}
}

func (s *Service) GetMe(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *Service) UpdateMe(ctx context.Context, id uuid.UUID, input UpdateInput) (*domainuser.User, error) {
	current, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	email := current.Email
	if input.Email != nil {
		email = *input.Email
	}

	username := current.Username
	if input.Username != nil {
		username = *input.Username
	}

	if err := current.UpdateProfile(email, username); err != nil {
		return nil, err
	}

	var updatedEmail *string
	if input.Email != nil {
		normalized := domainuser.NormalizeEmail(*input.Email)
		updatedEmail = &normalized
	}

	var updatedUsername *string
	if input.Username != nil {
		trimmed := strings.TrimSpace(*input.Username)
		updatedUsername = &trimmed
	}

	return s.users.UpdateProfile(ctx, id, updatedEmail, updatedUsername)
}
