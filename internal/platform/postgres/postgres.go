package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New Pool Creates a connection pool to the database at dsn and checks that the database is rechable before returning the pool. The pool is configured with a minimum of 2 connections and a maximum of 10 connections.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres:parse dsn : %w", err)
	}
	// Pool Sizing , little bit of warm connections and a hard limit to avoid overloading the database
	cfg.MinConns = 2
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnLifetimeJitter = 5 * time.Minute
	cfg.MaxConnIdleTime = 15 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: create pool : %w", err)
	}
	// The pool connection is lazy, so we need to ping the database to ensure that the connection is established and the database is reachable.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("postgres: ping database : %w", err)
	}
	return pool, nil
}
