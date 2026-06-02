package user

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

	domainuser "github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceuser "github.com/DaniilKalts/music-platform-api/internal/service/user"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

func TestHandlerGetMe(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		identity       *httpx.UserIdentity
		mockUser       *domainuser.User
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success",
			identity: &httpx.UserIdentity{ID: userID, Role: "USER"},
			mockUser: &domainuser.User{
				ID:           userID,
				Email:        "daniil.kalts@rbk.kz",
				Username:     "daniilkalts",
				Role:         domainuser.RoleUser,
				Subscription: domainuser.SubscriptionFree,
				CreatedAt:    now,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"` + userID.String() + `","email":"daniil.kalts@rbk.kz","username":"daniilkalts","role":"USER","subscription_type":"FREE","created_at":"2026-01-02T03:04:05Z"}`,
		},
		{
			name:           "unauthorized",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
		{
			name:           "not found",
			identity:       &httpx.UserIdentity{ID: userID, Role: "USER"},
			mockErr:        domainuser.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(&serviceMock{user: tt.mockUser, err: tt.mockErr})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
			if tt.identity != nil {
				r = r.WithContext(httpx.WithUser(r.Context(), *tt.identity))
			}

			h.GetMe(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandlerUpdateMe(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name           string
		identity       *httpx.UserIdentity
		mockUser       *domainuser.User
		mockErr        error
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success",
			identity: &httpx.UserIdentity{ID: userID, Role: "USER"},
			mockUser: &domainuser.User{
				ID:           userID,
				Email:        "new.email@rbk.kz",
				Username:     "newusername",
				Role:         domainuser.RoleUser,
				Subscription: domainuser.SubscriptionFree,
				CreatedAt:    now,
			},
			body:           `{"email":"new.email@rbk.kz","username":"newusername"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"` + userID.String() + `","email":"new.email@rbk.kz","username":"newusername","role":"USER","subscription_type":"FREE","created_at":"2026-01-02T03:04:05Z"}`,
		},
		{
			name:           "unauthorized",
			body:           `{"email":"new.email@rbk.kz"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
		{
			name:           "validation error",
			identity:       &httpx.UserIdentity{ID: userID, Role: "USER"},
			body:           `{"email":"bad"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"field email is invalid"}`,
		},
		{
			name:           "conflict",
			identity:       &httpx.UserIdentity{ID: userID, Role: "USER"},
			mockErr:        domainuser.ErrEmailAlreadyExists,
			body:           `{"email":"busy.email@rbk.kz"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"user with this email already exists"}`,
		},
		{
			name:           "internal error",
			identity:       &httpx.UserIdentity{ID: userID, Role: "USER"},
			mockErr:        errors.New("storage is down"),
			body:           `{"username":"newusername"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(&serviceMock{user: tt.mockUser, err: tt.mockErr})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me", bytes.NewBufferString(tt.body))
			if tt.identity != nil {
				r = r.WithContext(httpx.WithUser(r.Context(), *tt.identity))
			}

			h.UpdateMe(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

type serviceMock struct {
	user *domainuser.User
	err  error
}

func (m *serviceMock) GetMe(_ context.Context, _ uuid.UUID) (*domainuser.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func (m *serviceMock) UpdateMe(_ context.Context, _ uuid.UUID, _ serviceuser.UpdateInput) (*domainuser.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}
