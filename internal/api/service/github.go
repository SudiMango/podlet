package service

import (
	"context"
	"time"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/contract"
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
	result, err := s.githubRepository.CreateInstallationStateByUserId(ctx, userId, dto.CreateGithubInstallationStateRequest{
		Id:        userId,
		CreatedAt: createdAt,
		ExpiresAt: createdAt.Add(10 * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
