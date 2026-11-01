package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/jackc/pgx/v5"
)

const (
	jobPrefix string = "job"
)

var (
	ErrInvalidServerId    = errors.New("Server with given id doesn't exist.")
	ErrInvalidRepoName    = errors.New("Repository with given name doesn't exist. Please use the format `UserName/RepoName`.")
	ErrNoJobsAvailable    = errors.New("No job available.")
	ErrInvalidJobStatus   = errors.New("Job status must be ongoing, success, or failed.")
	ErrInvalidJobFilter   = errors.New("Job status must be pending, ongoing, success, or failed.")
	ErrJobNotFound        = errors.New("Job doesn't exist or cannot be updated.")
	ErrUnsupportedJobType = errors.New("Job type not supported.")
)

type JobService struct {
	jobRepository    *repository.JobRepository
	serverRepository *repository.ServerRepository
	githubService    *GithubService
}

func NewJobService(
	jobRepository *repository.JobRepository,
	serverRepository *repository.ServerRepository,
	githubService *GithubService,
) *JobService {
	return &JobService{
		jobRepository:    jobRepository,
		serverRepository: serverRepository,
		githubService:    githubService,
	}
}

func (s *JobService) CreatePendingJobForServer(ctx context.Context, userId string, request contract.CreateJobRequest) (*contract.CreateJobResponse, error) {

	// Check if server exists with id
	_, err := s.serverRepository.GetServerByUserIdAndServerId(ctx, userId, request.ServerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidServerId
		}

		return nil, err
	}

	// Check if repo exists with name
	githubRepository, err := s.githubService.GetGithubRepositoryByFullName(ctx, userId, request.RepoName)
	if err != nil {
		return nil, err
	}

	// Create job request
	result, err := s.jobRepository.CreatePendingJobForServer(ctx, userId, dto.CreateJobRequest{
		JobId:          util.GenerateID(jobPrefix),
		ServerId:       request.ServerId,
		JobType:        contract.JobTypeDeployNewRepository,
		RepoName:       githubRepository.FullName,
		RepositoryId:   githubRepository.Id,
		InstallationId: githubRepository.InstallationId,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *JobService) ClaimJob(ctx context.Context, serverId string) (*contract.ClaimJobResponse, error) {
	// Get earliest queued job
	job, err := s.jobRepository.GetOldestJobByServerId(ctx, serverId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoJobsAvailable
		}
		return nil, err
	}

	// Get job payload
	payload, err := s.prepareJobPayload(ctx, job)
	if err != nil {
		return nil, err
	}

	// Report job as claimed and return
	err = s.ReportJobByServerIdAndJobId(ctx, serverId, contract.ReportJobRequest{
		JobId:     job.Id,
		JobStatus: contract.StatusOngoing,
		Error:     "",
	})
	if err != nil {
		return nil, err
	}

	return &contract.ClaimJobResponse{
		Id:      job.Id,
		Type:    job.JobType,
		Payload: payload,
	}, nil
}

func (s *JobService) ReportJobByServerIdAndJobId(ctx context.Context, serverId string, request contract.ReportJobRequest) error {
	switch request.JobStatus {
	case contract.StatusOngoing, contract.StatusSuccess, contract.StatusFailed:
	default:
		return ErrInvalidJobStatus
	}

	if err := s.jobRepository.ReportJobByServerIdAndJobId(ctx, serverId, request); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJobNotFound
		}

		return err
	}

	return nil
}

func (s *JobService) GetAllJobsByUserId(ctx context.Context, userId string, status string) (*contract.JobsResponse, error) {
	jobs, err := s.jobRepository.GetAllJobsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return filterJobsByStatus(jobs, status)
}

func (s *JobService) GetAllJobsByServerIdAndUserId(ctx context.Context, userId string, serverId string, status string) (*contract.JobsResponse, error) {
	jobs, err := s.jobRepository.GetAllJobsByServerIdAndUserId(ctx, userId, serverId)
	if err != nil {
		return nil, err
	}

	return filterJobsByStatus(jobs, status)
}

func (s *JobService) GetJobByJobIdAndUserId(ctx context.Context, userId string, jobId string) (*contract.JobResponse, error) {
	job, err := s.jobRepository.GetJobByJobIdAndUserId(ctx, userId, jobId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}

		return nil, err
	}

	return job, nil
}

/*
	Private helpers
*/

func filterJobsByStatus(jobs *contract.JobsResponse, status string) (*contract.JobsResponse, error) {
	if status == "" {
		return jobs, nil
	}

	jobStatus := contract.JobStatus(status)
	switch jobStatus {
	case contract.StatusPending, contract.StatusOngoing, contract.StatusSuccess, contract.StatusFailed:
	default:
		return nil, ErrInvalidJobFilter
	}

	filtered := make(contract.JobsResponse, 0)
	for _, job := range *jobs {
		if job.JobStatus == jobStatus {
			filtered = append(filtered, job)
		}
	}

	return &filtered, nil
}

func (s *JobService) prepareJobPayload(ctx context.Context, job *dto.JobResponse) (json.RawMessage, error) {
	switch job.JobType {
	case contract.JobTypeDeployNewRepository:
		token, err := s.githubService.CreateRepositoryAccessToken(ctx, job.InstallationId, job.RepositoryId)
		if err != nil {
			return nil, err
		}

		return json.Marshal(contract.CloneRepositoryJob{
			CloneUrl:    "https://github.com/" + job.RepoName + ".git",
			GithubToken: token,
		})

	default:
		return nil, ErrUnsupportedJobType
	}
}
