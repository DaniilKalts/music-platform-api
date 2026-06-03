package popular

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

const (
	popularKey = "popular_tracks"
	popularTTL = 30 * time.Minute
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, tracks []*track.Track) error {
	data, err := json.Marshal(tracks)
	if err != nil {
		return fmt.Errorf("marshal popular tracks: %w", err)
	}

	if err := c.client.Set(ctx, popularKey, data, popularTTL).Err(); err != nil {
		return fmt.Errorf("set popular tracks in redis: %w", err)
	}

	return nil
}

func (c *Cache) Get(ctx context.Context) ([]*track.Track, error) {
	data, err := c.client.Get(ctx, popularKey).Bytes()
	if err != nil {
		return nil, err
	}

	var tracks []*track.Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("unmarshal popular tracks: %w", err)
	}

	return tracks, nil
}
