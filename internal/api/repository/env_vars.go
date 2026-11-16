package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type EnvVarRepository struct {
	store *db.Queries
}

func NewEnvVarRepository(store *db.Queries) *EnvVarRepository {
	return &EnvVarRepository{
		store: store,
	}
}

func (r *EnvVarRepository) CreateEnvVarByUserIdAndRepoId(ctx context.Context, userId string, request dto.CreateEnvVarRequest) (*contract.EnvVarResponse, error) {
	var pgUserId pgtype.UUID
	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.CreateEnvVarByUserIdAndRepoId(ctx, db.CreateEnvVarByUserIdAndRepoIdParams{
		ID:             request.Id,
		Key:            request.Key,
		ValueEncrypted: request.ValueEncrypted,
		UserID:         pgUserId,
		RepositoryID:   request.RepositoryId,
	})
	if err != nil {
		return nil, err
	}

	return &contract.EnvVarResponse{
		Id:        result.ID,
		Key:       result.Key,
		RepoName:  request.RepoName,
		CreatedAt: result.CreatedAt.Time,
	}, nil
}

func (r *EnvVarRepository) GetEnvVarsByUserIdAndRepoId(ctx context.Context, userId string, repoName string, repoId int64) (*contract.EnvVarsResponse, error) {
	var pgUserId pgtype.UUID
	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.GetEnvVarsByUserIdAndRepoId(ctx, db.GetEnvVarsByUserIdAndRepoIdParams{
		UserID:       pgUserId,
		RepositoryID: repoId,
	})
	if err != nil {
		return nil, err
	}

	envVars := make(contract.EnvVarsResponse, 0, len(result))
	for _, envVar := range result {
		envVars = append(envVars, contract.EnvVarResponse{
			Id:        envVar.ID,
			Key:       envVar.Key,
			RepoName:  repoName,
			CreatedAt: envVar.CreatedAt.Time,
		})
	}

	return &envVars, nil
}

func (r *EnvVarRepository) GetEnvVarsByUserIdAndRepoIdForAgent(ctx context.Context, userId string, repositoryId int64) (contract.EnvVarsForAgentResponse, error) {
	var pgUserId pgtype.UUID
	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.GetEnvVarsByUserIdAndRepoIdForAgent(ctx, db.GetEnvVarsByUserIdAndRepoIdForAgentParams{
		UserID:       pgUserId,
		RepositoryID: repositoryId,
	})
	if err != nil {
		return nil, err
	}

	envVars := make(contract.EnvVarsForAgentResponse, 0, len(result))
	for _, envVar := range result {
		envVars = append(envVars, contract.EnvVarForAgentResponse{
			Key:   envVar.Key,
			Value: envVar.ValueEncrypted,
		})
	}

	return envVars, nil
}

func (r *EnvVarRepository) DeleteEnvVarByUserIdAndRepoId(ctx context.Context, userId string, repoId int64, key string) (int64, error) {
	var pgUserId pgtype.UUID
	if err := pgUserId.Scan(userId); err != nil {
		return 0, err
	}

	return r.store.DeleteEnvVarByUserIdAndRepoId(ctx, db.DeleteEnvVarByUserIdAndRepoIdParams{
		UserID:       pgUserId,
		RepositoryID: repoId,
		Key:          key,
	})
}

func (r *EnvVarRepository) DeleteAllEnvVarsByUserIdAndRepoId(ctx context.Context, userId string, repoId int64) error {
	var pgUserId pgtype.UUID
	if err := pgUserId.Scan(userId); err != nil {
		return err
	}

	return r.store.DeleteAllEnvVarsByUserIdAndRepoId(ctx, db.DeleteAllEnvVarsByUserIdAndRepoIdParams{
		UserID:       pgUserId,
		RepositoryID: repoId,
	})
}
