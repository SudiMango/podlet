package service

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/cloudflare"
	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/contract"
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

func (s *CloudflareService) CreateDNSRecord(ctx context.Context, userId string, serverId string, request dto.CreateDNSRecordRequest) (*contract.DnsRecordResponse, error) {

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
