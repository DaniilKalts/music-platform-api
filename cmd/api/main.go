package main

import (
	"flag"
	"log"

	"go.uber.org/zap"

	"github.com/DaniilKalts/music-platform-api/internal/app"
	"github.com/DaniilKalts/music-platform-api/internal/config"
	"github.com/DaniilKalts/music-platform-api/pkg/logger"
)

func main() {
	configPath := flag.String("config-path", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	zapLogger, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()

	a, err := app.NewApp(&cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("build app", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		zapLogger.Fatal("run app", zap.Error(err))
	}
}
