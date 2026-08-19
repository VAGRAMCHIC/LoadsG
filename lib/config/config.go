package config

import (
	"os"
)

type Config struct {
	Address string // http://127.0.0.1:8080
	Token   string // access token
	Format  string // table | json | yaml
}

// Load читает env. Флаги cobra поверх этого уже перезаписывают значения.
func Load() *Config {
	return &Config{
		Address: envOr("LOADSG_ADDR", "http://127.0.0.1:8080"),
		Token:   os.Getenv("LOADSG_TOKEN"),
		Format:  envOr("LOADSG_FORMAT", "table"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}