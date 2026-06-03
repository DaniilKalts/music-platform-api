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
	track     *track.Cache
	genre     *genre.Cache
	search    *search.Cache
	popular   *popular.Cache
}

func (c *Caches) Track() *track.Cache {
	return c.track
}

func (c *Caches) Genre() *genre.Cache {
	return c.genre
}

func (c *Caches) Search() *search.Cache {
	return c.search
}

func (c *Caches) Popular() *popular.Cache {
	return c.popular
}

func NewCaches(client *redis.Client) *Caches {
	return &Caches{
		Blacklist: blacklist.NewCache(client),
		Refresh:   refresh.NewCache(client),
		track:     track.NewCache(client),
		genre:     genre.NewCache(client),
		search:    search.NewCache(client),
		popular:   popular.NewCache(client),
	}
}
