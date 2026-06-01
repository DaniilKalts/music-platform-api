package v1

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/auth"
)

type Dependencies struct {
	AuthService auth.Service
}

func RegisterRoutes(r chi.Router, deps Dependencies) {
	r.Route("/api/v1", func(r chi.Router) {
		auth.RegisterRoutes(r, deps.AuthService)
	})
}
