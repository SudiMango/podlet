package service

import (
	"context"
	"time"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
)

const (
	pendingServerIdPrefix string = "pnd"
	serverIdPrefix        string = "srv"
)

type ServerService struct {
	serverRepository *repository.ServerRepository
}

func NewServerService(serverRepository *repository.ServerRepository) *ServerService {
	return &ServerService{
		serverRepository: serverRepository,
	}
}

func (s *ServerService) CreatePendingServerByUserId(ctx context.Context, userId string) (*contract.PendingServerResponse, error) {
	token, err := util.GenerateToken()
	if err != nil {
		return nil, err
	}

	tokenHash := util.HashToken(token)

	createdAtTime := time.Now()
	request := dto.CreatePendingServerRequest{
		Id:                    util.GenerateID(pendingServerIdPrefix),
		CreatedAt:             createdAtTime,
		ExpiresAt:             createdAtTime.Add(10 * time.Minute),
		RegistrationTokenHash: tokenHash,
	}

	done, err := s.serverRepository.CreatePendingServerByUserId(ctx, userId, request)
	if err != nil {
		return nil, err
	}

	done.RegistrationToken = &token
	return done, nil
}

func (s *ServerService) CreateServerByUserId(ctx context.Context, request contract.CreateServerRequest) (*contract.ServerResponse, error) {

	// Check if the user owns this pending server and if it is not expired
	pendingServer, err := s.serverRepository.GetValidPendingServerRegistrationTokenHash(ctx, util.HashToken(request.RegistrationToken))
	if err != nil {
		return nil, err
	}

	// Register new server for user
	token, err := util.GenerateToken()
	if err != nil {
		return nil, err
	}

	tokenHash := util.HashToken(token)

	createdAtTime := time.Now()
	requestDb := dto.CreateServerRequest{
		Id:        util.GenerateID(serverIdPrefix),
		Name:      request.Name,
		CreatedAt: createdAtTime,
		UpdatedAt: createdAtTime,
		TokenHash: tokenHash,
	}

	done, err := s.serverRepository.CreateServerByUserId(ctx, pendingServer.UserId, requestDb)
	if err != nil {
		return nil, err
	}
	done.Token = &token

	err = s.serverRepository.DeletePendingServerByUserIdAndPendingServerId(ctx, pendingServer.UserId, pendingServer.Id)
	if err != nil {
		return done, err
	}

	return done, nil
}

func (s *ServerService) GetServersByUserId(ctx context.Context, userId string) (*contract.ServersResponse, error) {
	return s.serverRepository.GetServersByUserId(ctx, userId)
}

func (s *ServerService) GetServerByUserIdAndServerId(ctx context.Context, userId string, serverId string) (*contract.ServerResponse, error) {
	return s.serverRepository.GetServerByUserIdAndServerId(ctx, userId, serverId)
}
