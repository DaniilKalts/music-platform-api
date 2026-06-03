package userrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

const (
	emailUniqueConstraint    = "users_email_unique"
	usernameUniqueConstraint = "users_username_unique"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

func (r *Repository) Create(ctx context.Context, u user.User, password user.Password) (*user.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:               u.ID,
		Email:            u.Email,
		Username:         u.Username,
		PasswordHash:     password.Hash,
		Salt:             password.Salt,
		Role:             sqlc.UserRole(u.Role),
		SubscriptionType: sqlc.SubscriptionType(u.Subscription),
	})
	if err != nil {
		switch {
		case postgres.IsUniqueViolation(err, emailUniqueConstraint):
			return nil, user.ErrEmailAlreadyExists
		case postgres.IsUniqueViolation(err, usernameUniqueConstraint):
			return nil, user.ErrUsernameAlreadyExists
		default:
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	return toDomain(userRow{
		ID:               row.ID,
		Email:            row.Email,
		Username:         row.Username,
		Role:             row.Role,
		SubscriptionType: row.SubscriptionType,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}), nil
}

func (r *Repository) GetCredentialsByEmail(ctx context.Context, email string) (*user.User, user.Password, error) {
	row, err := r.queries.GetUserCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.Password{}, user.ErrNotFound
		}

		return nil, user.Password{}, fmt.Errorf("get user credentials by email: %w", err)
	}

	u := &user.User{
		ID:    row.ID,
		Email: row.Email,
		Role:  user.Role(row.Role),
	}
	password := user.Password{Hash: row.PasswordHash, Salt: row.Salt}

	return u, password, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return toDomain(userRow{
		ID:               row.ID,
		Email:            row.Email,
		Username:         row.Username,
		Role:             row.Role,
		SubscriptionType: row.SubscriptionType,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}), nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, email, username *string) (*user.User, error) {
	row, err := r.queries.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:       id,
		Email:    nullableText(email),
		Username: nullableText(username),
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, user.ErrNotFound
		case postgres.IsUniqueViolation(err, emailUniqueConstraint):
			return nil, user.ErrEmailAlreadyExists
		case postgres.IsUniqueViolation(err, usernameUniqueConstraint):
			return nil, user.ErrUsernameAlreadyExists
		default:
			return nil, fmt.Errorf("update user profile: %w", err)
		}
	}

	return toDomain(userRow{
		ID:               row.ID,
		Email:            row.Email,
		Username:         row.Username,
		Role:             row.Role,
		SubscriptionType: row.SubscriptionType,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}), nil
}

func (r *Repository) UpdateSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error) {
	row, err := r.queries.UpdateUserSubscription(ctx, sqlc.UpdateUserSubscriptionParams{
		ID:               id,
		SubscriptionType: sqlc.SubscriptionType(sub),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("update user subscription: %w", err)
	}

	return toDomain(userRow{
		ID:               row.ID,
		Email:            row.Email,
		Username:         row.Username,
		Role:             row.Role,
		SubscriptionType: row.SubscriptionType,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}), nil
}

func nullableText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *value, Valid: true}
}
