package middleware

import (
	"net/http"

	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, ok := httpx.UserFromContext(r.Context())
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if identity.Role != role {
				httpx.WriteError(w, http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
