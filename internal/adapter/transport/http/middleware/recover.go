package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
	"github.com/DaniilKalts/music-platform-api/pkg/logger"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.FromContext(r.Context()).Error("panic recovered", zap.Any("panic", recovered))
				httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
