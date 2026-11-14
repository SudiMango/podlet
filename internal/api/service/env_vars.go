package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
)

const (
	envVarPrefix   string = "env"
	maxEnvKeyLen   int    = 64
	maxEnvValueLen int    = 4096
)

var (
	ErrInvalidEnvKey   = errors.New("Env var key must start with a letter or underscore, contain only letters, numbers, and underscores, and be at most 64 characters.")
	ErrInvalidEnvValue = errors.New("Env var value must be between 1 and 4096 characters.")
	ErrEnvVarNotFound  = errors.New("Env var doesn't exist.")

	envKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

type EnvVarService struct {
	envVarRepository *repository.EnvVarRepository
	githubService    *GithubService
	gcm              cipher.AEAD
}

func NewEnvVarService(
	envVarRepository *repository.EnvVarRepository,
	githubService *GithubService,
) (*EnvVarService, error) {
	encoded := os.Getenv("ENV_ENCRYPTION_KEY")
	if encoded == "" {
		return nil, fmt.Errorf("ENV_ENCRYPTION_KEY is required")
	}

	key, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("ENV_ENCRYPTION_KEY must be 32 bytes, base64 encoded")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &EnvVarService{
		envVarRepository: envVarRepository,
		githubService:    githubService,
		gcm:              gcm,
	}, nil
}

func (s *EnvVarService) CreateEnvVarByUserIdAndRepoName(ctx context.Context, userId string, request contract.CreateEnvVarRequest) (*contract.EnvVarResponse, error) {
	if len(request.Key) > maxEnvKeyLen || !envKeyPattern.MatchString(request.Key) {
		return nil, ErrInvalidEnvKey
	}
	if request.Value == "" || len(request.Value) > maxEnvValueLen {
		return nil, ErrInvalidEnvValue
	}

	repo, err := s.githubService.GetGithubRepositoryByFullName(ctx, userId, request.RepoName)
	if err != nil {
		return nil, err
	}

	encrypted, err := s.encrypt(request.Value)
	if err != nil {
		return nil, err
	}

	return s.envVarRepository.CreateEnvVarByUserIdAndRepoId(ctx, userId, dto.CreateEnvVarRequest{
		Id:             util.GenerateID(envVarPrefix),
		Key:            request.Key,
		ValueEncrypted: encrypted,
		RepositoryId:   repo.Id,
		RepoName:       repo.FullName,
	})
}

func (s *EnvVarService) GetEnvVarsByUserIdAndRepoName(ctx context.Context, userId string, repoName string) (*contract.EnvVarsResponse, error) {
	repo, err := s.githubService.GetGithubRepositoryByFullName(ctx, userId, repoName)
	if err != nil {
		return nil, err
	}

	return s.envVarRepository.GetEnvVarsByUserIdAndRepoId(ctx, userId, repoName, repo.Id)
}

func (s *EnvVarService) GetDecryptedEnvVarsByUserIdAndRepoId(ctx context.Context, userId string, repositoryId int64) (map[string]string, error) {
	rows, err := s.envVarRepository.GetEnvVarsByUserIdAndRepoIdForAgent(ctx, userId, repositoryId)
	if err != nil {
		return nil, err
	}

	envVars := make(map[string]string, len(rows))
	for _, row := range rows {
		value, err := s.decrypt(row.Value)
		if err != nil {
			return nil, err
		}

		envVars[row.Key] = value
	}

	return envVars, nil
}

func (s *EnvVarService) DeleteEnvVarByUserIdAndRepoName(ctx context.Context, userId string, repoName string, key string) error {
	repo, err := s.githubService.GetGithubRepositoryByFullName(ctx, userId, repoName)
	if err != nil {
		return err
	}

	rowsAffected, err := s.envVarRepository.DeleteEnvVarByUserIdAndRepoId(ctx, userId, repo.Id, key)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEnvVarNotFound
	}

	return nil
}

func (s *EnvVarService) DeleteAllEnvVarsByUserIdAndRepoName(ctx context.Context, userId string, repoName string) error {
	repo, err := s.githubService.GetGithubRepositoryByFullName(ctx, userId, repoName)
	if err != nil {
		return err
	}

	return s.envVarRepository.DeleteAllEnvVarsByUserIdAndRepoId(ctx, userId, repo.Id)
}

/*
	Private helpers
*/

func (s *EnvVarService) encrypt(plaintext string) (string, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	sealed := s.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (s *EnvVarService) decrypt(encoded string) (string, error) {
	sealed, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	nonceSize := s.gcm.NonceSize()
	if len(sealed) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	plaintext, err := s.gcm.Open(nil, sealed[:nonceSize], sealed[nonceSize:], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
