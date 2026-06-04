package config

import "fmt"

type Config struct {
	Server   Server
	Postgres Postgres `envPrefix:"DB_"`
	Redis    Redis    `envPrefix:"REDIS_"`
	JWT      JWT      `envPrefix:"JWT_"`
	Logger   Logger   `envPrefix:"LOG_"`
	S3       S3
	Limits   Limits
}

type Limits struct {
	FreePlaylistLimit  int `env:"FREE_PLAYLIST_LIMIT" envDefault:"3"`
	FreeFavoritesLimit int `env:"FREE_FAVORITES_LIMIT" envDefault:"20"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis config: %w", err)
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config: %w", err)
	}

	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("logger config: %w", err)
	}

	if err := c.S3.Validate(); err != nil {
		return fmt.Errorf("s3 config: %w", err)
	}

	return nil
}
