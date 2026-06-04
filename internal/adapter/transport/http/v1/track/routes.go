package track

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/tracks", func(r chi.Router) {
		// Public catalog browsing — no auth required.
		r.Get("/", h.List)
		r.Get("/search", h.Search)
		r.Get("/genres", h.ListGenres)
		r.Get("/{id}", h.Get)

		// Recording a play is tied to a user, so it stays authenticated.
		r.With(middleware.RequireAuth).Post("/{id}/play", h.Play)
	})
}
