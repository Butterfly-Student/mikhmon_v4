//go:build modern

package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName  string
	HTTPPort string
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func Load() Config {
	return Config{
		AppName:  env("APP_NAME", "mikhmon-modern"),
		HTTPPort: env("HTTP_PORT", "8080"),
		Postgres: PostgresConfig{
			Host:     env("PG_HOST", "127.0.0.1"),
			Port:     env("PG_PORT", "5432"),
			User:     env("PG_USER", "postgres"),
			Password: env("PG_PASSWORD", "postgres"),
			DBName:   env("PG_DB", "mikhmon"),
			SSLMode:  env("PG_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     env("REDIS_ADDR", "127.0.0.1:6379"),
			Password: env("REDIS_PASSWORD", ""),
			DB:       0,
		},
	}
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
