package user

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/me", h.GetMe)
		r.Patch("/me", h.UpdateMe)
	})
}
