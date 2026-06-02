package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	userrepo "github.com/DaniilKalts/music-platform-api/internal/repository/user"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User, password user.Password) (*user.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*user.User, user.Password, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type Repositories struct {
	User UserRepository
}

func NewRepositories(db sqlc.DBTX) *Repositories {
	return &Repositories{
		User: userrepo.NewRepository(db),
	}
}
