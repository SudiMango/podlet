package service

import (
	"context"
	"errors"
	"time"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
)

var (
	ErrStateExpired = errors.New("Pending Github installation state has expired.")
)

type GithubService struct {
	githubRepository *repository.GithubRepository
}

func NewGithubService(githubRepository *repository.GithubRepository) *GithubService {
	return &GithubService{
		githubRepository: githubRepository,
	}
}

func (s *GithubService) CreateInstallationStateByUserId(ctx context.Context, userId string) (*contract.CreateGithubInstallationStateResponse, error) {
	createdAt := time.Now()
	token, err := util.GenerateToken()
	if err != nil {
		return nil, err
	}

	result, err := s.githubRepository.CreateInstallationStateByUserId(ctx, userId, dto.CreateGithubInstallationStateRequest{
		Id:        util.HashToken(token),
		CreatedAt: createdAt,
		ExpiresAt: createdAt.Add(10 * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	result.Id = token
	return result, nil
}

func (s *GithubService) GetStateByStateId(ctx context.Context, stateId string) (*dto.GetStateResponse, error) {
	result, err := s.githubRepository.GetStateByStateId(ctx, util.HashToken(stateId))
	if err != nil {
		return nil, err
	}

	if time.Now().After(result.ExpiresAt) {
		return nil, ErrStateExpired
	}

	return result, nil
}

func (s *GithubService) UpsertGithubInstallationByUserId(ctx context.Context, userId string, installation dto.GithubInstallationResponse) (*dto.GithubInstallation, error) {
	return s.githubRepository.UpsertGithubInstallation(ctx, userId, dto.UpsertGithubInstallationRequest{
		Id:                  installation.ID,
		AccountId:           installation.Account.ID,
		AccountLogin:        installation.Account.Login,
		AccountType:         installation.Account.Type,
		RepositorySelection: installation.RepositorySelection,
	})
}

func (s *GithubService) DeleteGithubInstallationStateByStateId(ctx context.Context, stateId string) error {
	return s.githubRepository.DeleteGithubInstallationState(ctx, util.HashToken(stateId))
}
