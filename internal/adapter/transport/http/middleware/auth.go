package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
	"github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

type TokenManager interface {
	ParseToken(tokenStr string, tokenType jwt.TokenType) (*jwt.Claims, error)
}

type Blacklist interface {
	IsRevoked(ctx context.Context, token string) (bool, error)
}

func Auth(manager TokenManager, blacklist Blacklist) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid auth header")
				return
			}

			token := parts[1]

			revoked, err := blacklist.IsRevoked(r.Context(), token)
			if err != nil {
				httpx.WriteInternalError(w, r, err)
				return
			}
			if revoked {
				httpx.WriteError(w, http.StatusUnauthorized, "token is revoked")
				return
			}

			claims, err := manager.ParseToken(token, jwt.TokenTypeAccess)
			if err != nil {
				httpx.WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := httpx.WithUser(r.Context(), httpx.UserIdentity{
				ID:   claims.UserID,
				Role: claims.Role,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := httpx.UserFromContext(r.Context()); !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
