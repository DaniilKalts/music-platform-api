package config

import "fmt"

type Config struct {
	Server   Server
	Postgres Postgres `envPrefix:"DB_"`
	JWT      JWT      `envPrefix:"JWT_"`
	Logger   Logger   `envPrefix:"LOG_"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config: %w", err)
	}

	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("logger config: %w", err)
	}

	return nil
}
