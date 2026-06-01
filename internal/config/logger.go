package config

import (
	"fmt"
	"strings"
)

type Logger struct {
	Level  string `env:"LEVEL" envDefault:"info"`
	Format string `env:"FORMAT" envDefault:"json"`
}

func (c Logger) Validate() error {
	switch strings.ToLower(strings.TrimSpace(c.Level)) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("unsupported log level %q", c.Level)
	}

	switch strings.ToLower(strings.TrimSpace(c.Format)) {
	case "json", "console":
	default:
		return fmt.Errorf("unsupported log format %q", c.Format)
	}

	return nil
}
