package config

import (
	"fmt"
	"strings"
	"time"
)

type Redis struct {
	Host         string        `env:"HOST" envDefault:"localhost"`
	Port         int           `env:"PORT" envDefault:"6379"`
	Password     string        `env:"PASSWORD"`
	DB           int           `env:"DB" envDefault:"0"`
	DialTimeout  time.Duration `env:"DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"3s"`
}

func (c Redis) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Redis) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("redis host is required")
	}

	if c.Port <= 0 {
		return fmt.Errorf("redis port must be positive")
	}

	if c.DB < 0 {
		return fmt.Errorf("redis db cannot be negative")
	}

	if c.DialTimeout <= 0 {
		return fmt.Errorf("dial timeout must be positive")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	return nil
}
