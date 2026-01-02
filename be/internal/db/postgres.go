package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	DatabaseURL string
}

func LoadPostgresConfigFromEnv() PostgresConfig {
	// Prefer DATABASE_URL if provided.
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return PostgresConfig{DatabaseURL: v}
	}

	// Build from PG* vars. Note: user didn't specify credentials; these should be supplied via env.
	host := envOr("PGHOST", "localhost")
	port := envOr("PGPORT", "5432")
	db := envOr("PGDATABASE", "budgeting")
	user := envOr("PGUSER", "postgres")
	pass := envOr("PGPASSWORD", "postgres")

	// If no user specified, produce empty userinfo (still useful to show in errors/README).
	if user == "" {
		return PostgresConfig{DatabaseURL: fmt.Sprintf("postgres://%s:%s/%s?sslmode=disable", host, port, db)}
	}
	if pass == "" {
		return PostgresConfig{DatabaseURL: fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", user, host, port, db)}
	}
	return PostgresConfig{DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)}
}

func (c PostgresConfig) Enabled() bool { return c.DatabaseURL != "" }

func OpenPostgresGorm(ctx context.Context, cfg PostgresConfig) (*gorm.DB, *sql.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, nil, fmt.Errorf("missing DatabaseURL")
	}
	gdb, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	// best-effort connectivity check
	_ = sqlDB.PingContext(ctx)

	return gdb, sqlDB, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
