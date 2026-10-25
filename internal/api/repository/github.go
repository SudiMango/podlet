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

func (r *GithubRepository) GetStateByStateId(ctx context.Context, stateId string) (*dto.GetStateResponse, error) {
	result, err := r.store.GetStateByStateId(ctx, stateId)
	if err != nil {
		return nil, err
	}

	return &dto.GetStateResponse{
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

func (r *GithubRepository) DeleteGithubInstallationState(ctx context.Context, stateId string) error {
	return r.store.DeleteGithubInstallationState(ctx, stateId)
}
