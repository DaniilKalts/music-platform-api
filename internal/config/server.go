package config

import (
	"fmt"
	"time"
)

type Server struct {
	Port               int           `env:"APP_PORT" envDefault:"8080"`
	HTTPTimeout        time.Duration `env:"SERVER_HTTP_TIMEOUT" envDefault:"15s"`
	HandlerTimeout     time.Duration `env:"SERVER_HANDLER_TIMEOUT" envDefault:"10s"`
	ShutdownTimeout    time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"15s"`
	CORSAllowedOrigins []string      `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000"`
}

func (c Server) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("app port must be in range 1..65535")
	}
	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("http timeout must be positive")
	}
	if c.HandlerTimeout <= 0 {
		return fmt.Errorf("handler timeout must be positive")
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}

	return nil
}

func (c Server) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
