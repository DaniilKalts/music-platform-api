package playlist

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/playlists", func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Get("/{id}/tracks", h.ListTracks)
		r.Post("/{playlistID}/tracks/{trackID}", h.AddTrack)
		r.Delete("/{playlistID}/tracks/{trackID}", h.RemoveTrack)
	})
}
