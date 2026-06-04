package v1

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/admin"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/auth"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/history"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/track"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/user"
	"github.com/DaniilKalts/music-platform-api/internal/service"
)

type Dependencies struct {
	AuthService     service.AuthService
	UserService     service.UserService
	TrackService    service.TrackService
	PlaylistService service.PlaylistService
	FavoriteService service.FavoriteService
	HistoryService  service.HistoryService
	AdminService    service.AdminService
}

func RegisterRoutes(r chi.Router, deps Dependencies) {
	r.Route("/api/v1", func(r chi.Router) {
		auth.RegisterRoutes(r, deps.AuthService)
		user.RegisterRoutes(r, deps.UserService)
		track.RegisterRoutes(r, deps.TrackService)
		favorite.RegisterRoutes(r, deps.FavoriteService)
		history.RegisterRoutes(r, deps.HistoryService)
		playlist.RegisterRoutes(r, deps.PlaylistService)
		admin.RegisterRoutes(r, deps.AdminService)
	})
}
