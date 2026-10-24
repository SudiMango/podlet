package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type ServerRepository struct {
	store *db.Queries
}

func NewServerRepository(store *db.Queries) *ServerRepository {
	return &ServerRepository{
		store: store,
	}
}

func (r *ServerRepository) CreatePendingServerByUserId(ctx context.Context, userId string, request contract.CreatePendingServerRequest) (*contract.PendingServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	done, err := r.store.CreatePendingServerByUserId(ctx, db.CreatePendingServerByUserIdParams{
		ID:                    request.Id,
		UserID:                pgUserId,
		RegistrationTokenHash: request.RegistrationTokenHash,
		CreatedAt: pgtype.Timestamptz{
			Time:  request.CreatedAt,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  request.ExpiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contract.PendingServerResponse{
		Id:                done.ID,
		UserId:            done.UserID.String(),
		CreatedAt:         done.CreatedAt.Time,
		ExpiresAt:         done.ExpiresAt.Time,
		RegistrationToken: nil,
	}, err
}

func (r *ServerRepository) GetPendingServersByUserId(ctx context.Context, userId string) (*contract.PendingServersResponse, error) {
	return nil, nil
}

func (r *ServerRepository) GetPendingServerByPendingServerId(ctx context.Context, pendingServerId string) (*contract.PendingServerResponse, error) {
	return nil, nil
}

func (r *ServerRepository) DeletePendingServerByPendingServerId(ctx context.Context, pendingServerId string) error {
	return nil
}

func (r *ServerRepository) CreateServerByUserId(ctx context.Context, userId string, request contract.CreateServerRequest) (*contract.ServerResponse, error) {
	return nil, nil
}

func (r *ServerRepository) GetServersByUserId(ctx context.Context, userId string) (*contract.ServersResponse, error) {
	return nil, nil
}

func (r *ServerRepository) GetServerByServerId(ctx context.Context, serverId string) (*contract.ServerResponse, error) {
	return nil, nil
}

func (r *ServerRepository) UpdateServerByServerId(ctx context.Context, serverId string, request contract.UpdateServerRequest) (*contract.ServerResponse, error) {
	return nil, nil
}

func (r *ServerRepository) DeleteServerByServerId(ctx context.Context, serverId string) error {
	return nil
}
