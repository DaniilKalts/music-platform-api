package genre

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

const (
	genresKey = "genres"
	genresTTL = 24 * time.Hour
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, genres []track.Genre) error {
	data, err := json.Marshal(genres)
	if err != nil {
		return fmt.Errorf("marshal genres: %w", err)
	}

	if err := c.client.Set(ctx, genresKey, data, genresTTL).Err(); err != nil {
		return fmt.Errorf("set genres in redis: %w", err)
	}

	return nil
}

func (c *Cache) Get(ctx context.Context) ([]track.Genre, error) {
	data, err := c.client.Get(ctx, genresKey).Bytes()
	if err != nil {
		return nil, err
	}

	var genres []track.Genre
	if err := json.Unmarshal(data, &genres); err != nil {
		return nil, fmt.Errorf("unmarshal genres: %w", err)
	}

	return genres, nil
}

func (c *Cache) Delete(ctx context.Context) error {
	if err := c.client.Del(ctx, genresKey).Err(); err != nil {
		return fmt.Errorf("delete genres from redis: %w", err)
	}

	return nil
}
