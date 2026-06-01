package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres"
	"github.com/DaniilKalts/music-platform-api/internal/config"
	"github.com/DaniilKalts/music-platform-api/internal/repository"
	"github.com/DaniilKalts/music-platform-api/internal/service"
	jwtpkg "github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

type Container struct {
	Config *config.Config
	Logger *zap.Logger

	DB *pgxpool.Pool

	Repositories *repository.Repositories
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

	repositories := repository.NewRepositories(db)
	tokenManager := jwtpkg.NewManager(
		[]byte(cfg.JWT.AccessSecret),
		[]byte(cfg.JWT.RefreshSecret),
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	services := service.NewServices(repositories, tokenManager)

	return &Container{
		Config:       cfg,
		Logger:       logger,
		DB:           db,
		Repositories: repositories,
		TokenManager: tokenManager,
		Services:     services,
	}, nil
}

func (c *Container) Close() {
	c.DB.Close()
}
