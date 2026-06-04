package admin

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Use(middleware.RequireRole("ADMIN"))

		r.Post("/tracks", h.CreateTrack)
		r.Put("/tracks/{id}", h.UpdateTrack)
		r.Delete("/tracks/{id}", h.DeleteTrack)
		r.Patch("/users/{id}/subscription", h.UpdateUserSubscription)
	})
}
