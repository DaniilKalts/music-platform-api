package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

type userRow struct {
	ID               uuid.UUID
	Email            string
	Username         string
	Role             sqlc.UserRole
	SubscriptionType sqlc.SubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func toDomain(r userRow) *user.User {
	return &user.User{
		ID:           r.ID,
		Email:        r.Email,
		Username:     r.Username,
		Role:         user.Role(r.Role),
		Subscription: user.Subscription(r.SubscriptionType),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
