package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
	"github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	tests := []struct {
		name           string
		authHeader     string
		claims         *jwt.Claims
		parseErr       error
		revoked        bool
		expectedStatus int
		expectedBody   string
		expectIdentity bool
	}{
		{
			name:           "missing header passes through",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid token injects identity",
			authHeader:     "Bearer access-token",
			claims:         &jwt.Claims{UserID: userID, Role: "USER"},
			expectedStatus: http.StatusOK,
			expectIdentity: true,
		},
		{
			name:           "bad header",
			authHeader:     "access-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid auth header"}`,
		},
		{
			name:           "revoked token",
			authHeader:     "Bearer access-token",
			revoked:        true,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"token is revoked"}`,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer access-token",
			parseErr:       jwt.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid or expired token"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				identity, ok := httpx.UserFromContext(r.Context())
				require.Equal(t, tt.expectIdentity, ok)
				if tt.expectIdentity {
					require.Equal(t, userID, identity.ID)
					require.Equal(t, "USER", identity.Role)
				}
				w.WriteHeader(http.StatusOK)
			})

			h := Auth(&tokenManagerMock{claims: tt.claims, err: tt.parseErr}, &blacklistMock{revoked: tt.revoked})(next)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				r.Header.Set("Authorization", tt.authHeader)
			}

			h.ServeHTTP(w, r)

			require.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				require.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestRequireAuth(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("missing identity", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		RequireAuth(next).ServeHTTP(w, r)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		require.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String())
	})

	t.Run("identity exists", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(httpx.WithUser(r.Context(), httpx.UserIdentity{ID: uuid.New(), Role: "USER"}))

		RequireAuth(next).ServeHTTP(w, r)

		require.Equal(t, http.StatusNoContent, w.Code)
	})
}

type tokenManagerMock struct {
	claims *jwt.Claims
	err    error
}

func (m *tokenManagerMock) ParseToken(_ string, _ jwt.TokenType) (*jwt.Claims, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.claims, nil
}

type blacklistMock struct {
	revoked bool
	err     error
}

func (m *blacklistMock) IsRevoked(_ context.Context, _ string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.revoked, nil
}
