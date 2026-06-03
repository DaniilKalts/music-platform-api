package favorite

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/favorites", func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Get("/tracks", h.List)
		r.Post("/tracks/{track_id}", h.Add)
		r.Delete("/tracks/{track_id}", h.Remove)
	})
}
