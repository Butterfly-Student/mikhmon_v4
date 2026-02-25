package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Redis        RedisConfig
	JWT          JWTConfig
	InternalWSKey string // Internal key for WebSocket authentication
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads configuration from environment variables or config file
func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set defaults
	setDefaults()

	// Read from environment variables
	viper.AutomaticEnv()

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is OK, we'll use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  viper.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetInt("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Database:        viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSL_MODE"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			AccessTokenTTL:  viper.GetDuration("JWT_ACCESS_TOKEN_TTL"),
			RefreshTokenTTL: viper.GetDuration("JWT_REFRESH_TOKEN_TTL"),
		},
		InternalWSKey: viper.GetString("INTERNAL_WS_KEY"),
	}

	return cfg, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_PORT", getEnv("SERVER_PORT", "8080"))
	viper.SetDefault("SERVER_READ_TIMEOUT", 30*time.Second)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 30*time.Second)
	viper.SetDefault("SERVER_IDLE_TIMEOUT", 120*time.Second)

	// Database defaults
	viper.SetDefault("DB_HOST", getEnv("DB_HOST", "localhost"))
	viper.SetDefault("DB_PORT", getEnvAsInt("DB_PORT", 5432))
	viper.SetDefault("DB_USER", getEnv("DB_USER", "mikhmon"))
	viper.SetDefault("DB_PASSWORD", getEnv("DB_PASSWORD", "mikhmon"))
	viper.SetDefault("DB_NAME", getEnv("DB_NAME", "mikhmon"))
	viper.SetDefault("DB_SSL_MODE", getEnv("DB_SSL_MODE", "disable"))
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Redis defaults
	viper.SetDefault("REDIS_HOST", getEnv("REDIS_HOST", "localhost"))
	viper.SetDefault("REDIS_PORT", getEnvAsInt("REDIS_PORT", 6379))
	viper.SetDefault("REDIS_PASSWORD", getEnv("REDIS_PASSWORD", ""))
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_POOL_SIZE", 10)

	// JWT defaults
	viper.SetDefault("JWT_SECRET", getEnv("JWT_SECRET", "your-secret-key-change-in-production"))
	viper.SetDefault("JWT_ACCESS_TOKEN_TTL", 24*time.Hour)
	viper.SetDefault("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour)

	// Internal WebSocket key
	viper.SetDefault("INTERNAL_WS_KEY", getEnv("INTERNAL_WS_KEY", "mikhmon-ws-internal-key"))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// DSN returns PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// Addr returns Redis address
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
