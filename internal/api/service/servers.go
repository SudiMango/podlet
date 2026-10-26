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

/*
	Pending server functions
*/

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

func (s *ServerService) GetValidPendingServersByUserId(ctx context.Context, userId string) (*contract.PendingServersResponse, error) {
	return s.serverRepository.GetValidPendingServersByUserId(ctx, userId)
}

func (s *ServerService) GetValidPendingServerByUserIdAndPendingServerId(ctx context.Context, userId string, pendingServerId string) (*contract.PendingServerResponse, error) {
	return s.serverRepository.GetValidPendingServerByUserIdAndPendingServerId(ctx, userId, pendingServerId)
}

func (s *ServerService) DeletePendingServerByUserIdAndPendingServerId(ctx context.Context, userId string, pendingServerId string) error {
	return s.serverRepository.DeletePendingServerByUserIdAndPendingServerId(ctx, userId, pendingServerId)
}

/*
	Server functions
*/

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

func (s *ServerService) UpdateServerLastSeenAtByUserIdAndServerId(ctx context.Context, userId string, serverId string) (*contract.ServerLifecycleState, error) {
	lifecycleState, err := s.serverRepository.UpdateServerLastSeenAtByUserIdAndServerId(ctx, userId, serverId)
	if err != nil {
		return nil, err
	}

	return lifecycleState, nil
}

func (s *ServerService) UpdateServerByUserIdAndServerId(ctx context.Context, userId string, serverId string, request contract.UpdateServerRequest) (*contract.ServerResponse, error) {
	if request.Name == nil {
		return s.serverRepository.GetServerByUserIdAndServerId(ctx, userId, serverId)
	}

	return s.serverRepository.UpdateServerByUserIdAndServerId(ctx, userId, serverId, *request.Name)
}

func (s *ServerService) DeleteServerByUserIdAndServerId(ctx context.Context, userId string, serverId string) error {
	return s.serverRepository.DeleteServerByUserIdAndServerId(ctx, userId, serverId)
}

func (s *ServerService) MarkServerForDeletionByUserIdAndServerId(ctx context.Context, userId string, serverId string) error {
	return s.serverRepository.MarkServerForDeletionByUserIdAndServerId(ctx, userId, serverId)
}
