package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
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

func (r *GithubRepository) CreateInstallationStateByUserId(ctx context.Context, userId string, request dto.CreateGithubInstallationStateRequest) (*dto.GithubInstallationState, error) {
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

	return &dto.GithubInstallationState{
		Id:        result.ID,
		UserId:    result.UserID.String(),
		CreatedAt: result.CreatedAt.Time,
		ExpiresAt: result.ExpiresAt.Time,
	}, nil
}

func (r *GithubRepository) GetStateByStateId(ctx context.Context, stateId string) (*dto.GithubInstallationState, error) {
	result, err := r.store.GetStateByStateId(ctx, stateId)
	if err != nil {
		return nil, err
	}

	return &dto.GithubInstallationState{
		Id:        result.ID,
		UserId:    result.UserID.String(),
		CreatedAt: result.CreatedAt.Time,
		ExpiresAt: result.ExpiresAt.Time,
	}, nil
}

func (r *GithubRepository) UpsertGithubInstallation(ctx context.Context, userId string, request dto.UpsertGithubInstallationRequest) (*dto.GithubInstallation, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.UpsertGithubInstallation(ctx, db.UpsertGithubInstallationParams{
		ID:                  request.Id,
		UserID:              pgUserId,
		AccountID:           request.AccountId,
		AccountLogin:        request.AccountLogin,
		AccountType:         request.AccountType,
		RepositorySelection: request.RepositorySelection,
	})
	if err != nil {
		return nil, err
	}

	return &dto.GithubInstallation{
		Id:                  result.ID,
		UserId:              result.UserID.String(),
		AccountId:           result.AccountID,
		AccountLogin:        result.AccountLogin,
		AccountType:         result.AccountType,
		RepositorySelection: result.RepositorySelection,
		CreatedAt:           result.CreatedAt.Time,
		UpdatedAt:           result.UpdatedAt.Time,
	}, nil
}

func (r *GithubRepository) GetGithubInstallationsByUserId(ctx context.Context, userId string) ([]dto.GithubInstallation, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	results, err := r.store.GetGithubInstallationsByUserId(ctx, pgUserId)
	if err != nil {
		return nil, err
	}

	installations := make([]dto.GithubInstallation, 0, len(results))
	for _, result := range results {
		installations = append(installations, dto.GithubInstallation{
			Id:                  result.ID,
			UserId:              result.UserID.String(),
			AccountId:           result.AccountID,
			AccountLogin:        result.AccountLogin,
			AccountType:         result.AccountType,
			RepositorySelection: result.RepositorySelection,
			CreatedAt:           result.CreatedAt.Time,
			UpdatedAt:           result.UpdatedAt.Time,
		})
	}

	return installations, nil
}

func (r *GithubRepository) DeleteGithubInstallationState(ctx context.Context, stateId string) error {
	return r.store.DeleteGithubInstallationState(ctx, stateId)
}
