package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	transporthttp "github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http"
	"github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1"
	"github.com/DaniilKalts/music-platform-api/internal/config"
)

type App struct {
	container *Container
	server    *http.Server
}

func NewApp(cfg *config.Config, logger *zap.Logger) (*App, error) {
	c, err := NewContainer(cfg, logger)
	if err != nil {
		return nil, err
	}

	router := transporthttp.NewRouter(
		c.Logger,
		v1.Dependencies{
			AuthService: c.Services.Auth,
			UserService: c.Services.User,
		},
		c.TokenManager,
		c.Caches.Blacklist,
		c.Config.Server.HandlerTimeout,
	)

	return &App{
		container: c,
		server: &http.Server{
			Addr:         c.Config.Server.Addr(),
			Handler:      router,
			ReadTimeout:  c.Config.Server.HTTPTimeout,
			WriteTimeout: c.Config.Server.HTTPTimeout,
			IdleTimeout:  c.Config.Server.HTTPTimeout,
		},
	}, nil
}

func (a *App) Run() error {
	defer a.container.Close()

	logger := a.container.Logger

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	logger.Info("server started", zap.String("addr", a.server.Addr))

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("server: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-sigCtx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.container.Config.Server.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	if err := <-errCh; err != nil {
		return err
	}

	logger.Info("server stopped")

	return nil
}
