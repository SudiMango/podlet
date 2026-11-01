package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5"
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
		ID:             request.JobId,
		UserID:         pgUserId,
		ServerID:       request.ServerId,
		JobType:        string(request.JobType),
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

func (r *JobRepository) GetOldestJobByServerId(ctx context.Context, serverId string) (*dto.JobResponse, error) {
	result, err := r.store.GetOldestJobByServerId(ctx, serverId)
	if err != nil {
		return nil, err
	}

	return &dto.JobResponse{
		Id:             result.ID,
		UserId:         result.UserID.String(),
		ServerId:       result.ServerID,
		JobStatus:      contract.JobStatus(result.JobStatus),
		ErrorMessage:   &result.ErrorMessage.String,
		JobType:        contract.JobType(result.JobType),
		RepoName:       result.RepoName,
		RepositoryId:   result.RepositoryID,
		InstallationId: result.InstallationID,
		CreatedAt:      result.CreatedAt.Time,
		UpdatedAt:      result.UpdatedAt.Time,
	}, err
}

func (r *JobRepository) ReportJobByServerIdAndJobId(ctx context.Context, serverId string, request contract.ReportJobRequest) error {
	rowsAffected, err := r.store.ReportJobByServerIdAndJobId(ctx, db.ReportJobByServerIdAndJobIdParams{
		ID:        request.JobId,
		ServerID:  serverId,
		JobStatus: db.ServerJobStatus(request.JobStatus),
		ErrorMessage: pgtype.Text{
			String: request.Error,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *JobRepository) GetAllJobsByUserId(ctx context.Context, userId string) (*contract.JobsResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.GetAllJobsByUserId(ctx, pgUserId)
	if err != nil {
		return nil, err
	}

	jobs := make(contract.JobsResponse, 0, len(result))
	for _, job := range result {
		var errorMessage *string
		if job.ErrorMessage.Valid {
			errorMessage = &job.ErrorMessage.String
		}

		j := contract.JobResponse{
			Id:           job.ID,
			ServerId:     job.ServerID,
			JobStatus:    contract.JobStatus(job.JobStatus),
			ErrorMessage: errorMessage,
			JobType:      contract.JobType(job.JobType),
			RepoName:     job.RepoName,
			CreatedAt:    job.CreatedAt.Time,
			UpdatedAt:    job.UpdatedAt.Time,
		}

		jobs = append(jobs, j)
	}

	return &jobs, err
}

func (r *JobRepository) GetAllJobsByServerIdAndUserId(ctx context.Context, userId string, serverId string) (*contract.JobsResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.GetAllJobsByServerIdAndUserId(ctx, db.GetAllJobsByServerIdAndUserIdParams{
		ServerID: serverId,
		UserID:   pgUserId,
	})
	if err != nil {
		return nil, err
	}

	jobs := make(contract.JobsResponse, 0, len(result))
	for _, job := range result {
		var errorMessage *string
		if job.ErrorMessage.Valid {
			errorMessage = &job.ErrorMessage.String
		}

		j := contract.JobResponse{
			Id:           job.ID,
			ServerId:     job.ServerID,
			JobStatus:    contract.JobStatus(job.JobStatus),
			ErrorMessage: errorMessage,
			JobType:      contract.JobType(job.JobType),
			RepoName:     job.RepoName,
			CreatedAt:    job.CreatedAt.Time,
			UpdatedAt:    job.UpdatedAt.Time,
		}

		jobs = append(jobs, j)
	}

	return &jobs, err
}

func (r *JobRepository) GetJobByJobIdAndUserId(ctx context.Context, userId string, jobId string) (*contract.JobResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	job, err := r.store.GetJobByJobIdAndUserId(ctx, db.GetJobByJobIdAndUserIdParams{
		ID:     jobId,
		UserID: pgUserId,
	})
	if err != nil {
		return nil, err
	}

	var errorMessage *string
	if job.ErrorMessage.Valid {
		errorMessage = &job.ErrorMessage.String
	}

	return &contract.JobResponse{
		Id:           job.ID,
		ServerId:     job.ServerID,
		JobStatus:    contract.JobStatus(job.JobStatus),
		ErrorMessage: errorMessage,
		JobType:      contract.JobType(job.JobType),
		RepoName:     job.RepoName,
		CreatedAt:    job.CreatedAt.Time,
		UpdatedAt:    job.UpdatedAt.Time,
	}, err
}
