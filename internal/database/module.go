package database

import (
	"context"
	"fmt"

	"github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"database",
	fx.Provide(
		providePool,
		provideQueries,
	),
)

func providePool(lifecycle fx.Lifecycle) (*pgxpool.Pool, error) {
	databaseURL, err := connectionString()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				return fmt.Errorf("connect to database: %w", err)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func provideQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
