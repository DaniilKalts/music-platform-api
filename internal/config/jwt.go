package config

import (
	"fmt"
	"strings"
	"time"
)

type JWT struct {
	AccessSecret  string        `env:"ACCESS_SECRET,required"`
	RefreshSecret string        `env:"REFRESH_SECRET,required"`
	AccessTTL     time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
	RefreshTTL    time.Duration `env:"REFRESH_TTL" envDefault:"720h"`
}

func (c JWT) Validate() error {
	if strings.TrimSpace(c.AccessSecret) == "" {
		return fmt.Errorf("access secret is required")
	}
	if strings.TrimSpace(c.RefreshSecret) == "" {
		return fmt.Errorf("refresh secret is required")
	}
	if c.AccessTTL <= 0 {
		return fmt.Errorf("access ttl must be positive")
	}
	if c.RefreshTTL <= 0 {
		return fmt.Errorf("refresh ttl must be positive")
	}

	return nil
}
