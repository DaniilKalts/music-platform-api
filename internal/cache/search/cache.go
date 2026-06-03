package search

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

const (
	searchPrefix = "search:"
	searchTTL    = 5 * time.Minute
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, query string, tracks []*track.Track) error {
	data, err := json.Marshal(tracks)
	if err != nil {
		return fmt.Errorf("marshal search results: %w", err)
	}

	if err := c.client.Set(ctx, key(query), data, searchTTL).Err(); err != nil {
		return fmt.Errorf("set search results in redis: %w", err)
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, query string) ([]*track.Track, error) {
	data, err := c.client.Get(ctx, key(query)).Bytes()
	if err != nil {
		return nil, err
	}

	var tracks []*track.Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("unmarshal search results: %w", err)
	}

	return tracks, nil
}

func key(query string) string {
	hash := sha256.Sum256([]byte(query))
	return searchPrefix + hex.EncodeToString(hash[:])
}
