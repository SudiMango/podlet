package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobRepository struct {
	store *db.Queries
}

func NewJobRepository(store *db.Queries) *JobRepository {
	return &JobRepository{
		store: store,
	}
}

func (r *JobRepository) CreatePendingJobForServer(ctx context.Context, userId string, request dto.CreateJobRequest) (*contract.CreateJobResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.CreatePendingJobForServer(ctx, db.CreatePendingJobForServerParams{
		ID:     request.JobId,
		UserID: pgUserId,
		ServerID: pgtype.Text{
			String: request.ServerId,
			Valid:  true,
		},
		JobType:        db.ServerJobType(request.JobType),
		RepoName:       request.RepoName,
		RepositoryID:   request.RepositoryId,
		InstallationID: request.InstallationId,
	})
	if err != nil {
		return nil, err
	}

	return &contract.CreateJobResponse{
		JobId:  request.JobId,
		Status: contract.JobStatus(result.JobStatus),
	}, nil
}

func (r *JobRepository) GetOldestJobByServerId(ctx context.Context, serverId string) (*contract.JobResponse, error) {
	result, err := r.store.GetOldestJobByServerId(ctx, pgtype.Text{
		String: serverId,
		Valid:  true,
	})
	if err != nil {
		return nil, err
	}

	return &contract.JobResponse{
		ID:             result.ID,
		UserID:         result.UserID.String(),
		ServerID:       result.ServerID.String,
		JobStatus:      contract.JobStatus(result.JobStatus),
		JobType:        contract.JobType(result.JobType),
		RepoName:       result.RepoName,
		RepositoryID:   result.RepositoryID,
		InstallationID: result.InstallationID,
		CreatedAt:      result.CreatedAt.Time,
		UpdatedAt:      result.UpdatedAt.Time,
	}, err
}
