package user

import (
	"time"

	domainuser "github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceuser "github.com/DaniilKalts/music-platform-api/internal/service/user"
)

type Response struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Username         string `json:"username"`
	Role             string `json:"role"`
	SubscriptionType string `json:"subscription_type"`
	CreatedAt        string `json:"created_at"`
}

type UpdateRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
}

func ToUpdateInput(r UpdateRequest) serviceuser.UpdateInput {
	return serviceuser.UpdateInput{Email: r.Email, Username: r.Username}
}

func ToResponse(u domainuser.User) Response {
	return Response{
		ID:               u.ID.String(),
		Email:            u.Email,
		Username:         u.Username,
		Role:             string(u.Role),
		SubscriptionType: string(u.Subscription),
		CreatedAt:        u.CreatedAt.UTC().Format(time.RFC3339),
	}
}
