package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/cache/blacklist"
)

type Caches struct {
	Blacklist *blacklist.Cache
}

func NewCaches(client *redis.Client) *Caches {
	return &Caches{
		Blacklist: blacklist.NewCache(client),
	}
}
