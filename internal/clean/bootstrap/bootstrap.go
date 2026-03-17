//go:build modern

package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	appcfg "mikhmon_v4/internal/clean/config"
	"mikhmon_v4/internal/clean/domain/entity"
	rediscache "mikhmon_v4/internal/clean/infrastructure/cache/redis"
	"mikhmon_v4/internal/clean/infrastructure/persistence/postgres"
	"mikhmon_v4/internal/clean/usecase/auth"
)

type App struct {
	Config      appcfg.Config
	Logger      *zap.Logger
	DB          *gorm.DB
	RedisClient *redis.Client
	Cache       *rediscache.Cache
	AuthService *auth.Service
}

func New(ctx context.Context, cfg appcfg.Config) (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	if err := db.AutoMigrate(&entity.User{}, &entity.Router{}); err != nil {
		return nil, fmt.Errorf("automigrate: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	userRepo := postgres.NewUserRepo(db)

	return &App{
		Config:      cfg,
		Logger:      logger,
		DB:          db,
		RedisClient: rdb,
		Cache:       rediscache.New(rdb),
		AuthService: auth.NewService(userRepo),
	}, nil
}
