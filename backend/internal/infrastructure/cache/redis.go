package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis.Client
type RedisClient struct {
	client *redis.Client
}

// NewRedis creates a new Redis client
func NewRedis(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// Get retrieves a value from cache
func (r *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Set stores a value in cache with TTL
func (r *RedisClient) Set(ctx context.Context, key string, value []byte, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

// Delete removes a key from cache
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Exists checks if a key exists in cache
func (r *RedisClient) Exists(ctx context.Context, key string) bool {
	n, err := r.client.Exists(ctx, key).Result()
	return err == nil && n > 0
}

// Client returns the underlying *redis.Client (needed by pubsub.New).
func (r *RedisClient) Client() *redis.Client { return r.client }

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Cache key helpers
func GetHotspotUsersKey(routerID uint, profile string) string {
	if profile == "" {
		return fmt.Sprintf("mikhmon:hotspot:users:%d:all", routerID)
	}
	return fmt.Sprintf("mikhmon:hotspot:users:%d:%s", routerID, profile)
}

func GetHotspotActiveKey(routerID uint) string {
	return fmt.Sprintf("mikhmon:hotspot:active:%d", routerID)
}

func GetHotspotHostsKey(routerID uint) string {
	return fmt.Sprintf("mikhmon:hotspot:hosts:%d", routerID)
}

func GetHotspotProfilesKey(routerID uint) string {
	return fmt.Sprintf("mikhmon:hotspot:profiles:%d", routerID)
}

func GetDashboardResourceKey(routerID uint) string {
	return fmt.Sprintf("mikhmon:dashboard:resource:%d", routerID)
}

func GetDashboardHotspotKey(routerID uint) string {
	return fmt.Sprintf("mikhmon:dashboard:hotspot:%d", routerID)
}

func GetDashboardTrafficKey(routerID uint, iface string) string {
	return fmt.Sprintf("mikhmon:dashboard:traffic:%d:%s", routerID, iface)
}

func GetDashboardLogsKey(routerID uint) string {
	return fmt.Sprintf("mikhmon:dashboard:logs:%d", routerID)
}

func GetLiveReportKey(routerID uint, month string) string {
	return fmt.Sprintf("mikhmon:livereport:%d:%s", routerID, month)
}

func GetSalesReportKey(routerID uint, owner string) string {
	return fmt.Sprintf("mikhmon:salesreport:%d:%s", routerID, owner)
}
