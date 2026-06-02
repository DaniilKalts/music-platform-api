package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/middleware"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.With(middleware.RequireAuth).Post("/logout", h.Logout)
		r.Post("/refresh", h.Refresh)
	})
}
