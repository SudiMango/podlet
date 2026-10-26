package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/github"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
)

var (
	ErrStateExpired               = errors.New("Pending Github installation state has expired")
	ErrInvalidInstallation        = errors.New("Installation does not belong to this GitHub App")
	ErrInstallationSuspended      = errors.New("Github installation is suspended")
	ErrContentsPermissionRequired = errors.New("Installation does not grant repository contents access")
)

type GithubService struct {
	githubRepository *repository.GithubRepository
	githubClient     *github.Client
	appID            int64
	appName          string
}

func NewGithubService(githubRepository *repository.GithubRepository, githubClient *github.Client) (*GithubService, error) {
	appID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		return nil, errors.New("invalid GITHUB_APP_ID")
	}

	appName := os.Getenv("GITHUB_APP_NAME")
	if appName == "" {
		return nil, errors.New("GITHUB_APP_NAME is required")
	}

	return &GithubService{
		githubRepository: githubRepository,
		githubClient:     githubClient,
		appID:            appID,
		appName:          appName,
	}, nil
}

func (s *GithubService) CreateInstallationStateByUserId(ctx context.Context, userId string) (*dto.GithubInstallationState, error) {
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

func (s *GithubService) VerifyInstallation(ctx context.Context, stateID string, installationId int64) (*dto.GithubInstallation, error) {
	state, err := s.GetStateByStateId(ctx, stateID)
	if err != nil {
		return nil, err
	}

	// Get installation and verify it gives proper perms to correct app
	installation, err := s.githubClient.GetInstallation(ctx, installationId)
	if err != nil {
		return nil, err
	}

	if installation.Id != installationId || installation.AppId != s.appID || installation.AppSlug != s.appName {
		return nil, ErrInvalidInstallation
	}

	if installation.SuspendedAt != nil {
		return nil, ErrInstallationSuspended
	}

	if installation.Permissions.Contents != "read" && installation.Permissions.Contents != "write" {
		return nil, ErrContentsPermissionRequired
	}

	// Save install info
	saved, err := s.UpsertGithubInstallationByUserId(ctx, state.UserId, *installation)
	if err != nil {
		return nil, err
	}

	// Delete old state
	if err := s.DeleteGithubInstallationStateByStateId(ctx, stateID); err != nil {
		return nil, err
	}

	return saved, nil
}

func (s *GithubService) CreateRepositoryAccessToken(ctx context.Context, installationId int64, repositoryId int64) (string, error) {
	token, err := s.githubClient.CreateInstallationToken(ctx, installationId, []int64{repositoryId})
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func (s *GithubService) GetStateByStateId(ctx context.Context, stateId string) (*dto.GithubInstallationState, error) {
	result, err := s.githubRepository.GetStateByStateId(ctx, util.HashToken(stateId))
	if err != nil {
		return nil, err
	}

	if time.Now().After(result.ExpiresAt) {
		return nil, ErrStateExpired
	}

	return result, nil
}

func (s *GithubService) UpsertGithubInstallationByUserId(ctx context.Context, userId string, installation github.Installation) (*dto.GithubInstallation, error) {
	return s.githubRepository.UpsertGithubInstallation(ctx, userId, dto.UpsertGithubInstallationRequest{
		Id:                  installation.Id,
		AccountId:           installation.Account.Id,
		AccountLogin:        installation.Account.Login,
		AccountType:         installation.Account.Type,
		RepositorySelection: installation.RepositorySelection,
	})
}

func (s *GithubService) GetGithubRepositoriesByUserId(ctx context.Context, userId string) (contract.GithubRepositoriesResponse, error) {
	installations, err := s.githubRepository.GetGithubInstallationsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := make(contract.GithubRepositoriesResponse, 0)
	seen := make(map[int64]struct{})

	for _, installation := range installations {
		token, err := s.githubClient.CreateInstallationToken(ctx, installation.Id, nil)
		if err != nil {
			return nil, err
		}

		result, err := s.githubClient.ListRepositories(ctx, token.Token)
		if err != nil {
			return nil, err
		}

		for _, repository := range result.Repositories {
			if _, exists := seen[repository.Id]; exists {
				continue
			}

			seen[repository.Id] = struct{}{}
			response = append(response, contract.GithubRepositoryResponse{
				Id:             repository.Id,
				InstallationId: installation.Id,
				Name:           repository.Name,
				FullName:       repository.FullName,
				Private:        repository.Private,
				HtmlUrl:        repository.HtmlUrl,
				CloneUrl:       repository.CloneUrl,
				DefaultBranch:  repository.DefaultBranch,
			})
		}
	}

	return response, nil
}

func (s *GithubService) GetGithubRepositoryByFullName(ctx context.Context, userId string, fullName string) (*contract.GithubRepositoryResponse, error) {
	owner, name, found := strings.Cut(fullName, "/")
	if !found || owner == "" || name == "" || strings.Contains(name, "/") {
		return nil, ErrInvalidRepoName
	}

	repositories, err := s.GetGithubRepositoriesByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, repository := range repositories {
		if strings.EqualFold(repository.FullName, fullName) {
			return &repository, nil
		}
	}

	return nil, ErrInvalidRepoName
}

func (s *GithubService) DeleteGithubInstallationStateByStateId(ctx context.Context, stateId string) error {
	return s.githubRepository.DeleteGithubInstallationState(ctx, util.HashToken(stateId))
}
