package refresh

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

const refreshTokenPrefix = "refresh_token:"
const rotateScript = `
if redis.call("EXISTS", KEYS[1]) == 0 then
  return 0
end
redis.call("DEL", KEYS[1])
redis.call("SET", KEYS[2], "1", "PX", ARGV[1])
return 1
`

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Add(ctx context.Context, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	if err := c.client.Set(ctx, key(token), "1", ttl).Err(); err != nil {
		return fmt.Errorf("add refresh token: %w", err)
	}

	return nil
}

func (c *Cache) Exists(ctx context.Context, token string) (bool, error) {
	count, err := c.client.Exists(ctx, key(token)).Result()
	if err != nil {
		return false, fmt.Errorf("check refresh token: %w", err)
	}

	return count > 0, nil
}

func (c *Cache) Remove(ctx context.Context, token string) error {
	if err := c.client.Del(ctx, key(token)).Err(); err != nil {
		return fmt.Errorf("remove refresh token: %w", err)
	}

	return nil
}

func (c *Cache) Rotate(ctx context.Context, oldToken, newToken string, newExpiresAt time.Time) (bool, error) {
	ttl := time.Until(newExpiresAt)
	if ttl <= 0 {
		return false, nil
	}

	rotated, err := c.client.Eval(ctx, rotateScript, []string{key(oldToken), key(newToken)}, ttl.Milliseconds()).Int()
	if err != nil {
		return false, fmt.Errorf("rotate refresh token: %w", err)
	}

	return rotated == 1, nil
}

func key(token string) string {
	return refreshTokenPrefix + jwt.HashToken(token)
}
