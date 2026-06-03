package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/cache/blacklist"
	"github.com/DaniilKalts/music-platform-api/internal/cache/genre"
	"github.com/DaniilKalts/music-platform-api/internal/cache/popular"
	"github.com/DaniilKalts/music-platform-api/internal/cache/refresh"
	"github.com/DaniilKalts/music-platform-api/internal/cache/search"
	"github.com/DaniilKalts/music-platform-api/internal/cache/track"
)

type Caches struct {
	Blacklist *blacklist.Cache
	Refresh   *refresh.Cache
	Track     *track.Cache
	Genre     *genre.Cache
	Search    *search.Cache
	Popular   *popular.Cache
}

func NewCaches(client *redis.Client) *Caches {
	return &Caches{
		Blacklist: blacklist.NewCache(client),
		Refresh:   refresh.NewCache(client),
		Track:     track.NewCache(client),
		Genre:     genre.NewCache(client),
		Search:    search.NewCache(client),
		Popular:   popular.NewCache(client),
	}
}
