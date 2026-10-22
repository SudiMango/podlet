package repository

import (
	"context"
	"fmt"

	db "github.com/SudiMango/podlet/internal/sql/database"
)

type HealthRepository struct {
	store *db.Queries
}

func NewHealthRepository(store *db.Queries) *HealthRepository {
	return &HealthRepository{
		store: store,
	}
}

func (r *HealthRepository) IsDatabaseHealthy(ctx context.Context) (bool, error) {
	result, err := r.store.PingDatabase(ctx)
	if err != nil {
		return false, fmt.Errorf("ping database: %w", err)
	}

	return result == 1, nil
}
