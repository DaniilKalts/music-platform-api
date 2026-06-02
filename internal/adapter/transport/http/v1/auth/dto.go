package auth

import (
	"time"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceauth "github.com/DaniilKalts/music-platform-api/internal/service/auth"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type RegisterResponse struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Username         string `json:"username"`
	Role             string `json:"role"`
	SubscriptionType string `json:"subscription_type"`
	CreatedAt        string `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

func ToRegisterInput(r RegisterRequest) serviceauth.RegisterInput {
	return serviceauth.RegisterInput{Email: r.Email, Username: r.Username, Password: r.Password}
}

func ToLoginInput(r LoginRequest) serviceauth.LoginInput {
	return serviceauth.LoginInput{Email: r.Email, Password: r.Password}
}

func ToRegisterResponse(u user.User) RegisterResponse {
	return RegisterResponse{
		ID:               u.ID.String(),
		Email:            u.Email,
		Username:         u.Username,
		Role:             string(u.Role),
		SubscriptionType: string(u.Subscription),
		CreatedAt:        formatTime(u.CreatedAt),
	}
}

func ToTokenResponse(token serviceauth.TokenPair) TokenResponse {
	return TokenResponse{
		AccessToken:           token.AccessToken,
		AccessTokenExpiresAt:  formatTime(token.AccessTokenExpiresAt),
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresAt: formatTime(token.RefreshTokenExpiresAt),
	}
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
