package service

import (
	"context"
	"errors"

	"github.com/SudiMango/podlet/internal/api/cloudflare"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/jackc/pgx/v5"
)

var (
	ErrCfTunnelNotFound = errors.New("Tunnel with given userId and serverId not found.")
)

type CloudflareService struct {
	cloudflareRepository *repository.CloudflareRepository
	cloudflareClient     *cloudflare.Client
}

func NewCloudflareService(cloudflareRepository *repository.CloudflareRepository, cloudflareClient *cloudflare.Client) *CloudflareService {
	return &CloudflareService{
		cloudflareRepository: cloudflareRepository,
		cloudflareClient:     cloudflareClient,
	}
}

/*
	Tunnels
*/

func (s *CloudflareService) CreateTunnel(ctx context.Context, userId string, serverId string) (*contract.TunnelResponse, error) {

	// Create tunnel from cf side
	cfResult, err := s.cloudflareClient.CreateTunnel(ctx)
	if err != nil {
		return nil, err
	}

	// Make db reference for tunnel
	result, err := s.cloudflareRepository.CreateTunnel(ctx, userId, serverId, *cfResult)
	if err != nil {
		return nil, err
	}

	// Get tunnel connector token
	token, err := s.cloudflareClient.GetTunnelConnectorToken(ctx, result.Id)
	if err != nil {
		return nil, err
	}
	result.ConnectorToken = token

	return result, nil
}

func (s *CloudflareService) GetTunnelByUserIdAndServerId(ctx context.Context, userId string, serverId string) (*contract.TunnelResponse, error) {
	result, err := s.cloudflareRepository.GetTunnelByUserIdAndServerId(ctx, userId, serverId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCfTunnelNotFound
		}

		return nil, err
	}

	token, err := s.cloudflareClient.GetTunnelConnectorToken(ctx, result.Id)
	if err != nil {
		return nil, err
	}
	result.ConnectorToken = token

	return result, nil
}

/*
	DNS Records
*/

func (s *CloudflareService) CreateDNSRecord(ctx context.Context, userId string, serverId string, request contract.CreateDnsRecordRequest) (*contract.DnsRecordResponse, error) {
	// Repoint dns record if it already exists
	existing, err := s.cloudflareRepository.GetDNSRecordByServerIdAndRepositoryId(ctx, userId, serverId, request.RepositoryId)
	if err == nil {
		return s.repointDnsRecord(ctx, existing, request)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Create new dns record if it doesn't exist
	return s.createNewDnsRecord(ctx, userId, serverId, request)
}

/*
	Private functions
*/

func (s *CloudflareService) repointDnsRecord(ctx context.Context, record *contract.DnsRecordResponse, request contract.CreateDnsRecordRequest) (*contract.DnsRecordResponse, error) {
	// Get current ingress rules for the tunnel
	configs, err := s.cloudflareClient.GetIngressConfigs(ctx, record.TunnelId)
	if err != nil {
		return nil, err
	}

	// Point the existing hostname's rule at the new service
	for i := range configs.Config.Ingress {
		rule := &configs.Config.Ingress[i]
		if rule.Hostname != nil && *rule.Hostname == record.Name {
			rule.Service = request.Service
			break
		}
	}

	// Push updated ingress rules to cf
	_, err = s.cloudflareClient.UpdateIngressConfig(
		ctx,
		record.TunnelId,
		cloudflare.PutIngressRequest{
			Config: configs.Config,
		},
	)
	if err != nil {
		return nil, err
	}

	// Update db reference with the new service
	if err := s.cloudflareRepository.UpdateDNSRecordService(ctx, record.Id, request.Service); err != nil {
		return nil, err
	}
	record.Service = request.Service

	return record, nil
}

func (s *CloudflareService) createNewDnsRecord(ctx context.Context, userId string, serverId string, request contract.CreateDnsRecordRequest) (*contract.DnsRecordResponse, error) {
	// Create dns record from cf side
	cfResult, err := s.cloudflareClient.CreateDNSRecord(ctx, request.TunnelId)
	if err != nil {
		return nil, err
	}

	request.Hostname = cfResult.Name

	// Update ingress records
	configs, err := s.cloudflareClient.GetIngressConfigs(ctx, request.TunnelId)
	if err != nil {
		return nil, err
	}

	configs.Config.Ingress = append(configs.Config.Ingress, contract.Ingress{
		Hostname: &request.Hostname,
		Service:  request.Service,
	})

	_, err = s.cloudflareClient.UpdateIngressConfig(
		ctx,
		request.TunnelId,
		cloudflare.PutIngressRequest{
			Config: configs.Config,
		},
	)
	if err != nil {
		return nil, err
	}

	request.DnsRecordId = cfResult.Id
	request.CreatedAt = cfResult.CreatedAt

	// Make db reference for dns record
	result, err := s.cloudflareRepository.CreateDNSRecord(ctx, userId, serverId, request)
	if err != nil {
		return nil, err
	}

	return result, nil
}
