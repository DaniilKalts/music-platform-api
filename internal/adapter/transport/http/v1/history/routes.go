package history

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/listening-history", func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Get("/", h.List)
	})
}
