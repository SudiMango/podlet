package repository

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/contract"
	db "github.com/SudiMango/podlet/internal/sql/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type CloudflareRepository struct {
	store *db.Queries
}

func NewCloudflareRepository(store *db.Queries) *CloudflareRepository {
	return &CloudflareRepository{
		store: store,
	}
}

/*
	Tunnels
*/

func (r *CloudflareRepository) CreateTunnel(ctx context.Context, userId string, serverId string, request contract.CreateTunnelResponse) (*contract.TunnelResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.CreateTunnel(ctx, db.CreateTunnelParams{
		ID:       request.Id,
		Name:     request.Name,
		UserID:   pgUserId,
		ServerID: serverId,
		CreatedAt: pgtype.Timestamptz{
			Time:  request.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contract.TunnelResponse{
		Id:        result.ID,
		Name:      result.Name,
		UserId:    result.UserID.String(),
		ServerId:  result.ServerID,
		CreatedAt: result.CreatedAt.Time,
	}, nil
}

/*
	DNS Records
*/

func (r *CloudflareRepository) CreateDNSRecord(ctx context.Context, userId string, serverId string, request dto.CreateDNSRecordRequest) (*contract.DnsRecordResponse, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}

	result, err := r.store.CreateDnsRecord(ctx, db.CreateDnsRecordParams{
		ID:       request.DnsRecordId,
		Name:     request.Hostname,
		Service:  request.Service,
		UserID:   pgUserId,
		ServerID: serverId,
		TunnelID: request.TunnelId,
		CreatedAt: pgtype.Timestamptz{
			Time:  request.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contract.DnsRecordResponse{
		Id:        result.ID,
		Name:      result.Name,
		Service:   result.Service,
		UserId:    result.UserID.String(),
		ServerId:  result.ServerID,
		TunnelId:  result.TunnelID,
		CreatedAt: request.CreatedAt,
	}, nil
}
