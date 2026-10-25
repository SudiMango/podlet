package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type GithubRepository struct {
	store *db.Queries
}

func NewGithubRepository(store *db.Queries) *GithubRepository {
	return &GithubRepository{
		store: store,
	}
}

func (r *GithubRepository) CreateInstallationStateByUserId(ctx context.Context, userId string, request dto.CreateGithubInstallationStateRequest) (*contract.CreateGithubInstallationStateResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.CreateInstallationStateByUserId(ctx, db.CreateInstallationStateByUserIdParams{
		ID:     request.Id,
		UserID: pgUserId,
		CreatedAt: pgtype.Timestamptz{
			Time:  request.CreatedAt,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  request.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contract.CreateGithubInstallationStateResponse{
		Id:        result.ID,
		UserId:    result.UserID.String(),
		CreatedAt: result.CreatedAt.Time,
		ExpiresAt: result.ExpiresAt.Time,
	}, nil
}
