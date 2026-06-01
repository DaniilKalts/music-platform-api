package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Postgres struct {
	Host             string        `env:"HOST" envDefault:"localhost"`
	Port             int           `env:"PORT" envDefault:"5432"`
	User             string        `env:"USER" envDefault:"postgres"`
	Password         string        `env:"PASSWORD" envDefault:"postgres"`
	Name             string        `env:"NAME" envDefault:"music_platform"`
	SSLMode          string        `env:"SSL_MODE" envDefault:"disable"`
	MaxConns         int32         `env:"MAX_CONNS" envDefault:"10"`
	MinConns         int32         `env:"MIN_CONNS" envDefault:"1"`
	MaxConnLifetime  time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime  time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"30m"`
	StatementTimeout time.Duration `env:"STATEMENT_TIMEOUT" envDefault:"3s"`
}

func (c Postgres) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("host is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be in range 1..65535")
	}
	if strings.TrimSpace(c.User) == "" {
		return fmt.Errorf("user is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if !isValidSSLMode(c.SSLMode) {
		return fmt.Errorf("ssl mode must be one of disable, allow, prefer, require, verify-ca, verify-full")
	}
	if c.MaxConns <= 0 {
		return fmt.Errorf("max conns must be positive")
	}
	if c.MinConns < 0 || c.MinConns > c.MaxConns {
		return fmt.Errorf("min conns must be in range 0..max conns")
	}
	if c.MaxConnLifetime <= 0 || c.MaxConnIdleTime <= 0 || c.StatementTimeout <= 0 {
		return fmt.Errorf("postgres timeouts must be positive")
	}

	return nil
}

func (c Postgres) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Name,
	}

	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func isValidSSLMode(value string) bool {
	switch value {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
		return true
	default:
		return false
	}
}
