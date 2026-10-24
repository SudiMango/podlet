package repository

import (
	"context"
	"time"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5"
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

/*
	Pending server functions
*/

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

	done, err := r.store.GetValidPendingServersByUserId(ctx, pgUserId)
	if err != nil {
		return nil, err
	}

	var response contract.PendingServersResponse
	for _, s := range done {
		sTemp := contract.PendingServerResponse{
			Id:                s.ID,
			UserId:            s.UserID.String(),
			CreatedAt:         s.CreatedAt.Time,
			ExpiresAt:         s.ExpiresAt.Time,
			RegistrationToken: nil,
		}
		response = append(response, sTemp)
	}

	return &response, nil
}

func (r *ServerRepository) GetValidPendingServerByUserIdAndPendingServerId(ctx context.Context, userId string, pendingServerId string) (*contract.PendingServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	done, err := r.store.GetValidPendingServerByUserIdAndPendingServerId(ctx, db.GetValidPendingServerByUserIdAndPendingServerIdParams{
		UserID: pgUserId,
		ID:     pendingServerId,
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
	}, nil
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

	rowsDeleted, err := r.store.DeletePendingServerByUserIdAndPendingServerId(ctx, db.DeletePendingServerByUserIdAndPendingServerIdParams{
		UserID: pgUserId,
		ID:     pendingServerId,
	})
	if err != nil {
		return err
	}
	if rowsDeleted == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

/*
	Server functions
*/

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
		Id:         done.ID,
		Name:       done.Name,
		UserId:     done.UserID.String(),
		CreatedAt:  done.CreatedAt.Time,
		UpdatedAt:  done.UpdatedAt.Time,
		LastSeenAt: nil,
		Token:      nil,
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
			Id:         s.ID,
			Name:       s.Name,
			UserId:     s.UserID.String(),
			CreatedAt:  s.CreatedAt.Time,
			UpdatedAt:  s.UpdatedAt.Time,
			LastSeenAt: &s.LastSeenAt.Time,
			Token:      nil,
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
		Id:         done.ID,
		Name:       done.Name,
		UserId:     done.UserID.String(),
		CreatedAt:  done.CreatedAt.Time,
		UpdatedAt:  done.UpdatedAt.Time,
		LastSeenAt: &done.LastSeenAt.Time,
		Token:      nil,
	}, nil
}

func (r *ServerRepository) GetServerByTokenHash(ctx context.Context, tokenHash string) (*contract.ServerResponse, error) {
	done, err := r.store.GetServerByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	return &contract.ServerResponse{
		Id:         done.ID,
		Name:       done.Name,
		UserId:     done.UserID.String(),
		CreatedAt:  done.CreatedAt.Time,
		UpdatedAt:  done.UpdatedAt.Time,
		LastSeenAt: &done.LastSeenAt.Time,
		Token:      nil,
	}, nil
}

func (r *ServerRepository) UpdateServerLastSeenAtByUserIdAndServerId(ctx context.Context, userId string, serverId string) (*dto.ServerLifecycleState, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	lifecycleState, err := r.store.UpdateServerLastSeenAtByUserIdAndServerId(ctx, db.UpdateServerLastSeenAtByUserIdAndServerIdParams{
		UserID: pgUserId,
		ID:     serverId,
		LastSeenAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if lifecycleState != "active" {
		state := dto.ServerLifecycleState(lifecycleState)
		return &state, nil
	}

	return nil, nil
}

func (r *ServerRepository) UpdateServerByUserIdAndServerId(ctx context.Context, userId string, serverId string, name string) (*contract.ServerResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	done, err := r.store.UpdateServerByUserIdAndServerId(ctx, db.UpdateServerByUserIdAndServerIdParams{
		UserID: pgUserId,
		ID:     serverId,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return &contract.ServerResponse{
		Id:         done.ID,
		Name:       done.Name,
		UserId:     done.UserID.String(),
		CreatedAt:  done.CreatedAt.Time,
		UpdatedAt:  done.UpdatedAt.Time,
		LastSeenAt: &done.LastSeenAt.Time,
		Token:      nil,
	}, nil
}

func (r *ServerRepository) DeleteServerByUserIdAndServerId(ctx context.Context, userId string, serverId string) error {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return err
	}

	rowsDeleted, err := r.store.DeleteServerByUserIdAndServerId(ctx, db.DeleteServerByUserIdAndServerIdParams{
		UserID: pgUserId,
		ID:     serverId,
	})
	if err != nil {
		return err
	}
	if rowsDeleted == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ServerRepository) MarkServerForDeletionByUserIdAndServerId(ctx context.Context, userId string, serverId string) error {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return err
	}

	rowsDeleted, err := r.store.MarkServerForDeletionByUserIdAndServerId(ctx, db.MarkServerForDeletionByUserIdAndServerIdParams{
		UserID:         pgUserId,
		ID:             serverId,
		LifecycleState: db.ServerLifecycleState(dto.ServerLifecyclePendingDeletion),
	})
	if err != nil {
		return err
	}
	if rowsDeleted == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
