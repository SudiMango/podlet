package service

import (
	"context"
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
	ErrInvalidServerId = errors.New("Server with given id doesn't exist.")
	ErrInvalidRepoName = errors.New("Repository with given name doesn't exist. Please use the format `UserName/RepoName`.")
	ErrNoJobsAvailable = errors.New("No job available.")
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
		JobType:        contract.JobTypeCloneRepository,
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
	job, err := s.jobRepository.GetOldestJobByServerId(ctx, serverId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoJobsAvailable
		}
		return nil, err
	}

	githubToken, err := s.githubService.CreateRepositoryAccessToken(ctx, job.InstallationID, job.RepositoryID)
	if err != nil {
		return nil, err
	}

	return &contract.ClaimJobResponse{
		JobId:       job.ID,
		JobType:     job.JobType,
		CloneUrl:    "https://github.com/" + job.RepoName + ".git",
		GithubToken: githubToken,
	}, nil
}
