package auth

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceauth "github.com/DaniilKalts/music-platform-api/internal/service/auth"
	jwtpkg "github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

func TestHandlerRegister(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		mockUser       *user.User
		mockErr        error
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			mockUser: &user.User{
				ID:           userID,
				Email:        "daniil.kalts@rbk.kz",
				Username:     "daniilkalts",
				Role:         user.RoleUser,
				Subscription: user.SubscriptionFree,
				CreatedAt:    now,
			},
			body:           `{"email":"daniil.kalts@rbk.kz","username":"daniilkalts","password":"12345678"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"` + userID.String() + `","email":"daniil.kalts@rbk.kz","username":"daniilkalts","role":"USER","subscription_type":"FREE","created_at":"2026-01-02T03:04:05Z"}`,
		},
		{
			name:           "validation error",
			mockErr:        user.ErrInvalidEmail,
			body:           `{"email":"bad","username":"daniilkalts","password":"12345678"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"field email is invalid"}`,
		},
		{
			name:           "conflict",
			mockErr:        user.ErrEmailAlreadyExists,
			body:           `{"email":"daniil.kalts@rbk.kz","username":"daniilkalts","password":"12345678"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"user with this email already exists"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(&serviceMock{
				registerUser: tt.mockUser,
				registerErr:  tt.mockErr,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(tt.body))

			h.Register(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandlerLogout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		logoutErr      error
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			authHeader:     "Bearer access-token",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"empty auth header"}`,
		},
		{
			name:           "invalid token",
			logoutErr:      jwtpkg.ErrInvalidToken,
			authHeader:     "Bearer bad-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid or expired token"}`,
		},
		{
			name:           "internal error",
			logoutErr:      errors.New("redis is down"),
			authHeader:     "Bearer access-token",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(&serviceMock{logoutErr: tt.logoutErr})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
			if tt.authHeader != "" {
				r.Header.Set("Authorization", tt.authHeader)
			}

			h.Logout(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				require.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestHandlerRefresh(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	refreshExp := time.Date(2026, 1, 3, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		mockToken      *serviceauth.TokenPair
		mockErr        error
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			mockToken: &serviceauth.TokenPair{
				AccessToken:           "new-access-token",
				AccessTokenExpiresAt:  now,
				RefreshToken:          "new-refresh-token",
				RefreshTokenExpiresAt: refreshExp,
			},
			body:           `{"refresh_token":"refresh-token"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"new-access-token","access_token_expires_at":"2026-01-02T03:04:05Z","refresh_token":"new-refresh-token","refresh_token_expires_at":"2026-01-03T03:04:05Z"}`,
		},
		{
			name:           "invalid refresh token",
			mockErr:        jwtpkg.ErrInvalidToken,
			body:           `{"refresh_token":"bad-token"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid refresh token"}`,
		},
		{
			name:           "internal error",
			mockErr:        errors.New("storage is down"),
			body:           `{"refresh_token":"refresh-token"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(&serviceMock{refreshToken: tt.mockToken, refreshErr: tt.mockErr})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(tt.body))

			h.Refresh(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandlerLogin(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	refreshExp := time.Date(2026, 1, 3, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		mockToken      *serviceauth.TokenPair
		mockErr        error
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			mockToken: &serviceauth.TokenPair{
				AccessToken:           "access-token",
				AccessTokenExpiresAt:  now,
				RefreshToken:          "refresh-token",
				RefreshTokenExpiresAt: refreshExp,
			},
			body:           `{"email":"daniil.kalts@rbk.kz","password":"12345678"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"access_token":"access-token","access_token_expires_at":"2026-01-02T03:04:05Z","refresh_token":"refresh-token","refresh_token_expires_at":"2026-01-03T03:04:05Z"}`,
		},
		{
			name:           "invalid credentials",
			mockErr:        serviceauth.ErrInvalidCredentials,
			body:           `{"email":"daniil.kalts@rbk.kz","password":"wrong-password"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid email or password"}`,
		},
		{
			name:           "internal error",
			mockErr:        errors.New("storage is down"),
			body:           `{"email":"daniil.kalts@rbk.kz","password":"12345678"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(&serviceMock{
				loginToken: tt.mockToken,
				loginErr:   tt.mockErr,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(tt.body))

			h.Login(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

type serviceMock struct {
	registerUser *user.User
	registerErr  error
	loginToken   *serviceauth.TokenPair
	loginErr     error
	logoutErr    error
	refreshToken *serviceauth.TokenPair
	refreshErr   error
}

func (m *serviceMock) Register(_ context.Context, _ serviceauth.RegisterInput) (*user.User, error) {
	if m.registerErr != nil {
		return nil, m.registerErr
	}
	return m.registerUser, nil
}

func (m *serviceMock) Login(_ context.Context, _ serviceauth.LoginInput) (*serviceauth.TokenPair, error) {
	if m.loginErr != nil {
		return nil, m.loginErr
	}
	return m.loginToken, nil
}

func (m *serviceMock) Logout(_ context.Context, _ string) error {
	return m.logoutErr
}

func (m *serviceMock) Refresh(_ context.Context, _ string) (*serviceauth.TokenPair, error) {
	if m.refreshErr != nil {
		return nil, m.refreshErr
	}
	return m.refreshToken, nil
}
