package track

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/tracks", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/search", h.Search)
		r.Get("/genres", h.ListGenres)
		r.Get("/{id}", h.Get)

		r.With(middleware.RequireAuth).Post("/{id}/play", h.Play)
	})
}
