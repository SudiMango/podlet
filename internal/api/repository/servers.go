package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
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

func (r *ServerRepository) CreatePendingServerByUserId(ctx context.Context, userId string, request dto.CreatePendingServerRequest) (*contract.PendingServerResponse, error) {
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

func (r *ServerRepository) GetValidPendingServersByUserId(ctx context.Context, userId string) (*contract.PendingServersResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *ServerRepository) GetValidPendingServerByUserIdAndPendingServerId(ctx context.Context, userId string, pendingServerId string) (*contract.PendingServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *ServerRepository) GetValidPendingServerRegistrationTokenHash(ctx context.Context, registrationTokenHash string) (*contract.PendingServerResponse, error) {
	done, err := r.store.GetValidPendingServerRegistrationTokenHash(ctx, registrationTokenHash)
	if err != nil {
		return nil, err
	}

	return &contract.PendingServerResponse{
		Id:                done.ID,
		UserId:            done.UserID.String(),
		CreatedAt:         done.CreatedAt.Time,
		ExpiresAt:         done.ExpiresAt.Time,
		RegistrationToken: &done.RegistrationTokenHash,
	}, nil
}

func (r *ServerRepository) DeletePendingServerByUserIdAndPendingServerId(ctx context.Context, userId string, pendingServerId string) error {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return err
	}

	return r.store.DeletePendingServerByUserIdAndPendingServerId(ctx, db.DeletePendingServerByUserIdAndPendingServerIdParams{
		UserID: pgUserId,
		ID:     pendingServerId,
	})
}

func (r *ServerRepository) CreateServerByUserId(ctx context.Context, userId string, request dto.CreateServerRequest) (*contract.ServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	done, err := r.store.CreateServerByUserId(ctx, db.CreateServerByUserIdParams{
		ID:        request.Id,
		Name:      request.Name,
		UserID:    pgUserId,
		TokenHash: request.TokenHash,
		CreatedAt: pgtype.Timestamptz{
			Time:  request.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  request.UpdatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contract.ServerResponse{
		Id:        done.ID,
		Name:      done.Name,
		UserId:    done.UserID.String(),
		CreatedAt: done.CreatedAt.Time,
		UpdatedAt: done.UpdatedAt.Time,
		Token:     nil,
	}, nil
}

func (r *ServerRepository) GetServersByUserId(ctx context.Context, userId string) (*contract.ServersResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	done, err := r.store.GetServersByUserId(ctx, pgUserId)
	if err != nil {
		return nil, err
	}

	var response contract.ServersResponse
	for _, s := range done {
		sTemp := contract.ServerResponse{
			Id:        s.ID,
			Name:      s.Name,
			UserId:    s.UserID.String(),
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
			Token:     nil,
		}
		response = append(response, sTemp)
	}

	return &response, nil
}

func (r *ServerRepository) GetServerByUserIdAndServerId(ctx context.Context, userId string, serverId string) (*contract.ServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	done, err := r.store.GetServerByUserIdAndServerId(ctx, db.GetServerByUserIdAndServerIdParams{
		UserID: pgUserId,
		ID:     serverId,
	})
	if err != nil {
		return nil, err
	}

	return &contract.ServerResponse{
		Id:        done.ID,
		Name:      done.Name,
		UserId:    done.UserID.String(),
		CreatedAt: done.CreatedAt.Time,
		UpdatedAt: done.UpdatedAt.Time,
		Token:     nil,
	}, nil
}

func (r *ServerRepository) UpdateServerByUserIdAndServerId(ctx context.Context, userId string, serverId string, request contract.UpdateServerRequest) (*contract.ServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *ServerRepository) DeleteServerByUserIdAndServerId(ctx context.Context, userId string, serverId string) error {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return err
	}

	return nil
}
