package track

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

const (
	trackPrefix = "track:"
	trackTTL    = 15 * time.Minute
	notFoundTTL = 5 * time.Minute
)

var errTrackNotFoundCached = errors.New("track not found (cached)")

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, t *track.Track) error {
	ttl := withJitter(trackTTL, 0.1)

	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshal track: %w", err)
	}

	if err := c.client.Set(ctx, key(t.ID), data, ttl).Err(); err != nil {
		return fmt.Errorf("set track in redis: %w", err)
	}

	return nil
}

func (c *Cache) SetNotFound(ctx context.Context, id uuid.UUID) error {
	if err := c.client.Set(ctx, key(id), "null", notFoundTTL).Err(); err != nil {
		return fmt.Errorf("set not found in redis: %w", err)
	}
	return nil
}

func (c *Cache) Get(ctx context.Context, id uuid.UUID) (*track.Track, error) {
	data, err := c.client.Get(ctx, key(id)).Bytes()
	if err != nil {
		return nil, err
	}

	if string(data) == "null" {
		return nil, errTrackNotFoundCached
	}

	var t track.Track
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("unmarshal track: %w", err)
	}

	return &t, nil
}

func (c *Cache) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.client.Del(ctx, key(id)).Err(); err != nil {
		return fmt.Errorf("delete track from redis: %w", err)
	}

	return nil
}

func key(id uuid.UUID) string {
	return trackPrefix + id.String()
}

func withJitter(base time.Duration, jitter float64) time.Duration {
	if jitter <= 0 {
		return base
	}
	f := float64(base)
	min := f * (1 - jitter)
	max := f * (1 + jitter)
	return time.Duration(min + rand.Float64()*(max-min))
}
