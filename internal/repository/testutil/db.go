package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres"
	"github.com/DaniilKalts/music-platform-api/internal/config"
)

func SetupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	// Находим корень проекта для миграций
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../..")
	
	migrationsPath := filepath.Join(root, "database/migrations")
	os.Setenv("MIGRATIONS_DIR", migrationsPath)

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/music_platform_test?sslmode=disable"
	}

	cfg := &config.Postgres{
		Host:             "localhost",
		Port:             5432,
		User:             "postgres",
		Password:         "postgres",
		Name:             "music_platform_test",
		SSLMode:          "disable",
		MaxConns:         10,
		MinConns:         1,
		MaxConnLifetime:  time.Hour,
		MaxConnIdleTime:  time.Minute * 30,
		StatementTimeout: time.Second * 3,
	}

	ctx := context.Background()
	pool, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		t.Skipf("Skipping integration test: failed to connect to test database. Ensure 'music_platform_test' exists and Postgres is running on 5432 with password 'postgres'. Error: %v", err)
		return nil, func() {}
	}

	cleanup := func() {
		TruncateTables(t, pool)
		pool.Close()
	}

	return pool, cleanup
}

func TruncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	if pool == nil {
		return
	}
	tables := []string{
		"listening_history",
		"favorites",
		"playlist_tracks",
		"playlists",
		"tracks",
		"albums",
		"artists",
		"genres",
		"users",
	}

	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))
	_, err := pool.Exec(context.Background(), query)
	require.NoError(t, err, "Failed to truncate tables")
}
