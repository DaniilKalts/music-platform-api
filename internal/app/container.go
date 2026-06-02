package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	redisadapter "github.com/DaniilKalts/music-platform-api/internal/adapter/cache/redis"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres"
	"github.com/DaniilKalts/music-platform-api/internal/cache"
	"github.com/DaniilKalts/music-platform-api/internal/config"
	"github.com/DaniilKalts/music-platform-api/internal/repository"
	"github.com/DaniilKalts/music-platform-api/internal/service"
	jwtpkg "github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

type Container struct {
	Config *config.Config
	Logger *zap.Logger

	DB    *pgxpool.Pool
	Redis *redis.Client

	Repositories *repository.Repositories
	Caches       *cache.Caches
	TokenManager *jwtpkg.Manager
	Services     *service.Services
}

func NewContainer(cfg *config.Config, logger *zap.Logger) (_ *Container, err error) {
	ctx := context.Background()

	db, err := postgres.NewClient(ctx, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	redisClient, err := redisadapter.NewClient(ctx, &cfg.Redis)

	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	defer func() {
		if err != nil {
			_ = redisClient.Close()
		}
	}()

	repositories := repository.NewRepositories(db)
	caches := cache.NewCaches(redisClient)
	tokenManager := jwtpkg.NewManager(
		[]byte(cfg.JWT.AccessSecret),
		[]byte(cfg.JWT.RefreshSecret),
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	services := service.NewServices(repositories, tokenManager, caches.Blacklist, caches.Refresh)

	return &Container{
		Config:       cfg,
		Logger:       logger,
		DB:           db,
		Redis:        redisClient,
		Repositories: repositories,
		Caches:       caches,
		TokenManager: tokenManager,
		Services:     services,
	}, nil
}

func (c *Container) Close() {
	_ = c.Redis.Close()
	c.DB.Close()
}
