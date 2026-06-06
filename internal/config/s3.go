package config

import "fmt"

type S3 struct {
	Endpoint  string `env:"S3_ENDPOINT" envDefault:"localhost:9000"`
	PublicURL string `env:"S3_PUBLIC_URL" envDefault:"http://localhost:9000"`
	AccessKey string `env:"S3_ACCESS_KEY" envDefault:"admin"`
	SecretKey string `env:"S3_SECRET_KEY" envDefault:"password"`
	Bucket    string `env:"S3_BUCKET" envDefault:"tracks"`
	UseSSL    bool   `env:"S3_USE_SSL" envDefault:"false"`
}

func (s S3) Validate() error {
	if s.Endpoint == "" {
		return fmt.Errorf("S3_ENDPOINT is required")
	}
	if s.PublicURL == "" {
		return fmt.Errorf("S3_PUBLIC_URL is required")
	}
	if s.AccessKey == "" {
		return fmt.Errorf("S3_ACCESS_KEY is required")
	}
	if s.SecretKey == "" {
		return fmt.Errorf("S3_SECRET_KEY is required")
	}
	if s.Bucket == "" {
		return fmt.Errorf("S3_BUCKET is required")
	}
	return nil
}
