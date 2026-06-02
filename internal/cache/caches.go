package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/cache/blacklist"
	"github.com/DaniilKalts/music-platform-api/internal/cache/refresh"
)

type Caches struct {
	Blacklist *blacklist.Cache
	Refresh   *refresh.Cache
}

func NewCaches(client *redis.Client) *Caches {
	return &Caches{
		Blacklist: blacklist.NewCache(client),
		Refresh:   refresh.NewCache(client),
	}
}
